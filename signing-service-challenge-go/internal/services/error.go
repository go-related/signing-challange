package services

import "fmt"

type ServiceError struct {
	msg    string
	Status int
}

func (r *ServiceError) Error() string {
	return fmt.Sprintf("validation error: %v", r.msg)
}

func NewServiceError(msg string, status int) *ServiceError {
	return &ServiceError{msg: msg, Status: status}
}
