`int16` reads the same 16 bits with the top bit as a sign, so `ff 83` is −125 rather than 65411. Converting to `float64` before dividing keeps the tenths: −12.5 °C.
