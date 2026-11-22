## What makes `XBandRadio` satisfy `Radio`?

- [ ] A line declaring that `XBandRadio` implements `Radio`
  > Go has no such declaration and no `implements` keyword.
- [x] Its method `RateMbps() float64`
  > Having every method the interface lists is all it takes.
- [ ] A field `RateMbps float64`
  > An interface lists methods, and a field never satisfies one.

## What does `fmt.Printf("%T", r)` print after these lines?

- [ ] `main.Radio`
  > `%T` prints the type of the value inside the interface, never the interface type itself.
- [ ] `main.XBandRadio`
  > `r = UHFRadio{}` replaced the `XBandRadio`: an interface variable isn't tied to the first type it held.
- [x] `main.UHFRadio`
  > `r` now holds a `UHFRadio`, and `%T` prints the type of the value inside.

## Which values can be stored in `v`?

- [x] Values of every type
  > `any` lists no methods, so every type satisfies it.
- [ ] Only values of built-in types, such as `437.25`
  > Your own types satisfy `any` too, because it asks for no methods.
- [ ] Only values whose type has at least one method
  > `any` asks for no methods, so a type without any satisfies it as well.
