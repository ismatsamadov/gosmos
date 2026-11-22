Make one `csv.Reader` for `r` before the loop, then call its `Read` method each time round and append the record it returns.

Check for `io.EOF` first: it means the sheet ended, and the records you have collected are the answer.

Any other error is a broken sheet, such as a quote that never closes. Return `nil` and that error rather than the rows read before it.
