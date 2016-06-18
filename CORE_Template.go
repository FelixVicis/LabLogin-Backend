package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	funcMap := template.FuncMap{}
	var err error
	tpl, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.*")
	if err != nil {
		panic(err)
	}
}

// Internal Function
// Passes along any information to templates and then executes them.
func ServeTemplateWithParams(res http.ResponseWriter, templateName string, params interface{}) {
	err := tpl.ExecuteTemplate(res, templateName, &params)
	HandleError(res, err)
}
