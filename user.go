package acmeserverless

import "encoding/json"

// User is a single user in the ACME Serverless Fitness Shop
type User struct {
	// ID is the unique identifier of the user in the shop
	ID string `json:"id"`

	// Username is the username of the user
	Username string `json:"username"`

	// Password is the password of the user
	Password string `json:"password"`

	// Firstname is the firstname of the user
	Firstname string `json:"firstname"`

	// Lastname is the lastname of the user
	Lastname string `json:"lastname"`

	// Email is where we need to send spam ;-)
	Email string `json:"email"`
}

// UnmarshalUser parses the JSON-encoded data and stores the result in a User
func UnmarshalUser(data string) (User, error) {
	var r User
	err := json.Unmarshal([]byte(data), &r)
	return r, err
}

// Marshal returns the JSON encoding of User
func (r *User) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// RegisterUserResponse is sent back to the front-end service after a new user tries to register for the
// ACME Serverless Fitness Shop
type RegisterUserResponse struct {
	// Message is a status message indicating success or failure
	Message string `json:"message"`

	// ResourceID represents the user that has been created together with the new user ID
	ResourceID string `json:"resourceId"`

	// Status is the HTTP status code indicating success or failure
	Status int64 `json:"status"`
}

// UnmarshalRegisterUserResponse parses the JSON-encoded data and stores the result
// in a RegisterUserResponse
func UnmarshalRegisterUserResponse(data string) (RegisterUserResponse, error) {
	var r RegisterUserResponse
	err := json.Unmarshal([]byte(data), &r)
	return r, err
}

// Marshal returns the JSON encoding of RegisterUserResponse
func (r *RegisterUserResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// AllUsers is the response struct for the reply to the API call to get all users.
type AllUsers struct {
	Data []User `json:"data"`
}

// Marshal returns the JSON encoding of AllUsers
func (r *AllUsers) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// LoginRequest is the request to log in to be able to buy new products
type LoginRequest struct {
	// Username is the username of the user
	Username string `json:"username"`

	// Password is the password of the user
	Password string `json:"password"`
}

// UnmarshalLoginRequest parses the JSON-encoded data and stores the result in a LoginRequest
func UnmarshalLoginRequest(data string) (LoginRequest, error) {
	var r LoginRequest
	err := json.Unmarshal([]byte(data), &r)
	return r, err
}

// Marshal returns the JSON encoding of AllUsers
func (r *LoginRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// LoginResponse is the response sent to the client to indicate whether or not the user has successfully
// authenticated
type LoginResponse struct {
	// Accesstoken is the JWT access token containing information and claims
	AccessToken string `json:"access_token"`

	// RefreshToken is the token needed to refresh the access token
	RefreshToken string `json:"refresh_token"`

	// Status is the HTTP status code indicating success or failure
	Status int64 `json:"status"`
}

// UnmarshalLoginResponse parses the JSON-encoded data and stores the result in a LoginResponse
func UnmarshalLoginResponse(data string) (LoginResponse, error) {
	var r LoginResponse
	err := json.Unmarshal([]byte(data), &r)
	return r, err
}

// Marshal returns the JSON encoding of LoginResponse
func (r *LoginResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// VerifyTokenResponse is sent to the service requesting a validation of the access token
type VerifyTokenResponse struct {
	// Message is a status message indicating success or failure
	Message string `json:"message"`

	// Status is the HTTP status code indicating success or failure
	Status int `json:"status"`
}

// UnmarshalVerifyTokenResponse parses the JSON-encoded data and stores the result in a VerifyTokenResponse
func UnmarshalVerifyTokenResponse(data string) (VerifyTokenResponse, error) {
	var r VerifyTokenResponse
	err := json.Unmarshal([]byte(data), &r)
	return r, err
}

// Marshal returns the JSON encoding of VerifyTokenResponse
func (r *VerifyTokenResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UserDetailsResponse is sent back to the user with details about that specific user
type UserDetailsResponse struct {
	// User are the details about the user
	User User `json:"data"`

	// Status is the HTTP status code indicating success or failure
	Status int `json:"status"`
}

// Marshal returns the JSON encoding of UserDetailsResponse
func (r *UserDetailsResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
