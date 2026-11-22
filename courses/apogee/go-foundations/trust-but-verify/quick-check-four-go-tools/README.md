---
slug: quick-check-four-go-tools
title: Quick Check: Four Go Tools
type: quiz
position: 7
difficulty: 5
topics: modules, testing
created: 2026-11-13T11:32:35Z
modified: 2026-11-13T16:25:47Z
---
Before her sign-off, Mirela runs four commands over the uplink code in a terminal. In the app, Format Code runs `gofmt`, and every test run starts with part of `go vet`; the rest needs a terminal.

```text
gofmt -l .         layout gofmt would change
go vet ./...       suspicious code, wrong verbs
go test -cover     how much code the tests ran
govulncheck ./...  known vulnerabilities
```

`govulncheck` comes from the Go team but is installed on its own. It checks a program's modules, the code it uses from other projects, and the Go release that builds it, against a list of reported vulnerabilities.

## Your task

Answer Mirela's three questions.
