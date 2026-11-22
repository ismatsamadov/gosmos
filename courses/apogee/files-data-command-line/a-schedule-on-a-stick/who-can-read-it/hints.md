Inspect the path first and return any error from `Stat`. Its successful result holds the file mode.

Apply `Perm` before masking the group and other bits. Zero after that mask means none of those permissions are set.
