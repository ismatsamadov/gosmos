#!/bin/bash
# pkgcheck.sh <folder>...
#
# Checks lesson packages before they are imported (docs/authoring.md, "Checks"), with the Go on
# PATH. Each <folder> is a package (it holds README.md) or any folder above packages, such as a
# chapter's or a whole export. For every package it checks:
#   - front matter: type, difficulty, topics; the readme's length (60-200 words, quizzes
#     60-120) and its "## Your task" heading
#   - code lessons: 2 or 3 hints, a 1-3 sentence solution note, gofmt
#   - solution: passes 30 runs in a row (tests: go test -count=30; code: 30 runs printing
#     expected_output.txt byte for byte), and tests also pass under -race (-count=3)
#   - starter: builds and fails for the lesson's reason (tests fail; code prints something else,
#     compared as the app compares: line endings, trailing spaces and final newlines ignored)
# It prints one line per problem and a summary. Exit status 1 when there were problems.
#
# Set QUICK=1 to run the solution 3 times instead of 30 and skip -race (for a fast first pass).

set -u
export GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOFLAGS=
# The go.mod a lesson without one gets, as the app writes it: the Go's major and minor version.
goversion=$(go env GOVERSION | sed -E 's/^go([0-9]+\.[0-9]+).*/\1/')
runs=30; racecount=3
if [ "${QUICK:-}" = 1 ]; then runs=3; racecount=0; fi

problems=0
checked=0
say() { echo "$1: $2"; problems=$((problems + 1)); }

body() { awk 'f>=2{print} /^---$/{f++}' "$1/README.md"; }
field() { awk -v k="$2" 'NR>1 && /^---$/{exit} index($0, k": ")==1{print substr($0, length(k)+3)} $0==k":"{print ""}' "$1/README.md"; }
paragraphs() { awk 'BEGIN{RS=""} END{print NR}' "$1"; }

# lay out <pkg> <dir> <solved>: the lesson's files as the app writes them, with a go.mod.
layout() {
	local pkg=$1 dir=$2 solved=$3
	mkdir -p "$dir"
	cp -R "$pkg/starter/." "$dir/"
	if [ "$solved" = 1 ]; then cp -R "$pkg/solution/." "$dir/"; fi
	[ -f "$dir/go.mod" ] || printf 'module lesson\n\ngo %s\n' "$goversion" > "$dir/go.mod"
}

normalize() { sed -e 's/\r$//' -e 's/[[:space:]]*$//' | awk '{lines[NR]=$0} END{n=NR; while (n>0 && lines[n]=="") n--; for (i=1;i<=n;i++) print lines[i]}'; }

