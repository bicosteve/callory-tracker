package main

import (
	"github.com/bicosteve/callory-tracker/pkg/models"
	"html/template"
	"path/filepath"
)

/*
This file contains all the methods related to template functionalities.
*/

/*
functions -> object that contains string keyed map for look up between the name
of custom template functions
*/
var functions = template.FuncMap{}

// templateData struct will hold any dynamic data that we want to pass to HTM Templates
type templateData struct {
	User        *models.User
	Food        *models.Food
	Foods       []*models.Food
	CurrentYear int
}

// Used for template caching when the application starts
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get slice of all the file paths with extension .html
	// use filepath.Glob()

	// This gives a slice of all the '.page' templates for the app
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract the filename e.g home.page.html
		// assign it to the name variable
		name := filepath.Base(page)

		// Parse the page template file in to a template set
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err

		}

		// Use ParseGlob() method to add any file with .base.html to template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.base.html"))
		if err != nil {
			return nil, err
		}

		// Use ParseGlob() to add any *.partial.html file to template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache using the name of the page
		// example 'home.page.html' as the key in the map
		cache[name] = ts

	}

	return cache, nil
}
