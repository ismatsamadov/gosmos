Decode into an `Envelope` with `decodeOne` first. Its errors, an unknown field or a second report, go back to the caller as they are.

Switch on the envelope's version: version 1 sets the kind to `image`, version 2 needs nothing, and every other version returns the unsupported-version error, with `%d` for both numbers.

Pass the kind and the raw payload to `decodePayload`, and wrap its error with `%w` after `report <id>: ` so the decoding error stays inside.
