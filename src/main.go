package main

import (
	"html/template"
	"monitor/home"
	"monitor/service"
	"net/http"
	"time"
)

var registry = service.NewRegistry()

func main() {
	go registry.ConsoleMonitor()

	home := &home.Handler{
		Template: template.Must(template.ParseFiles("./templates/index.html")),
		Status: func() []service.Data {
			now := time.Now()
			status := make([]service.Data, 0)

			for _, state := range registry.States() {
				status = append(status, service.Data{
					ServiceName:   state.Name,
					MissedCheckIn: state.WasCheckinMissed(&now),
					LastHeartbeat: state.LastHeartbeat.Format(time.RFC3339),
				})
			}

			return status
		},
	}
	http.HandleFunc("/", home.ServeHTTP)

	http.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		serviceName := r.FormValue("site")

		if s, ok := registry.TryGet(serviceName); ok {
			s.LastHeartbeat = time.Now()
			registry.Set(s)
		} else {
			registry.Set(service.State{
				Name:          serviceName,
				LastHeartbeat: time.Now(),
				HasAlerted:    false,
			})
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		// w.WriteHeader(200)

	})

	http.ListenAndServe(":1911", nil)
}
