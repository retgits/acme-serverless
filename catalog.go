package acmeserverless

import "encoding/json"

// CatalogItem represents the products as they are stored in the data store
type CatalogItem struct {
	// ID is the unique identifier of the product
	ID string `json:"id"`

	// Name is the name of the product
	Name string `json:"name"`

	// ShortDescription is a short description of the product suited for Point of Sales or mobile apps
	ShortDescription string `json:"shortDescription"`

	// Description is a longer description of the product suited for websites
	Description string `json:"description"`

	// ImageURL1 is the location of the first image
	ImageURL1 string `json:"imageUrl1"`

	// ImageURL2 is the location of the second image
	ImageURL2 string `json:"imageUrl2"`

	// ImageURL3 is the location of the third image
	ImageURL3 string `json:"imageUrl3"`

	// Price is the monetary value of the product
	Price float32 `json:"price"`

	// Tags are keys that represent additional sorting information for front-end displays
	Tags []string `json:"tags"`
}

// UnmarshalCatalogItem parses the JSON-encoded data and stores the result
// in a CatalogItem
func UnmarshalCatalogItem(data string) (CatalogItem, error) {
	var r CatalogItem
	err := json.Unmarshal([]byte(data), &r)
	return r, err
}

// Marshal returns the JSON encoding of Product
func (r *CatalogItem) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// CreateCatalogItemResponse is the respons that is sent back to the API after a new item
// has been added to the catalog.
type CreateCatalogItemResponse struct {
	// Message is a status message indicating success or failure
	Message string `json:"message"`

	// ResourceID represents the product that has been created
	// together with the new product ID
	ResourceID CatalogItem `json:"resourceId"`

	// Status is the HTTP status code indicating success or failure
	Status int `json:"status"`
}

// Marshal returns the JSON encoding of CreateCatalogItemResponse
func (r *CreateCatalogItemResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// AllCatalogItemsResponse is the response struct for the reply to
// the API call to get all products.
type AllCatalogItemsResponse struct {
	Data []CatalogItem `json:"data"`
}

// Marshal returns the JSON encoding of AllCatalogItemsResponse
func (r *AllCatalogItemsResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
