package main

import (
	"dbconn"
	"time"
)

func preloginInitializer() {
	<-time.After(1 * time.Second)
	dbconn.Initialization(mainEventBus)
	log.Debugln("Spawning login...")
	currwnd.Hide()
	chanTransist <- WND_LOGIN
}
