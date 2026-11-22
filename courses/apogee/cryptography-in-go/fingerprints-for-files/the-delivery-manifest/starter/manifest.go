package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"strings"
)

// checkManifest checks every file the manifest lists against fsys and reports one line per file.
func checkManifest(fsys fs.FS, manifest io.Reader) ([]string, error) {
	_, _, _ = bufio.NewScanner, strconv.ParseInt, strings.Fields
	_, _ = errors.Is, fmt.Errorf
	return nil, nil
}
