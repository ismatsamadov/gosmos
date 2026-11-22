Check the limit and number of targets before wrapping or consuming the source. An empty target list would otherwise discard the recording.

Combine the target slice with `MultiWriter(targets...)`, and pass a limited source to `Copy`. Return both of Copy's results.
