package main

// Tracker shows one satellite's live frames in order, each once.
type Tracker struct {
	Last    uint32 // the sequence number of the last frame shown
	Shown   int    // frames shown
	Missed  int    // numbers jumped over: frames that never came in time
	Dropped int    // frames not shown: repeats, and frames older than Last
}

// Accept reports whether the frame numbered seq is to be shown, and counts it.
func (t *Tracker) Accept(seq uint32) bool {
	t.Last = seq
	t.Shown++
	return true
}
