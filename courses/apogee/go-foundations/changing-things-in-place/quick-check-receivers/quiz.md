## Which of the two methods needs a pointer receiver?

- [ ] Both: a value receiver can't read the struct's fields
  > A value receiver can read every field, as `Duration` does; what it can't do is change the original.
- [x] Only `Stow`, because it changes the station
  > A method that changes its struct needs the struct's address; with a value receiver, `Stow` would change a copy.
- [ ] Only `Duration`, because it returns a value
  > Returning a result needs no pointer: `Duration` only reads `p`, so a copy serves it well.

## `cwr` is a `Station`, not a `*Station`. What happens at `cwr.Stow("wind 75 km/h")`?

- [ ] It doesn't compile until you write `(&cwr).Stow("wind 75 km/h")`
  > It compiles: when `cwr` is a variable, Go writes the `&` for you.
- [ ] It compiles, but `Stow` changes a copy, so `cwr` stays unstowed
  > No copy is made: the method gets `&cwr`, the address of `cwr` itself.
- [x] Go takes the address of `cwr` for you, so `cwr` is stowed
  > Calling a pointer method on a variable is shorthand for `(&cwr).Stow(...)`, so `cwr.Stowed` becomes true.
