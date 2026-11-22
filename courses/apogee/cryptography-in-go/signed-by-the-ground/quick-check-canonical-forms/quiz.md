## The two commands share one canonical text. What does that mean for their signatures?

- [x] A signature for one is a signature for the other
  > The signature covers only the canonical text, and the two texts are the same bytes.
- [ ] Each still gets its own signature, since the commands differ
  > The signature is computed from the canonical text alone, which is the same for both.
- [ ] The second command's signature fails to verify
  > It verifies, which is the trouble: it is the first command's signature too.

## Which change gives every canonical text exactly one meaning?

- [ ] Join the fields with `,` instead of `|`
  > A comma inside an argument would cause the same trouble as a bar.
- [ ] Sign the arguments in alphabetical order
  > Order changes nothing here: the bar inside `SAFE|REBOOT` still reads as two fields.
- [x] Refuse `|` inside fields, or write each field's length before it
  > Either way every text splits into fields in exactly one way, so it has one meaning.
