Call `decodeEnvelope` first and return its error: a broken envelope has no kind to switch on.

In the switch, give a variable of type `any` a pointer to a new, empty report of the right type. `json.Unmarshal` fills whatever that pointer points to.

Put an unknown kind into the error with `%q`, and return the error from decoding the payload too, such as a chunk count written as text.
