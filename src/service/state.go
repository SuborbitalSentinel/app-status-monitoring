package service

import "time"

type State struct {
	Name          string
	LastHeartbeat time.Time
	HasAlerted    bool
}

func (s State) WasCheckinMissed(currentTime *time.Time) bool {
	return s.LastHeartbeat.Add(1 * time.Minute).Before(*currentTime)
}
