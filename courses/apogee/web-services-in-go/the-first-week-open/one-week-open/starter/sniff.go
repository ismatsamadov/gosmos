package main

import "net/http"

// checkCommands refuses a command file whose bytes are not plain text, whatever the form said it was.
func checkCommands(data []byte) error {
	if kind := http.DetectContentType(data); kind != "text/plain; charset=utf-8" {
		return &APIError{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "command file must be plain text, not " + kind}
	}
	return nil
}
