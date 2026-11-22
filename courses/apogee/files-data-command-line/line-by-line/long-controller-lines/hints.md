The TODO sits where the call belongs: after the scanner is made and before `Scan` runs for the first time.

`sc.Buffer` takes a starting buffer, such as an empty byte slice with some capacity, and the largest size the buffer may grow to. Pass `maxRecordBytes` as that size, not a bigger number: the 300 KiB record must still fail with `bufio.ErrTooLong`.
