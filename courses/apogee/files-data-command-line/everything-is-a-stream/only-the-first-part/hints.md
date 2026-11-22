Put the limit around the reader before consuming bytes. Slicing afterwards cannot undo reads from the source.

Passing the limited reader to `ReadAll` preserves the usual error handling while stopping consumption at the limit.
