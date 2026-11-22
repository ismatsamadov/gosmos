`main` gives `runPass` a context that ends in 30 ms. Follow it down the calls: where does it stop being passed along?

`context.Background()` in `sendCommand` starts a new context with no deadline. Replace it with a `ctx` parameter, first in the list, which `runPass` fills in.
