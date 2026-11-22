Let a `csv.Reader` read the whole sheet: a quoted note with a line break is still one record with six cells, however many lines it covers.

Skip the first record, the header, by its index rather than by slicing: an empty sheet has no first record to slice past.

Take index 5 of every other record. The reader has already removed the surrounding quotes and turned each doubled quote back into one.
