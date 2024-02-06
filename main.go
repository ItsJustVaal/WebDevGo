package main

import (
	"log"
	"net/http"

	"github.com/ItsJustVaal/WebDevGo/controllers"
	"github.com/ItsJustVaal/WebDevGo/models"
	"github.com/ItsJustVaal/WebDevGo/templates"
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

	// DB init
	cfg := models.DefaultConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	
	// Users controller
	userService := models.UserService{
		DB: db,
	}
	usersC := controllers.Users{
		UserService: &userService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	mainRouter.Get("/signup", usersC.New)
	mainRouter.Get("/signin", usersC.SignIn)
	mainRouter.Post("/users", usersC.Create)
	mainRouter.Post("/signin", usersC.ProcessSignIn)




	mainRouter.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	mainRouter.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	mainRouter.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))






	log.Fatal(http.ListenAndServe(":3000", mainRouter))
}
