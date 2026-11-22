## What does `fmt.Println(seen["thistle"])` print?

- [ ] `nil`
  > `nil` is the zero value of pointers, maps and slices, but this map holds `int` values, and an `int` is never `nil`.
- [x] `0`
  > A missing key gives the zero value of the map's value type, and for `int` that is 0.
- [ ] Nothing: the program panics, because `thistle` isn't a key
  > Reading a missing key never panics in Go; it returns the zero value, and the comma-ok form tells you the key was missing.

## In what order does `for id := range seen` visit the keys?

- [ ] In the order the keys were written in the literal
  > A map keeps no insertion order, and Go even varies the order on purpose.
- [ ] In sorted order, as `fmt.Println` shows a map
  > `fmt` sorts the keys when it prints a whole map, but `range` doesn't sort them.
- [x] In no fixed order: it can change from one run to the next
  > Map iteration order is unspecified, so when order matters, sort the keys yourself.

## What happens at `pending["thistle"] = 1`?

- [x] The program panics: `pending` is a nil map, and you can't write to one
  > `var pending map[string]int` declares a nil map: reading from it gives zero values, but writing panics. Give it a map first, such as `map[string]int{}`.
- [ ] Go creates the map and stores 1 under `thistle`
  > Go never creates a map for you on assignment; writing to a nil map panics when the line runs.
- [ ] It doesn't compile, because `pending` was never given a map
  > It compiles: the compiler doesn't track whether a map variable is nil. The mistake shows up only when the line runs.
