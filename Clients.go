package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	// "fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/nu7hatch/gouuid"
	"net/http"
)

func NewUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

func RegisterClient(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.
	ServeJsonOfStruct(res, JsonOptions{
		Status: "Success",
	}, NewUUID())
}
