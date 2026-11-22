package main

// offer puts s on the Board's queue q, which has room for at least one, without ever waiting: while q is full it drops the oldest schedule waiting, and it returns how many it dropped.
func offer(q chan *Schedule, s *Schedule) (dropped int) {
	for {
		select {
		case q <- s:
			return dropped
		default:
		}
		select {
		case <-q: // the oldest schedule waiting
			dropped++
		default:
		}
	}
}
