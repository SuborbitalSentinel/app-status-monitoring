package home

import (
	"monitor/util"
	"net/http"
)

type ServiceData struct {
	ServiceName   string
	MissedCheckIn bool
	LastHeartbeat string
}

type Handler struct {
	Template          util.TemplateExecutor
	CreateServiceData func() []ServiceData
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Services": h.CreateServiceData(),
	}
	h.Template.Execute(w, data)
}
