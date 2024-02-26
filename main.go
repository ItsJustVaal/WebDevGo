package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ItsJustVaal/WebDevGo/controllers"
	"github.com/ItsJustVaal/WebDevGo/migrations"
	"github.com/ItsJustVaal/WebDevGo/models"
	"github.com/ItsJustVaal/WebDevGo/templates"
	"github.com/ItsJustVaal/WebDevGo/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.PSQL = models.DefaultConfig()

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = "apwfWAAw0f8AWfafwareaweaAfwWAg9W2Lkf"
	cfg.CSRF.Secure = false
	cfg.Server.Address = ":3000"

	return cfg, nil
}

func main() {
	// DB init
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Services init
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	passwordResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)

	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "forgot-pw.gohtml", "tailwind.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "check-email.gohtml", "tailwind.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))

	// Middlewear
	umw := controllers.UserMiddlewear{
		SessionService: sessionService,
	}

	csrfMW := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Secure(cfg.CSRF.Secure))

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
	mainRouter.Get("/forgot-pw", usersC.ForgotPassword)
	mainRouter.Post("/forgot-pw", usersC.ProcessForgotPassword)
	mainRouter.Get("/reset-pw", usersC.ResetPassword)
	mainRouter.Post("/reset-pw", usersC.ProcessResetPassword)
	mainRouter.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	mainRouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	log.Fatal(http.ListenAndServe(cfg.Server.Address, mainRouter))
}
