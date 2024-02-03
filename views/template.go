package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err.Error())
		http.Error(w, "Invalid data passed to execute", http.StatusInternalServerError)
		return
	}
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error){
	tpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template error: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

func Parse(fp string) (Template, error) {
	tpl, err := template.ParseFiles(fp)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template error: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}