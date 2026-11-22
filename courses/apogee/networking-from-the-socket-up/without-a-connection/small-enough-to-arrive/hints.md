Build the datagram being filled as a `[]byte`. Before adding a line to one that already holds lines, check its length plus one byte for the newline plus the line's length against `maxDatagram`.

When the line doesn't fit, append the datagram to the result and start an empty one; a line goes in with a newline before it only when the datagram already holds something. After the loop, a datagram still being filled belongs in the result too.

If the 513-byte case comes back as one datagram, the newline between the two lines wasn't counted. Check each line's own length against `maxDatagram` first, and wrap `ErrTooBig` with `%w`.
