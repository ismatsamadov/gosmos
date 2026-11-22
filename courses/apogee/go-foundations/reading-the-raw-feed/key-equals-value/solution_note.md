`strings.Cut` splits each field once, at its first `=`, so `note=a=b` keeps its value `a=b`, and checking found stops a word without any `=`, like `thistle`, from passing for a key.
