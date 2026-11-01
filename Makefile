# Gosmos: build and install the app, run every test, check the library.
.PHONY: app install test check verify

APP := build/Gosmos.app
LIBRARY := courses/apogee.sqlite
PACKAGES := CodeEditor JQ LessonHTTP

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

# Runs every code lesson of the library: each reference solution passes, each starter fails.
# Also wants 2 or 3 hints and a solution note for every code lesson, and an explanation for
# every quiz answer.
verify:
	go run ./cmd/library verify --teaching $(LIBRARY)
