package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gofrs/uuid"
	acmeserverless "github.com/retgits/acme-serverless"
)

var (
	region    string
	table     string
	dynamoURL string
)

// Create a single instance of the dynamoDB service which can be reused if the container stays warm.
var dbs *dynamodb.DynamoDB

// initialize creates the connection to dynamoDB. If the environment variable DYNAMO_URL is set, the
// connection is made to that URL instead of relying on the AWS SDK to provide the URL.
func initialize() {
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	if len(dynamoURL) > 0 {
		awsSession.Config.Endpoint = aws.String(dynamoURL)
	}

	dbs = dynamodb.New(awsSession)

	return
}

// AddUser stores a new user in Amazon DynamoDB
func AddUser(user acmeserverless.User) error {
	// Create a JSON encoded string of the user
	payload, err := user.Marshal()
	if err != nil {
		return err
	}

	// Create a map of DynamoDB Attribute Values containing the table keys
	km := make(map[string]*dynamodb.AttributeValue)
	km["PK"] = &dynamodb.AttributeValue{
		S: aws.String("USER"),
	}
	km["SK"] = &dynamodb.AttributeValue{
		S: aws.String(user.ID),
	}

	// Create a map of DynamoDB Attribute Values containing the table data elements
	em := make(map[string]*dynamodb.AttributeValue)
	em[":keyid"] = &dynamodb.AttributeValue{
		S: aws.String(user.Username),
	}
	em[":payload"] = &dynamodb.AttributeValue{
		S: aws.String(string(payload)),
	}

	uii := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(table),
		Key:                       km,
		ExpressionAttributeValues: em,
		UpdateExpression:          aws.String("SET Payload = :payload, KeyID = :keyid"),
	}

	_, err = dbs.UpdateItem(uii)
	if err != nil {
		return err
	}

	return nil
}

// AddCatalogItem stores a new product in Amazon DynamoDB
func AddCatalogItem(product acmeserverless.CatalogItem) error {
	// Marshal the newly updated product struct
	payload, err := product.Marshal()
	if err != nil {
		return err
	}

	// Create a map of DynamoDB Attribute Values containing the table keys
	km := make(map[string]*dynamodb.AttributeValue)
	km["PK"] = &dynamodb.AttributeValue{
		S: aws.String("PRODUCT"),
	}
	km["SK"] = &dynamodb.AttributeValue{
		S: aws.String(product.ID),
	}

	// Create a map of DynamoDB Attribute Values containing the table data elements
	em := make(map[string]*dynamodb.AttributeValue)
	em[":payload"] = &dynamodb.AttributeValue{
		S: aws.String(string(payload)),
	}

	uii := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(table),
		Key:                       km,
		ExpressionAttributeValues: em,
		UpdateExpression:          aws.String("SET Payload = :payload"),
	}

	_, err = dbs.UpdateItem(uii)
	if err != nil {
		return err
	}

	return nil
}

// AddOrder stores a new order in Amazon DynamoDB
func AddOrder(order acmeserverless.Order) error {
	// Generate and assign a new orderID
	o.OrderID = uuid.Must(uuid.NewV4()).String()
	o.Status = aws.String("Pending Payment")

	// Marshal the newly updated product struct
	payload, err := order.Marshal()
	if err != nil {
		return fmt.Errorf("error marshalling order: %s", err.Error())
	}

	// Create a map of DynamoDB Attribute Values containing the table keys
	km := make(map[string]*dynamodb.AttributeValue)
	km["PK"] = &dynamodb.AttributeValue{
		S: aws.String("ORDER"),
	}
	km["SK"] = &dynamodb.AttributeValue{
		S: aws.String(order.OrderID),
	}

	// Create a map of DynamoDB Attribute Values containing the table data elements
	em := make(map[string]*dynamodb.AttributeValue)
	em[":keyid"] = &dynamodb.AttributeValue{
		S: aws.String(order.UserID),
	}
	em[":payload"] = &dynamodb.AttributeValue{
		S: aws.String(string(payload)),
	}

	uii := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(table),
		Key:                       km,
		ExpressionAttributeValues: em,
		UpdateExpression:          aws.String("SET Payload = :payload, KeyID = :keyid"),
	}

	_, err = dbs.UpdateItem(uii)
	if err != nil {
		return fmt.Errorf("error updating dynamodb: %s", err.Error())
	}

	return nil
}

// StoreItems saves the cart items from a single user into Amazon DynamoDB
func StoreItems(userID string, item acmeserverless.CartItem) error {
	payload, err := item.Marshal()
	if err != nil {
		return err
	}

	// Create a map of DynamoDB Attribute Values containing the table keys
	// for the access pattern PK = CART SK = ID
	km := make(map[string]*dynamodb.AttributeValue)
	km["PK"] = &dynamodb.AttributeValue{
		S: aws.String("CART"),
	}
	km["SK"] = &dynamodb.AttributeValue{
		S: aws.String(userID),
	}

	em := make(map[string]*dynamodb.AttributeValue)
	em[":payload"] = &dynamodb.AttributeValue{
		S: aws.String(string(payload)),
	}

	uii := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(table),
		Key:                       km,
		ExpressionAttributeValues: em,
		UpdateExpression:          aws.String("SET Payload = :payload"),
	}

	_, err = dbs.UpdateItem(uii)
	return err
}

func main() {
	// Read flags
	flag.StringVar(&region, "region", "", "The region to send requests to (required)")
	flag.StringVar(&table, "table", "", "The Amazon DynamoDB table to use (required)")
	flag.StringVar(&dynamoURL, "endpoint", "", "An optional endpoint URL (optional, hostname only or fully qualified URI)")
	flag.Parse()

	// Make sure region and table are set
	if len(region) < 1 {
		log.Fatal("Error: the 'region' flag must be set")
	}

	if len(table) < 1 {
		log.Fatal("Error: the 'table' flag must be set")
	}

	// Initialize the database connection
	initialize()

	// Read all files.
	// if any of the files are not read successfully the function panics
	userData, err := ioutil.ReadFile("./user-data.json")
	catalogData, err = ioutil.ReadFile("./catalog-data.json")
	orderData, err = ioutil.ReadFile("./order-data.json")
	cartData, err = ioutil.ReadFile("./cart-data.json")
	if err != nil {
		panic(err)
	}

	var users []acmeserverless.User

	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Println(err)
	}

	for _, usr := range users {
		err = AddUser(usr)
		if err != nil {
			log.Println(err)
		}
	}

	var products []acmeserverless.CatalogItem

	err = json.Unmarshal(data, &products)
	if err != nil {
		log.Println(err)
	}

	for _, product := range products {
		err = AddProduct(product)
		if err != nil {
			log.Println(err)
		}
	}

	var orders acmeserverless.Orders

	err = json.Unmarshal(data, &orders)
	if err != nil {
		log.Println(err)
	}

	for _, ord := range orders {
		ord, err = AddOrder(ord)
		if err != nil {
			log.Println(err)
		}
	}

	var carts acmeserverless.Carts

	err = json.Unmarshal(data, &carts)
	if err != nil {
		log.Println(err)
	}

	for _, crt := range carts {
		err = StoreItems(crt.UserID, crt.Items)
		if err != nil {
			log.Println(err)
		}
	}
}
