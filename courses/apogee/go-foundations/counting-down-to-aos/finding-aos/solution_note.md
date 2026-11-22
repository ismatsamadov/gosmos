`break` ends the search at the first minute that clears the mask, so later minutes that also clear it never overwrite `aos`; when none does, the loop runs out and `aos` keeps -1.
