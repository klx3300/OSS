package evbus

// Subscriber is entities listening on the event bus
type Subscriber struct {
	SubscriberId int
	// Priority ranged from -999 ~ 999.
	SubscriberPriority int
	// Not subscribing anything means subscribing everything.
	// Keep that in mind..
	SubscribedTypes []int
	MessageChannel  chan Event
	// The bool value returned in this functor will determine
	// whether this event is consumed, preventing further subscribers
	// to receive it. true -> consumed
	// Anyway, an event will be consumed after iterating through every subscribers resulted in
	// any matching subscriber found.
	// Warning: dont run time-consuming stuff in it, return as soon as possible!
	MessageFunctor func(Subscriber, Event) bool
	Extra          interface{}
}
