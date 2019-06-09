package candies

import (
	"evbus"
	"logger"
	"os"
)

var log = logger.Logger{
	LogLevel: 0,
	Name:     "Candy:",
}

// Initialization finishes the initialization process
func Initialization(bus *evbus.EventBus) {
	log.Debugln("Initializing...")
	var everr error
	var tmpevtid int32
	tmpevtid, everr = bus.RegisterEventType("UnrecoverableFailure")
	if everr != nil {
		log.FailLn("Unable to register event type UnrecoverableFailure:", everr.Error())
		log.FailLn("This failure itself is unrecoverable. ::QUIT.")
		os.Exit(-1)
	}
	EventTypeUnrecoverableFail = int(tmpevtid)
	tmpevtid, everr = bus.RegisterEventType("UserStatusReport")
	if everr != nil {
		log.FailLn("Unable to register event type UserStatusReport:", everr.Error())
		log.FailLn("This failure itself is unrecoverable. ::QUIT.")
		os.Exit(-1)
	}
	EventTypeStatusReport = int(tmpevtid)
	log.Debugln("Initialized.")
}
