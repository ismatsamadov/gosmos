A pointer field stays nil when the key is missing or `null`, and points at the value when it is present, even 0. Checking for nil keeps the old mask for every update that never mentioned it.
