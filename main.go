package main

import (
	"log"
	"net/http"

	"github.com/ItsJustVaal/WebDevGo/controllers"
	"github.com/ItsJustVaal/WebDevGo/migrations"
	"github.com/ItsJustVaal/WebDevGo/models"
	"github.com/ItsJustVaal/WebDevGo/templates"
	"github.com/ItsJustVaal/WebDevGo/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

func main() {
	// DB init
	cfg := models.DefaultConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Services init
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

	// Middlewear
	umw := controllers.UserMiddlewear{
		SessionService: &sessionService,
	}

	csrft := "apwfWAAw0f8AWfafwareaweaAfwWAg9W2Lkf"
	csrfMW := csrf.Protect([]byte(csrft), csrf.Secure(false))

	mainRouter := chi.NewRouter()
	mainRouter.Use(csrfMW, umw.SetUser, middleware.Logger)

	mainRouter.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	mainRouter.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	mainRouter.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	// User Routes
	mainRouter.Get("/signup", usersC.New)
	mainRouter.Get("/signin", usersC.SignIn)
	mainRouter.Post("/users", usersC.Create)
	mainRouter.Post("/signin", usersC.ProcessSignIn)
	mainRouter.Post("/signout", usersC.ProcessSignOut)
	mainRouter.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	mainRouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	log.Fatal(http.ListenAndServe(":3000", mainRouter))
}
