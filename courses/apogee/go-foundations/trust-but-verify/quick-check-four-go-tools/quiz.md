## A check calls `t.Errorf("seq %d", c.Op)`, and `c.Op` is a string such as `PING`. What happens when you run the tests?

- [ ] The tests run, and the message prints `%!d(string=PING)`
  > That is what `fmt` would print, but `go test` runs vet's check of format strings first, and the tests never start.
- [x] `go test` reports the wrong verb and runs no test
  > Part of `go vet` runs before every `go test`, and a verb that doesn't fit its value stops the run with `[build failed]`.
- [ ] Nothing, until someone runs `gofmt`
  > `gofmt` changes only how code is laid out, its spaces, tabs and line breaks, never what it means.

## `go test -cover` prints `coverage: 80.0% of statements`. What does that tell Mirela?

- [x] One statement in five never ran during the tests
  > Those statements could be wrong and every test would still pass; `go test -coverprofile=c.out` and `go tool cover -html=c.out` show which ones.
- [ ] One test in five failed
  > A failing test makes `go test` print FAIL; coverage counts the code the tests ran, not how they ended.
- [ ] The code is right four times in five
  > Coverage says a statement ran, not that anything checked it: a test with no checks at all can cover every line.

## The uplink code uses only the standard library. What can `govulncheck` still report?

- [ ] Nothing: with no modules there is nothing to check
  > It also checks the Go release that builds the program, and the standard library can have reported vulnerabilities too.
- [ ] Mistakes in the uplink code itself, as `go vet` does
  > `govulncheck` looks only for vulnerabilities already reported in what a program uses; new mistakes are for vet and your tests to find.
- [x] A reported vulnerability in the Go release that builds it
  > Its list covers the standard library too, so a program built with an old Go release can be reported, and the fix is a newer release.
