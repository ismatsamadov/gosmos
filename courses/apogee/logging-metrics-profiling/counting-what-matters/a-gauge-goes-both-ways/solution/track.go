package main

// track counts run in g while it runs, and returns run's error.
func track(g *Gauge, run func() error) error {
	g.Add(1)
	defer g.Add(-1)
	return run()
}
