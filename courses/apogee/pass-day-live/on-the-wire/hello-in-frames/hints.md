`binary.BigEndian`'s `PutUint16` and `PutUint32` write the length and the seq into a byte slice, and `crc32.ChecksumIEEE` gives the CRC of the bytes you hand it: the type, the seq and the body, never the length. The length counts every byte after itself, 9 plus the body's.

Read the frames through the `bufio.Reader` that read the banner: a radio that sends its banner and WELCOME in one write leaves the WELCOME in that reader's buffer. `io.ReadFull` reads exactly the 2 bytes of a length, then exactly as many as it counts, however the station splits them. Keep in the pass whether a WELCOME, a NAK or a BYE came, and read until the pass's context closes the call.

A NAK's body starts with its code, and code 2 is `NakStowed`. If a test says a call hung up before its LOS, look at what happens when a read fails: Varsk Island's line damages a frame now and then, and the call stays open until LOS all the same. If HELLO's bytes differ, the message names the part that does.
