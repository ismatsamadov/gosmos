---
slug: hello-pipit
title: Hello, Pipit
type: code
position: 0
difficulty: 1
topics: printing, programs
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
Welcome to Apogee. Your first program writes a message for Pipit-1, the company's test satellite. The operations team, or ops, will send it from their control room to a ground station, a radio dish that talks to satellites. Pipit-1 will echo it back during its next pass: the few minutes it is in the station's view.

A Go program starts in `func main()`. There, `fmt.Println` prints the text between the quotes and ends the line:

```go
package main

import "fmt"

func main() {
	fmt.Println("Port Ellery standing by")
}
```

Leave `package main` and `import "fmt"` as they are; every program in this chapter starts with them.

## Your task

Change the message so the program prints exactly:

```text
Apogee ground calling Pipit-1
```

Press **Run** to see the output, then **Submit** to check it.
