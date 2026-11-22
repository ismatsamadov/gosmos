Lock `l.mu` first and `defer` the unlock, so the map and the bucket stay yours until `Allow` returns. A customer with no entry in `l.buckets` gets `&Bucket{}`, stored in the map before you use it; a zero bucket fills itself on its first `Take`.

After `Take`, the whole tokens left are `b.Tokens` turned into an `int`, which drops the fraction. When `b.Tokens` is under 1, the missing part of a token, `1 - b.Tokens`, divided by `refillPerSec`, is how many seconds until one is back.

Round those seconds up with `math.Ceil` before turning them into a `time.Duration` of whole seconds, so 0.4 s becomes 1s. If the 20 goroutines let more than 10 through, the lock doesn't cover both finding the bucket and taking from it.
