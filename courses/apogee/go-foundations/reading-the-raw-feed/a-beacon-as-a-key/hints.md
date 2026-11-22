Look up the full array as a map key before recording it. Changing the caller's array later must not change an already stored key.

The slice's upper bound is excluded. Slice the pointed-to array directly to share its storage instead of copying bytes into a new slice.
