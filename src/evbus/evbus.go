package evbus

import (
	"errtypes"
	"logger"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

var log = logger.Logger{
	LogLevel: 0,
	Name:     "EventBus:",
}

// EventBus is a bus for events dispatched to it.
type EventBus struct {
	Events                  sync.Pool
	EventCounter            uint32
	PendingCounter          uint32
	PendingThreshold        uint32
	Subscribers             *SubscriberMap
	WildcardSubscribers     *SubscriberMap
	SubscriberCounter       int32
	DispatchList            *DispatchMap
	EventTypeList           *EventTypeMap
	EventTypeCounter        int32
	shutdownSignal          bool
	EachEventBackfileBase   time.Duration
	EachEventBackfireThresh time.Duration
	EventLoopBackfireBase   time.Duration
	EventLoopBackfireThresh time.Duration
}

// NewEventBus returns a brand new event bus instance.
func NewEventBus() *EventBus {
	evbus := new(EventBus)
	evbus.Subscribers = NewSubscriberMap()
	evbus.WildcardSubscribers = NewSubscriberMap()
	evbus.DispatchList = NewDispatchMap()
	evbus.EventTypeList = NewEventTypeMap()
	evbus.SubscriberCounter = 0
	evbus.EventTypeCounter = 0
	evbus.EventCounter = 0
	evbus.PendingCounter = 0
	// This is pretty magic, but you can set by urself.
	// Just don't mess things up.
	evbus.PendingThreshold = 256
	evbus.shutdownSignal = false
	evbus.EachEventBackfileBase = 10 * time.Millisecond
	evbus.EachEventBackfireThresh = 500 * time.Millisecond
	evbus.EventLoopBackfireBase = 1 * time.Millisecond
	evbus.EventLoopBackfireThresh = 10 * time.Millisecond
	go evbus.eventDispatcher()
	return evbus
}

// Finalize stops the event dispatching process after that, dont use it!
func (evbus *EventBus) Finalize() {
	evbus.shutdownSignal = true
}

// RegisterEventType register the given name, and return the event type id on succeed.
func (evbus *EventBus) RegisterEventType(typename string) (int32, error) {
	if evbus.shutdownSignal {
		return 0, errtypes.EBADF
	}
	_, ldok := evbus.EventTypeList.Access(typename)
	if ldok {
		return 0, errtypes.EEXIST
	}
	evid := atomic.AddInt32(&(evbus.EventTypeCounter), 1) - 1
	evbus.EventTypeList.SetValue(typename, evid)
	evbus.DispatchList.SetValue(evid, NewSubscriberMap())
	return evid, nil
}

// GetEventTypeId doesn't need to be access from outside
func (evbus *EventBus) GetEventTypeId(typename string) (int32, error) {
	if evbus.shutdownSignal {
		return 0, errtypes.EBADF
	}
	ld, ldok := evbus.EventTypeList.Access(typename)
	if ldok {
		return ld, nil
	}
	return ld, errtypes.ENEXIST
}

// Subscribe to some events. parameter sb.id will be ignored and the returned integer
// will be the final IDs.
func (evbus *EventBus) Subscribe(sb Subscriber) int {
	if evbus.shutdownSignal {
		log.Warnln("EventBus: Subscribing into finalized event bus.")
		return -1
	}
	if sb.SubscriberPriority > 999 {
		log.Warnln("EventBus: Subscriber priority", sb.SubscriberPriority, "larger than 999. Capped.")
		sb.SubscriberPriority = 999
	}
	if sb.SubscriberPriority < -999 {
		log.Warnln("EventBus: Subscriber priority", sb.SubscriberPriority, "lesser than -999. Capped.")
		sb.SubscriberPriority = -999
	}
	sbid := atomic.AddInt32(&(evbus.SubscriberCounter), 1) - 1
	sb.SubscriberId = int(sbid)
	evbus.Subscribers.SetValue(sbid, sb)
	if len(sb.SubscribedTypes) == 0 {
		evbus.WildcardSubscribers.SetValue(sbid, sb)
	} else {
		for i := 0; i < len(sb.SubscribedTypes); i++ {
			evtmap, evtok := evbus.DispatchList.Access(int32(sb.SubscribedTypes[i]))
			if !evtok {
				log.Infoln("EventBus: Registering subscriber", sbid, "reached invalid event type", sb.SubscribedTypes[i])
			} else {
				evtmap.SetValue(sbid, sb)
			}
		}
	}
	return int(sbid)
}

// Unsubscribe to events. The sbid must be provided.
func (evbus *EventBus) Unsubscribe(sbid int) error {
	if evbus.shutdownSignal {
		return errtypes.EBADF
	}
	sb, existok := evbus.Subscribers.Access(int32(sbid))
	if !existok {
		return errtypes.ENEXIST
	}
	if len(sb.SubscribedTypes) == 0 {
		evbus.WildcardSubscribers.Remove(int32(sbid))
	} else {
		for i := 0; i < len(sb.SubscribedTypes); i++ {
			evtmap, evtok := evbus.DispatchList.Access(int32(sb.SubscribedTypes[i]))
			if !evtok {
				log.Warnln("EventBus: Unregistering subscriber", sbid, "reached invalid event type", sb.SubscribedTypes[i])
			} else {
				evtmap.Remove(int32(sb.SubscribedTypes[i]))
			}
		}
	}
	evbus.Subscribers.Remove(int32(sbid))
	return nil
}

// PublishEvent will ignore the event id in parameter.
func (evbus *EventBus) PublishEvent(event Event) error {
	if evbus.shutdownSignal {
		return errtypes.EBADF
	}
	tryup := atomic.AddUint32(&(evbus.PendingCounter), 1) - 1
	if tryup >= evbus.PendingThreshold {
		atomic.AddUint32(&(evbus.PendingCounter), ^uint32(0))
		return errtypes.EAGAIN
	}
	// Generate an id for it
	evid := atomic.AddUint32(&(evbus.EventCounter), 1)
	event.EventId = evid
	evbus.Events.Put(event)
	return nil
}

type subscriberByPriority []Subscriber

func (s subscriberByPriority) Len() int {
	return len(s)
}
func (s subscriberByPriority) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s subscriberByPriority) Less(i, j int) bool {
	return s[i].SubscriberPriority < s[j].SubscriberPriority
}

