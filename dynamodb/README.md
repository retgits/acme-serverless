# DynamoDB

The ACME Serverless Fitness Shop leverages [Amazon DynamoDB](https://aws.amazon.com/dynamodb/) to store data.

## Create the table

To create the table you'll need a [Pulumi account](https://app.pulumi.com/signup). Once you have your Pulumi account and configured the [Pulumi CLI](https://www.pulumi.com/docs/get-started/aws/install-pulumi/), you can initialize a new stack using the Pulumi templates in the [pulumi](./pulumi) folder.

```bash
cd pulumi
pulumi stack init <your pulumi org>/acmeserverless-dynamodb/dev
```

To create the Pulumi stack, and create the Amazon DynamoDB table, run `pulumi up`

Pulumi is configured using a file called `Pulumi.dev.yaml`. A sample configuration is available in the Pulumi directory. You can rename [`Pulumi.dev.yaml.sample`](./pulumi/Pulumi.dev.yaml.sample) to `Pulumi.dev.yaml` and update the variables accordingly. Alternatively, you can change variables directly in the [main.go](./pulumi/main.go) file in the pulumi directory.

If you want to keep track of the resources in Pulumi, you can add tags to your stack as well.

```bash
pulumi stack tag set app:name acmeserverless
pulumi stack tag set app:feature acmeserverless-dynamodb
pulumi stack tag set app:domain infra
```

## Seed the table

To seed the DynamoDB table with random data, you can use the Go app in the [seed](./seed) directory. The app has two required flags and one optional one:

* `region`: The region to send requests to (required)
* `table`: The Amazon DynamoDB table to use (required, the name is part of the output shown by `pulumi up`)
* `endpoint`: An optional endpoint URL (optional, hostname only or fully qualified URI in case you're using [DynamoDB Local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html))

As an example, using the default settings, you can run

```bash
cd seed
go run main.go -region=us-west-2 -table=dev-acmeserverless-dynamodb
```

To generate your own data, you can use [Mockaroo](https://www.mockaroo.com/) and import the `schema.json` files to start off.
