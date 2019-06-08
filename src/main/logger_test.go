package main

import "testing"
import "logger"

func TestLogger(t *testing.T) {
	logger.Log.Debugln("Test debug.")
	logger.Log.Infoln("Test information.")
	logger.Log.Succln("Test success.")
	logger.Log.Warnln("Test warning.")
	logger.Log.FailLn("Test failure.")
}

// func TestLoggerPerf(t *testing.T) {
// 	completionCh := make(chan int)
// 	for i := 0; i < 32; i++ {
// 		go func(i int) {
// 			for j := 0; j < 10000; j++ {
// 				logger.Log.Warnln("Thread", i, "Epoch", j)
// 			}
// 			completionCh <- 0
// 		}(i)
// 	}
// 	for i := 0; i < 32; i++ {
// 		<-completionCh
// 	}
// }
