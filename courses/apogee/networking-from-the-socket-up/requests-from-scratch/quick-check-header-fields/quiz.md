## How does `readHeader` keep `Content-Type: Text/Plain` and `content-type: text/plain`?

- [ ] As two fields, one for each spelling of the name
  > Names match in any letter case, so both lines go under one name, `content-type`.
- [x] As one field, `content-type`, with both values as they were sent
  > Names are stored in lower case, while values keep theirs: `Text/Plain` stays `Text/Plain`.
- [ ] As one field holding `text/plain` twice, since the server lower-cases both
  > Only the name is lower-cased; a value is kept exactly as sent, spaces around it aside.

## What does `readHeader` keep from the two X-Pass lines?

- [ ] Only the last value, `cirrine-1/YST/13094`
  > Nothing is overwritten: each value is appended to the name's slice.
- [ ] Only the first value, `furrow-a/YST/18218`
  > The second line isn't dropped; it adds a second value.
- [x] Both values, in the order they came
  > A Header keeps every value of a name in order, and Get joins them with ", ".

## Why is `Host : 127.0.0.1` refused rather than read as Host?

- [x] Another program may read that line differently, and refusing is safer than guessing which way it went
  > If the station server took it as Host and the next program didn't, the two would act on different requests; a refusal leaves nothing to guess, and HTTP itself requires this 400.
- [ ] Because no space is allowed anywhere in a header line
  > Values may hold spaces, and spaces around a value are trimmed; only the name must be exactly a token.
- [ ] Because Go's own client sends that line and must be told
  > Go's client never puts a space before the colon; the refusal is for programs that do.
