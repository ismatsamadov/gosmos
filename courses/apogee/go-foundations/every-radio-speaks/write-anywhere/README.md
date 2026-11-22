---
slug: write-anywhere
title: Write Anywhere
type: code_tests
position: 4
difficulty: 5
topics: interfaces, printing
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
The Board draws its list of next passes on the Fishbowl's big screen. Nkem wants a test that reads what it drew, but `fmt.Println` only reaches the screen.

`io.Writer` is an interface with one method, `Write`, which takes bytes. The terminal (`os.Stdout`), files, network connections and a `*strings.Builder` all satisfy it. `fmt.Fprintln` and `fmt.Fprintf` work like `Println` and `Printf` but write to the writer you pass first:

```go
func writeBanner(w io.Writer, sat string) {
	fmt.Fprintf(w, "Tracking %s\n", sat)
}

writeBanner(os.Stdout, "Pipit-2") // to the screen
var b strings.Builder
writeBanner(&b, "Pipit-2") // into b.String()
```

## Your task

Make `writeBoard` write its header and every pass ID to `w` instead of the terminal. `main` still passes `os.Stdout`; the test passes `&b`, a `*strings.Builder`.

> [!NOTE]
> Why `&b`? A builder's `Write` changes the builder, so it has a pointer receiver, and only the pointer satisfies `io.Writer`.
