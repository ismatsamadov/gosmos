package main

// Station holds the current dish state.
type Station struct {
	Code       string
	Stowed     bool
	StowReason string
}

// Stow parks this station's dishes and records the reason.
func (s *Station) Stow(reason string) { s.Stowed = true; s.StowReason = reason }

// stowAction returns this station's Stow method for the checklist to call later.
func stowAction(s *Station) func(string) {
	return func(reason string) {}
}
