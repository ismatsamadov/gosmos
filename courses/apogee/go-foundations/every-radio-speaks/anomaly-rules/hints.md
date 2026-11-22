Each `Check` compares one field of `f` with its rule's limit. What should it return when the reading is fine, and what when it isn't?

Build each message with `fmt.Sprintf`, using `%.2f` for volts and `%.1f` for degrees, and return it together with `true`.

Compare with `<` for the battery and `>` for the temperature, so a reading exactly at a limit still returns `"", false`.
