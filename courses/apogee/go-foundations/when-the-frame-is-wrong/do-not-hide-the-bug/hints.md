The panic value has type `any`. A comma-ok type assertion tells you whether it is the known `decoderFault`.

Use the fault's reason for the error only when the assertion succeeds. Otherwise panic again with the original value, without converting it to text.
