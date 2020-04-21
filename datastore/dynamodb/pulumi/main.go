package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/go/aws/dynamodb"
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

// DynamoConfig contains the key-value pairs for the configuration of Amazon DynamoDB in this stack
type DynamoConfig struct {
	// Controls how you are charged for read and write throughput and how you manage capacity
	BillingMode pulumi.String `json:"billingmode"`

	// The number of write units for this table
	WriteCapacity pulumi.Int `json:"writecapacity"`

	// The number of read units for this table
	ReadCapacity pulumi.Int `json:"readcapacity"`
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Read the configuration data from Pulumi.<stack>.yaml
		conf := config.New(ctx, "awsconfig")

		// Create a new Tags object with the data from the configuration
		var tags Tags
		conf.RequireObject("tags", &tags)

		// Create a new DynamoConfig object with the data from the configuration
		var dynamoConfig DynamoConfig
		conf.RequireObject("dynamodb", &dynamoConfig)

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

		// The table attributes represent a list of attributes that describe the key schema for the table and indexes
		tableAttributeInput := []dynamodb.TableAttributeInput{
			dynamodb.TableAttributeArgs{
				Name: pulumi.String("PK"),
				Type: pulumi.String("S"),
			}, dynamodb.TableAttributeArgs{
				Name: pulumi.String("SK"),
				Type: pulumi.String("S"),
			},
		}

		// The set of arguments for constructing an Amazon DynamoDB Table resource
		tableArgs := &dynamodb.TableArgs{
			Attributes:    dynamodb.TableAttributeArray(tableAttributeInput),
			BillingMode:   pulumi.StringPtrInput(dynamoConfig.BillingMode),
			HashKey:       pulumi.String("PK"),
			RangeKey:      pulumi.String("SK"),
			Tags:          pulumi.Map(tagMap),
			Name:          pulumi.String(fmt.Sprintf("%s-%s", ctx.Stack(), ctx.Project())),
			ReadCapacity:  dynamoConfig.ReadCapacity,
			WriteCapacity: dynamoConfig.WriteCapacity,
		}

		// NewTable registers a new resource with the given unique name, arguments, and options
		table, err := dynamodb.NewTable(ctx, fmt.Sprintf("%s-%s", ctx.Stack(), ctx.Project()), tableArgs)
		if err != nil {
			return err
		}

		// Export the ARN and Name of the table
		ctx.Export("Table::Arn", table.Arn)
		ctx.Export("Table::Name", table.Name)

		return nil
	})
}
