Array equality checks every byte. Passing a `[12]byte` copies its contents, so changing the satellite number in the local copy leaves the caller's original intact.
