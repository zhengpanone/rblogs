package main

import log "github.com/sirupsen/logrus"

func main() {
	requestLogger := log.WithFields(log.Fields{
		"request_id": "request_id",
		"user_ip":    "user_ip",
	})
	requestLogger.Info("something happened on that request") // will log request_id and user_ip
	requestLogger.Warn("something not great happened")
}

/*
INFO[0000] something happened on that request            request_id=request_id user_ip=user_ip
WARN[0000] something not great happened                  request_id=request_id user_ip=user_ip
*/
