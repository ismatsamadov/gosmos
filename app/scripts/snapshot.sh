#!/bin/bash
# snapshot.sh <out.png> <light|dark> [launch arguments...]
# snapshot.sh --clean
#
# Draws one screen of the app to a PNG, to look at lessons the way a learner meets them
# (docs/authoring.md, "Seeing it in the app"). It runs a copy of build/Gosmos.app (make app)
# under its own bundle identifier, app.gosmos.Gosmos.qa, with its own learner data, so the real
# app's settings, window layout and progress stay untouched. The app draws its own window
# (-GosmosSnapshot), so this works without screen recording permission. The copy, made again
# whenever build/Gosmos.app is newer, and its data live in $QA (default $TMPDIR/gosmos-qa).
#
#   QA_LIBRARY=<file>   the library the copy shows in place of courses/apogee.sqlite, such as
#                       a scratch copy with new lessons imported
#   QA_DATA=<folder>    the learner data folder (default $QA/data); a new one starts in the
#                       first course's opening
#   QA_SIZE=1400x900    the window's size (taller for a long readme)
#   QA_DELAY=3          seconds before drawing (a failed Submit needs about 8)
#   QA_ID=qa            the last part of the copy's bundle identifier; two copies that run at
#                       once need two, each with its own QA folder
#
# The launch arguments pick the screen: -GosmosOpenLesson <lesson uuid> and
# -GosmosLaunchAction <action>; docs/authoring.md lists them. --clean removes the copy, its data
# and its settings.
set -euo pipefail
here="$(cd "$(dirname "$0")" && pwd)"
repo="$(cd "$here/../.." && pwd)"
QA="${QA:-${TMPDIR:-/tmp}/gosmos-qa}"
id=app.gosmos.Gosmos.${QA_ID:-qa}
app="$QA/Gosmos QA.app"
lsregister=/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister

if [ "${1:-}" = --clean ]; then
	[ -d "$app" ] && "$lsregister" -u "$app" 2>/dev/null || true
	rm -rf "$QA" "$HOME/Library/Saved Application State/$id.savedState"
	defaults delete "$id" 2>/dev/null || true
	rm -f "$HOME/Library/Preferences/$id.plist"
	exit 0
fi
if [ $# -lt 2 ] || { [ "$2" != light ] && [ "$2" != dark ]; }; then
	echo "usage: snapshot.sh <out.png> <light|dark> [launch arguments...] | --clean" >&2
	exit 2
fi
out=$1 appearance=$2
shift 2
built="$repo/build/Gosmos.app"
[ -d "$built" ] || { echo "snapshot.sh: no $built; run make app first" >&2; exit 1; }

if [ ! -d "$app" ] || [ "$built" -nt "$app" ] || [ -n "${QA_LIBRARY:-}" ]; then
	rm -rf "$app"
	mkdir -p "$QA"
	ditto "$built" "$app"
	/usr/libexec/PlistBuddy -c "Set :CFBundleIdentifier $id" -c "Set :CFBundleName Gosmos QA" "$app/Contents/Info.plist"
	if [ -n "${QA_LIBRARY:-}" ]; then
		cp "$QA_LIBRARY" "$app/Contents/Resources/Courses/apogee.sqlite"
	fi
	codesign --force --sign - --timestamp=none "$app"
fi

data="${QA_DATA:-$QA/data}"
mkdir -p "$data" "$(dirname "$out")"
rm -f "$out"
"$app/Contents/MacOS/Gosmos" -GosmosDataFolder "$data" -GosmosAppearance "$appearance" \
	-GosmosWindowSize "${QA_SIZE:-1400x900}" -GosmosSnapshotDelay "${QA_DELAY:-3}" -GosmosSnapshot "$out" "$@" \
	>/dev/null 2>&1 &
pid=$!
# A snapshot run never quits by itself: wait for the PNG to stop growing, then stop the app.
last=-1
for _ in $(seq 1 240); do
	sleep 0.25
	if [ -s "$out" ]; then
		size=$(stat -f %z "$out")
		[ "$size" = "$last" ] && break
		last=$size
	fi
done
{
	kill "$pid" || true
	for _ in $(seq 1 20); do kill -0 "$pid" || break; sleep 0.1; done
	kill -9 "$pid" || true
	wait "$pid" || true
} 2>/dev/null # the shell reports the stopped app on stderr
[ -s "$out" ] || { echo "snapshot.sh: the app drew nothing to $out" >&2; exit 1; }
echo "$out"
