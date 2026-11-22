`ReadFile` returns a byte slice followed by an error. Return that error before converting anything if the read failed.

The task asks for exact text. Do not split, trim or add a newline to the bytes you read.
