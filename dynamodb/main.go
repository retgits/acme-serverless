package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	cart "github.com/retgits/acme-serverless-cart"
	cartDB "github.com/retgits/acme-serverless-cart/internal/datastore/dynamodb"
	catalog "github.com/retgits/acme-serverless-catalog"
	catalogDB "github.com/retgits/acme-serverless-catalog/internal/datastore/dynamodb"
	order "github.com/retgits/acme-serverless-order"
	orderDB "github.com/retgits/acme-serverless-order/internal/datastore/dynamodb"
	user "github.com/retgits/acme-serverless-user"
	userDB "github.com/retgits/acme-serverless-user/internal/datastore/dynamodb"
)

const (
	region = "us-west-2"
	table  = "acmeserverless"
)

func main() {
	os.Setenv("REGION", region)
	os.Setenv("TABLE", table)

	data, err := ioutil.ReadFile("./user-data.json")
	if err != nil {
		log.Println(err)
	}

	var users []user.User

	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Println(err)
	}

	userdb := userDB.New()

	for _, usr := range users {
		err = userdb.AddUser(usr)
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

	catalogdb := catalogDB.New()

	for _, product := range products {
		err = catalogdb.AddProduct(product)
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

	orderdb := orderDB.New()

	for _, ord := range orders {
		ord, err = orderdb.AddOrder(ord)
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

	cartdb := cartDB.New()

	for _, crt := range cartdb {
		err = dynamoStore.StoreItems(crt.Userid, crt.Items)
		if err != nil {
			log.Println(err)
		}
	}
}
