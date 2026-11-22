The loop already stops at a failure, so nothing inside it needs to change. The answer comes after the loop.

Return `sc.Err()` as the second result instead of nil. It is nil after a clean end, so a whole log still comes back without an error, and the lines collected before a failure come back with its reason.
