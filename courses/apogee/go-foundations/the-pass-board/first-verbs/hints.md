Which parts of `Next passes for Pipit-1: 4` never change, and which two come from `name` and `passes`?

The fixed text becomes the format string, with the verb for a string where the name goes and the verb for a whole number where the count goes. `Printf` doesn't end the line, so without `\n` the next line runs on after the 4.

The shape is `fmt.Printf("Next passes for __: __\n", ____, ____)`, with a verb in each short blank and a variable in each long one.
