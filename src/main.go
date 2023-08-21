package main

import (
	"log"
	"monitor/home"
	"monitor/service"
	"monitor/util"
	"net/http"
	"sort"
	"time"
)

var registry = service.NewRegistry()

func main() {
	go registry.ConsoleMonitor()

	home := &home.Handler{
		Template: util.DebugTemplateExecutor{Filepath: "./templates/index.gotmpl"},
		CreateServiceData: func() []home.ServiceData {
			now := time.Now()
			serviceDto := make([]home.ServiceData, 0)
			states := registry.States()
			sort.Slice(states, func(lhs, rhs int) bool {
				if states[lhs].WasCheckinMissed(&now) != states[rhs].WasCheckinMissed(&now) {
					return states[lhs].WasCheckinMissed(&now)
				}
				return states[lhs].Name < states[rhs].Name
			})

			for _, state := range states {
				childServices := make([]home.ServiceData, 0)
				if state.Relationship == service.Parent {
					for _, childState := range states {
						if childState.Relationship == service.Child && childState.ParentKey == state.ParentKey {
							childServices = append(childServices, home.ServiceData{
								ServiceName:   childState.Name,
								MissedCheckIn: childState.WasCheckinMissed(&now),
								LastHeartbeat: childState.LastHeartbeat.Format("Jan 2 15:04:05"),
							})
						}
					}
				}
				offlineChildCount := 0
				for _, childState := range childServices {
					if childState.MissedCheckIn {
						offlineChildCount += 1
					}
				}
				if state.Relationship != service.Child {
					serviceDto = append(serviceDto, home.ServiceData{
						ServiceName:   state.Name,
						MissedCheckIn: state.WasCheckinMissed(&now),
						LastHeartbeat: state.LastHeartbeat.Format("Mon Jan 2 15:04:05 MST 2006"),
						OnlineChildCount: len(childServices) - offlineChildCount,
						ChildServices: childServices,
					})

				}
			}

			return serviceDto
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
