Look the user up as the starter does, but don't return when the name is missing: use `dummyRecord` as the stored record instead, and remember that the user was unknown.

Call `c.matches` once, whoever is logging in. Then return `errInvalidLogin` if it gave an error, if the password didn't match, or if the user was unknown, so the dummy record can never let anyone in.

If the test counts no `kdf` call for the unknown user, something still returns before `c.matches`. If the two errors differ, an `errors.New` is left over: remove it, and the `errors` import with it.
