package evbus

import "time"

// Event represents each individual event
type Event struct {
	EventId uint32
	// a no type event will only dispatch to wildcard subscribers.
	EventType []int
	EventTime time.Time
	Payload   interface{}
}
