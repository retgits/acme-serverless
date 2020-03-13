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
* [A Pulumi Account](https://app.pulumi.com/signup) for deployments if you choose SQS for communication;
* The services use [Sentry.io](https://sentry.io) for tracing and error reporting

## Supported AWS Services

### Data Store

The ACME Serverless Fitness Shop needs to store data. In the list below you can find the supported data store services with a link to the deployment instructions.

* [Amazon DynamoDB](./dynamodb)

### Eventing

Wherever possible, the ACME Serverless Fitness Shop uses event-driven communication. In the list below you can find the supported eventing solutions with a link to the deployment instructions.

* [Amazon EventBridge](./eventbridge)
* [Amazon Simple Queue Service](./sqs)

### APIs

All APIs will be accessible through Amazon API Gateway

### Hosting

The Point-of-Sales app can be hosted on [Amazon S3](https://aws.amazon.com/s3).

## Overview

![architecture](./overview-sqs.png)

The diagram shows how the services in the different domains work together. The architecture above shows the Amazon Simple Queue Service deployment option

## Contributing

Pull requests are welcome in their individual repositories. For major changes or questions, please open [an issue](https://github.com/retgits/acme-serverless/issues) first to discuss what you would like to change.

## License

See the [LICENSE](./LICENSE) file in the repository.
