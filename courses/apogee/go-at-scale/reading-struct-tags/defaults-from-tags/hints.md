`reflect.Indirect` follows a pointer to what it points to, and only a struct reached that way can be set: check the result's `Kind()` and `CanSet()` before anything else.

Range over the struct's `Fields()`; read each tag with `Lookup("default")`, and skip a field when the key is missing, when `CanSet()` is false or when `IsZero()` is false.

If the test reports a panic from `SetString`, the loop reached the loader's unexported password: skip any field whose `Value` can't be set.
