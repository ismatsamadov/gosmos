Change only the field's type to a pointer; its tag stays as it is.

After decoding, a nil `MaskDeg` means the update said nothing about the mask, so return `oldDeg`. Otherwise return the value it points to, reading it with `*`.
