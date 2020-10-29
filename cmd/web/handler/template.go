package handler

import (
	"html/template"
	"log"
	"net/http"
)

func render(w http.ResponseWriter, r *http.Request, tmplName string,
	funcMap template.FuncMap, data interface{}) {
	baseFile := "layout"
	tmpls := []string{
		"./ui/html/layout.html",
	}
	tmpls = append(tmpls, tmplName)

	tmpl, err := template.New(tmplName).Funcs(funcMap).ParseFiles(tmpls...)
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

// function mapping
func rolesChecked(roleName string, userRoles []string) string {
	for _, ur := range userRoles {
		if ur == roleName {
			return "checked"
		}
	}
	return ""
}

func safe(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}
