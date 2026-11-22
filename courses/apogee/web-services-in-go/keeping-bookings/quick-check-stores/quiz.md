## `MemStore`'s methods all have the receiver `(m *MemStore)`. Which values can be used as a `Store`?

- [x] A `*MemStore`, such as `NewMemStore()` returns
  > Methods with a pointer receiver belong to the pointer type, so `*MemStore` has all four and satisfies `Store`.
- [ ] A `MemStore` value, and a `*MemStore` too
  > A plain `MemStore` value doesn't have the pointer-receiver methods, so it isn't a `Store`.
- [ ] Neither, until `MemStore` declares that it implements `Store`
  > Go has no such declaration: having the methods is all a type needs to satisfy an interface.

## Why do `Create`'s conflict check and its insert happen under one lock?

- [ ] A mutex can't be locked twice in one method
  > It can be locked, unlocked and locked again; the trouble is what happens in between.
- [x] Otherwise two customers booking one pass at once could both pass the check, and both get the pass
  > Between an unlock after the check and a lock for the insert, another `Create` can check and insert too.
- [ ] Holding the lock longer makes `Create` faster
  > A lock held longer never makes anything faster; it makes the check and the insert one step.

## Why does `List` copy the bookings into a new slice instead of returning one the store keeps?

- [ ] Go can't return a slice that a struct's field holds
  > It can; the question is what a caller could then do to the store through it.
- [ ] So that copying the bookings needs no lock
  > Copying reads the map, so it needs the lock like everything else.
- [x] So that no caller can change the store's bookings through what it got back
  > A caller may sort, change or keep its copy, and the store's own bookings stay as they were.
