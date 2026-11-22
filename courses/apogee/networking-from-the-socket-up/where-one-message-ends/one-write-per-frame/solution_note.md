One `Write` per frame means a connection shared by goroutines carries each frame whole, since nothing can land inside a single call. Returning its error is how the caller learns the line has gone.
