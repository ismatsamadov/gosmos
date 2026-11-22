## Which `pointDish` follows Go's convention for work that can be canceled?

- [ ] `func pointDish(code string, azDeg float64, ctx context.Context) error`
  > It works, but by convention `ctx` comes first, so every reader finds it in the same place.
- [x] `func pointDish(ctx context.Context, code string, azDeg float64) error`
  > The context is the first parameter, named `ctx`, ahead of everything else the function needs.
- [ ] `func pointDish(code string, azDeg float64) error`, calling `context.Background()` inside
  > A context made inside never hears the caller's deadline; that is how Tuesday's uplink held La Grande for an hour.

## The pass ends at LOS on its own. Do you still call the cancel function `passContext` returned?

- [x] Yes, with `defer cancel()`, to release the context as soon as the work is done
  > Until it ends, the context keeps resources such as its LOS timer; calling `cancel` releases them as soon as the work is done, instead of at LOS.
- [ ] No: the deadline will cancel the context, so `cancel` adds nothing
  > If the work finishes before LOS, the context and its timer stay alive until LOS; `cancel` releases them right away.
- [ ] Only when the work returns an error
  > The context holds the same resources whether the work succeeds or fails, so cancel is called either way.

## Where should the uplink service keep the `ctx` for a pass?

- [ ] In a field of the `Uplink` struct, so every method can reach it
  > A struct outlives the call its context belongs to, so its methods could end up using a context that has ended or belongs to another pass.
- [ ] In a package-level variable that every pass shares
  > Each pass has its own deadline, and one shared variable can hold only one of them.
- [x] Nowhere: each method takes it as its first parameter
  > Passing `ctx` with each call keeps every piece of work tied to the right pass and its deadline.
