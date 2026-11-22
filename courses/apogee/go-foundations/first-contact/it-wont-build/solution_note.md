Go refuses to build a program that imports a package it never uses. Without the `"os"` import only `fmt` is left, `main` uses it, and the banner builds and prints.
