package main

import (
	"context"
	"fmt"
	"time"
)

// Command is one telecommand on its way up.
type Command struct {
	SatID string
	Seq   uint64
	Op    string
}

// transmit sends cmd over the radio, which takes a second, unless ctx ends first.
func transmit(ctx context.Context, cmd Command) error {
	select {
	case <-time.After(time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// sendCommand transmits cmd, naming the satellite and sequence number in any error.
func sendCommand(ctx context.Context, cmd Command) error {
	if err := transmit(ctx, cmd); err != nil {
		return fmt.Errorf("uplink %s %d: %w", cmd.SatID, cmd.Seq, err)
	}
	return nil
}

// runPass sends the pass's command and logs the outcome.
func runPass(ctx context.Context, cmd Command) {
	if err := sendCommand(ctx, cmd); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("uplink %s %d: sent\n", cmd.SatID, cmd.Seq)
}

func main() {
	// The pass has 30 ms left before LOS.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	runPass(ctx, Command{SatID: "pipit-2", Seq: 118, Op: "SET_MODE"})
}
