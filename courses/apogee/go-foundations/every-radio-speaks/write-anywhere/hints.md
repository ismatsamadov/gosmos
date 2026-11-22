Which lines in `writeBoard` print without ever mentioning `w`?

`fmt.Fprintln` takes the writer first, then the values, and ends the line the way `Println` does.
