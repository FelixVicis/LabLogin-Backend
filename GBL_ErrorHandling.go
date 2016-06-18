package main

import (
	"errors"
	"net/http"
)

var (
	ErrNotImplemented = errors.New("Structure Error: Function Not Implemented!")
)

// Internal Function
// generic error handling for any error we encounter.
func HandleError(res http.ResponseWriter, e error) {
	if e != nil {
		http.Error(res, e.Error(), http.StatusInternalServerError)
	}
}
