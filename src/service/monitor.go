package service

import (
	"log"
	"time"
)

// If the current pattern is followed, only one Monitor can be used at any given time.
// To fix this, we would need to move the state of if we have or haven't alerted to something else so the Monitors don't fight

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
