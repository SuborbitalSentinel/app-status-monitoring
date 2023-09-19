package service

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Registry struct {
	services         map[string]State
	idToName         map[string]string
	parentToChildren map[string][]string
	mutex            sync.Mutex
}

func NewRegistry() Registry {
	return Registry{
		services:         make(map[string]State),
		idToName:         make(map[string]string),
		parentToChildren: make(map[string][]string),
		mutex:            sync.Mutex{},
	}
}

func (r *Registry) Reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.services = make(map[string]State)
	r.idToName = make(map[string]string)
	r.parentToChildren = make(map[string][]string)
}

 func(r *Registry) SetServiceName(serviceId string, serviceName string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.idToName[serviceId] = serviceName
}

func (r *Registry) GetServiceName(serviceId string) string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if val, ok := r.idToName[serviceId]; ok {
		return val
	}

	var hosts []string
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "HOST_") {
			hosts = append(hosts, strings.Split(v, "=")[1])
		}
	}

	for _, host := range hosts {
		res, err := http.Get(host)
		if err != nil {
			log.Printf("Failed to complete request for container names from Host: %s\n", host)
			log.Println(err)
			continue
		}
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		lines := strings.Split(string(body), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				parts := strings.Split(line, " ")
				serviceId := strings.TrimSpace(parts[0])
				serviceName := strings.TrimSpace(parts[1])

				r.idToName[serviceId] = serviceName
			}
		}
	}

	return r.idToName[serviceId]
}

func (r *Registry) States() []State {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	states := []State{}
	for _, state := range r.services {
		states = append(states, state)
	}
	return states
}

func (r *Registry) Set(state State) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.services[state.Id] = state

	if state.Relationship == Child && !contains(r.parentToChildren[state.ParentKey], state.Id) {
		r.parentToChildren[state.ParentKey] = append(r.parentToChildren[state.ParentKey], state.Id)
	}
}

func (r *Registry) TryGet(id string) (State, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if s, ok := r.services[id]; ok {
		return s, ok
	}
	return State{}, false
}

func contains(slice []string, element string) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}
