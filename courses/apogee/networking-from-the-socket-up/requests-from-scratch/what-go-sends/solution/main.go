package main

import (
	"bufio"
	"fmt"
	"io"
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
	done := askForList(ln.Addr().String()) // Go's own client asks for the list

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("accept:", err)
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second)) // no read waits longer than a second
	br := bufio.NewReader(conn)

	for {
		line, err := br.ReadString('\n')
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		fmt.Printf("%q\n", line)
		if line == "\r\n" {
			break
		}
	}

	io.WriteString(conn, "HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n")
	<-done
}
