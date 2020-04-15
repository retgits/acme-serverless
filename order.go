package acmeserverless

import (
	"encoding/json"

	"github.com/retgits/creditcard"
)

// Orders is a slide of Order objects
type Orders []Order

// Marshal returns the JSON encoding of Orders.
func (r *Orders) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Order represents an order that is made by a user in the ACME Serverless Fitness Shop.
type Order struct {
	// OrderID uniquely identifies the order
	OrderID string `json:"_id"`

	// Status represents the current status of the order
	Status *string `json:"status,omitempty"`

	// UserID represents the user who placed the order
	UserID string `json:"userid,omitempty"`

	// Firstname is the firstname of the user who placed the order
	Firstname *string `json:"firstname,omitempty"`

	// Lastname is the lastname of the user who placed the order
	Lastname *string `json:"lastname,omitempty"`

	// Address is an address where the order must be delivered
	Address *Address `json:"address,omitempty"`

	// Email is the email address to send spam to ;-)
	Email *string `json:"email,omitempty"`

	// Delivery is the delivery method the shipment must use
	Delivery string `json:"delivery"`

	// Creditcard isthe creditcard used to pay the order
	Card creditcard.Card `json:"card,omitempty"`

	// Cart contains all items part of the order
	Cart []CartItem `json:"cart"`

	// Total represents the monetary value of the order
	Total string `json:"total,omitempty"`
}

// Marshal returns the JSON encoding of an Order
func (r *Order) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UnmarshalOrder parses the JSON-encoded data and stores the result in a
// Order object.
func UnmarshalOrder(data string) (Order, error) {
	var r Order
	err := json.Unmarshal([]byte(data), &r)
	return r, err
}

// Address is an address where the order must be delivered
type Address struct {
	// Street is the streetname
	Street *string `json:"street,omitempty"`

	// City is the city name
	City *string `json:"city,omitempty"`

	// Zip is the zip or postal code
	Zip *string `json:"zip,omitempty"`

	// State is the state
	State *string `json:"state,omitempty"`

	// Country is the country where the shipment must be sent
	Country *string `json:"country,omitempty"`
}

// OrderStatus represents the current status of an order and is sent by the Order service
type OrderStatus struct {
	// OrderID uniquely represents an order
	OrderID string `json:"order_id"`

	// Payment is the data that is emitted by the payment service
	Payment CreditCardValidationDetails `json:"payment"`

	// UserID is the unique representation of the user
	UserID string `json:"userid"`
}

// Marshal returns the JSON encoding of OrderStatus
func (r *OrderStatus) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
