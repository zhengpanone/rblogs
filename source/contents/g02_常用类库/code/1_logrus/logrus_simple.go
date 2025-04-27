package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithFields(log.Fields{
		"name": "test",
	}).Info("Test logrus info logs")
}

// Output INFO[0000] Test logrus info logs                         name=test
