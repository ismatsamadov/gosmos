# golang/website's files for the Go reference

`make reference` builds `reference/go.sqlite` from the Go that `go env GOROOT` names and from
these files of [golang/website](https://github.com/golang/website), kept here so the reference
builds with nothing fetched. They are byte for byte as they are at commit
`b8d34bbe2d9e230cc70ce3e4272dc651f494a188`, the Makefile's `WEBSITE_COMMIT`, under the same
paths. They are The Go Authors' work: the articles are published on go.dev under the Creative
Commons Attribution 4.0 International License, and the code under the BSD license in
golang/website's `LICENSE`, kept beside them (`docs/reference-format.md`, "Credit").

- `_content/doc/effective_go.html`: Effective Go, with the programs it includes from
  `_content/doc/progs/`: `eff_bytesize.go`, `eff_qr.go`, `eff_sequence.go`, `eff_unused1.go`
  and `eff_unused2.go`.
- `_content/doc/faq.md`: the Go FAQ.
- `_content/doc/code.html`: How to Write Go Code.

The go command leaves out folders whose names start with `_`, so the programs are not part of
the repository's module. To move to a newer commit, fetch the same paths at it, set
`WEBSITE_COMMIT`, run `make reference`, and read the articles in the app.
