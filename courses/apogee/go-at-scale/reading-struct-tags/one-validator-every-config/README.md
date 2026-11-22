---
slug: one-validator-every-config
title: One Validator for Every Config
type: code_tests
position: 9
difficulty: 9
topics: errors, reflection, structs
created: 2026-11-12T03:14:42Z
modified: 2026-11-12T03:14:42Z
---
Friday 8 June. The third wave has flown since Monday on forms checked one rule at a time; from today one call checks a whole form, and passplan 2's own settings each time it starts.

`Validate` walks a struct's exported fields in order, splits each `check` tag at its commas and hands each rule to `checkRule`, in `rules.go`; a missing or empty tag has no rules. A struct field, embedded or not, is walked the same way, and its problems name its own type: `HostConfig.Dishes`. `errors.Join` makes one error of them all, a problem per line, or `nil` when all are `nil`.

```go
var errs []error
for _, kmh := range []int{64, 81, 74} {
	if kmh > 72 {
		err := fmt.Errorf("%d km/h", kmh)
		errs = append(errs, err)
	}
}
fmt.Println(errors.Join(errs...)) // two lines
```

## Your task

`Validate` returns `nil`. Check `v`, a struct or a pointer to one: every exported field's rules, in order, named `<Type>.<Field>`, and the fields of any struct field; join every problem. Anything else gives `validate: want a struct, got %T`.

> [!NOTE]
> The tests check the fifteen hosts' forms in `forms.go`, as first sent.
