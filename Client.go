package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
)

func NewUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

func registerClient(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.
	fmt.Sprint(NewUUID())
}
