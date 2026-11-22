## What does the first line of the request above hold?

- [x] The method, the path with its query, and the version
  > That is the request line: `GET`, then `/passes?sat=furrow-a&station=SQT`, then `HTTP/1.1`.
- [ ] The method and the whole URL, `http://127.0.0.1` included
  > The address travels in the `Host` header; the first line holds only the path and its query.
- [ ] The method and the path, without the query
  > The query rides on the first line too, right after the `?` that ends the path.

## Which method asks the service to cancel booking `bk-4494`?

- [ ] `GET`
  > `GET` reads the booking and changes nothing.
- [ ] `POST`
  > `POST` creates something new, such as a booking; it doesn't remove one.
- [x] `DELETE`
  > `DELETE` removes what the path names, here the booking `bk-4494`.

## What follows the `?` in `/passes?sat=furrow-a&station=SQT`?

- [ ] More of the path, one folder for each `&`
  > The path ends at the `?`; nothing after it is part of the path.
- [x] The query: `key=value` pairs joined by `&`
  > `sat=furrow-a` and `station=SQT` are two pairs, and `&` separates them.
- [ ] Headers, written into the URL
  > Headers travel on lines of their own; the query is part of the request line.
