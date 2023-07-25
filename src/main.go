package main

import (
	"html/template"
	"log"
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
		CreateServiceData: func() []home.ServiceData {
			now := time.Now()
			status := make([]home.ServiceData, 0)

			for _, state := range registry.States() {
				status = append(status, home.ServiceData{
					ServiceName:   state.Name,
					MissedCheckIn: state.WasCheckinMissed(&now),
					LastHeartbeat: state.LastHeartbeat.Format(time.RFC3339),
				})
			}

			return status
		},
	}
	http.HandleFunc("/", home.ServeHTTP)

	http.HandleFunc("/reset", func(w http.ResponseWriter, _ *http.Request) {
		registry.Reset()
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
		}
		serviceId := r.FormValue("service")
		relationship := r.FormValue("relationship")
		parentKey := r.FormValue("parent-key")
		log.Printf("Received /alive request: '%s', '%s', '%s'", serviceId, relationship, parentKey)

		if s, ok := registry.TryGet(serviceId); ok {
			s.Name = registry.GetServiceName(serviceId)
			s.LastHeartbeat = time.Now()
			registry.Set(s)
		} else {
			registry.Set(service.State{
				Id:            serviceId,
				Name:          registry.GetServiceName(serviceId),
				LastHeartbeat: time.Now(),
				HasAlerted:    false,
				Relationship:  service.ToRelationship(relationship),
				ParentKey:     parentKey,
			})
		}

		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":1911", nil)
}
