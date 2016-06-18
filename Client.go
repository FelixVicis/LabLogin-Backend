package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/nu7hatch/gouuid"
)

func NewUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

func RegisterClient(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.
	fmt.Fprint(res,NewUUID())
}