# Gosmos: build and install the app, run every test, check the library.
.PHONY: app install test check verify reference library notarize

APP := build/Gosmos.app
LIBRARY := courses/apogee.sqlite
REFERENCE := reference/go.sqlite
PACKAGES := CodeEditor

# The Apogee library, built from its text in courses/apogee/: the same text makes the same
# file, and a file already the same is left alone.
library:
	go run ./cmd/author build courses/apogee $(LIBRARY)

# The app, with the Go toolchain lessons run with built in: the one `go env GOROOT` names, or
# GO_ROOT=<folder>. app/scripts/bundle.sh trims and signs it.
app: library
	swift build -c release --package-path app --product Gosmos
	app/scripts/bundle.sh "$$(swift build -c release --package-path app --show-bin-path)/Gosmos" $(APP)

# Sends a Developer ID build (SIGN_IDENTITY="Developer ID Application: …" make app) to Apple's
# notary service, staples the ticket to the app, and zips it as build/Gosmos.zip to hand out.
# NOTARY_PROFILE names the notarytool credentials kept in the keychain, which
# `xcrun notarytool store-credentials <profile>` saves once (docs/architecture.md, "Signing").
NOTARY_PROFILE ?= gosmos
notarize:
	@codesign -dv --verbose=2 $(APP) 2>&1 | grep -q '^Authority=Developer ID Application' || \
		{ echo "notarize: $(APP) is not signed with a Developer ID; make app with SIGN_IDENTITY set" >&2; exit 1; }
	@codesign -dv $(APP) 2>&1 | grep -q '^CodeDirectory.*flags=.*runtime' || \
		{ echo "notarize: $(APP) is not signed with the hardened runtime" >&2; exit 1; }
	codesign --verify --deep --strict $(APP)
	rm -f build/Gosmos.zip
	ditto -c -k --keepParent $(APP) build/Gosmos.zip
	xcrun notarytool submit build/Gosmos.zip --keychain-profile $(NOTARY_PROFILE) --wait
	xcrun stapler staple $(APP)
	rm -f build/Gosmos.zip
	ditto -c -k --keepParent $(APP) build/Gosmos.zip
	spctl --assess --type execute -vv $(APP)

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

# The library's checks: what the app checks when it opens it, and given code that changes only
# where a lesson changes it or courses/apogee-drift.txt says why; then the Go reference's.
check: library
	go run ./cmd/library check $(LIBRARY)
	go run ./cmd/library drift --accept courses/apogee-drift.txt $(LIBRARY)
	go run ./cmd/reference check $(REFERENCE)

# Runs every code lesson of the library: each reference solution passes, each starter fails.
# Also wants 2 or 3 hints and a solution note for every code lesson, and an explanation for
# every quiz answer.
verify: library
	go run ./cmd/library verify --teaching $(LIBRARY)

# Makes the Go reference anew: the standard library of the Go that `go env GOROOT` names, its
# spec and memory model, and Effective Go, the FAQ and How to Write Go Code from golang/website
# at WEBSITE_COMMIT, whose files reference/website/ keeps (its README.md). Nothing is fetched,
# and the app reads the reference offline.
WEBSITE_COMMIT := b8d34bbe2d9e230cc70ce3e4272dc651f494a188

reference:
	go run ./cmd/reference build --goroot "$$(go env GOROOT)" --website reference/website \
		--website-commit $(WEBSITE_COMMIT) --out $(REFERENCE)
