`newClient` keeps its reader in the client's `br` field. Whatever it read past the banner, perhaps the WELCOME and the STATUS, is in that buffer now and nowhere else.

`readFrame` takes a `*bufio.Reader`: give it `c.br` instead of making a new one.

If the second `recv` waits until the test's one-second guard ends it, the STATUS is sitting in the buffer of a reader that was thrown away after the WELCOME.
