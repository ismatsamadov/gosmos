Both branches of the draft return a different text: the answer must not depend on the name at all, so it can come after the `if`, once, as `resetSent`.

The `if` decides only whether anything is sent. When `r.known(user)` is true, `r.send(user, r.Start(user))` starts the reset and queues its link; otherwise nothing happens.
