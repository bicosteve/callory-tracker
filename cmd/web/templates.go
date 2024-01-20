package main

import (
	"html/template"
	"path/filepath"
)

/*
This file contains all the methods related to template functionalities.
*/

/*
functions -> object that contains string keyed map for look up between the name of custom template functions
*/
var functions = template.FuncMap{}

func loadTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get slice of all the filepaths with extension .html
	// use filepath.Glob()

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract the filename eg home.html
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err

		}

		// ParseGlob() method add any file with .base to template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.base.html"))
		if err != nil {
			return nil, err
		}

		// Adding *.partial.html to template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}

	return cache, nil
}
