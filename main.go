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
	r.POST("/create table user", CreateTableUser)
	r.POST("/insert into user", CreateUser)
	r.POST("/drop from user", DropUser)
	// Client Commands
	r.POST("/get new client", RegisterClient)
	// Record Commands
	r.POST("/select state from record", SelectStateFromRecord)
	r.POST("/toggle state from record", ToggleStateFromRecord)
	r.POST("/select all from record", SelectAllFromRecord)
	r.POST("/select current from record", SelectCurrentFromRecord)
	r.POST("/drop table record", DropAllFromRecord)
}

func Index(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, "index.html", nil)
}
