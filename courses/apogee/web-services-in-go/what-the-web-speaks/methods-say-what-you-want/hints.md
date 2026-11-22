Each line of `emails` pairs a method with a URL, and its comment says which email it is. The first email only asks a question, so it keeps `http.MethodGet`; the other two change something.

Booking a pass creates something new, and canceling removes something. `net/http` has a constant for each of those methods, and each one replaces an `http.MethodGet` in the list.

If the last line still starts with `GET`, look at the line whose URL ends in `bk-4494`: that is the cancellation, and its method is the one to change.
