`poller` sends three values, but `main` receives only once. How many receives does it take to get all three?

Each `<-frames` takes exactly one value off the channel, so the receive and the `Println` belong in a loop that runs once for each frame `poller` sends.
