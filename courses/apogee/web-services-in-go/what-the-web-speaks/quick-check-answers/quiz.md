## The booking service will answer a new booking with `201`. What does it say that `200` doesn't?

- [ ] That the answer has no body
  > That is `204 No Content`; a `201` can carry a body, such as the booking it made.
- [x] That something new was created
  > `201 Created` says the request worked and made something, such as a booking.
- [ ] That the program should send it again later
  > Waiting and trying again is what `429` and the 5xx statuses ask for; `201` is a success.

## After the code above runs, what does `ct` hold?

- [x] `application/json`
  > `Get` matches header names in any letter case, so `content-type` finds the stand-in's `Content-Type`.
- [ ] An empty string, because the stand-in wrote `Content-Type`
  > Header names don't depend on letter case: `Get` finds the header however you spell its name.
- [ ] `text/plain`, a default the recorder adds
  > The stand-in set `application/json` itself, and the recorder keeps what the handler wrote.

## How does the request in the code above reach the stand-in?

- [ ] Over a network connection that stays inside your Mac
  > No connection is opened at all; a real connection comes tomorrow, with a server.
- [ ] Over the network, to the computer where Nkem wrote it
  > Nothing leaves your Mac: the stand-in's code runs inside your own program.
- [x] Through a plain function call, in memory
  > `ServeHTTP` is called like any function, and the recorder keeps what it writes.
