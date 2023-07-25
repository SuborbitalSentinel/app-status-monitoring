package service

import (
	"log"
	"time"
)

func (r *Registry) ConsoleMonitor() {
	for {
		now := time.Now()

		for _, state := range r.States() {
			checkinMissed := state.WasCheckinMissed(&now)
			if checkinMissed && !state.HasAlerted {
				state.HasAlerted = true
				r.Set(state)
				log.Printf("Service: %s, appears to be down!\n", state.Name)
			} else if !checkinMissed && state.HasAlerted {
				state.HasAlerted = false
				r.Set(state)
				log.Printf("Service: %s, is back online!\n", state.Name)
			}
		}

		time.Sleep(10 * time.Second)
	}
}
