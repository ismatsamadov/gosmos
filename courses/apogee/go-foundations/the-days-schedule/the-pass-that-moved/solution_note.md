`live[:2]` shared `live`'s array, so writing `draft[1]` changed the live Board too. `slices.Clone` copies the two times into a new array, so the draft can change without touching `live`.
