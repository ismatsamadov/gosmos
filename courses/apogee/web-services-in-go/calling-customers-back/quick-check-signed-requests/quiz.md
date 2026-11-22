## Why does the signature cover the timestamp as well as the body?

- [ ] So a receiver can sort the events it gets by time
  > Sorting needs no signature: the time in the header would do. The signature is about trust.
- [x] So a request someone copied can't be sent again later under a new time
  > A new timestamp breaks the signature, and the old one soon falls outside the five minutes.
- [ ] So two events with the same body get signatures of different lengths
  > An HMAC-SHA256 is always 64 hex digits, whatever it covers.

## When may the receiver trust the `booking_id` in the body?

- [ ] As soon as the body decodes as JSON
  > Anyone can send well-formed JSON; decoding it proves nothing about who sent it.
- [ ] Once it has answered the request with `204 No Content`
  > The answer comes after every check; it can't make the body any more trustworthy.
- [x] Once the signature it computes from the bytes that arrived matches the header
  > Only the service and the customer know the secret, so a match proves the service sent exactly these bytes.

## Why does the receiver accept a time up to 300 seconds from its own clock, either way?

- [x] Because the service's clock and the receiver's can be a little apart
  > No two clocks agree exactly; five minutes allows for that while an old request is still refused.
- [ ] Because a retry can reach the receiver five minutes after the first attempt
  > Every attempt is signed again with its own timestamp, so a retry is never older than its own sending.
- [ ] Because checking a signature can take up to five minutes
  > Computing and comparing an HMAC takes microseconds.
