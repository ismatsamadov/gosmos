Loop over `privateRoutes` and `call` each route with `""` for the header; every status must be 401. Use `t.Errorf`, so that one failed check doesn't hide the others.

`token` trades a key for a token, which goes in the header as `"Bearer "` and the token. Every booking that `bookingsOf` returns for Furrowlight's token must be of `furrow-a` or `furrow-b`. For the next check, put the API key itself after `Bearer `.

If a planted service survives, find it in `plant` in `main_test.go` and ask which of your checks it would pass. Get Furrowlight's token before calling `advance`, use it again afterwards, and check the message `call` returns, not only the status.
