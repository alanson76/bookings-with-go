package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/alanson76/playground/web/02_firstApp/pkg/config"
	"github.com/alanson76/playground/web/02_firstApp/pkg/models"
)

// map functions to template
var functions = template.FuncMap{}

// structs for cached templates, etc
var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// addDefaultData add default data to all the templates
func addDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// reders template
// relative directory from the go.mod for files
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	// Develop mode / Production mode
	// use cached template if production mode
	// load new cache every time if development mode
	var tc map[string]*template.Template
	if app.UseCache {
		// get templates cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get the current template
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template the from templates cache")
	}

	// allocates the memory as bytes
	buf := new(bytes.Buffer)

	// wrting template to the buffer
	td = addDefaultData(td)
	_ = t.Execute(buf, td)

	// drain the buffer to browser
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	// cash the templates
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// fmt.Println("building template cache for ", name)

		// template sets
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// check if there are any layout templates
		matches, err := filepath.Glob("./templates/layouts/*.layout.html")
		if err != nil {
			return myCache, err
		}

		// if there are layouts, parse the .layout.html
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/layouts/*.layout.html")
			if err != nil {
				return myCache, err
			}

		}

		// assign the parsed template to map, myCache
		myCache[name] = ts
	}

	return myCache, nil
}
