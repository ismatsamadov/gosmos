A tag goes after the field's type, on the same line, in backquotes. The key inside the quotes is the export's spelling, letter for letter.

Tag all five fields, not only `SatID`. Decoding forgives case, but encoding writes an untagged field under its Go name, such as `ID` or `AOS`, and the encoded booking must match Yusra's exactly.
