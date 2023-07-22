package home

import (
	"html/template"
	"net/http"
)

type Data struct {
	ServiceName   string
	MissedCheckIn bool
	LastHeartbeat string
}

type Handler struct {
	Template *template.Template
	Status   func() []Data
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Services": h.Status(),
	}
	h.Template.Execute(w, data)
}
