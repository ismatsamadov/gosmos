## What does sc.Scan() return once the input has ended?

- [ ] true, with an empty string in sc.Text()
  > An empty line is still a line; at the end there is no line left, and Scan says so.
- [x] false, so the for loop ends
  > Scan returns false when no line is left, or when reading fails, and the loop stops there.
- [ ] It returns io.EOF as an error
  > Scan returns only a bool; reaching the end is not an error, so sc.Err() stays nil.

## What does strings.Fields return for a line of only spaces and tabs?

- [x] An empty slice, so the loop skips the line
  > Fields drops every space and tab, nothing is left, and len(words) is 0.
- [ ] A slice holding one empty string
  > Fields never returns an empty word; strings.Split can, but Fields drops every space and tab.
- [ ] An error saying the line is blank
  > Fields returns only a slice: a blank line is ordinary input for it.

## In c, ok := s.commands[words[0]], what do c and ok hold for a name not in the table?

- [ ] A nil command and an error
  > A map lookup never returns an error; ok alone says whether the name was there.
- [ ] Nothing: the lookup panics
  > Reading a missing key never panics; only writing to a nil map does.
- [x] The zero command, and ok is false
  > A missing key gives the zero value of command, and ok is false, so dispatch can answer it as unknown.
