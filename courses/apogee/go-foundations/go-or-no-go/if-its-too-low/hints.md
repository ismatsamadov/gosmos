Compare your output with the expected one: every `Pass` line still prints, but only one warning is left. Which pass keeps it, and why that one?

An `if` runs the lines between its braces only when its condition is true. Each warning's `fmt.Println` goes inside one, and the `Printf` above it stays outside, so every pass still prints.
