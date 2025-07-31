package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/arynkh/snippetbox/internal/models"
)

// holdiing structure for any dynamic data we want to pass to the template
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
	Form        any
}

// returns nicely formatted string representation of a time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate, //add the humanDate function to the template.FuncMap (stored in a global variable)
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	//get a slice of all filepaths that match the pattern below. This gives a slice of all the filepaths for our app 'page' templates
	// like  [ui/html/pages/home.html ui/html/pages/view.html]
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//extract the file name (like 'home.html') from the full filepath
		name := filepath.Base(page)

		//use template.New() to create an empty template set with the name of the page and add the functions from the functions variable
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		//call parseglob *on this template set* to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		//add the template set to the map, using the name of the page as the key
		cache[name] = ts
	}

	return cache, nil //return map
}
