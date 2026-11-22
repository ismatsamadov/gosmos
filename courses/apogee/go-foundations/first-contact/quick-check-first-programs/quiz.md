## Where does this program start running?

- [ ] At the first line of the file
  > The first line, `package main`, says which package the file belongs to; no statement runs there.
- [ ] At `import "fmt"`
  > An import lets the file use the `fmt` package; it isn't a statement, so the program doesn't start there.
- [x] At the start of `func main()`
  > This program's code starts in `main`: Go calls it and runs its statements from the first one down.

## What does this program print?

- [ ] `DEBUG: radio on`, then `AOS in 2 minutes`
  > The first call sits behind `//`, so Go treats the whole line as a comment and never runs it.
- [x] `AOS in 2 minutes`
  > The commented call never runs, so only the second `fmt.Println` prints.
- [ ] Nothing, because the comment stops the program
  > This comment has no effect: the line after it still runs.

## You can't edit `main_test.go`. What do you change to make its test pass?

- [ ] The `want` value inside the test
  > `want` belongs to the read-only test: it states what the function must return.
- [x] The function in `main.go` that the test calls
  > The test checks what that function returns, so you change the function until it returns `want`.
- [ ] The name of the test
  > The test file is read-only, and a new name wouldn't change what the test checks anyway.
