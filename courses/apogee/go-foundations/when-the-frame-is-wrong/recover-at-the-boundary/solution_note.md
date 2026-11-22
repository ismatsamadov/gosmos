The directly deferred recovery catches the decoder panic after cleanup runs. It sets the named error result, while a normal decoder return preserves its reading and nil error.
