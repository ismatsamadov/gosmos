package main

import (
	"fmt"
	"net"
	"time"
)

// frame is Pipit-2's telemetry at 06:37:15, as the recorder sends it.
const frame = "pipit-2 30412 2028-06-12T06:37:15Z batt_v=8.04 temp_c=9.5 mode=NOMINAL"

func main() {
	view, err := net.ListenPacket("udp", "127.0.0.1:0") // the live view, at Port Ellery
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	defer view.Close()
	rec, err := net.ListenPacket("udp", "127.0.0.1:0") // the recorder, in the hut
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	defer rec.Close()
	view.SetReadDeadline(time.Now().Add(time.Second)) // no read waits longer than a second

	buf := make([]byte, 512)
	n, fromRecorder := 0, false
	// TODO: send frame from rec to view, then read it on view into buf
	fmt.Printf("%d bytes, from the recorder: %t\n", n, fromRecorder)
	fmt.Printf("%q\n", buf[:n])
}
