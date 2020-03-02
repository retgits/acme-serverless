package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	eventType    string
	fileLocation string
	sqsQueue     string
	region       string
	service      string
)

func main() {
	flag.StringVar(&eventType, "event", "", "The name of event to send (required)")
	flag.StringVar(&sqsQueue, "queue", "", "The URL of the Amazon SQS queue to send events to (required)")
	flag.StringVar(&service, "service", "", "The ACME Serverless Fitness Shop service to send events to (required)")
	flag.StringVar(&fileLocation, "location", "", "The folder containing the JSON files (required)")
	flag.StringVar(&region, "region", "us-west-2", "The region to send requests to (defaults to us-west-2)")
	flag.Parse()

	// Make sure the required flags are set
	if len(eventType) < 1 {
		log.Fatal("Error: the 'event' flag must be set")
	}

	if len(sqsQueue) < 1 {
		log.Fatal("Error: the 'queue' flag must be set")
	}

	if len(service) < 1 {
		log.Fatal("Error: the 'service' flag must be set")
	}

	if len(fileLocation) < 1 {
		log.Fatal("Error: the 'location' flag must be set")
	}

	bytes, err := ioutil.ReadFile(path.Join(fileLocation, service, fmt.Sprintf("%s.json", eventType)))
	if err != nil {
		log.Fatalf("error reading file: %s", err.Error())
	}
	payload := string(bytes)

	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	svc := sqs.New(awsSession)

	sendMessageInput := &sqs.SendMessageInput{
		QueueUrl:    aws.String(sqsQueue),
		MessageBody: aws.String(payload),
	}

	sendMessageOutput, err := svc.SendMessage(sendMessageInput)
	if err != nil {
		log.Fatalf("error reading file: %s", err.Error())
	}

	log.Printf("MessageID: %s\n", *sendMessageOutput.MessageId)
}
