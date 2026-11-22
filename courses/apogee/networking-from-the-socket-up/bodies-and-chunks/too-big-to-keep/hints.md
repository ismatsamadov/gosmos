The check belongs between `chunkSize` and the `make` that holds the chunk: by then you know its size, and nothing of it has been read.

The room left is `max` minus the bytes the body holds so far, `int64(len(body))`. A size larger than that room takes the body past the limit; a size equal to it fills the body exactly and is fine.

If the test reports `recording cut off` where it wants 413, the chunk is read before the check: its bytes never come, so the read meets the end of the stream. Skip the check when `max` is 0.
