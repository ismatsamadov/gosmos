## What does the first line print?

- [ ] `[7.9]`
  > The 5 is a width: `7.9` is only three characters, so two spaces are added.
- [x] `[  7.9]`
  > One decimal turns 7.91 into `7.9`, and the width of 5 adds two spaces on the left.
- [ ] `[7.9  ]`
  > Padding goes on the left unless the verb has a `-`, so the spaces come before the number.

## What does the second line print?

- [x] `[YST ]`
  > The `-` puts the padding on the right: `YST` and one space make 4 characters.
- [ ] `[ YST]`
  > That is what `%4s` prints; the `-` moves the padding to the right.
- [ ] `[YST]`
  > A width is a minimum, so a three-letter code gets one space to reach 4.

## What does the third line print?

- [ ] `[  7]`
  > Spaces are the default padding; the `0` before the width asks for zeros.
- [ ] `[700]`
  > Zeros pad on the left, where they don't change the value; `700` would be a different number.
- [x] `[007]`
  > The `0` pads with zeros up to the width of 3.
