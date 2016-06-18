package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	r := httprouter.New()
	http.Handle("/", r)
	r.GET("/", index) // <user> Root page
}

func index(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, "index.html", nil)
}

func documentation(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, "documentation.html", nil)
}
