package main

// Run is one planning run as the rest of passplan 2 sees it: Stop ends it, and Done is closed once it has ended.
type Run struct {
	done chan struct{}
	stop func()
}

// newRun returns a run that is still going.
func newRun() *Run {
	r := &Run{done: make(chan struct{})}
	r.stop = func() { close(r.done) } // TODO: make stop safe to call more than once
	return r
}

// Stop ends the run; the Board's shutdown and a newer plan may both call it.
func (r *Run) Stop() { r.stop() }

// Done returns a channel that is closed once the run has been stopped.
func (r *Run) Done() <-chan struct{} { return r.done }
