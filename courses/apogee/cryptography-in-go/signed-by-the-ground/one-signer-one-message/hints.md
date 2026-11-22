In `Sign`, make the HMAC with `hmac.New(sha256.New, s.key)`, write the command into it and sum it, all in that one call: a fresh hash for every command.

Delete the `mac` field from `Signer` and the `mac:` part of the composite literal in `NewSigner`. The `hash` import then has nothing left to do, so remove it.

If the second signature is right but some of the 400 are wrong, something is still shared between calls: every value `Sign` writes into must be made inside `Sign`.
