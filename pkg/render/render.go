package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Anka-Abdullah/bookings/pkg/config"
	"github.com/Anka-Abdullah/bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// renders a template
func RenderTemplate(w http.ResponseWriter, html string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if app.UseChace {
		// get the emplate chace from the app config
		tc = app.TemplateChace
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[html]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myChace := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myChace, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myChace, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myChace, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*layout.html")

			if err != nil {
				return myChace, err
			}
		}

		myChace[name] = ts
	}

	return myChace, nil
}
