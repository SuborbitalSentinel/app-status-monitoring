package service

import "time"

type State struct {
	Id            string
	Name          string
	LastHeartbeat time.Time
	HasAlerted    bool
	ParentKey     string
	Relationship  Relationship
}

func (s State) WasCheckinMissed(currentTime *time.Time) bool {
	return s.LastHeartbeat.Add(1 * time.Minute).Before(*currentTime)
}

type Relationship int

func (r Relationship) ToString() string {
	switch r {
	case Child:
		return "child"
	case Parent:
		return "parent"
	default:
		return "standalone"
	}
}

const (
	Parent Relationship = iota
	Child
	Standalone
)

func ToRelationship(relationship string) Relationship {
	switch relationship {
	case "parent":
		return Parent
	case "child":
		return Child
	default:
		return Standalone
	}
}
