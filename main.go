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
	"github.com/gorilla/csrf"
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

	// Root
	mainRouter.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	// Users controller
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))

	// User Gets
	mainRouter.Get("/signup", usersC.New)
	mainRouter.Get("/signin", usersC.SignIn)
	mainRouter.Get("/users/me", usersC.CurrentUser)

	// User Posts
	mainRouter.Post("/users", usersC.Create)
	mainRouter.Post("/signin", usersC.ProcessSignIn)
	mainRouter.Post("/signout", usersC.ProcessSignOut)

	mainRouter.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	mainRouter.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	csrft := "apwfWAAw0f8AWfafwareaweaAfwWAg9W2Lkf"
	csrfMW := csrf.Protect([]byte(csrft), csrf.Secure(false))
	log.Fatal(http.ListenAndServe(":3000", csrfMW(mainRouter)))
}
