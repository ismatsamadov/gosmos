`strings.Split(line, " ")` gives exactly three parts only when the line has exactly two spaces. Check the count, then that the method isn't empty, that the target starts with `/` and that `versionShaped` accepts the version.

Every failed check gives the same refusal, a `RequestError` with status 400 and the message `bad request line`. A good shape with another version gets 505, with the version quoted by `%q` in the message.

If the two-spaces or the tab case passes when it shouldn't, the line went through `strings.Fields`, which treats any run of spaces and tabs as one break.
