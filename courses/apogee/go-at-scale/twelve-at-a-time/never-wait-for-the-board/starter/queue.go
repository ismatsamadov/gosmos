package main

// offer puts s on the Board's queue q, which has room for at least one, without ever waiting: while q is full it drops the oldest schedule waiting, and it returns how many it dropped.
func offer(q chan *Schedule, s *Schedule) (dropped int) {
	q <- s // TODO: never wait: while q is full, drop its oldest schedule
	return 0
}
