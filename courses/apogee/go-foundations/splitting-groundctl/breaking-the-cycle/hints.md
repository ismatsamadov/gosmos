Follow the arrows: `alert` needs `telemetry` for `Frame`, and `telemetry` needs `alert` for `Severity`. Which of the two types could live in a package that imports neither?

Cut the `Severity` type, its constants and its `String` method from `alert/alert.go`, and paste them under the package line of `severity/severity.go`. They need no imports there.

In `telemetry/telemetry.go` and `alert/alert.go`, import `"apogee/groundctl/severity"` and put `severity.` in front of every `Severity`, `Info`, `Warning` and `Critical` they use; `telemetry` no longer needs `alert` at all.
