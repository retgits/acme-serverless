package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

var (
	eventType    string
	fileLocation string
	eventBus     string
	region       string
	service      string
)

func main() {
	flag.StringVar(&eventType, "event", "", "The name of event to send (required)")
	flag.StringVar(&eventBus, "bus", "", "The name of the Amazon EventBridge bus to send events to (required)")
	flag.StringVar(&service, "service", "", "The ACME Serverless Fitness Shop service to send events to (required)")
	flag.StringVar(&fileLocation, "location", "", "The folder containing the JSON files (required)")
	flag.StringVar(&region, "region", "us-west-2", "The region to send requests to (defaults to us-west-2)")
	flag.Parse()

	// Make sure the required flags are set
	if len(eventType) < 1 {
		log.Fatal("Error: the 'event' flag must be set")
	}

	if len(eventBus) < 1 {
		log.Fatal("Error: the 'bus' flag must be set")
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

	svc := eventbridge.New(awsSession)

	entries := make([]*eventbridge.PutEventsRequestEntry, 1)
	entries[0] = &eventbridge.PutEventsRequestEntry{
		Detail:       aws.String(payload),
		DetailType:   aws.String("myDetailType"),
		EventBusName: aws.String(eventBus),
		Resources:    []*string{aws.String("TestMessage"), aws.String(eventBus)},
		Source:       aws.String("cli"),
	}
	event := &eventbridge.PutEventsInput{
		Entries: entries,
	}

	output, err := svc.PutEvents(event)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range output.Entries {
		log.Printf("EventID: %s", *entry.EventId)
	}
}
