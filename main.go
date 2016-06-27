package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	r := httprouter.New()
	http.Handle("/", r)
	r.GET("/", Index)
	// User Commands
	r.GET("/create table user", CreateTableUser)
	r.GET("/insert into user", CreateUser)
	r.GET("/drop from user", DropUser)
	// Client Commands
	r.GET("/get new client", RegisterClient)
	// Record Commands
	r.GET("/select state from record", SelectStateFromRecord)
	r.GET("/toggle state from record", ToggleStateFromRecord)
	r.GET("/select all from record", SelectAllFromRecord)
	r.GET("/select current from record", SelectCurrentFromRecord)
	r.GET("/drop table record", DropAllFromRecord)
}

func Index(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, "index.html", nil)
}
