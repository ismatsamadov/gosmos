`append` doesn't change `passes` itself. Where does the slice it returns go right now?

`_` is the blank identifier: whatever you assign to it is gone. Assign the result to the variable that should grow instead.
