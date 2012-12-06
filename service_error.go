package restful

import "fmt"

// ServiceError is a transport object to pass information about a non-Http error occurred in a WebService while processing a request.
type ServiceError struct {
	Code    int
	Message string
}

// Return a new ServiceError using the code and reason
func NewError(code int, message string) ServiceError {
	return ServiceError{Code: code, Message: message}
}

// Return a text representation of the service error
func (self ServiceError) Error() string {
	return fmt.Sprintf("[ServiceError:%v] %v", self.Code, self.Message)
}
