package logging

import (
	"fmt"
	"os"

	setting "logrus_demo/gin_logrus_packing/config"

	"github.com/sirupsen/logrus"
)

var WebLog *logrus.Logger

func Init() {
	initWebLog()
}

func initWebLog() {

	WebLog = initLog(setting.Conf.LogConfig.WebLogName)
}

// 初始化日志句柄
func initLog(logFileName string) *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}

	logFilePath := setting.Conf.LogFilePath
	logName := logFilePath + logFileName
	var f *os.File
	var err error
	//判断日志文件夹是否存在，不存在则创建
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		os.MkdirAll(logFilePath, os.ModePerm)
	}
	//判断日志文件是否存在，不存在则创建，否则就直接打开
	if _, err := os.Stat(logName); os.IsNotExist(err) {
		f, err = os.Create(logName)
	} else {
		f, err = os.OpenFile(logName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	}

	if err != nil {
		fmt.Println("open log file failed")
	}

	log.Out = f
	log.Level = logrus.InfoLevel
	return log
}
