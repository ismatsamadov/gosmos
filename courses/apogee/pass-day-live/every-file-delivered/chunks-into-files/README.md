---
slug: chunks-into-files
title: Chunks Into Files
type: code_tests
position: 0
difficulty: 9
topics: bytes, files, maps
created: 2026-11-12T05:36:26Z
modified: 2026-11-12T11:16:23Z
---
Hollis has heard about the investors' week, and his email to Yusra is one subject line: "EVERY CHUNK OR NOTHING."

On a booked pass the satellite sends a file down in numbered pieces of 256 bytes, called chunks: a FILE frame with the file's chunk count and name, a CHUNK frame for each chunk, then END with the CRC-32 of the whole file. Chunks come out of order now and then, and now and then twice; a CHUNK or END belongs to the last FILE on its call. A file lands when every chunk is there and their CRC-32, in number order, matches END's:

```text
Thu 16:04:25 LANDED thistle-2-20290726-1603.dat
```

## Your task

Put each file back together. At its END, if it is whole, save it as `files/<name>`, never partial, and print `LANDED <name>`. After the week, write `reports/incomplete.txt`, one line per file still incomplete, by name:

```text
furrow-k-20290724-0740.tif missing 21 22 23
thistle-2-20290726-1603.dat crc mismatch
```

The tests compare every file of the week byte for byte, Furrow-K's Tuesday image three chunks short among them, and fly Thursday at Port Ellery with a line lost mid-file, a wrong CRC and two chunks missing.
