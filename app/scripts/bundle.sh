#!/bin/sh
# Assembles Gosmos.app from a built binary, with the Go toolchain it runs lessons with, and
# signs it ad hoc.
# Usage: app/scripts/bundle.sh <binary> <app path>   (make app runs it)
# GO_ROOT=<folder> is the Go to build in; the default is the one `go env GOROOT` names.
set -eu

binary=$1
app=$2
here=$(cd "$(dirname "$0")/.." && pwd)
repo=$(cd "$here/.." && pwd)

if [ ! -x "$binary" ]; then
	echo "bundle.sh: $binary is not an executable; build it first (make app)" >&2
	exit 1
fi

# The Go built into the app is darwin/arm64, all macOS 27 runs on.
go_root=${GO_ROOT:-$(cd / && GOTOOLCHAIN=local go env GOROOT)}
if [ ! -x "$go_root/bin/go" ] || [ ! -x "$go_root/bin/gofmt" ] || [ ! -d "$go_root/pkg/tool/darwin_arm64" ]; then
	echo "bundle.sh: $go_root is not a darwin/arm64 Go toolchain; set GO_ROOT to one" >&2
	exit 1
fi

if [ ! -f "$repo/db/learner.sql" ]; then
	echo "bundle.sh: $repo/db/learner.sql is missing" >&2
	exit 1
fi

# The app is put together in a folder of its own under $TMPDIR, then moved into place whole,
# the old one out first: build/ is inside this repository's Go module, and a go command walking
# ./... (go vet, go test) must never meet a half-copied Go there.
stage=$(mktemp -d "${TMPDIR:-/tmp}/gosmos-bundle.XXXXXX")
trap 'rm -rf "$stage"' EXIT
new="$stage/$(basename "$app")"

mkdir -p "$new/Contents/MacOS" "$new/Contents/Resources"
cp "$binary" "$new/Contents/MacOS/Gosmos"
cp "$here/Resources/Info.plist" "$new/Contents/Info.plist"
cp "$here/Resources/AppIcon.icns" "$new/Contents/Resources/AppIcon.icns"

# Gosmos's license travels with the app; the Go toolchain keeps its own license below.
cp "$repo/LICENSE" "$new/Contents/Resources/LICENSE"

# The learner database's schema, from which the app makes ~/Library/Application Support/Gosmos/
# learner.sqlite on first launch.
cp "$repo/db/learner.sql" "$new/Contents/Resources/learner.sql"

# The course libraries built into the app, opened read-only and immutable.
mkdir -p "$new/Contents/Resources/Courses"
for course in "$repo"/courses/*.sqlite; do
	if [ -e "$course" ]; then cp "$course" "$new/Contents/Resources/Courses/"; fi
done

# The Go documentation the app shows offline, when the repository has it.
mkdir -p "$new/Contents/Resources/Reference"
for reference in "$repo"/reference/*.sqlite; do
	if [ -e "$reference" ]; then cp "$reference" "$new/Contents/Resources/Reference/"; fi
done

# GRDB's privacy manifest, a resource bundle SwiftPM puts next to the binary. Only the bundles
# named here go in: the build folder keeps those of packages the app no longer uses.
for name in GRDB_GRDB; do
	bundle="$(dirname "$binary")/$name.bundle"
	if [ ! -d "$bundle" ]; then
		echo "bundle.sh: $bundle is missing" >&2
		exit 1
	fi
	cp -R "$bundle" "$new/Contents/Resources/"
done

# The Go toolchain, so lessons build, run, test and format with nothing installed: GO_ROOT
# without what building and testing standard-library code never uses (docs/architecture.md
# says why each goes): Go's own tests and test data, its API lists, documents and extras, the
# WebAssembly support, the time zone file's generator, the sources of the go command and its
# tools (their programs stay), and the fix and cover tools.
go="$new/Contents/Resources/Go"
ditto "$go_root" "$go"
rm -rf "$go/api" "$go/doc" "$go/misc" "$go/test" "$go/lib/wasm" "$go/lib/time/mkzip.go" "$go/src/cmd" \
	"$go/pkg/tool/darwin_arm64/fix" "$go/pkg/tool/darwin_arm64/cover"
find "$go/src" -type d -name testdata -prune -exec rm -rf {} +
# Every .go file left is in Go's own module (src/go.mod), which ./... in this repository skips.
if [ ! -f "$go/src/go.mod" ] || [ -n "$(find "$go" -name '*.go' ! -path "$go/src/*" -print)" ]; then
	echo "bundle.sh: the Go in $go_root has .go files outside the module in its src folder" >&2
	exit 1
fi
# Go's license, and its patent grant when there is one, go with it: from GOROOT, else from the
# folder above, where Homebrew keeps the license. Without the license there is no app.
for file in LICENSE PATENTS; do
	for dir in "$go_root" "$go_root/.."; do
		if [ -f "$dir/$file" ]; then
			cp "$dir/$file" "$go/$file"
			break
		fi
	done
done
if [ ! -f "$go/LICENSE" ]; then
	echo "bundle.sh: found no LICENSE for the Go in $go_root or in the folder above it" >&2
	exit 1
fi

# Every program in the toolchain is signed like the app, then the app. A Developer ID build
# signs each with its identity, the hardened runtime and a timestamp (docs/architecture.md).
find "$go" -type f -perm -u+x -exec sh -c '
	set -e
	for f; do
		case $(file -b "$f") in
		Mach-O*executable* | Mach-O*dynamically\ linked* | Mach-O*bundle*)
			codesign --remove-signature "$f"
			codesign --sign - --timestamp=none "$f" ;;
		esac
	done' sh {} +
codesign --force --sign - --timestamp=none "$new"

mkdir -p "$(dirname "$app")"
if [ -e "$app" ]; then mv "$app" "$stage/previous"; fi
mv "$new" "$app"
echo "Built $app with $("$app/Contents/Resources/Go/bin/go" version | cut -d' ' -f3), $(du -sh "$app" | cut -f1)"
