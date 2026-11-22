A method goes inside the pattern string itself: the method, one space, then the path. The handlers stay as they are.

All three patterns need it, the guide's `/guide/` too. A pattern without a method still matches every method, so one forgotten pattern still answers a POST.

If `GET /stations` gets 404, a pattern is missing its space: without one, the mux reads `GET` as a host name. If a POST still gets 200, one of the three patterns has no method yet.
