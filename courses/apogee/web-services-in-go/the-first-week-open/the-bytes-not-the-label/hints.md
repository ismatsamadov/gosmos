Call `http.DetectContentType(data)` once and keep its answer: you need it for the comparison and for the message.

Compare the whole string, parameters included: `text/html` and `text/xml` are text too, and `text/plain; charset=utf-16le` is plain text the station can't read.

The status is `http.StatusUnsupportedMediaType`, and the message ends with the type exactly as the sniffer gave it, such as `command file must be plain text, not application/pdf`.
