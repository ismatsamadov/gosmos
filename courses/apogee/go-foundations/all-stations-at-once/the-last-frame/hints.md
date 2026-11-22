`range frames` ends only when the channel is closed and empty. Which function knows that the last frame has been sent?

The built-in `close` takes the channel to close. Call it in `poller`, after the loop that sends the frames.
