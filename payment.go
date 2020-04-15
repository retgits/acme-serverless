package acmeserverless

import (
	"encoding/json"
	"strconv"

	"github.com/retgits/creditcard"
)

// CreditCardValidatedEvent is sent by the payment service when the creditcard has been validated.
type CreditCardValidatedEvent struct {
	// Metadata for the event.
	Metadata Metadata `json:"metadata"`

	// Data contains the payload data for the event.
	Data CreditCardValidationDetails `json:"data"`
}

// UnmarshalCreditCardValidatedEvent parses the JSON-encoded data and stores the result in a
// CreditCardValidatedEvent.
func UnmarshalCreditCardValidatedEvent(data []byte) (CreditCardValidatedEvent, error) {
	var r CreditCardValidatedEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal returns the JSON encoding of CreditCardValidatedEvent.
func (e *CreditCardValidatedEvent) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// CreditCardValidationDetails contain the details of the validation by the payment service.
type CreditCardValidationDetails struct {
	// Indicates whether the transaction was a success or not.
	Success bool `json:"success"`

	// The HTTP statuscode of the event.
	Status int `json:"status"`

	// A string containing the result of the service.
	Message string `json:"message"`

	// The monetary amount of the transaction.
	Amount string `json:"amount,omitempty"`

	// The unique identifier of the transaction.
	TransactionID string `json:"transactionID"`

	// The unique identifier of the order.
	OrderID string `json:"orderID"`
}

// UnmarshalCreditCardValidationDetails parses the JSON-encoded data and stores the result in a
// CreditCardValidationDetails.
func UnmarshalCreditCardValidationDetails(data []byte) (CreditCardValidationDetails, error) {
	var r CreditCardValidationDetails
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal returns the JSON encoding of CreditCardValidationDetails.
func (e *CreditCardValidationDetails) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// PaymentRequestDetails contains the data that is needed to validate the payment.
type PaymentRequestDetails struct {
	// The unique identifier of the order.
	OrderID string `json:"orderID"`

	// Card used for the transaction.
	Card creditcard.Card `json:"card"`

	// Total monetary value of the transaction.
	Total string `json:"total"`
}

// UnmarshalPaymentRequestDetails parses the JSON-encoded data and stores the result in a
// PaymentRequestDetails.
func UnmarshalPaymentRequestDetails(data []byte) (PaymentRequestDetails, error) {
	var r PaymentRequestDetails
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal returns the JSON encoding of PaymentRequestDetails.
func (e *PaymentRequestDetails) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// PaymentRequestedEvent is sent by the Order service when the creditcard for the order should be
// validated and charged.
type PaymentRequestedEvent struct {
	// Metadata for the event.
	Metadata Metadata `json:"metadata"`

	// Data contains the payload data for the event.
	Data PaymentRequestDetails `json:"data"`
}

// UnmarshalPaymentRequestedEvent parses the JSON-encoded data and stores the result in a
// PaymentRequestedEvent.
func UnmarshalPaymentRequestedEvent(data []byte) (PaymentRequestedEvent, error) {
	var r PaymentRequestedEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal returns the JSON encoding of PaymentRequestedEvent.
func (e *PaymentRequestedEvent) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// UnmarshalShopPayment parses the JSON-encoded data and stores the result in a
// ShopPayment.
func UnmarshalShopPayment(data []byte) (ShopPayment, error) {
	var r ShopPayment
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal returns the JSON encoding of ShopPayment.
func (r *ShopPayment) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// ShopPayment allows for backward compatibility between the serverless version and the containerized version
// of the ACME Fitness Shop. This format should only be used for HTTP based interactions.
type ShopPayment struct {
	Card  ShopCard `json:"card"`
	Total string   `json:"total"`
}

// ShopCard allows for backward compatibility between the serverless version and the containerized version
// of the ACME Fitness Shop. This format should only be used for HTTP based interactions.
type ShopCard struct {
	Number   string `json:"number"`
	ExpYear  string `json:"expYear"`
	ExpMonth string `json:"expMonth"`
	Ccv      string `json:"ccv"`
}

// ToCreditCard maps the old shopcard data to the new creditcard and allows for backward compatibility between
// the serverless version and the containerized version of the ACME Fitness Shop. This format should only be used
// for HTTP based interactions.
func (s *ShopCard) ToCreditCard() creditcard.Card {
	em, _ := strconv.Atoi(s.ExpMonth)
	ey, _ := strconv.Atoi(s.ExpYear)
	return creditcard.Card{
		CVV:         s.Ccv,
		Number:      s.Number,
		ExpiryMonth: em,
		ExpiryYear:  ey,
	}
}
