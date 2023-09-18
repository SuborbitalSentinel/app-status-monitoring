package util

import (
	"html/template"
	"net/http"
)

// This whole file exists to facilitate the use of a Release or Debug html template.
// If the Debug template is used, changes can be made to the template without restarting the server.
// If the Release template is used, changes to the template will not be reflected until the server is restarted.

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
