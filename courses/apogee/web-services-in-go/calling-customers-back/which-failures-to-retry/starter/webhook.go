package main

// AnswerError is a receiver's answer outside 200 to 299: the event did not arrive.
type AnswerError struct {
	StatusCode int    // 503
	Status     string // "503 Service Unavailable"
}

// Error says what the receiver answered.
func (e *AnswerError) Error() string { return "receiver answered " + e.Status }
