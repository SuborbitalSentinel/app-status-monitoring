package service

import "sync"

type Registry struct {
	services map[string]State
	mutex    sync.Mutex
}

func NewRegistry() Registry {
	return Registry{
		services: make(map[string]State),
		mutex:    sync.Mutex{},
	}
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
	r.services[state.Name] = state
}

func (r *Registry) TryGet(name string) (State, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, ok := r.services[name]; ok {
		return r.services[name], ok
	}
	return State{}, false
}
