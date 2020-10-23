package handler

import (
	"log"
	"net/http"
	"text/template"
)

func render(w http.ResponseWriter, r *http.Request, tmplName string,
	data interface{}) {
	baseFile := "layout"
	tmpls := []string{
		"./ui/html/layout.html",
	}
	tmpls = append(tmpls, tmplName)

	tmpl, err := template.ParseFiles(tmpls...)
	if err != nil {
		log.Println("parse files error:", err)
		w.Write([]byte(err.Error()))
	}

	err = tmpl.ExecuteTemplate(w, baseFile, data)
	if err != nil {
		log.Println("execute error:", err)
		w.Write([]byte(err.Error()))
	}
}
