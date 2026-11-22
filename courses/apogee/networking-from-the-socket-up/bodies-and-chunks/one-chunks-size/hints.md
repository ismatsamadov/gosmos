`strings.Cut(line, ";")` gives the part before the first semicolon, or the whole line when there is none, and `strings.TrimRight` takes spaces and tabs off its end only: a space before the digits is still refused.

Check the digits yourself: none, more than 16, or any character outside `0123456789abcdefABCDEF` is a bad size, and `strings.Trim` with that set leaves an empty string only when every character belongs. Then parse them with `strconv.ParseInt` in base 16.

If `1000` comes back as 1000, the size is still read in base 10. The refusal quotes the whole line as it came, with `%q`.
