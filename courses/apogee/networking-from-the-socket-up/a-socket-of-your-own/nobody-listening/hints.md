The refused dial's error wraps `syscall.ECONNREFUSED` a few layers down, and `errors.Is` finds it there whatever the layers above it say.

Check for the refusal before the general `err != nil` case, or inside it, and return the name followed by `: not listening`.

An address without a port fails before any dial is tried, with an error that isn't a refusal, so it should still come back as its own text.
