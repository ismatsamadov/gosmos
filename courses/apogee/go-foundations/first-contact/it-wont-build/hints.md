Read the error from the left: `./main.go:5:2` means file `main.go`, line 5, column 2. What is on line 5, and what does the rest of the error say about it?

A package counts as used only when the code refers to something in it, the way `main` calls `fmt.Println`. Does anything in `main` come from `os`?
