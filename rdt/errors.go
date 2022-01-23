package rdt

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
	Message string `json:"message"`
}

// ErrorList defines model for ErrorList.
type ErrorList []Error

// ErrorResponse defines model for GetOrderAddressResponse.
type ErrorResponse struct {
	// A list of error responses returned when a request is unsuccessful.
	Errors ErrorList `json:"errors,omitempty"`
}
