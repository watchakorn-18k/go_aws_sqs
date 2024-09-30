package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var sqsClient *sqs.SQS
var queueURL string
var mongoClient *mongo.Client

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("err load ENV")
	}
	// AWS SQS configuration
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		return
	}
	sqsClient = sqs.New(sess)
	queueURL = os.Getenv("QUEUE_URL") // URL ของคิว SQS

	// MongoDB configuration
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	mongoClient = client

	// Worker loop
	fmt.Println("Starting worker loop...")
	for {
		receiveAndProcessMessages()
		fmt.Println("Waiting for new messages...")
		time.Sleep(1 * time.Second) // ตรวจสอบคิวทุกๆ 10 วินาที
	}
}

func receiveAndProcessMessages() {
	result, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(10), // จำนวนข้อความสูงสุดที่ต้องการรับในครั้งเดียว
		WaitTimeSeconds:     aws.Int64(5),  // ระยะเวลารอข้อความใหม่
	})

	if err != nil {
		log.Printf("Failed to receive messages: %v", err)
		return
	}
	for _, message := range result.Messages {
		fmt.Println("Start Message ID:", *message.MessageId)
		processMessage(*message.Body)

		// ลบข้อความออกจากคิวหลังจากดำเนินการเสร็จสิ้น
		_, err = sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queueURL),
			ReceiptHandle: message.ReceiptHandle,
		})

		if err != nil {
			log.Printf("Failed to delete message from SQS: %v", err)
		}
	}
}

func processMessage(messageBody string) {
	var user User
	if err := json.Unmarshal([]byte(messageBody), &user); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	// Insert data into MongoDB
	collection := mongoClient.Database("quueu_test").Collection("users")
	_, err := collection.InsertOne(context.Background(), bson.M{
		"name":  user.Name,
		"email": user.Email,
	})

	if err != nil {
		log.Printf("Failed to insert data into MongoDB: %v", err)
	} else {
		log.Printf("Successfully inserted data: %v", user)
	}
}
