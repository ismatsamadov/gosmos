`json.Unmarshal` needs somewhere to write: declare a `Booking` variable and pass its address, with `&`, not the value itself.

When `Unmarshal` returns an error, some fields may already be filled: a wrong type for the id still lets the station through. Return `Booking{}` with the error rather than the variable.
