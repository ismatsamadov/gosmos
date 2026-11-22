`strings.Fields` treats any run of white space as one separator, so the doubled space no longer leaves an empty field, and `fields[1]` is the sequence number again.
