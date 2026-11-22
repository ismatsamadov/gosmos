## What is the type of `saved`?

- [ ] `bool`
  > That is the type returned by a call, not the type of the function itself.
- [x] `func(float64) bool`
  > The function accepts one `float64` and returns one `bool`; assigning it preserves that type.
- [ ] `float64`
  > That is its parameter type, not the type of the function value.

## Which line calls the rule?

- [x] `result := saved(25)`
  > Parentheses and an argument call the function, producing `true` here.
- [ ] `saved := ready`
  > This keeps the function value without running it.
- [ ] The declaration of the function literal
  > Writing a function literal creates a function value; its body runs when it is called.

## What must happen before calling a variable declared `var rule func(float64) bool`?

- [ ] Nothing; its default rule returns `false`
  > The zero value is `nil`, not a function that returns a default result.
- [ ] Convert it to `bool`
  > A function cannot be converted to its result type.
- [x] Assign a function to it
  > Calling a `nil` function panics; assigning a function gives the variable something to call.
