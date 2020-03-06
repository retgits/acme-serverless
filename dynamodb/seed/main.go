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

	cart "github.com/retgits/acme-serverless-cart"
	catalog "github.com/retgits/acme-serverless-catalog"
	order "github.com/retgits/acme-serverless-order"
	user "github.com/retgits/acme-serverless-user"
)

var (
	region    string
	table     string
	dynamoURL string
)

// Create a single instance of the dynamoDB service
// which can be reused if the container stays warm
var dbs *dynamodb.DynamoDB

// initialize creates the connection to dynamoDB. If the environment
// variable DYNAMO_URL is set, the connection is made to that URL
// instead of relying on the AWS SDK to provide the URL
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
func AddUser(usr user.User) error {
	// Create a JSON encoded string of the user
	payload, err := usr.Marshal()
	if err != nil {
		return err
	}

	// Create a map of DynamoDB Attribute Values containing the table keys
	km := make(map[string]*dynamodb.AttributeValue)
	km["PK"] = &dynamodb.AttributeValue{
		S: aws.String("USER"),
	}
	km["SK"] = &dynamodb.AttributeValue{
		S: aws.String(usr.ID),
	}

	// Create a map of DynamoDB Attribute Values containing the table data elements
	em := make(map[string]*dynamodb.AttributeValue)
	em[":keyid"] = &dynamodb.AttributeValue{
		S: aws.String(usr.Username),
	}
	em[":payload"] = &dynamodb.AttributeValue{
		S: aws.String(payload),
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

// AddProduct stores a new product in Amazon DynamoDB
func AddProduct(p catalog.Product) error {
	// Marshal the newly updated product struct
	payload, err := p.Marshal()
	if err != nil {
		return err
	}

	// Create a map of DynamoDB Attribute Values containing the table keys
	km := make(map[string]*dynamodb.AttributeValue)
	km["PK"] = &dynamodb.AttributeValue{
		S: aws.String("PRODUCT"),
	}
	km["SK"] = &dynamodb.AttributeValue{
		S: aws.String(p.ID),
	}

	// Create a map of DynamoDB Attribute Values containing the table data elements
	em := make(map[string]*dynamodb.AttributeValue)
	em[":payload"] = &dynamodb.AttributeValue{
		S: aws.String(payload),
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
func AddOrder(o order.Order) (order.Order, error) {
	// Generate and assign a new orderID
	o.OrderID = uuid.Must(uuid.NewV4()).String()
	o.Status = aws.String("Pending Payment")

	// Marshal the newly updated product struct
	payload, err := o.Marshal()
	if err != nil {
		return o, fmt.Errorf("error marshalling order: %s", err.Error())
	}

	// Create a map of DynamoDB Attribute Values containing the table keys
	km := make(map[string]*dynamodb.AttributeValue)
	km["PK"] = &dynamodb.AttributeValue{
		S: aws.String("ORDER"),
	}
	km["SK"] = &dynamodb.AttributeValue{
		S: aws.String(o.OrderID),
	}

	// Create a map of DynamoDB Attribute Values containing the table data elements
	em := make(map[string]*dynamodb.AttributeValue)
	em[":keyid"] = &dynamodb.AttributeValue{
		S: aws.String(o.UserID),
	}
	em[":payload"] = &dynamodb.AttributeValue{
		S: aws.String(payload),
	}

	uii := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(table),
		Key:                       km,
		ExpressionAttributeValues: em,
		UpdateExpression:          aws.String("SET Payload = :payload, KeyID = :keyid"),
	}

	_, err = dbs.UpdateItem(uii)
	if err != nil {
		return o, fmt.Errorf("error updating dynamodb: %s", err.Error())
	}

	return o, nil
}

// StoreItems saves the cart items from a single user into Amazon DynamoDB
func StoreItems(userID string, i cart.Items) error {
	payload, err := i.Marshal()
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

	initialize()

	data, err := ioutil.ReadFile("./user-data.json")
	if err != nil {
		log.Println(err)
	}

	var users []user.User

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

	data, err = ioutil.ReadFile("./catalog-data.json")
	if err != nil {
		log.Println(err)
	}

	var products []catalog.Product

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

	data, err = ioutil.ReadFile("./order-data.json")
	if err != nil {
		log.Println(err)
	}

	var orders order.Orders

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

	data, err = ioutil.ReadFile("./cart-data.json")
	if err != nil {
		log.Println(err)
	}

	var carts cart.Carts

	err = json.Unmarshal(data, &carts)
	if err != nil {
		log.Println(err)
	}

	for _, crt := range carts {
		err = StoreItems(crt.Userid, crt.Items)
		if err != nil {
			log.Println(err)
		}
	}
}
