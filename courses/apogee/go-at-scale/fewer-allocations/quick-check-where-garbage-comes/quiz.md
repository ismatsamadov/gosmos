## Which of the four lines allocates nothing?

- [ ] Line 1, since `p` only points at a copy `note` already has
  > `p` is returned, so `c` must outlive `note`: the compiler moves `c` to the heap, and `-gcflags=-m` says `moved to heap: c`.
- [x] Line 2
  > An `int` copied into a local stays in `note`'s stack frame, which is gone when `note` returns, with nothing for the collector to clear.
- [ ] Line 3, since it only prints a number
  > `Println` takes `any`, so the elevation is copied into an interface value, and that copy goes on the heap.
- [ ] Line 4, since it only sets a field
  > `string(b)` copies `b`'s bytes into a new string, and since `p` outlives `note`, that string goes on the heap.

## Why does line 3 allocate, though `c` is already in memory?

- [x] `Println` takes `any`, and the compiler can't prove the number's copy inside it stays in the frame
  > So that copy goes on the heap. Only a few values have ready-made copies, such as a zero or an integer from 0 to 255, so a value such as 74.9 is copied on every call.
- [ ] `Println` allocates a new buffer for every line it prints
  > `fmt` reuses its buffers from one call to the next; what this call allocates is the copy of the number inside the `any`.
- [ ] A float64 always lives on the heap
  > A float64 lives wherever its variable does, often on the stack; what line 3 allocates is a new copy of it, inside the `any`.

## After line 4, `b` is reused for the next frame. What does `p.ID` hold then?

- [ ] The next frame's bytes, since `p.ID` shares `b`'s memory
  > A string can't change while a slice can, so the conversion can't share the slice's memory: it copies the bytes.
- [ ] Nothing, since reusing `b` frees the string
  > The string has its own copy of the bytes, which stays as long as something refers to it.
- [x] The bytes `b` held at line 4
  > `string(b)` copied them into the new string, so whatever `b` holds later can't change `p.ID`.
