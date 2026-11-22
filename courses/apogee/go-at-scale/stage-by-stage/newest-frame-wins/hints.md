Before storing a frame, look up what `h` already holds for its satellite with `old, ok := h[f.SatID]`.

Store the frame when the satellite isn't there yet, or when the new frame's `Seq` is higher than `old.Seq`.

If only the frame numbered 0 goes missing, the comparison runs against the map's zero value: check `ok` before comparing.
