`WriteTo` is a method of the socket that sends, `rec`, and takes the bytes, `[]byte(frame)`, and where they go: the live view's own address, `view.LocalAddr()`. `ReadFrom` is a method of the socket that receives, `view`.

`ReadFrom(buf)` returns three values: how many bytes it copied, the sender's address and an error. `n` and `err` exist already, so `n, from, err :=` declares only `from`. When an error isn't nil, print it and return.

Two addresses are the same when their `String()` results are equal. If the program prints an i/o timeout after a second, the datagram went to an address `view` isn't on.
