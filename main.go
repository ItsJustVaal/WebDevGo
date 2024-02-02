package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/ItsJustVaal/WebDevGo/controllers"
	"github.com/ItsJustVaal/WebDevGo/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.Logger)
	mainRouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	
	// Gets
	mainRouter.Get("/", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))))
	mainRouter.Get("/contact", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))))
	mainRouter.Get("/faq", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))))






	log.Fatal(http.ListenAndServe(":3000", mainRouter))
}
