Each `f` the loop gets is a `reflect.StructField`, and `f.Tag` is a `reflect.StructTag` with `Get` and `Lookup` methods.

`Lookup("check")` returns two values: the rules, and whether the field has a `check` key at all.

When the second value is false, the line ends in `no check`; otherwise it ends in the rules, without `check:` and without quotes.
