package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/gofrs/uuid"
	acmeserverless "github.com/retgits/acme-serverless"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	username string
	password string
	hostname string
	port     string
)

// The pointer to MongoDB provides the API operation methods for making requests to MongoDB.
var dbs *mongo.Database

// initialize creates the connection to MongoDB.
func initialize() {
	connString := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, hostname, port)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %s", err.Error())
	}
	dbs = client.Database("acmeserverless")
}

// AddCatalogItem stores a new product in Amazon DynamoDB
func AddCatalogItem(p acmeserverless.CatalogItem) error {
	coll := dbs.Collection("catalog")
	payload, err := p.Marshal()
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = coll.InsertOne(ctx, bson.D{{"SK", p.ID}, {"PK", "PRODUCT"}, {"Payload", string(payload)}})

	return err
}

// AddUser stores a new user in Amazon DynamoDB
func AddUser(usr acmeserverless.User) error {
	coll := dbs.Collection("user")

	payload, err := usr.Marshal()
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = coll.InsertOne(ctx, bson.D{{"SK", usr.ID}, {"KeyID", usr.Username}, {"PK", "USER"}, {"Payload", string(payload)}})

	return err
}

func ptrString(p string) *string {
	return &p
}

// AddOrder stores a new order in Amazon DynamoDB
func AddOrder(o acmeserverless.Order) error {
	coll := dbs.Collection("order")

	// Generate and assign a new orderID
	o.OrderID = uuid.Must(uuid.NewV4()).String()
	o.Status = ptrString("Pending Payment")

	// Marshal the newly updated product struct
	payload, err := o.Marshal()
	if err != nil {
		return fmt.Errorf("error marshalling order: %s", err.Error())
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = coll.InsertOne(ctx, bson.D{{"SK", o.OrderID}, {"KeyID", o.UserID}, {"PK", "ORDER"}, {"Payload", string(payload)}})

	return nil
}

// StoreItems saves the cart items from a single user into Amazon DynamoDB
func StoreItems(userID string, i acmeserverless.CartItems) error {
	coll := dbs.Collection("cart")

	payload, err := i.Marshal()
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = coll.InsertOne(ctx, bson.D{{"SK", userID}, {"PK", "CART"}, {"Payload", string(payload)}})

	return err
}

func main() {
	// Read flags
	flag.StringVar(&username, "username", "", "")
	flag.StringVar(&password, "password", "", "")
	flag.StringVar(&hostname, "hostname", "", "")
	flag.StringVar(&port, "port", "", "")
	flag.Parse()

	initialize()

	// Read all files.
	// if any of the files are not read successfully the function panics
	userData, err := ioutil.ReadFile("../../dynamodb/seed/user-data.json")
	catalogData, err := ioutil.ReadFile("../../dynamodb/seed/catalog-data.json")
	orderData, err := ioutil.ReadFile("../../dynamodb/seed/order-data.json")
	cartData, err := ioutil.ReadFile("../../dynamodb/seed/cart-data.json")
	if err != nil {
		panic(err)
	}

	var products []acmeserverless.CatalogItem

	err = json.Unmarshal(catalogData, &products)
	if err != nil {
		log.Println(err)
	}

	for _, product := range products {
		err = AddCatalogItem(product)
		if err != nil {
			log.Println(err)
		}
	}

	var users []acmeserverless.User

	err = json.Unmarshal(userData, &users)
	if err != nil {
		log.Println(err)
	}

	for _, usr := range users {
		err = AddUser(usr)
		if err != nil {
			log.Println(err)
		}
	}

	var orders acmeserverless.Orders

	err = json.Unmarshal(orderData, &orders)
	if err != nil {
		log.Println(err)
	}

	for _, ord := range orders {
		err = AddOrder(ord)
		if err != nil {
			log.Println(err)
		}
	}

	var carts acmeserverless.Carts

	err = json.Unmarshal(cartData, &carts)
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
