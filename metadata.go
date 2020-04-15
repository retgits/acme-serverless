package acmeserverless

// Metadata contains information on the domain, source, type, and status
// of the event.
type Metadata struct {
	// Domain represents the the event came from
	// like Payment or Order.
	Domain string `json:"domain"`

	// Source represents the function the event came from
	// like ValidateCreditCard or SubmitOrder.
	Source string `json:"source"`

	// Type respresents the type of event this is
	// like CreditCardValidated.
	Type string `json:"type"`

	// Status represents the current status of the event
	// like Success or Failure.
	Status string `json:"status"`
}
