---
slug: a-function-behind-main
title: A Function Behind Main
type: code_tests
position: 0
difficulty: 6
topics: programs, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Four of the van's eleven scripts print the station table, each its own way. The team wants one command, gctl, and Nkem asks one question before it grows: "How will a test run it without starting a program?"

Give the program's work to a function that receives the process's arguments and three streams: standard input, standard output for results and standard error for messages. It returns the exit status, the number a finished program hands back, 0 for success, so a test can call it with a slice and buffers.

```go
func main() {
	os.Exit(count(os.Args[1:], os.Stdout))
}

func count(args []string, out io.Writer) int {
	fmt.Fprintln(out, len(args), "arguments")
	return 0
}
```

## Your task

Implement `run`: write each station's code and name to `stdout`, one line per station in the table's order, and return 0:

```text
YST Ystmark
PEL Port Ellery
SQT Salar Quieto
LAN Lantern Atoll
CWR Cape Wren
```

> [!WARNING]
> `fmt.Println` writes to the process's own output, not to the test's buffer. Write to the `stdout` parameter.
