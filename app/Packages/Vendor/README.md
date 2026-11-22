# The app's vendored packages

The Swift packages the app builds with are kept here, so a clone of the repository builds the
app with nothing fetched. `app/Package.swift` names them by path. Each folder holds the
package's own sources, byte for byte as released, with its license, and a trimmed copy of its
own manifest: only the targets the app builds, a first comment saying what was left out, and
`-suppress-warnings`, since their warnings are not Gosmos's to fix.

| Folder | Package | Release | Commit | License |
|---|---|---|---|---|
| `GRDB.swift/` | [groue/GRDB.swift](https://github.com/groue/GRDB.swift) | 7.11.1 | `b83108d10f42680d78f23fe4d4d80fc88dab3212` | MIT, `LICENSE` |
| `swift-markdown/` | [swiftlang/swift-markdown](https://github.com/swiftlang/swift-markdown) | 0.9.0 | `25cb61d3482054b09ae76ca4f281b1bfe7fe5a43` | Apache 2.0 with Runtime Library Exception, `LICENSE.txt` and `NOTICE.txt` |
| `swift-cmark/` | [swiftlang/swift-cmark](https://github.com/swiftlang/swift-cmark) | 0.9.0 | `08ddb528923cc1a6527e02b7a1aee9e516ca749a` | BSD 2-clause and others, `COPYING` |

What each keeps:

- **GRDB.swift**: `GRDB/`, its sources, without `Documentation.docc`; `Sources/GRDBSQLite`,
  the system SQLite it builds on. Left out: the SQLCipher and custom-SQLite variants, tests,
  documentation and project files.
- **swift-markdown**: `Sources/Markdown` and `Sources/CAtomic`, on the swift-cmark beside it
  (`.package(path: "../swift-cmark")`) rather than fetched. Left out: tests, snippets, the
  command-line tool (`Tools/`), the DocC plugin and the CMake files.
- **swift-cmark**: `src/` and `extensions/`, the `cmark-gfm` and `cmark-gfm-extensions`
  libraries swift-markdown uses. Left out: the command-line program, tests, benchmarks and the
  build files for other systems.

To move to a new release:

1. Check out the release of each package and copy the folders above over these, leaving
   out what is listed as left out.
2. Bring each trimmed `Package.swift` up to date with the release's own: its tools version,
   platforms, settings and targets, keeping only the targets above.
3. Update this table and the versions in `docs/architecture.md`, then run `make test` and
   `make app`.
