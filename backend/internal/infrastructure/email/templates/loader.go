package templates

import (
	"embed"
	"fmt"
	"text/template"
)

//go:embed *.txt
var templateFS embed.FS

// TemplateProvider provides TemplateProvider related functionality.
type TemplateProvider interface {
	Get(name string) *template.Template
}

// TemplateLoader provides TemplateLoader related functionality.
type TemplateLoader struct {
	templates *template.Template
}

// NewTemplateLoader creates a new TemplateLoader instance.
func NewTemplateLoader() (*TemplateLoader, error) {
	tmpl, err := template.ParseFS(templateFS, "*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}
	return &TemplateLoader{templates: tmpl}, nil
}

// Get provides Get related functionality.
func (l *TemplateLoader) Get(name string) *template.Template {
	return l.templates.Lookup(name)
}
