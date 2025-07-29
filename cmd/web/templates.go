package main

import "github.com/arynkh/snippetbox/internal/models"

// holdiing structure for any dynamic data we want to pass to the template
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