func (evbus *EventBus) handleEachEvent(event Event) {
	// use backfire to reduce contention
	backfireTime := evbus.EachEventBackfileBase
	for true {
		if evbus.shutdownSignal {
			return
		}
		subscriberList := make([]Subscriber, 0)
		slPtr := &subscriberList
		evbus.WildcardSubscribers.Iterate(func(i int32, sb Subscriber) bool {
			*(slPtr) = append(*slPtr, sb)
			return true
		})
		for i := 0; i < len(event.EventType); i++ {
			etmap, etok := evbus.DispatchList.Access(int32(event.EventType[i]))
			if !etok {
				log.Warnln("EventBus: Handling event with invalid event type", event.EventType[i])
			} else {
				etmap.Iterate(func(i int32, sb Subscriber) bool {
					*(slPtr) = append(*slPtr, sb)
					return true
				})
			}
		}
		if len(subscriberList) == 0 {
			log.Infoln("EventBus: event type", event.EventType, "has no subscriber.")
			time.Sleep(backfireTime)
			backfireTime = backfireTime * 2
			if backfireTime >= evbus.EachEventBackfireThresh {
				backfireTime = evbus.EachEventBackfireThresh
			}
			continue
		}
		// dispatch it
		if evbus.shutdownSignal {
			return
		}
		sort.Sort(subscriberByPriority(subscriberList))
		for i := len(subscriberList) - 1; i >= 0; i-- {
			rv := subscriberList[i].MessageFunctor(subscriberList[i], event)
			if rv == true {
				break
			}
			if evbus.shutdownSignal {
				return
			}
		}
		break
	}
}

func (evbus *EventBus) eventDispatcher() {
	backfireTime := evbus.EventLoopBackfireBase
	for true {
		if evbus.shutdownSignal {
			return
		}
		pendEv := evbus.Events.Get()
		if pendEv == nil {
			time.Sleep(backfireTime)
			backfireTime += time.Millisecond
			if backfireTime >= evbus.EventLoopBackfireThresh {
				backfireTime = evbus.EventLoopBackfireThresh
			}
			continue
		}
		backfireTime = evbus.EventLoopBackfireBase
		atomic.AddUint32(&(evbus.PendingCounter), ^uint32(0))
		ev, ok := pendEv.(Event)
		if !ok {
			log.Warnln("EventBus: Garbage", pendEv, "in pending event pool.")
			continue
		}
		go evbus.handleEachEvent(ev)
	}
}
