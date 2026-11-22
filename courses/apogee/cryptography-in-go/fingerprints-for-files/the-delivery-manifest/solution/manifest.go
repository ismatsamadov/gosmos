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
	var report []string
	sc := bufio.NewScanner(manifest)
	for n := 1; sc.Scan(); n++ {
		l, ok := parseManifestLine(sc.Text())
		if !ok {
			return nil, fmt.Errorf("manifest line %d: want <digest> <bytes> <file>", n)
		}
		state, err := checkManifestFile(fsys, l)
		if err != nil {
			return nil, err
		}
		report = append(report, state+" "+l.name)
	}
	return report, sc.Err()
}

// manifestLine is one line of a manifest: a file's digest, the hash that made it, its size and its name.
type manifestLine struct {
	digest, alg, name string
	size              int64
}

// parseManifestLine reads <digest> <bytes> <file>, and reports whether the text is such a line.
func parseManifestLine(text string) (manifestLine, bool) {
	f := strings.Fields(text)
	if len(f) != 3 {
		return manifestLine{}, false
	}
	alg, _, err := parseDigest(f[0])
	size, sizeErr := strconv.ParseInt(f[1], 10, 64)
	if err != nil || sizeErr != nil || size < 0 {
		return manifestLine{}, false
	}
	return manifestLine{digest: f[0], alg: alg, name: f[2], size: size}, true
}

// checkManifestFile returns ok, wrong size, changed or missing for the file l lists.
func checkManifestFile(fsys fs.FS, l manifestLine) (string, error) {
	f, err := fsys.Open(l.name)
	if errors.Is(err, fs.ErrNotExist) {
		return "missing", nil
	} else if err != nil {
		return "", err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return "", err
	}
	if info.Size() != l.size {
		return "wrong size", nil
	}
	got, err := fingerprintAs(f, l.alg)
	if err != nil {
		return "", err
	}
	if got != l.digest {
		return "changed", nil
	}
	return "ok", nil
}
