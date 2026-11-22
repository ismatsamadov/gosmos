The message is already right. What does `%v` keep of `err`, and what does it lose?

`fmt.Errorf` has one verb that wraps an error instead of copying its text; the readme's example uses it.
