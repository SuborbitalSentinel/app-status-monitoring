package home

import (
	"html/template"
	"monitor/service"
	"net/http"
)

type Handler struct {
	Template *template.Template
	Status   func() []service.Data
}

func (self *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Services": self.Status(),
	}
	self.Template.Execute(w, data)
}
