package common

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	UnknownErrorName   = "UNKNOWN_ERROR"
	DefaaultHTTPStatus = http.StatusInternalServerError
)

// Error indicates a domain error
type Error interface {
	Error() string
	ClientMsg() string
}

// DomainError used for expressing errors occuring in application.
type DomainError struct {
	code         ErrorCode              // code indicates an ErrorCode customized for domain logic.
	err          error                  // err contains a native error. It will be logged in system logs.
	clientMsg    string                 // clientMsg contains a message that will return to clients
	remoteStatus int                    // remoteStatus contains proxy HTTP status code. It is used for remote process related errors.
	detail       map[string]interface{} // detail contains some details that clients may need. It is business-driven.
}

type ErrorOption func(*DomainError) error

func WithMsg(msg string) ErrorOption {
	return func(de *DomainError) error {
		de.clientMsg = msg
		return nil
	}
}

func WithStatus(status int) ErrorOption {
	return func(de *DomainError) error {
		if status < 100 || status > 599 {
			return fmt.Errorf("invalid HTTP status code: %d", status)
		}
		de.remoteStatus = status
		return nil
	}
}

func WithDetail(detail map[string]interface{}) ErrorOption {
	return func(de *DomainError) error {
		de.detail = detail
		return nil
	}
}

func NewError(code ErrorCode, err error, opts ...ErrorOption) Error {
	if err, ok := err.(Error); ok {
		return err
	}

	e := DomainError{code: code, err: err}
	for _, o := range opts {
		o(&e)
	}
	return e
}

func (e DomainError) Error() string {
	var msgs []string
	if e.remoteStatus != 0 {
		msgs = append(msgs, strconv.Itoa(e.remoteStatus))
	}
	if e.err != nil {
		msgs = append(msgs, e.err.Error())
	}
	if e.clientMsg != "" {
		msgs = append(msgs, e.clientMsg)
	}

	return strings.Join(msgs, ": ")
}

func (e DomainError) Name() string {
	if e.code.Name == "" {
		return UnknownErrorName
	}
	return e.code.Name
}

func (e DomainError) ClientMsg() string {
	return e.clientMsg
}

func (e DomainError) HTTPStatus() int {
	if e.code.StatusCode == 0 {
		return DefaaultHTTPStatus
	}
	return e.code.StatusCode
}

func (e DomainError) RemoteHTTPStatus() int {
	return e.remoteStatus
}

func (e DomainError) Detail() map[string]interface{} {
	return e.detail
}
