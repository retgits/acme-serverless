# EventBridge

The ACME Serverless Fitness Shop leverages [Amazon EventBridge](https://aws.amazon.com/eventbridge/) as a serverless event bus to connect microservices together.

## Create the eventbus

To create the eventbus you'll need the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) installed and configured. After that you can run

```bash
aws events create-event-bus --name acmeserverless
```

## Sending events to test

To send a message to the EventBus using the test data from the files in any of the subfolders, you can use the Go app in the test folder. The app has three required flags and one optional one:

* `event`: The name of event to send (required)
* `bus`: The name of the Amazon EventBridge bus to send events to (required)
* `service`: The ACME Serverless Fitness Shop service to send events to (required)
* `region`: The region to send requests to (optional, defaults to us-west-2)

```bash
go run main.go -event=<any of the files existing in the folder of the specific service> -bus=<name of the custom bus> -service=<name of the service>
```
