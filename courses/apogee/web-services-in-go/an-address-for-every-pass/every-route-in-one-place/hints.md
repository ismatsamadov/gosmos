Compare `routes` with the task's list line by line: one route is missing, and one pattern has no method.

A pattern without a method answers every method with the same handler, so `/bookings/{id}` has to become two routes, one per method, each with its own handler from `bookings.go`.

If `DELETE /bookings/bk-4497` gets 200 and a booking, it is still reaching `getBooking`. If `GET /passes` gets `404 page not found`, its route is missing. Each `Allow` header names exactly the methods registered on that path.