check() {
	local pkg=$1 name type words min max tmp
	name=${pkg#"$ROOT"/}
	checked=$((checked + 1))
	type=$(field "$pkg" type)
	case "$type" in code|code_tests|quiz) ;; *) say "$name" "type is \"$type\""; return ;; esac
	words=$(body "$pkg" | wc -w | tr -d ' ')
	if [ "$type" = quiz ]; then min=60; max=120; else min=60; max=200; fi
	if [ "$words" -lt $min ] || [ "$words" -gt $max ]; then say "$name" "readme has $words words, not $min-$max"; fi
	body "$pkg" | grep -q '^## Your task$' || say "$name" "readme has no \"## Your task\" heading"
	body "$pkg" | grep -q '^# ' && say "$name" "readme has an H1 heading"
	if body "$pkg" | grep -qiE '!\[|<img|<video|<iframe'; then say "$name" "readme has a picture or an embed"; fi
	local long
	long=$(body "$pkg" | expand -t 4 | awk '/^```/{c=!c; next} c && length($0) > 50 {n++} END{print n+0}')
	[ "$long" -eq 0 ] || say "$name" "readme has $long code lines over 50 characters (the app's readme column wraps them)"
	if [ "$type" = quiz ]; then
		[ -f "$pkg/quiz.md" ] || say "$name" "has no quiz.md"
		return
	fi
	[ -f "$pkg/expected_output.txt" ] && [ "$type" = code_tests ] && say "$name" "code_tests lesson has expected_output.txt"
	[ -f "$pkg/expected_output.txt" ] || [ "$type" = code_tests ] || { say "$name" "code lesson has no expected_output.txt"; return; }
	local hints sentences
	hints=$(paragraphs "$pkg/hints.md" 2>/dev/null || echo 0)
	{ [ "$hints" -ge 2 ] && [ "$hints" -le 3 ]; } || say "$name" "has $hints hints, not 2 or 3"
	sentences=$(tr '\n' ' ' < "$pkg/solution_note.md" | grep -oE '[.!?](\s|$)' | wc -l | tr -d ' ')
	{ [ "$sentences" -ge 1 ] && [ "$sentences" -le 3 ]; } || say "$name" "solution note has $sentences sentences, not 1-3"
	local unformatted
	unformatted=$(gofmt -l "$pkg/starter" "$pkg/solution" 2>&1)
	[ -z "$unformatted" ] || say "$name" "not gofmt'd: $(echo "$unformatted" | tr '\n' ' ')"

	tmp=$(mktemp -d)
	trap 'rm -rf "$tmp"' RETURN
	layout "$pkg" "$tmp/solved" 1
	layout "$pkg" "$tmp/starter" 0
	local out status
	if [ "$type" = code_tests ]; then
		out=$(cd "$tmp/solved" && go test -count=$runs -timeout 120s . 2>&1); status=$?
		[ $status -eq 0 ] || say "$name" "solution fails go test -count=$runs: $(echo "$out" | grep -m3 -E 'FAIL|panic|error|\.go:[0-9]+' | tr '\n' ' ')"
		if [ $racecount -gt 0 ] && [ $status -eq 0 ]; then
			out=$(cd "$tmp/solved" && CGO_ENABLED=1 go test -race -count=$racecount -timeout 120s . 2>&1) || say "$name" "solution fails under -race: $(echo "$out" | grep -m3 -E 'FAIL|DATA RACE|panic' | tr '\n' ' ')"
		fi
		out=$(cd "$tmp/starter" && go vet . 2>&1) || {
			grep -q 'compiler-errors' <<<"$(field "$pkg" topics)" || say "$name" "starter does not build or vet: $(echo "$out" | head -3 | tr '\n' ' ')"
			return
		}
		out=$(cd "$tmp/starter" && go test -count=1 -timeout 20s . 2>&1) && say "$name" "starter passes its tests"
		if grep -q '^panic:' <<<"$out"; then say "$name" "starter panics: $(echo "$out" | grep -m1 '^panic:')"; fi
		if grep -q 'test timed out' <<<"$out"; then say "$name" "starter's tests time out"; fi
	else
		out=$(cd "$tmp/solved" && go vet . 2>&1 && go build -o "$tmp/solved.bin" . 2>&1) || { say "$name" "solution does not build: $(echo "$out" | head -3 | tr '\n' ' ')"; return; }
		local i
		for ((i = 0; i < runs; i++)); do
			if ! (cd "$tmp/solved" && "$tmp/solved.bin" > "$tmp/out.txt" 2>/dev/null); then say "$name" "solution exits nonzero on run $((i + 1))"; break; fi
			cmp -s "$tmp/out.txt" "$pkg/expected_output.txt" || { say "$name" "solution output differs from expected_output.txt on run $((i + 1)): $(diff "$tmp/out.txt" "$pkg/expected_output.txt" | head -4 | tr '\n' ' ')"; break; }
		done
		out=$(cd "$tmp/starter" && go vet . 2>&1 && go build -o "$tmp/starter.bin" . 2>&1) || {
			grep -q 'compiler-errors' <<<"$(field "$pkg" topics)" || say "$name" "starter does not build or vet: $(echo "$out" | head -3 | tr '\n' ' ')"
			return
		}
		if (cd "$tmp/starter" && "$tmp/starter.bin" > "$tmp/start.txt" 2>"$tmp/start.err"); then
			[ "$(normalize < "$tmp/start.txt")" = "$(normalize < "$pkg/expected_output.txt")" ] && say "$name" "starter already prints the expected output"
		fi
		grep -q '^panic:' "$tmp/start.err" && say "$name" "starter panics"
	fi
}

for arg in "$@"; do
	ROOT=$(cd "$arg" && pwd)
	if [ -f "$ROOT/README.md" ]; then ROOT=$(dirname "$ROOT"); check "$(cd "$arg" && pwd)"; continue; fi
	while IFS= read -r readme; do check "$(dirname "$readme")"; done < <(find "$ROOT" -name README.md -not -path '*/starter/*' -not -path '*/solution/*' | sort)
done
if [ $problems -eq 0 ]; then
	echo "ok: $checked packages"
else
	echo "$problems problems in $checked packages"
	exit 1
fi
