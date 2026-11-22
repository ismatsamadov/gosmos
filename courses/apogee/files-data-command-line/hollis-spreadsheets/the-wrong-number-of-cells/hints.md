Set `FieldsPerRecord` on the reader before the first `Read`. The header is checked like every other record, so a five-column sheet fails from line 1.

When `errors.Is(err, csv.ErrFieldCount)`, add the error to a slice and go round again. The record that came with it is the wrong width, so leave it out.

At `io.EOF`, return the records with `errors.Join` of the collected errors; with none collected, it returns `nil`. Any other error, such as a broken quote, ends the reading at once.
