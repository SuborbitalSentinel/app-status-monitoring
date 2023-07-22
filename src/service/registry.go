package service

import "sync"

type Registry struct {
	services         map[string]State
	parentToChildren map[string][]string
	mutex            sync.Mutex
}

func NewRegistry() Registry {
	return Registry{
		services:         make(map[string]State),
		parentToChildren: make(map[string][]string),
		mutex:            sync.Mutex{},
	}
}

func contains(slice []string, element string) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
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
