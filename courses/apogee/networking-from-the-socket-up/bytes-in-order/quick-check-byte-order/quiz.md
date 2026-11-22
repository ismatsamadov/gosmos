## What does `az` hold?

- [ ] 38215
  > That is `95 47`: the bytes read little-endian, as page 14 of the manual would have it.
- [x] 18325
  > Big-endian puts the most significant byte first: 0x4795 is 18325, which is 183.25°.
- [ ] 4795
  > `47 95` are hexadecimal digits; read as a number they are 0x4795, not the decimal 4795.

## Why does `state` need no byte order?

- [x] It is a single byte, so there is no order to choose
  > Byte order only says how the bytes of one number are arranged; one byte has nothing to arrange.
- [ ] Go reads single bytes big-endian unless told otherwise
  > There is no order to apply to one byte, in Go or anywhere else.
- [ ] Its values are small, 0 to 4, so the order can't change them
  > Small values in a two-byte field still change with the order: 1 is `00 01` one way and `01 00` the other.

## What does `temp` hold?

- [ ] 65411
  > That is the same bits read as a `uint16`; `int16` reads them with a sign.
- [ ] −12.5
  > That is the temperature in degrees, after dividing by 10; `temp` holds the tenths.
- [x] −125
  > 0xff83 has its top bit set, so `int16` reads it as negative: −125 tenths of a degree.
