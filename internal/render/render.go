package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/TobiasSchnier/bookings/internal/config"
	"github.com/TobiasSchnier/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// sets the config vor the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// get TemaplteCache from Appconfig
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested Template from Cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get TemplateCache")
	}
	buf := new(bytes.Buffer)

	// set default data
	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render Template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	//get all files named *.page.tmpl from Templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
