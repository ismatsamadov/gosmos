package main

// reasons are the words after each status code, the same words net/http uses.
var reasons = map[int]string{
	200: "OK",
	201: "Created",
	400: "Bad Request",
	404: "Not Found",
	405: "Method Not Allowed",
	408: "Request Timeout",
	409: "Conflict",
	411: "Length Required",
	413: "Request Entity Too Large",
	414: "Request URI Too Long",
	431: "Request Header Fields Too Large",
	501: "Not Implemented",
	505: "HTTP Version Not Supported",
}

// dateLayout is how the Date field writes a time, in UTC.
const dateLayout = "Mon, 02 Jan 2006 15:04:05 GMT"

// Field is one header field of an answer, written as given.
type Field struct{ Name, Value string }
