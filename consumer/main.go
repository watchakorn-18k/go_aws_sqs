// consumer/main.go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

var sqsClient *sqs.SQS
var queueURL string

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("err load ENV")
	}
	app := fiber.New()
	app.Use(logger.New())

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

	app.Post("/api/insert", insertToQueue)

	log.Fatal(app.Listen(":3000"))
}

// insertToQueue ส่งข้อมูลไปยัง AWS SQS
func insertToQueue(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid request")
	}
	data["name"] = fmt.Sprintf("%v_%v", data["name"].(string), rand.Intn(100000))
	// ส่งข้อมูลไปยัง SQS
	messageBody := `{"name": "` + data["name"].(string) + `", "email": "` + data["email"].(string) + `"}`

	res, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody:    aws.String(messageBody),
		QueueUrl:       aws.String(queueURL),
		MessageGroupId: aws.String("insert_user"), // กำหนด MessageGroupId ที่ไม่ซ้ำกัน
	})

	if err != nil {
		log.Printf("Error sending message to SQS: %v", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to enqueue")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Data inserted in queue successfully",
		"data":    map[string]interface{}{"message_id": res.MessageId},
	})
}
