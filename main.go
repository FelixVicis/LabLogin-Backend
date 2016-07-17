package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	r := httprouter.New()
	dbug := false

	http.Handle("/", r)
	r.GET("/", Index)
	Init_UserRoutes(r, dbug)   // User Commands
	Init_ClientRoutes(r, dbug) // Client Commands
	Init_RecordRoutes(r, dbug) // Record Commands
}

func Index(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, "index.html", nil)
}
