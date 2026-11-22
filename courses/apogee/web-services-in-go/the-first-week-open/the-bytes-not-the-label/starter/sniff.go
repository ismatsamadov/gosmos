package main

// checkCommands refuses a command file whose bytes are not plain text, whatever the form said it was.
func checkCommands(data []byte) error {
	return nil // TODO: refuse a file that http.DetectContentType doesn't call plain text
}
