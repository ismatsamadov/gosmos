Multiply each angle by 100 and round it with `math.Round`, and keep both results as `float64` until they have passed the range checks.

Check the rounded values: 0 to 35999 for the azimuth, 0 to 9000 for the elevation. 359.996 rounds to 36000, which is past the last hundredth before north.

`binary.BigEndian.AppendUint16` takes a `uint16`: convert the rounded azimuth to one, and the elevation to `int16` first, the type the protocol gives it. If 1.15 comes out as `00 72`, the conversion happened without rounding.
