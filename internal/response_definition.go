package internal

// Payload is the payload for every response
type Payload struct {
	Transaction interface{}            `json:"transaction,omitempty"`
	Token       interface{}            `json:"token,omitempty"`
	Error       map[string]interface{} `json:"error,omitempty"`
}

// Properties are the predefined set of properties for each response
type Properties struct {
	APIVersion string
}

// Base is the basic definition of every response
type Base struct {
	// The API version
	//
	// required: true
	// example: 0.0.2
	APIVersion string `json:"apiVersion"`
	// The request ID
	//
	// required: true
	// example: dc380b72-41c9-47bf-8be5-f3a7a493f4ca
	ID string `json:"id,omitempty"`
	// The request method
	//
	// required: true
	Method string `json:"method,omitempty"`
}

// Definition is the structure of every http response
type Definition struct {
	Base
	Payload
}
