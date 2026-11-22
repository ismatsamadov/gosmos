Each step is a call and an `if` that returns `nil` and the error as soon as it fails: `checkRequest`, then `readRSA`, then `keyFingerprint`, then the comparison, then `encryptTo` for each secret. Build the lines only at the end.

Compare `string(r.Text)` with `"brinewatch handover " + fp` as a whole: a line naming the key by a few digits, or naming another key, must not match. `hex.EncodeToString` writes each locked secret for its line.

If the impostor's key swapped in beside the request gets a secret, the comparison isn't using the fingerprint of the key read from `r.Key`. If Brinewatch can't open a line, check that each secret is locked under its own label: `apiKeyLabel` for the API key, `webhookLabel` for the webhook secret.
