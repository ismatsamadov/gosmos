Loop over `files` and look at each `f.Name`; the first name that fails decides the error, and the rest need not be looked at.

Let `filepath.IsLocal` judge each name rather than searching the text yourself: `log/../correction.xml` holds `..` and still stays inside, while `log/../../.profile` does not.

Build the error with `fmt.Errorf` and the `%q` verb, so the name appears in quotes. Return nil once the loop has passed every name.
