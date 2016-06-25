package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	r := httprouter.New()
	http.Handle("/", r)
	r.GET("/", Index)
	r.GET("/doc", Documentation)
	r.POST("/state", UserState)
	r.POST("/newStudent", RegisterStudent)
	r.POST("/newClient", RegisterClient)
	r.POST("/loginUser", LoginUser)
	r.POST("/logoutUser", LogoutUser)
}

func Index(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, "index.html", nil)
}

func Documentation(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, "documentation.html", nil)
}
