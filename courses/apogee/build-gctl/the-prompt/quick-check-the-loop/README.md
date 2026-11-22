---
slug: quick-check-the-loop
title: Quick Check: The Loop
type: quiz
position: 4
difficulty: 5
topics: loops, maps, strings
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
Nkem reviews the prompt before the night shift sees it on Friday. Each of her three questions is about one line of this loop:

```go
sc := bufio.NewScanner(in)
for sc.Scan() {
	words := strings.Fields(sc.Text())
	if len(words) == 0 {
		continue
	}
	c, ok := s.commands[words[0]]
}
```

## Your task

Answer three questions: what `Scan` does once the input has ended, what `strings.Fields` gives for a line of spaces and tabs, and what the lookup gives for a name the table doesn't have.
