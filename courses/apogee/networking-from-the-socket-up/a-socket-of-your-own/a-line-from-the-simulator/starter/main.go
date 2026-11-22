package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	defer ln.Close()
	go NewSim("YST-1", 1).Serve(ln) // the maker's simulator, on this Mac only

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(time.Second)) // no read waits longer than a second

	buf := make([]byte, 64)
	var n int
	// TODO: read what the controller sends first into buf, setting n
	fmt.Printf("%d bytes: %q\n", n, buf[:n])
}
