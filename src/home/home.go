package home

import (
	"monitor/util"
	"net/http"
)

// This is the data structure that will be used to fill in the html template for the home page.
type ServiceData struct {
	ServiceName      string
	MissedCheckIn    bool
	LastHeartbeat    string
	OnlineChildCount int
	ChildServices    []ServiceData
}

// This type is a bit of indirection to allow us to specify the Template to use and how we derive the data for the Home page
// The primary reason is so this module doesn't need to care about what's required to create the data structs
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
