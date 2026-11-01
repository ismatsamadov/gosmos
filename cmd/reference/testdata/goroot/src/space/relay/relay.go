// Package relay passes signals between stations.
package relay

// A Signal travels between stations.
type Signal interface {
	// Send sends the signal.
	Send() error
}
