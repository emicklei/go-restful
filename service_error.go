package restful

import "fmt"

type ServiceError struct {
	Code    int
	Message string
}

func NewError(code int, message string) ServiceError {
	return ServiceError{Code: code, Message: message}
}

func (self ServiceError) Error() string {
	return fmt.Sprintf("[ServiceError:%v] %v", self.Code, self.Message)
}
