package acmeserverless

import "encoding/json"

// Carts is a slice of Cart objects
type Carts []Cart

// Marshal returns the JSON encoding of Carts
func (r *Carts) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Cart represents a shoppingcart for a user of the ACME Serverless Fitness Shop
type Cart struct {
	// Items is a slice of Item objects, each being a single object in the cart of the user
	Items []CartItem `json:"cart"`

	// UserID is the unique identifier of the user in the ACME Serverless Fitness Shop
	UserID string `json:"userid"`
}

// Marshal returns the JSON encoding of Cart
func (r *Cart) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UnmarshalCart parses the JSON-encoded data and stores the result in a Cart
func UnmarshalCart(data string) (Cart, error) {
	var r Cart
	err := json.Unmarshal([]byte(data), &r)
	return r, err
}

// Items is a slice of Item objects
type CartItems []CartItem

// Marshal returns the JSON encoding of Items
func (r *CartItems) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UnmarshalItems parses the JSON-encoded data and stores the result in an CartItems object
func UnmarshalItems(data string) (CartItems, error) {
	var r CartItems
	err := json.Unmarshal([]byte(data), &r)
	return r, err
}

// CartItem represents the items that the ACME Serverless Fitness Shop user has in their shopping cart
type CartItem struct {
	// Description is a description of the items
	Description string `json:"description"`

	// ItemID is the unique identifier of the item
	// This ID is set when the item originates in the cart domain
	ItemID *string `json:"itemid,omitempty"`

	// ID is the unique representation of the item
	// This ID is set when the item originates in the order domain
	ID *string `json:"id,omitempty"`

	// Name is the name of the item
	Name string `json:"name"`

	// Price is the monetairy value of the item
	Price float64 `json:"price"`

	// Quantity is how many of the item the user has in their cart
	Quantity int64 `json:"quantity"`
}

// Marshal returns the JSON encoding of CartItem
func (r *CartItem) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UnmarshalItem parses the JSON-encoded data and stores the result in an Item
func UnmarshalItem(data []byte) (CartItem, error) {
	var r CartItem
	err := json.Unmarshal(data, &r)
	return r, err
}

// CartItemTotal represents how many items the user currently has in their cart
type CartItemTotal struct {
	// CartItemTotal is the number of items
	CartItemTotal int64 `json:"cartitemtotal"`

	// UserID is the unique identifier of the user in the ACME Serverless Fitness Shop
	UserID string `json:"userid"`
}

// Marshal returns the JSON encoding of CartTotal
func (r *CartItemTotal) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// CartValueTotal represents the total value of all items currently in the cart of the iser
type CartValueTotal struct {
	// CartTotal is the value of items
	CartTotal float64 `json:"carttotal"`

	// UserID is the unique identifier of the user in the
	// ACME Serverless Fitness Shop
	UserID string `json:"userid"`
}

// Marshal returns the JSON encoding of CartValue
func (r *CartValueTotal) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UserIDResponse returns the UserID
type UserIDResponse struct {
	// UserID is the unique identifier of the user in the ACME Serverless Fitness Shop
	UserID string `json:"userid"`
}

// Marshal returns the JSON encoding on UserIDResponse
func (r *UserIDResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
