package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zepyrshut/breakfast-n-go/internal/config"
	"github.com/zepyrshut/breakfast-n-go/internal/driver"
	"github.com/zepyrshut/breakfast-n-go/internal/handlers"
	"github.com/zepyrshut/breakfast-n-go/internal/helpers"
	"github.com/zepyrshut/breakfast-n-go/internal/models"
	"github.com/zepyrshut/breakfast-n-go/internal/render"

	"github.com/alexedwards/scs/v2"
)

// Define the port number, localhost is valid too.
const portNumber = "127.0.0.1:8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// Main function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	fmt.Printf("server is running on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this is to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("connecting to database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=root")
	if err != nil {
		log.Fatal("cannont connect to db")
	}
	log.Println("connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	// repository pattern
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
