package main

import log "github.com/sirupsen/logrus"

func main() {
	log.Trace("Something very low level.")
	log.Debug("Useful debugging information.")
	log.Info("Something noteworthy happened!")
	log.Warn("You should probably take a look at this.")
	log.Error("Something failed but I'm not quitting.")
	// 记完日志后会调用os.Exit(1)
	log.Fatal("Bye.")
	// 记完日志后会调用 panic()
	log.Panic("I'm bailing.")
}

/**
INFO[0000] Something noteworthy happened!
WARN[0000] You should probably take a look at this.
ERRO[0000] Something failed but I'm not quitting.
FATA[0000] Bye.
*/
