package main

import log "github.com/sirupsen/logrus"

func main() {
	log.WithFields(log.Fields{
		"event": "event",
		"topic": "topic",
		"key":   "key",
	}).Fatal("Failed to send event")
}

// FATA[0000] Failed to send event                          event=event key=key topic=topic
