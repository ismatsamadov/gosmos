Check `len(body)` before anything else and return the error with the length you got. `binary.Decode` alone would accept 17 bytes and ignore the last one.

Declare a `rawStatus`, pass its address to `binary.Decode` with `binary.BigEndian`, and return its error if there is one.

Then build the `Status` from the raw fields: hundredths of a degree divided by 100, tenths of a km/h and of a °C divided by 10, each converted to `float64` first, and the state as a `DishState`.
