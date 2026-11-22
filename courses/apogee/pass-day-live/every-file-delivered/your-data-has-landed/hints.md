`customers.json` holds one JSON array of customers, each with four strings: `id`, `name`, `webhook_url` and `webhook_secret`, whose text, as bytes, is the signing key. A satellite's `Operator` names its customer, and the booking that matched the contact gives the event its booking ID, satellite and station.

Make the event as the file lands, numbering it under a lock, and send it from a goroutine of its own, so the pass keeps reading its call; wait for those goroutines before `week done`. `http.NewRequestWithContext` with the rehearsal's context builds the request, `w.HTTP.Do` sends it, and the answer's body must be closed.

A receiver refusing a `bad signature` checked other bytes than the ones signed: sign exactly what `json.Marshal` gave you and send those bytes. `hmac.New(sha256.New, secret)` signs the timestamp, a dot and the body, in lowercase hex. A refused timestamp came in milliseconds or from another instant than `w.Now().Unix()` at the attempt.
