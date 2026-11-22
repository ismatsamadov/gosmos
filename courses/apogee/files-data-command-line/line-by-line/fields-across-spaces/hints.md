`strings.Fields` does all of the splitting: tabs, repeated spaces and white space at the ends need no handling of their own.

Check how many fields came back before returning them. The error's text must match exactly, so build it with `fmt.Errorf` and the length of the slice.
