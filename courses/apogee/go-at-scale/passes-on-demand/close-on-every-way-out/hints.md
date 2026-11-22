The close at the end runs only when the loop over `ReadPasses` finishes. A loop over `Archive` that breaks, or stops at a damaged line, makes `yield` return false, and the iterator returns before it gets there.

A deferred call runs when the function it's in returns, whichever `return` that is. Defer the file's `Close` right after the `if err != nil` block, where `f` is a real file.

Remove the close at the end too, or a full read closes the file twice. The test counts closes: one after a full read, a loop that breaks after one pass and a loop that stops at an error.
