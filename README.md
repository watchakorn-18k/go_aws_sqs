# Go AWS SQS

This project demonstrates how to use AWS Simple Queue Service (SQS) with a Go application, structured into two main components: **producer** and **consumer**. Each component is contained in its respective folder, performing distinct roles in the messaging workflow.

## Project Structure

```
/go-aws-sqs
|-- /producer
|   |-- main.go
|   |-- go.mod
|   |-- .env
|   |-- ...
|-- /consumer
|   |-- main.go
|   |-- go.mod
|   |-- .env
|   |-- ...
```

### Producer

The `producer` folder contains the application responsible for receiving API requests from clients and sending messages to an SQS queue.

- **Functionality**:
  - The producer listens for incoming HTTP requests (using the Fiber framework).
  - When a request is received, it extracts the necessary data and sends it as a message to the specified SQS queue.
  - It ensures proper error handling and logs the status of message sending.

### Consumer

The `consumer` folder contains the application that processes messages from the SQS queue.

- **Functionality**:
  - The consumer continuously polls the SQS queue for new messages.
  - Upon receiving a message, it processes the data (e.g., inserts it into a MongoDB database).
  - After successful processing, the consumer deletes the message from the queue to prevent reprocessing.

## Prerequisites

- Go installed on your machine
- AWS account with SQS service enabled
- AWS credentials configured (using IAM roles or AWS CLI)

## Environment Variables
`/go-aws-sqs/producer/.env`
```env
AWS_REGION=ap-southeast-1
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
MONGODB_URI=
QUEUE_URL=
```

`/go-aws-sqs/consumer/.env`
```env
AWS_REGION=ap-southeast-1
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
QUEUE_URL=
```

## How to Run

1. **Producer**:
   - Navigate to the `producer` directory and run the application:
     ```bash
     cd producer
     go mod tidy
     go run main.go
     ```

2. **Consumer**:
   - Open another terminal window, navigate to the `consumer` directory, and run the application:
     ```bash
     cd consumer
     go mod tidy
     go run main.go
     ```

## Conclusion

This project provides a basic implementation of a producer-consumer pattern using AWS SQS in Go. You can extend the functionality by adding more complex processing logic, error handling, and integration with other AWS services.
```