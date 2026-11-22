Start from `reflect.TypeOf(v)` and range over its `Fields()`: each value is a `reflect.StructField`.

A field's `Type` prints as `time.Duration`, and its `Type.Kind()` prints as `int64`. `fmt.Sprintf` with three `%s` verbs builds the line.

If a line shows `Duration` without `time.`, you printed the type's `Name()`; print the type itself, whose `String` method includes the package.
