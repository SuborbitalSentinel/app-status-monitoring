package service

import "time"

type State struct {
	Id            string       // The unique identifier of the service.
	Name          string       // The human-readable name of the service
	LastHeartbeat time.Time    // The last time we have heard from this service
	HasAlerted    bool         // Whether or not this services has been alerted to a source
	ParentKey     string       // The unique identifier of the parent service; If this service is a child of another service, this will be the unique identifier of the parent service
	Relationship  Relationship // What type of service this is [Child, Parent, Standalone]
}

func (s State) WasCheckinMissed(currentTime *time.Time) bool {
	return s.LastHeartbeat.Add(30 * time.Second).Before(*currentTime)
}

// This is how we do enums in golang
type Relationship int

const (
	Parent Relationship = iota
	Child
	Standalone
)

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
