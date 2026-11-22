A datagram that fits leaves at least one byte of a `maxDatagram+1` slice unused, so `n` reaches 513 only when the datagram was longer than 512 bytes.

After the error check, compare `n` with `maxDatagram`; when it is more, return `nil`, `from` and `ErrTooBig` instead of the bytes.

If the 513-byte datagram comes back as 512 bytes and no error, the slice is one byte short of telling whole from cut. Nothing needs reading to throw a cut datagram's rest away: the next `ReadFrom` starts at the next datagram by itself.
