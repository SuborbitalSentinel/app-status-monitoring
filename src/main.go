/*
Name: Service Monitor
Description: A service monitor is used to monitor the health of various in-house services.
	It's designed to be fairly flexable and unopinionated, except for dockerized services.
	For dockerized services, it will attempt to asertain the friendly name by querying
	host machines for their container id / friendly name pairs. It assumes that it's a dockerized service
	if the service-name is missing from the /alive Post body.

	The host machines that will be queried are any environment variables beginning with: HOST_
Endpoints:
	/:
		Returns a visualized dashboard of the services that are being monitored.
	/reset:
		Resets the service registry (good to call after redeploying services)
	/alive form values:
		1) service-id: [Required] [String] unique identifier. (recomend using some sort of UUID).
		2) relationship: [Optional] ["Parent" | "Child" | ""] Defines whether the service is a parent, child, or standalone.
		3) parent-key: [Optional] Child and Parent services should share the same parent-key. Ignored for standalone services.
		4) service-name: [Optional] [String] human-readable name. If not provided, assumed that service is dockerized; container name resolved automagically
*/

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

	// TODO: Replace DebugTemplateExecutor with a ReleaseTemplateExecutor
	// releaseTemplates := util.NewReleaseTemplateExecutor("./templates/index.gotmpl")

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
						ServiceName:      state.Name,
						MissedCheckIn:    state.WasCheckinMissed(&now),
						LastHeartbeat:    state.LastHeartbeat.Format("Mon Jan 2 15:04:05 MST 2006"),
						OnlineChildCount: len(childServices) - offlineChildCount,
						ChildServices:    childServices,
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
		serviceId := r.FormValue("service-id")
		serviceName := r.FormValue("service-name")
		relationship := r.FormValue("relationship")
		parentKey := r.FormValue("parent-key")
		log.Printf("Received /alive request: '%s', '%s', '%s', '%s'", serviceId, serviceName, relationship, parentKey)

		if s, ok := registry.TryGet(serviceId); ok {
			// TODO: Try to remove this line by fixing the dirty read issue.
			s.Name = registry.GetServiceName(serviceId)

			s.LastHeartbeat = time.Now()
			registry.Set(s)
		} else {
			var name string
			if serviceName != "" {
				name = serviceName
			} else {
				name = registry.GetServiceName(serviceId)
			}

			registry.Set(service.State{
				Id:            serviceId,
				Name:          name,
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
