Make one decoder for the whole stream before the loop, and inside the loop declare a fresh `Booking` for each `Decode` call to fill.

Check for `io.EOF` before any other error: it is the one error that means the log ended cleanly, and the bookings collected so far are the answer.

On any other error, return the slice you have together with that error. The test with a line cut off expects the three complete bookings and `io.ErrUnexpectedEOF`.
