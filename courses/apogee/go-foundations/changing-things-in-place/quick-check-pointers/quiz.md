## What does `p` hold after `p := &x`?

- [ ] A copy of `x`, holding 3
  > A copy would come from `p := x`; the `&` takes the address of `x` instead.
- [x] The address of `x`
  > `&x` is a pointer to `x`, so `p` leads to the variable `x` itself.
- [ ] The value that `x` points to
  > `x` is an `int`, not a pointer, so it points to nothing; `&` gives its address.

## After `*p = 5`, what does `x` hold?

- [x] 5
  > `*p` is the variable `p` points to, which is `x`, so storing 5 there changes `x`.
- [ ] 3, because `p` has a copy of its own
  > `p` holds an address, not a copy, so writing through it reaches `x`.
- [ ] Nothing: it doesn't compile, because `p` holds an address, not an `int`
  > `*p` is the `int` that `p` points to, so assigning 5 to it compiles.

## What does `q` hold after `var q *int`?

- [ ] 0
  > 0 is the zero value of an `int`; `q` is a pointer, and a pointer's zero value is different.
- [ ] The address of a new `int` holding 0
  > `var` makes no `int` for `q` to point to; it gives `q` the zero value of its type.
- [x] `nil`
  > A pointer's zero value is `nil`: it points to nothing until it's given an address.
