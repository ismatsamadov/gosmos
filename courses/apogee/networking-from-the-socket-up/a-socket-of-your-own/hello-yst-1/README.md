---
slug: hello-yst-1
title: Hello, YST-1
type: code_tests
position: 7
difficulty: 7
topics: defer, networking, time
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T07:44:27Z
---
Friday 10 March. Before Dish 1's old controller comes out on the 20th, Sunniva wants dishd's first real conversation with the simulator.

The given `handshake` reads the banner, sends HELLO, the first message a controller expects, and reads the answer, WELCOME: the dish, its firmware (the program built into it) and `LastSeq`, the number of the last command it ran. `net.DialTimeout` gives up after the time you give it; `Dial` waits as long as the system does:

```go
conn, err := net.DialTimeout("tcp",
	"127.0.0.1:7300", 2*time.Second)
if err != nil {
	return err // refused, or no answer in 2s
}
```

## Your task

Make `firstContact` dial `addr` with a one-second timeout, close the connection when it returns and call `handshake`. Return the banner's name, then the firmware and last command:

```text
YST-1 firmware 2.4, last command 0
```

> [!NOTE]
> Across the island's line, dishd will give a dial 5 seconds. On the bench, one is plenty.
