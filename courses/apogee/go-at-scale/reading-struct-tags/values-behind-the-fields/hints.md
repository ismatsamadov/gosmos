Take `reflect.ValueOf(v)`, and when its `Kind()` is `reflect.Pointer`, use its `Elem()` instead; then return `nil` unless the kind is `reflect.Struct`.

Range over the struct value's `Fields()` with two loop variables: the `StructField` gives the name and `IsExported`, and the field's `Value` goes to `fmt.Sprintf` as it is.

If the test reports a panic from `Interface`, or a `password=` line, the loop reaches the unexported field: skip each field whose `StructField` isn't exported.
