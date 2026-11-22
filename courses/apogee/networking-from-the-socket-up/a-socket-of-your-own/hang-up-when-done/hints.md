A connection stays open until someone closes it, so `bannerOf` needs `conn.Close()` to run however it returns. A `defer` does exactly that, when `bannerOf` returns.

The `defer` goes after the `if err != nil` block. Above it, a failed dial leaves `conn` nil, and deferring its `Close` panics on the spot.

If the test counts fewer than 3 connections closed by dishd, the `Close` never ran. If it reports a panic, the `defer` comes before the error check.
