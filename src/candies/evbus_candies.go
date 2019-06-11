package candies

import (
	"evbus"
	"time"
)

// EventTypeUnrecoverableFail represents some unrecoverable failure happened.
// os.Exit() will be called after clean-ups done.
// payload is its reason.
var EventTypeUnrecoverableFail int

// EventTypeStatusReport represents some status change that should be shown to user.
// payload is the string.
var EventTypeStatusReport int

// CandyEvent is the candy-wrapped event bus event
type CandyEvent evbus.Event

// NewEvent return a candy event object
func NewEvent(payload interface{}) *CandyEvent {
	d := new(CandyEvent)
	d.Payload = payload
	return d
}

// Type append a new type to the candy event
func (c *CandyEvent) Type(evtype int) *CandyEvent {
	c.EventType = append(c.EventType, evtype)
	return c
}

// Fin return a acceptable obj
func (c *CandyEvent) Fin() evbus.Event {
	return evbus.Event(*c)
}

// NewChannelSubscriber is a simple subscriber that just put things into that channel
func NewChannelSubscriber(evtype int, pri int, consumed bool, theChannel chan evbus.Event) evbus.Subscriber {
	var s evbus.Subscriber
	s.SubscribedTypes = []int{evtype}
	s.SubscriberPriority = pri
	s.MessageChannel = theChannel
	s.MessageFunctor = func(sub evbus.Subscriber, ev evbus.Event) bool {
		sub.MessageChannel <- ev
		return consumed
	}
	return s
}

// FastGG simply declares unrecoverable happened.
func FastGG(info string, ebus *evbus.EventBus) {
	ebus.PublishEvent(NewEvent(info).Type(EventTypeUnrecoverableFail).Fin())
	// wait for GG to happen
	<-time.After(5 * time.Second)
}

// FastRegister frees you from handling errors
func FastRegister(name string, ebus *evbus.EventBus) int {
	r, err := ebus.RegisterEventType(name)
	if err != nil {
		FastGG("Unable to register "+name+":"+err.Error(), ebus)
		return 0
	}
	return int(r)
}
