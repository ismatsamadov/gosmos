`os.OpenFile` takes the path, the flags and the mode. The flags are `os.O_APPEND`, `os.O_CREATE` and `os.O_WRONLY`, joined with `|` into one number.

The mode is the last argument, written in octal as Go writes it: `0640`. Return what `os.OpenFile` returns; it is already a file and an error.

If the test finds only the second line, the second open wrote from the file's start: `os.Create` or `O_TRUNC` empties the file, and without `O_APPEND` the new line lands over the old one. If only the mode is wrong, check the last argument: `0644` lets every user read the log.
