package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
	"github.com/pulumi/pulumi/sdk/go/pulumi/config"
)

// Tags are key-value pairs to apply to the resources created by this stack
type Tags struct {
	// Author is the person who created the code, or performed the deployment
	Author pulumi.String

	// Feature is the project that this resource belongs to
	Feature pulumi.String

	// Team is the team that is responsible to manage this resource
	Team pulumi.String

	// Version is the version of the code for this resource
	Version pulumi.String
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Read the configuration data from Pulumi.<stack>.yaml
		conf := config.New(ctx, "awsconfig")

		// Create a new Tags object with the data from the configuration
		var tags Tags
		conf.RequireObject("tags", &tags)

		// Create a map[string]pulumi.Input of the tags
		// the first four tags come from the configuration file
		// the last two are derived from this deployment
		tagMap := make(map[string]pulumi.Input)
		tagMap["Author"] = tags.Author
		tagMap["Feature"] = tags.Feature
		tagMap["Team"] = tags.Team
		tagMap["Version"] = tags.Version
		tagMap["ManagedBy"] = pulumi.String("Pulumi")
		tagMap["Stage"] = pulumi.String(ctx.Stack())

		// Create the Payment Error Queue
		paymentErrQueue, err := sqs.NewQueue(ctx, fmt.Sprintf("%s-%s-payment-error", ctx.Stack(), ctx.Project()), &sqs.QueueArgs{
			MessageRetentionSeconds:  pulumi.Int(120),
			Name:                     pulumi.String(fmt.Sprintf("%s-%s-payment-error", ctx.Stack(), ctx.Project())),
			VisibilityTimeoutSeconds: pulumi.Int(30),
			Tags:                     pulumi.Map(tagMap),
		})
		if err != nil {
			return err
		}

		// The redrive policy for the other Payment Request and Payment Response are based on the
		// ARN of the Payment Error Queue
		paymentRedrivePolicy := paymentErrQueue.Arn.ApplyString(func(name string) string {
			return fmt.Sprintf("{\"maxReceiveCount\":1,\"deadLetterTargetArn\":\"%s\"}", name)
		})

		// Create the Payment Request Queue
		paymentRequestQueue, err := sqs.NewQueue(ctx, fmt.Sprintf("%s-%s-payment-request", ctx.Stack(), ctx.Project()), &sqs.QueueArgs{
			MessageRetentionSeconds:  pulumi.Int(120),
			Name:                     pulumi.String(fmt.Sprintf("%s-%s-payment-request", ctx.Stack(), ctx.Project())),
			VisibilityTimeoutSeconds: pulumi.Int(30),
			Tags:                     pulumi.Map(tagMap),
			RedrivePolicy:            paymentRedrivePolicy,
		})
		if err != nil {
			return err
		}

		// Create the Payment Response Queue
		paymentResponseQueue, err := sqs.NewQueue(ctx, fmt.Sprintf("%s-%s-payment-response", ctx.Stack(), ctx.Project()), &sqs.QueueArgs{
			MessageRetentionSeconds:  pulumi.Int(120),
			Name:                     pulumi.String(fmt.Sprintf("%s-%s-payment-response", ctx.Stack(), ctx.Project())),
			VisibilityTimeoutSeconds: pulumi.Int(30),
			Tags:                     pulumi.Map(tagMap),
			RedrivePolicy:            paymentRedrivePolicy,
		})
		if err != nil {
			return err
		}

		// Create the Shipment Error Queue
		shipmentErrQueue, err := sqs.NewQueue(ctx, fmt.Sprintf("%s-%s-shipment-error", ctx.Stack(), ctx.Project()), &sqs.QueueArgs{
			MessageRetentionSeconds:  pulumi.Int(120),
			Name:                     pulumi.String(fmt.Sprintf("%s-%s-shipment-error", ctx.Stack(), ctx.Project())),
			VisibilityTimeoutSeconds: pulumi.Int(30),
			Tags:                     pulumi.Map(tagMap),
		})
		if err != nil {
			return err
		}

		// The redrive policy for the other Shipment Request and Shipment Response are based on the
		// ARN of the Shipment Error Queue
		shipmentRedrivePolicy := shipmentErrQueue.Arn.ApplyString(func(name string) string {
			return fmt.Sprintf("{\"maxReceiveCount\":1,\"deadLetterTargetArn\":\"%s\"}", name)
		})

		// Create the Shipment Request Queue
		shipmentRequestQueue, err := sqs.NewQueue(ctx, fmt.Sprintf("%s-%s-shipment-request", ctx.Stack(), ctx.Project()), &sqs.QueueArgs{
			MessageRetentionSeconds:  pulumi.Int(120),
			Name:                     pulumi.String(fmt.Sprintf("%s-%s-shipment-request", ctx.Stack(), ctx.Project())),
			VisibilityTimeoutSeconds: pulumi.Int(30),
			Tags:                     pulumi.Map(tagMap),
			RedrivePolicy:            shipmentRedrivePolicy,
		})
		if err != nil {
			return err
		}

		// Create the Shipment Response Queue
		shipmentResponseQueue, err := sqs.NewQueue(ctx, fmt.Sprintf("%s-%s-shipment-response", ctx.Stack(), ctx.Project()), &sqs.QueueArgs{
			MessageRetentionSeconds:  pulumi.Int(120),
			Name:                     pulumi.String(fmt.Sprintf("%s-%s-shipment-response", ctx.Stack(), ctx.Project())),
			VisibilityTimeoutSeconds: pulumi.Int(30),
			Tags:                     pulumi.Map(tagMap),
			RedrivePolicy:            shipmentRedrivePolicy,
		})
		if err != nil {
			return err
		}

		// Export the ARNs and Names of the queues
		ctx.Export("PaymentErrorQueue::Arn", paymentErrQueue.Arn)
		ctx.Export("PaymentErrorQueue::Name", paymentErrQueue.Name)
		ctx.Export("PaymentRequestQueue::Arn", paymentRequestQueue.Arn)
		ctx.Export("PaymentRequestQueue::Name", paymentRequestQueue.Name)
		ctx.Export("PaymentResponseQueue::Arn", paymentResponseQueue.Arn)
		ctx.Export("PaymentResponseQueue::Name", paymentResponseQueue.Name)
		ctx.Export("ShipmentErrorQueue::Arn", shipmentErrQueue.Arn)
		ctx.Export("ShipmentErrorQueue::Name", shipmentErrQueue.Name)
		ctx.Export("ShipmentRequestQueue::Arn", shipmentRequestQueue.Arn)
		ctx.Export("ShipmentRequestQueue::Name", shipmentRequestQueue.Name)
		ctx.Export("ShipmentResponseQueue::Arn", shipmentResponseQueue.Arn)
		ctx.Export("ShipmentResponseQueue::Name", shipmentResponseQueue.Name)

		return nil
	})
}
