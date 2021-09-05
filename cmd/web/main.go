package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alanson76/bookings-with-go/pkg/config"
	"github.com/alanson76/bookings-with-go/pkg/handlers"
	"github.com/alanson76/bookings-with-go/pkg/render"
	"github.com/alexedwards/scs/v2"
)

// Variables
const portNumber string = ":3000"

// configuration variable
var app config.AppConfig

// access to session
var session *scs.SessionManager

// main goroutine
func main() {

	// change this to true when in production
	app.InProduction = false

	// session configuration
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session // saves the session to the app configuration

	// create template cache and assign to the global config variable, config.AppConfig
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	// assign the template cache to configuration so that others can use
	app.TemplateCache = tc
	// true for production, false for development
	app.UseCache = false

	// creates new repository and assign the repository
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// pass the app to render
	render.NewTemplates(&app)

	// routers
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	// run server
	fmt.Println("Server is running on port ", portNumber)
	// err = http.ListenAndServe(portNumber, nil)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
