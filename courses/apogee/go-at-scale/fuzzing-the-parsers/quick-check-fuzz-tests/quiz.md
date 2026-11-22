## What does plain `go test` run `FuzzParseLine` on?

- [ ] Invented lines for a few seconds, then the seeds
  > Invented lines come only with `-fuzz`: plain `go test` never makes one up.
- [x] Every seed in `lineSeeds` and every file in `testdata/fuzz/FuzzParseLine/`, nothing else
  > Together they are the seed corpus, and plain `go test` runs each one as a subtest.
- [ ] Only the seeds that `f.Add` adds
  > The corpus files run too: that is how Saturday's line runs on every `go test`.

## What does `go test -fuzz=FuzzParseLine -fuzztime=30s` add?

- [x] New lines made from the seeds, until one fails or 30 seconds pass
  > The fuzzer changes the inputs it has, keeps those that reach new code, and stops at the first failure.
- [ ] Every seed, run again 30 times
  > Running the same seeds again finds nothing new; the fuzzer invents inputs instead.
- [ ] Every possible line of up to 30 characters
  > There are far too many: the fuzzer changes the inputs it has, led by the code each one reaches.

## The fuzzer finds a line that makes `FuzzParseLine` fail. What happens to it?

- [ ] It is printed once, then lost when the run ends
  > It is written to disk, so the failure can be run again.
- [ ] It is added to `lineSeeds` in `seeds.go`
  > The fuzzer never edits your code: it writes a file of its own.
- [x] It is saved in `testdata/fuzz/FuzzParseLine/`, and plain `go test` runs it from then on
  > So the failure stays in the seed corpus: it fails every run until the code is fixed, and guards the fix after.
