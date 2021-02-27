package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPError struct {
	status  int
	detail  string
	err     error
	payload interface{}
}

func NewHTTPError(status int, detail string, err error, payload interface{}) error {
	return &HTTPError{
		status:  status,
		detail:  detail,
		err:     err,
		payload: payload,
	}
}

func (e *HTTPError) Error() string {
	if e.err == nil {
		return fmt.Sprintf("response with status %d: %s", e.status, e.detail)
	}
	return fmt.Sprintf("response with status %d: %s", e.status, e.err)
}

func (e *HTTPError) Unwrap() error {
	return e.err
}

func (e *HTTPError) ResponseBody() interface{} {
	if e.payload != nil {
		return e.payload
	}
	body := ErrorBody{
		Error: e.detail,
	}
	if e.err != nil {
		body.Error = e.err.Error()
	}
	return body
}

func (e *HTTPError) Response(w http.ResponseWriter) error {
	w.WriteHeader(e.Status())
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body := e.ResponseBody()
	var err error
	if bytes, ok := body.(Raw); ok {
		_, err = w.Write(bytes)
	} else {
		err = json.NewEncoder(w).Encode(body)
	}
	return err
}

func (e *HTTPError) Status() int {
	return e.status
}

func (e *HTTPError) Err() error {
	return e.err
}

func (e *HTTPError) Detail() string {
	return e.detail
}

func (e *HTTPError) Data() interface{} {
	return e.payload
}

type ErrorBody struct {
	Error string
}

// ErrorBadRequest - 400.
func ErrorBadRequest(msg string, payload interface{}) error {
	return NewHTTPError(http.StatusBadRequest, msg, nil, payload)
}

// ErrorUnauthorized - 401.
func ErrorUnauthorized(msg string, payload interface{}) error {
	return NewHTTPError(http.StatusUnauthorized, msg, nil, payload)
}

// ErrorForbidden - 403.
func ErrorForbidden(msg string, payload interface{}) error {
	return NewHTTPError(http.StatusForbidden, msg, nil, payload)
}

// ErrorNotFound - 404.
func ErrorNotFound(msg string, payload interface{}) error {
	return NewHTTPError(http.StatusNotFound, msg, nil, payload)
}
