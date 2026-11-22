Follow a pointer with `reflect.Indirect`, refuse anything that isn't a struct, then range over the struct's `Fields()`, skipping a field that isn't exported.

A field's name for `checkRule` is the struct type's `Name()`, a dot and the field's name. When a field's `Kind()` is `reflect.Struct`, call `Validate` on its value, from `Interface()`, and keep what it returns with the other problems.

If a problem reads `unknown rule ""`, a field with no `check` tag gave `strings.Split` an empty string, and it returned one empty rule: skip empty rules.
