package lifecycle

import (
	"container/list"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
)

var (
	lockMut  sync.RWMutex
	sigList  map[os.Signal]*list.List
	stopWait sync.WaitGroup
	sigChan  chan os.Signal
)

func init() {
	// We want to always be able to wait for done, so allready boot that one, ensure it gets executed last
	// stopWait gets decremented after closing of channel
	stopWait.Add(1)
	Finally(func() {
		// Close channel, therefor exit goroutine, and clean the lists to prevent double exists
		close(sigChan)
	})
}

// PanicHandler must be called as a deferred method ensuring that cleanup tasks
// are handled when a panic occurs. It attempts to log the generated panic and
// calls all registered callbacks when a panic occurs.
func PanicHandler() {
	defer func() {
		if err := recover(); err != nil && sigChan != nil {
			fmt.Println("Panic received, attempting gracefull exit; Error:", err)
			defer os.Exit(1)
		}
	}()

	if err := recover(); err != nil && sigChan != nil {
		if errC, ok := err.(error); ok {
			log.Error().Stack().Err(errC).Msg("panic received, attempting gracefull exit")
		} else {
			log.Error().Msgf("panic received, attempting gracefull exit; Error: %v", err)
		}

		Stop()
	}
}

// Stop can be used to manually signal a stop. It awaits all final callbacks to be called before returning.
func Stop() {
	if len(sigList) > 0 {
		sigChan <- os.Interrupt
	}

	StopListener()
}

// StopListener waits for the exit
func StopListener() {
	stopWait.Wait()
}

// Finally adds callbacks which are called on both graceful and forceful exit.
func Finally(cb func()) {
	listenSignal(os.Kill, cb)
	listenSignal(syscall.SIGQUIT, cb)
	listenSignal(os.Interrupt, cb)
}

func FinallyClose[T any](cs chan T) {
	Finally(func() { close(cs) })
}

func EFinally(cbs ...func() error) {
	for _, cb := range cbs {
		Finally(func() {
			if err := cb(); err != nil {
				log.Warn().Err(err).Msg("finaly callback error")
			}
		})
	}
}

// listenSignal Adds callbacks for some arbitrary signal
func listenSignal(s os.Signal, cb func()) {
	if sigList == nil {
		sigList = make(map[os.Signal]*list.List)
		sigChan = make(chan os.Signal, 3)

		go func() {
			defer stopWait.Done()
			var s os.Signal
			var l *list.List
			var ok bool
			for {
				s, ok = <-sigChan
				if !ok {
					return
				}

				// Try and select the list, continue otherwise
				if l, ok = sigList[s]; !ok {
					continue
				}

				lockMut.RLock()
				defer lockMut.RUnlock()
				// Iterate back first, to ensure oldest get executed last
				for e := l.Front(); e != nil; e = e.Prev() {
					e.Value.(func())()
				}
			}
		}()
	}

	lockMut.Lock()
	defer lockMut.Unlock()
	sl, slok := sigList[s]
	if !slok {
		signal.Notify(sigChan, s)
		sl = list.New()
		sigList[s] = sl
	}

	// Append callback
	sl.PushBack(cb)
}
