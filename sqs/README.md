# EventBridge

The ACME Serverless Fitness Shop leverages Amazon Simple Queue Service (SQS) as a queueing service between microservices.

## Create the queues

To create the table you'll need a [Pulumi account](https://app.pulumi.com/signup). Once you have your Pulumi account and configured the [Pulumi CLI](https://www.pulumi.com/docs/get-started/aws/install-pulumi/), you can initialize a new stack using the Pulumi templates in the [pulumi](./pulumi) folder.

```bash
cd pulumi
pulumi stack init <your pulumi org>/acmeserverless-sqs/dev
```

To create the Pulumi stack, and create the Amazon SQS queues, run `pulumi up`

If you want to keep track of the resources in Pulumi, you can add tags to your stack as well.

```bash
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature acmeserverless-sqs
pulumi stack tag set app:domain infra
```

## Sending events to test

To send a message to an SQS queue using the test data from the files in the [test](./test/) folder, you can use the Go app in the test folder. The app has three required flags and one optional one:

* `event`: The name of event to send (required)
* `queue`: The URL of the Amazon SQS queue to send events to (required)
* `service`: The ACME Serverless Fitness Shop service to send events to (required)
* `region`: The region to send requests to (optional, defaults to us-west-2)

```bash
go run main.go -event=<any of the files existing in the folder of the specific service> -queue=<url of the sqs queue> -service=<name of the service>
```

If you want to test from the [AWS Lambda Console](https://console.aws.amazon.com/lambda/home), you'll have to wrap the test data in a SQS record envelop:

```json
{
  "Records": [
    {
      "messageId": "19dd0b57-b21e-4ac1-bd88-01bbb068cb78",
      "receiptHandle": "MessageReceiptHandle",
      "body": "", // This is where the data, an escaped JSON string, should be pasted
      "attributes": {
        "ApproximateReceiveCount": "1",
        "SentTimestamp": "1523232000000",
        "SenderId": "123456789012",
        "ApproximateFirstReceiveTimestamp": "1523232000001"
      },
      "messageAttributes": {},
      "md5OfBody": "7b270e59b47ff90a553787216d55d91d",
      "eventSource": "aws:sqs",
      "eventSourceARN": "arn:aws:sqs:us-east-1:123456789012:MyQueue",
      "awsRegion": "us-east-1"
    }
  ]
}
```
