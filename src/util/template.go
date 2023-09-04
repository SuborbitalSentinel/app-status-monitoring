package util

import (
	"html/template"
	"net/http"
)

type TemplateExecutor interface {
	Execute(w http.ResponseWriter, data interface{}) error
}

type DebugTemplateExecutor struct {
	Filepath string
}

type ReleaseTemplateExecutor struct {
	Template *template.Template
}

func NewReleaseTemplateExecutor(filepath string) ReleaseTemplateExecutor {
	return ReleaseTemplateExecutor{
		Template: template.Must(template.ParseFiles(filepath)),
	}
}

func (e DebugTemplateExecutor) Execute(w http.ResponseWriter, data interface{}) error {
	t := template.Must(template.ParseFiles(e.Filepath))

	return t.Execute(w, data)
}

func (e ReleaseTemplateExecutor) Execute(w http.ResponseWriter, data interface{}) error {
	return e.Template.Execute(w, data)
}
