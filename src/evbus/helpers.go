package evbus

import (
	"container/list"
	"sync"
)

// SubscriberMap is the map of subscriber id(int) -> subscriber(Subscriber)
type SubscriberMap sync.Map

// EventTypeMap is the map of name(string) -> event type id(int32)
type EventTypeMap sync.Map

// DispatchMap is the shortcut when bus need to dispatch every event.
// Mapping is event type id(int) -> subscriber map(*SubscriberMap)
type DispatchMap sync.Map

func asSubscriber(liter *list.Element) *Subscriber {
	return liter.Value.(*Subscriber)
}

// NewEventTypeMap creates a new instance.
func NewEventTypeMap() *EventTypeMap {
	return (*EventTypeMap)(new(sync.Map))
}

// Access the event type map.
func (m *EventTypeMap) Access(key string) (int32, bool) {
	tmpM := ((*sync.Map)(m))
	ld, ldok := tmpM.Load(key)
	if !ldok {
		return 0, false
	}
	return ld.(int32), true
}

// AccessOrInit explains itself, return false on initialized
func (m *EventTypeMap) AccessOrInit(key string, value int32) (int32, bool) {
	tmpM := ((*sync.Map)(m))
	ld, ldok := tmpM.LoadOrStore(key, value)
	if !ldok {
		return value, false
	}
	return ld.(int32), true
}

// SetValue setup the mapping
func (m *EventTypeMap) SetValue(key string, value int32) {
	((*sync.Map)(m)).Store(key, value)
}

// Iterate requires you to read sync.Map.Range
func (m *EventTypeMap) Iterate(itfunc func(string, int32) bool) {
	((*sync.Map)(m)).Range(func(x, y interface{}) bool {
		return itfunc(x.(string), y.(int32))
	})
}

// Remove the given key from the map.
func (m *EventTypeMap) Remove(key string) {
	((*sync.Map)(m)).Delete(key)
}

// NewDispatchMap creates a new instance.
func NewDispatchMap() *DispatchMap {
	return (*DispatchMap)(new(sync.Map))
}

// Access the event type map.
func (m *DispatchMap) Access(key int32) (*SubscriberMap, bool) {
	tmpM := ((*sync.Map)(m))
	ld, ldok := tmpM.Load(key)
	if !ldok {
		return nil, false
	}
	return ld.(*SubscriberMap), true
}

// AccessOrInit explains itself, return false on initialized
func (m *DispatchMap) AccessOrInit(key int32, value *SubscriberMap) (*SubscriberMap, bool) {
	tmpM := ((*sync.Map)(m))
	ld, ldok := tmpM.LoadOrStore(key, value)
	if !ldok {
		return value, false
	}
	return ld.(*SubscriberMap), true
}

// SetValue setup the mapping
func (m *DispatchMap) SetValue(key int32, value *SubscriberMap) {
	((*sync.Map)(m)).Store(key, value)
}

// Iterate requires you to read sync.Map.Range
func (m *DispatchMap) Iterate(itfunc func(int32, *SubscriberMap) bool) {
	((*sync.Map)(m)).Range(func(x, y interface{}) bool {
		return itfunc(x.(int32), y.(*SubscriberMap))
	})
}

// Remove the given key from the map.
func (m *DispatchMap) Remove(key int32) {
	((*sync.Map)(m)).Delete(key)
}

// NewSubscriberMap creates a new instance.
func NewSubscriberMap() *SubscriberMap {
	return (*SubscriberMap)(new(sync.Map))
}

// Access the event type map.
func (m *SubscriberMap) Access(key int32) (Subscriber, bool) {
	tmpM := ((*sync.Map)(m))
	ld, ldok := tmpM.Load(key)
	if !ldok {
		return Subscriber{}, false
	}
	return ld.(Subscriber), true
}

// AccessOrInit explains itself, return false on initialized
func (m *SubscriberMap) AccessOrInit(key int32, value Subscriber) (Subscriber, bool) {
	tmpM := ((*sync.Map)(m))
	ld, ldok := tmpM.LoadOrStore(key, value)
	if !ldok {
		return value, false
	}
	return ld.(Subscriber), true
}

// SetValue setup the mapping
func (m *SubscriberMap) SetValue(key int32, value Subscriber) {
	((*sync.Map)(m)).Store(key, value)
}

// Iterate requires you to read sync.Map.Range
func (m *SubscriberMap) Iterate(itfunc func(int32, Subscriber) bool) {
	((*sync.Map)(m)).Range(func(x, y interface{}) bool {
		return itfunc(x.(int32), y.(Subscriber))
	})
}

// Remove the given key from the map.
func (m *SubscriberMap) Remove(key int32) {
	((*sync.Map)(m)).Delete(key)
}
