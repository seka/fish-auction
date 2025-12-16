package templates

import (
	"embed"
	"fmt"
	"text/template"
)

//go:embed *.txt
var templateFS embed.FS

type TemplateProvider interface {
	Get(name string) *template.Template
}

type TemplateLoader struct {
	templates *template.Template
}

func NewTemplateLoader() (*TemplateLoader, error) {
	tmpl, err := template.ParseFS(templateFS, "*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}
	return &TemplateLoader{templates: tmpl}, nil
}

func (l *TemplateLoader) Get(name string) *template.Template {
	return l.templates.Lookup(name)
}
