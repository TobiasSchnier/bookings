package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TobiasSchnier/bookings/internal/config"
	"github.com/TobiasSchnier/bookings/internal/handlers"
	"github.com/TobiasSchnier/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

// const cant be changed
const portNumber = ":8081"

// outside func main, so middleware can access it as well (same package)
var app config.AppConfig

var session *scs.SessionManager

func main() {
	// change this to true in PROD
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create TemplateCache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println("Starting app on Port:", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
