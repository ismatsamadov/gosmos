Call `run` first and keep its error. When it is nil, return nil at once: a flown pass leaves no file behind.

For a failure, create the file with `os.Create`; if that fails, return `errors.Join` of the pass's error and the file's. Otherwise `fr.WriteTo(f)` writes the snapshot and returns a count you don't need and an error you do. Join the pass's error, that one and `f.Close()`'s: `errors.Join` drops the nil ones.

Don't call `fr.Stop`: the recorder must keep recording for the next pass. If the test can't find the pass's error, return the joined error, not the snapshot's alone.
