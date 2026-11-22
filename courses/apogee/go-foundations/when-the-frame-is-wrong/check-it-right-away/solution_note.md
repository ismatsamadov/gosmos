The loop keeps each error and checks it at once: a rejected reading prints its reason, and `continue` skips both the stored line and the count, so only the two good readings are stored.
