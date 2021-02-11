package sigx

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Listen(fn func(os.Signal)) {
	go func() {
		// we use buffered to mitigate losing the signal
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

		sig := <-sigchan
		if fn != nil {
			fn(sig)
		}
	}()
}

// ListenExit will listen to OS signals (currently SIGINT, SIGKILL, SIGTERM)
// and will trigger the callback when signal are received from OS
func ListenExit(fn func(os.Signal)) {
	go func() {
		// we use buffered to mitigate losing the signal
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, os.Kill, syscall.SIGTERM)

		sig := <-sigchan
		if fn != nil {
			fn(sig)
		}
	}()
}

// AwaitReload waits for
func AwaitReload(lg *zap.SugaredLogger, sigChan chan os.Signal, closerWait chan bool, closer func()) chan bool {
	// The blocking signal handler that main() waits on.
	out := make(chan bool)

	// Respawn a new process and exit the running one.
	respawn := func() {
		if err := syscall.Exec(os.Args[0], os.Args, os.Environ()); err != nil {
			lg.Fatalf("error spawning process: %v", err)
		}
		os.Exit(0)
	}

	// Listen for reload signal.
	go func() {
		for range sigChan {
			lg.Infof("reloading on signal ...")

			go closer()
			select {
			case <-closerWait:
				// Wait for the closer to finish.
				respawn()
			case <-time.After(time.Second * 3):
				// Or timeout and force close.
				respawn()
			}
		}
	}()

	return out
}
