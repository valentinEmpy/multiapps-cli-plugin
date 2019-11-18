package commands

import (
	"time"
)

// Worker will do its Action once every interval, making up for lost time that
// happened during the Action by only waiting the time left in the interval.
type Worker struct {
	Stopped         bool          // A flag determining the state of the worker
	ShutdownChannel chan string   // A channel to communicate to the routine
	Interval        time.Duration // The interval with which to run the Action
	period          time.Duration // The actual period of the wait
}

// NewWorker creates a new worker and instantiates all the data structures required.
func NewWorker(interval time.Duration) *Worker {
	return &Worker{
		Stopped:         false,
		ShutdownChannel: make(chan string),
		Interval:        interval,
		period:          interval,
	}
}

// Run starts the worker and listens for a shutdown call.
func (w *Worker) Run(action func()) {

	// Loop that runs forever
	for {
		select {
		case <-w.ShutdownChannel:
			w.ShutdownChannel <- "Down"
			return
		case <-time.After(w.period):
			// This breaks out of the select, not the for loop.
			break
		}

		started := time.Now()
		action()
		finished := time.Now()

		duration := finished.Sub(started)
		w.period = w.Interval - duration

	}

}

// Shutdown is a graceful shutdown mechanism
func (w *Worker) Shutdown() {
	w.Stopped = true

	w.ShutdownChannel <- "Down"
	<-w.ShutdownChannel

	close(w.ShutdownChannel)
}
