package main

import "github.com/sirupsen/logrus"

var log = logrus.New()

func main() {
	log.Formatter = &logrus.JSONFormatter{}
	log.SetReportCaller(true) // 可以开启记录函数名，但是会消耗性能
	log.WithFields(logrus.Fields{
		"event": "event",
		"topic": "topic",
		"key":   "key",
	}).Info("Failed to send event")
}

/**
{"event":"event",
"file":"/Users/mac/code/blog/gblogs/source/contents/g02_常用类库/code/logrus/format/logrus_format.go:14",
"func":"main.main",
"key":"key",
"level":"info",
"msg":"Failed to send event",
"time":"2024-05-24T12:42:19+08:00",
"topic":"topic"
}
*/
