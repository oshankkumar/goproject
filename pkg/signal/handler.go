package signal

import (
	"os"
	"os/signal"
	"syscall"
)

var closeOnce = make(chan struct{})

var signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

func SetupHandler() <-chan struct{} {
	close(closeOnce)

	stop := make(chan struct{})

	sigChan := make(chan os.Signal, len(signals))
	signal.Notify(sigChan, signals...)

	go func() {
		<-sigChan
		close(stop)
		<-sigChan
		os.Exit(1)
	}()

	return stop
}
