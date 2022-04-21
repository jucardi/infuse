package paths

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jucardi/go-osx/log"
)

func listenForShutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit

		log.Get().Infof("signal captured: %s", sig.String())
		cleanup()
		os.Exit(0)
	}()
}

func cleanup() {
	for _, v := range created {
		if err := os.RemoveAll(v); err != nil {
			log.Get().Warnf("Unable to remove temporary directory, %s", v)
		} else {
			log.Get().Debugf("Temp Dir '%s' deleted.", v)
		}
	}
}
