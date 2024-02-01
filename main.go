package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1> Welcome to my site </h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1> Contact Page </h1><p> To get in touch, email me at <a href=\"mailto:robmmar90@gmail.com\">Robmmar90@gmail.com</a></p>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
	<h1>FAQs</h1> 
	<ul>
		<li>Q: Is there a free version?</li> 
		<li>A: Yes! We offer a free trial for 30 days on any paid plans</li>
		&nbsp 
		<li>Q: What are your support hours?</li> 
		<li>A: We support staff answering emails 24/7, though response times may be a bit slower on weekends</li>
		&nbsp
		<li>Q: How do I contact support?</li> 
		<li>A: Email us - <a href=/contact>Robmmar90@gmail.com</a></li>
	`)
}

func faqHandlerTwo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	id := chi.URLParam(r, "id")
	log.Println(id)
	fmt.Fprint(w, id)
}

func main() {

	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.Logger)
	// Gets
	mainRouter.Get("/", homeHandler)
	mainRouter.Get("/contact", contactHandler)
	mainRouter.Get("/faq", faqHandler)
	mainRouter.Get("/test/{id}", faqHandlerTwo)






	log.Fatal(http.ListenAndServe(":3000", mainRouter))
}
