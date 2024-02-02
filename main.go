package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func executeTemplate(w http.ResponseWriter, fp string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Printf("Error parsing template: %v", err.Error())
		http.Error(w, "There was an error Parsing", http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err.Error())
		http.Error(w, "Invalid data passed to execute", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "home.gohtml"))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "contact.gohtml"))
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "faq.gohtml"))
}


func main() {
	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.Logger)
	mainRouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	
	// Gets
	mainRouter.Get("/", homeHandler)
	mainRouter.Get("/contact", contactHandler)
	mainRouter.Get("/faq", faqHandler)






	log.Fatal(http.ListenAndServe(":3000", mainRouter))
}
