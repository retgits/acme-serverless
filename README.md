# ACME Serverless Fitness Shop

> Serverless and Fitness, because combining two amazing things can only lead to more amazing things

## Getting Started

These instructions will allow you to run entire ACME Serverless Fitness Shop

The ACME Serverless Fitness Shop contains seven different domains of service:

* [Shipment](https://github.com/retgits/acme-serverless-shipment)
* [Payment](https://github.com/retgits/acme-serverless-payment)
* [Order](https://github.com/retgits/acme-serverless-order)
* [Cart](https://github.com/retgits/acme-serverless-cart)
* [Catalog](https://github.com/retgits/acme-serverless-catalog)
* [User](https://github.com/retgits/acme-serverless-user)
* [Point-of-Sales](https://github.com/retgits/acme-serverless-pos)

To get started you'll need:

* [Go (at least Go 1.12)](https://golang.org/dl/);
* [An AWS Account](https://portal.aws.amazon.com/billing/signup);
* The _vuln_ targets for Make and Mage rely on the [Snyk](http://snyk.io/) CLI.
* The services use [Sentry.io](https://sentry.io) for tracing and error reporting

## Supported AWS Services

### Data Store

* [Amazon DynamoDB](https://aws.amazon.com/dynamodb/): You can use the makefile in the [dynamodb](./dynamodb) folder to create the DynamoDB table. The command you need to run is `make -f Makefile.dynamodb deploy`.

To start your journey off with random data, you can use the [`Makefile.dynamodb`](./dynamodb/Makefile.dynamodb) as well. The `seed` target will add seed data (from the various `data.json` files) into Amazon DynamoDB. To generate your own data, you can use [Mockaroo](https://www.mockaroo.com/) and import the `schema.json` files to start off.

### Eventing

* [Amazon EventBridge](https://aws.amazon.com/eventbridge/): Each of the domains that supports Amazon EventBridge will have instructions and sources how to run the ACME Serverless Fitness Shop using Amazon EventBridge;
* [Amazon Simple Queue Service](https://aws.amazon.com/sqs/): Each of the domains that supports Amazon Simple Queue Service will have instructions and sources how to run the ACME Serverless Fitness Shop using Amazon Simple Queue Service;
* [Amazon API Gateway](https://aws.amazon.com/api-gateway/) _(only when there are public APIs available)_

#### Prerequisites for EventBridge

* [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) installed and configured
* [Custom EventBus](https://docs.aws.amazon.com/eventbridge/latest/userguide/create-event-bus.html) configured, the name of the configured event bus should be set as the `feature` parameter in the `template.yaml` file of the service you want to deploy.

#### Build and deploy for EventBridge

Clone this repository

```bash
git clone https://github.com/retgits/acme-serverless
cd acme-serverless
```

Change directories to the [deploy/cloudformation](./deploy/cloudformation) folder of the service you want to deploy

```bash
cd ./deploy/cloudformation/<service>
## like cd ./deploy/cloudformation/payment
```

Download the sources of the service you want to deploy

```bash
make -f Makefile.lambda get
```

If your event bus is not called _acmeserverless_, update the name of the `feature` parameter in the `template.yaml` file. Now you can build and deploy the Lambda function:

```bash
make -f Makefile.lambda build TYPE=eventbridge
make -f Makefile.lambda deploy TYPE=eventbridge
```

#### Testing with EventBridge

You can test the function from the [AWS Lambda Console](https://console.aws.amazon.com/lambda/home) using the test data from the files in [eventbridge](./eventbridge/). To send a message to the event bus, you can use the Go app in `./eventbridge` and run

```bash
go run main.go -event=<any of the files existing in the folder of the specific service> -location=<location on disk of the eventbridge folder> -bus=<name of the custom bus> -service=<name of the service>
```

#### Prerequisites for SQS

* [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) installed and configured

#### Build and deploy for SQS

Clone this repository

```bash
git clone https://github.com/retgits/acme-serverless
cd acme-serverless
```

Change directories to the [deploy/cloudformation](./deploy/cloudformation) folder of the service you want to deploy

```bash
cd ./deploy/cloudformation/<service>
## like cd ./deploy/cloudformation/payment
```

Download the sources of the service you want to deploy

```bash
make -f Makefile.lambda get
```

Now you can build and deploy the Lambda function:

```bash
make -f Makefile.lambda build TYPE=sqs
make -f Makefile.lambda deploy TYPE=sqs
```

#### Testing with SQS

To send a message to an SQS queue using the test data from the files in [sqs](./sqs/), you can use the Go app in `./sqs` and run

```bash
go run main.go -event=<any of the files existing in the folder of the specific service> -location=<location on disk of the sqs folder> -queue=<url of the sqs queue> -service=<name of the service>
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

### Hosting

* The Point-of-Sales app can be hosted on [Amazon S3](https://aws.amazon.com/s3).

## Using Make

The Makefiles for the services have a few a bunch of options available:

| Target  | Description                                                |
|---------|------------------------------------------------------------|
| build   | Build the executable for Lambda                            |
| get     | Performs a git clone to get the sources for the service    |
| clean   | Remove all generated files                                 |
| deploy  | Deploy the app to AWS Lambda                               |
| destroy | Deletes the CloudFormation stack and all created resources |
| help    | Displays the help for each target (this message)           |
| vuln    | Scans the Go.mod file for known vulnerabilities using Snyk |

The targets `build` and `deploy` need a variable **`TYPE`** set to either `eventbridge` or `sqs` to build and deploy the correct Lambda functions.

## Using Mage

If you want to "go all Go" (_pun intended_) and write plain-old go functions to build and deploy, you can use [Mage](https://magefile.org/) (which still leverages the CloudFormation templates). Mage is a make/rake-like build tool using Go so Mage automatically uses the functions you create as Makefile-like runnable targets.

### Prerequisites for Mage

To use Mage, you'll need to install it first:

```bash
go get -u -d github.com/magefile/mage
cd $GOPATH/src/github.com/magefile/mage
go run bootstrap.go
```

Instructions curtesy of Mage

### Targets

The Magefile in this repository has a bunch of targets available:

| Target | Description                                                                                              |
|--------|----------------------------------------------------------------------------------------------------------|
| build  | compiles the individual commands in the cmd folder, along with their dependencies.                       |
| clean  | removes object files from package source directories.                                                    |
| deploy | packages, deploys, and returns all outputs of your stack.                                                |
| deps   | resolves and downloads dependencies to the current development module and then builds and installs them. |
| get    | performs a git clone of the source code from GitHub for the service specified.                           |
| test   | 'Go test' automates testing the packages named by the import paths.                                      |
| vuln   | uses Snyk to test for any known vulnerabilities in go.mod.                                               |

Mage relies on a few environment variables to complete the work:

* `STAGE`: The stage to deploy to (defaults to dev)
* `SERVICE`: The service to deploy (defaults to payment)
* `TYPE`: The service type to deploy (defaults to sqs)
* `AUTHOR`: The author of the project (defaults to retgits)
* `TEAM`: The name of the team (defaults to vcs)
* `AWS_S3_BUCKET`: The Amazon S3 bucket to upload files to (defaults to myS3Bucket)

## Overview

![architecture](./overview-sqs.png)

The diagram shows how the services in the different domains work together. The architecture above shows the Amazon Simple Queue Service deployment option

## Contributing

Pull requests are welcome in their individual repositories. For major changes or questions, please open [an issue](https://github.com/retgits/acme-serverless/issues) first to discuss what you would like to change.

## License

See the [LICENSE](./LICENSE) file in the repository.
