package main

import "fmt"

// Session holds a station's dish for one pass.
type Session struct {
	Station string
}

// openSession books station's dish for a pass.
func openSession(station string) *Session {
	fmt.Println("open", station)
	return &Session{Station: station}
}

// Close frees the dish for the next pass.
func (s *Session) Close() {
	fmt.Println("close", s.Station)
}

// track follows the satellite through the pass, and fails in wind of 72 km/h or more.
func track(windKmh float64) error {
	if windKmh >= 72 {
		return fmt.Errorf("wind %.0f km/h", windKmh)
	}
	return nil
}

// runPass opens a session on station, tracks the pass and closes the session.
func runPass(station string, windKmh float64) {
	s := openSession(station)
	if err := track(windKmh); err != nil {
		fmt.Println("pass failed:", err)
		return
	}
	fmt.Println("pass ok")
	s.Close()
}

func main() {
	runPass("YST", 40)
	runPass("CWR", 81)
}
