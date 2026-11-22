package main

import (
	"bufio"
	"fmt"
	"strings"
)

// readBanner reads the controller's banner line and returns the controller's name.
func readBanner(br *bufio.Reader) (string, error) {
	line, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSuffix(line, "\n")
	// TODO: also remove a \r that comes right before the \n
	parts := strings.Split(line, " ")
	if len(parts) != 4 || parts[1] != "READY" || parts[2] != "PROTO" || !validName(parts[0]) {
		return "", fmt.Errorf("not a banner: %q", line)
	}
	if parts[3] != "1" {
		return "", fmt.Errorf("unsupported protocol %q", parts[3])
	}
	return parts[0], nil
}

// validName reports whether name is 1 to 16 capital letters, digits and hyphens.
func validName(name string) bool {
	return len(name) >= 1 && len(name) <= 16 && strings.Trim(name, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-") == ""
}
