package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zepyrshut/breakfast-n-go/internal/config"
	"github.com/zepyrshut/breakfast-n-go/internal/handlers"
	"github.com/zepyrshut/breakfast-n-go/internal/render"

	"github.com/alexedwards/scs/v2"
)

// Define the port number, localhost is valid too.
const portNumber = "127.0.0.1:8080"

var app config.AppConfig
var session *scs.SessionManager

// Main function
func main() {

	// change this is to true when in production

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Printf("Server is running on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
