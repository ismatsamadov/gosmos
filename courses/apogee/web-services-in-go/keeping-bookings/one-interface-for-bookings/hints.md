An interface's methods are written one per line, each as a name and a signature with no body, like the `Get` already in `Store`. Add the three from the readme inside its braces.

A field's type comes after its name, so the field called `Store` can have the type `Store`. Put `Store` where `*MemStore` is, and the test's failing store fits.

If the build now says a type doesn't implement `Store`, compare that method in `Store` with the readme's list: an interface accepts only methods with exactly the same parameter and result types.
