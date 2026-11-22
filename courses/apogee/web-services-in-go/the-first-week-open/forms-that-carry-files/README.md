---
slug: forms-that-carry-files
title: Forms That Carry Files
type: code_tests
position: 0
difficulty: 9
topics: errors, http, streams
created: 2026-11-13T06:35:18Z
modified: 2026-11-13T16:15:49Z
---
Wednesday 2 February, 09:10. Customers now attach a command file to a booking: the commands they want sent up during the pass. Cirrine's engineer attached Cirrine-1's whole software image instead, 1.9 GB, and the draft read every byte into memory before looking at any of them.

A browser or a script sends a file as a form, `multipart/form-data`: parts between boundary lines, each with headers of its own. This body holds one part, named `file`:

```text
--b7
Content-Disposition: form-data; name="file"

PING
SET_MODE SAFE
--b7--
```

`r.MultipartReader()` reads such a body one part at a time. Its `NextPart` returns the next part, a reader of that part's bytes whose `FormName()` is its name, and `io.EOF` after the last one.

## Your task

`readCommands` returns the whole body. Wrap `r.Body` in `http.MaxBytesReader` with `maxBodyBytes`, then return the bytes of the part named `file`. Return any error as it comes, and `http.ErrMissingFile` when no part has that name.
