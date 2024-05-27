package cron

import (
	"log"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	duration := time.Second * 5
	handlerFunc := func() {
		log.Println("hello world")
	}
	cron := New(WithDuration(duration), WithHandlerFunc(handlerFunc))
	err := cron.Start()
	if err != nil {
		return
	}
	defer cron.Stop()
}
