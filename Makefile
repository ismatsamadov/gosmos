# Gosmos: build and install the app, run every test, check the library.
.PHONY: app install test check verify reference

APP := build/Gosmos.app
LIBRARY := courses/apogee.sqlite
REFERENCE := reference/go.sqlite
PACKAGES := CodeEditor

# The app, with the Go toolchain lessons run with built in: the one `go env GOROOT` names, or
# GO_ROOT=<folder>. app/scripts/bundle.sh trims and signs it.
app:
	swift build -c release --package-path app --product Gosmos
	app/scripts/bundle.sh "$$(swift build -c release --package-path app --show-bin-path)/Gosmos" $(APP)

# Checks this repository's library and installs the app. The app's settings are left alone.
install: app check
	rm -rf ~/Applications/Gosmos.app
	mkdir -p ~/Applications
	ditto $(APP) ~/Applications/Gosmos.app

# Every test: the Go tools, the app, and each of the app's own packages.
test:
	go vet ./...
	go test ./...
	swift test --package-path app
	for p in $(PACKAGES); do swift test --package-path app/Packages/$$p || exit 1; done

check:
	go run ./cmd/library check $(LIBRARY)
	go run ./cmd/reference check $(REFERENCE)

# Runs every code lesson of the library: each reference solution passes, each starter fails.
# Also wants 2 or 3 hints and a solution note for every code lesson, and an explanation for
# every quiz answer.
verify:
	go run ./cmd/library verify --teaching $(LIBRARY)

# Makes the Go reference anew: the standard library of the Go that `go env GOROOT` names, its
# spec and memory model, and Effective Go, the FAQ and How to Write Go Code from golang/website
# at WEBSITE_COMMIT. Those files are fetched once into $TMPDIR/godocs/website, the only network
# the reference ever needs; the app reads it offline.
WEBSITE_COMMIT := b8d34bbe2d9e230cc70ce3e4272dc651f494a188
WEBSITE_FILES := doc/effective_go.html doc/faq.md doc/code.html doc/progs/eff_bytesize.go \
	doc/progs/eff_qr.go doc/progs/eff_sequence.go doc/progs/eff_unused1.go doc/progs/eff_unused2.go

reference:
	website="$${TMPDIR:-/tmp}/godocs/website/$(WEBSITE_COMMIT)"; \
	for f in $(WEBSITE_FILES); do \
		test -s "$$website/_content/$$f" || { \
			curl -fsSL --create-dirs -o "$$website/_content/$$f.part" \
				"https://raw.githubusercontent.com/golang/website/$(WEBSITE_COMMIT)/_content/$$f" && \
			mv "$$website/_content/$$f.part" "$$website/_content/$$f"; } || exit 1; \
	done; \
	go run ./cmd/reference build --goroot "$$(go env GOROOT)" --website "$$website" \
		--website-commit $(WEBSITE_COMMIT) --out $(REFERENCE)
