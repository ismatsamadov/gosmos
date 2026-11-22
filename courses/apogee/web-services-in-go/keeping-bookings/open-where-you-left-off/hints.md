`os.Open` returns an error that `errors.Is(err, fs.ErrNotExist)` recognizes when there is no file: return the empty store then. Any other error from opening is a real failure; otherwise defer closing the file.

Make a `json.Decoder` on the file, call its `DisallowUnknownFields`, and decode into a `[]Booking`. A second `Decode` must return `io.EOF`, or something follows the array.

Pass the decoded bookings to `f.mem.load` and return its error too: it refuses an ID saved twice. Every failure returns `nil` for the store, so no caller can save over a file it couldn't read.
