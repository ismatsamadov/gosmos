The file frames come on the call like the others. FILE's body is the chunk count, a big-endian uint16, then the name; CHUNK's is the chunk's number, a uint16, then its bytes; END's is the CRC-32, a uint32. When a line is lost and dialed again, the station sends FILE again for the file in progress, then the chunks still to come.

Keep each file in one place every pass can reach, its chunks in a map by number under a mutex, and keep a number you already have only once. At END, join the chunks in number order, compare `crc32.ChecksumIEEE` of the result with END's, and save the file through a temporary file and a rename, as the plan's files are.

A file in `files/` with too many bytes kept a chunk that came twice; a file that never lands, with a chunk that came after a higher one, had its CRC taken over the chunks as they came. A file still incomplete after the week lists its missing numbers ascending, or `crc mismatch` when every chunk came and END's CRC didn't match.
