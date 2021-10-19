package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Anka-Abdullah/bookings/pkg/config"
	"github.com/Anka-Abdullah/bookings/pkg/handlers"
	"github.com/Anka-Abdullah/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

var port = ":4567"
var app config.AppConfig
var session *scs.SessionManager

func main() {

	//change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Println(err)
		log.Fatal("cannot create template cache")
	}

	app.TemplateChace = tc
	app.UseChace = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("starting application on http://localhost%s \n", port)
	fmt.Println("http://localhost:4567/about")

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
