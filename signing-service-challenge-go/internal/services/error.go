package services

import (
	"fmt"
)

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

type DBError struct {
	msg string
}

func (r *DBError) Error() string {
	return fmt.Sprintf("error: %v", r.msg)
}

func NewDBError(msg string) *DBError {
	return &DBError{msg: msg}
}
