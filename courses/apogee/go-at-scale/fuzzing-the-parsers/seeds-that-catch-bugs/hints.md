Each mutant breaks on one kind of line, so each needs its own seed: a line of five fields, a line whose battery reads `batt_v=nan`, and a line whose mode reads `mode=nominal`.

Copy the good line three times. Drop the mode field from the end of one, put `nan` in place of the battery's `8.02` in another, and write the mode in lowercase in the third.

If a mutant still survives, its seed has another fault that the parser refuses first: keep the rest of the line exactly as the good line has it.
