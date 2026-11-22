---
slug: a-property-not-an-answer
title: A Property, Not an Answer
type: code_tests
position: 3
difficulty: 9
topics: functions, testing
created: 2026-11-11T12:51:43Z
modified: 2026-11-11T12:51:43Z
---
Mirela, on what a check must do once the fuzzer runs all night: "It invents lines nobody wrote an answer for. Tell me what's true of every answer."

A fuzz test can't compare a frame with the one it expects, since nobody has seen the line before. It checks a property instead: something true of every frame the parser accepts, such as a battery from 0 to 20 volts. NaN, the float that is not a number, is never below, above or equal to any number, so a check tests for it with `math.IsNaN`:

```go
kmh := math.NaN()
fmt.Println(kmh > 72, kmh <= 72) // false false
fmt.Println(math.IsNaN(kmh))     // true
```

## Your task

`checkLine`, in `checks.go`, only calls `parse`. When `parse` accepts `line`, fail `t` unless the frame has a satellite ID, a battery that is a number from 0 to 20, a temperature that is a number from −100 to 150, and one of the four modes: `NOMINAL`, `SAFE`, `DETUMBLE` or `LEOP`.

> [!NOTE]
> A refused line passes, whatever its error says: refusing is the parser's job.
