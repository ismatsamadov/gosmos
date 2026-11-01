-- The Go reference: one SQLite file of documentation for the Go that ships inside Gosmos,
-- reference/go.sqlite in this repository, so that looking something up never needs the
-- network. It holds the standard library (every public package `go list std` lists, with its
-- symbols, their declarations and docs, and the examples from its tests) and five articles:
-- the language specification and the memory model from that same Go, and Effective Go, the FAQ
-- and How to Write Go Code from golang/website. cmd/reference makes it (make reference); the app
-- only reads it. docs/reference-format.md says how it is made.
--
-- Conventions, as in schema.sql:
--   * Every table is STRICT. Booleans are INTEGER 0 or 1. Timestamps are RFC 3339 in UTC with a
--     trailing Z. An optional value that is missing is NULL, not ''. Lists keep their order in a
--     `position` column counted from 0.
--   * Every connection sets PRAGMA foreign_keys = ON. The file uses the rollback journal.
--   * Docs and article text are Markdown of the kind the app renders: no pictures, no raw HTML.
--     Code blocks are fenced, ```go for Go and a plain fence for anything else, such as the
--     spec's grammar.
--   * Links are of two kinds only, and none leads out of the app; any other link in the sources
--     is left as its text:
--       godoc:<import path>               a package: godoc:strings
--       godoc:<import path>.<Name>        a const, var, func or type: godoc:strings.Builder
--       godoc:<import path>.<Type>.<Name> a method, or a field or interface method of the type,
--                                         which have no symbol row of their own: godoc:io.Reader.Read
--       godoc-article:<slug>              an article: godoc-article:spec
--       godoc-article:<slug>#<anchor>     a section of one: godoc-article:spec#Struct_types
--   * cmd/reference makes the whole file anew each time from this schema; nothing updates it.

PRAGMA foreign_keys = ON;

------------------------------------------------------------------------------------------------
-- Where it came from
------------------------------------------------------------------------------------------------

-- What the file was made from, and the credit the app shows for it: exactly one row.
CREATE TABLE source (
  singleton      INTEGER PRIMARY KEY CHECK (singleton = 1),
  go_version     TEXT NOT NULL CHECK (go_version GLOB 'go[0-9]*'),  -- such as go1.27.3, from $GOROOT/VERSION
  platform       TEXT NOT NULL CHECK (platform GLOB '?*/?*'),       -- the build constraints applied: darwin/arm64
  website_commit TEXT NOT NULL CHECK (length(website_commit) = 40 AND website_commit NOT GLOB '*[^0-9a-f]*'),
  generated_at   TEXT NOT NULL CHECK (generated_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(generated_at) IS NOT NULL),
  credit         TEXT NOT NULL CHECK (credit <> ''),   -- Markdown: whose work each part is, and its license
  license        TEXT NOT NULL CHECK (license <> '')   -- Go's LICENSE file (BSD-3-Clause), as it is
) STRICT;

------------------------------------------------------------------------------------------------
-- The standard library
------------------------------------------------------------------------------------------------

-- A package, as go/doc documents it for the platform in source: exported names only.
CREATE TABLE package (
  id          INTEGER PRIMARY KEY,
  import_path TEXT NOT NULL UNIQUE CHECK (import_path <> ''),  -- such as math/rand/v2
  name        TEXT NOT NULL CHECK (name <> ''),                -- such as rand
  synopsis    TEXT CHECK (synopsis <> ''),     -- the doc's first sentence, plain text
  doc         TEXT CHECK (doc <> ''),          -- Markdown
  deprecated  INTEGER NOT NULL CHECK (deprecated IN (0, 1))  -- its doc has a "Deprecated: " paragraph
) STRICT;

-- A declaration as go doc shows it: one func, method or type, or one whole const or var block
-- (a single const or var is a block of one). Its symbols share it, so a block is kept once.
CREATE TABLE declaration (
  id         INTEGER PRIMARY KEY,
  package_id INTEGER NOT NULL REFERENCES package(id),
  kind       TEXT NOT NULL CHECK (kind IN ('const', 'var', 'func', 'type', 'method')),
  source     TEXT NOT NULL CHECK (source <> ''),  -- Go, gofmt-style: no doc comment, no function body
  doc        TEXT CHECK (doc <> ''),              -- Markdown
  UNIQUE (id, package_id, kind)                   -- what symbol references
) STRICT;

-- An exported name, in the order go doc lists them: the package's consts, vars and funcs, then
-- each type followed by the consts, vars and funcs grouped with it and its methods. A type's
-- funcs, such as strings.NewReader, are funcs with no receiver that come after it.
CREATE TABLE symbol (
  id             INTEGER PRIMARY KEY,
  package_id     INTEGER NOT NULL REFERENCES package(id),
  declaration_id INTEGER NOT NULL,
  kind           TEXT NOT NULL CHECK (kind IN ('const', 'var', 'func', 'type', 'method')),
  receiver       TEXT CHECK (receiver <> ''),  -- a method's type, without * or type parameters
  name           TEXT NOT NULL CHECK (name <> ''),
  synopsis       TEXT CHECK (synopsis <> ''),  -- the first sentence of its doc, or of its block's
  deprecated     INTEGER NOT NULL CHECK (deprecated IN (0, 1)),
  position       INTEGER NOT NULL CHECK (position >= 0),
  CHECK ((kind = 'method') = (receiver IS NOT NULL)),
  UNIQUE (package_id, receiver, name),
  UNIQUE (package_id, position),
  UNIQUE (id, package_id),                     -- what example references
  FOREIGN KEY (declaration_id, package_id, kind) REFERENCES declaration(id, package_id, kind)
) STRICT;
-- UNIQUE above counts every NULL receiver as different; this index covers the names that have none.
CREATE UNIQUE INDEX symbol_top_level ON symbol(package_id, name) WHERE receiver IS NULL;
CREATE INDEX symbol_declaration ON symbol(declaration_id);

-- An example from the package's tests: of the package itself, or of one of its symbols.
CREATE TABLE example (
  id         INTEGER PRIMARY KEY,
  package_id INTEGER NOT NULL REFERENCES package(id),
  symbol_id  INTEGER,                                   -- NULL: an example of the package
  suffix     TEXT CHECK (suffix GLOB '[a-z]*'),         -- ExampleBuilder_grow's grow; NULL if none
  position   INTEGER NOT NULL CHECK (position >= 0),    -- order within the package
  doc        TEXT CHECK (doc <> ''),                    -- Markdown
  code       TEXT NOT NULL,      -- the function's body, or the whole test file when whole_file is 1
  output     TEXT,               -- what it must print; NULL when it has no output comment and is
                                 -- only compiled, '' when it must print nothing
  unordered  INTEGER NOT NULL CHECK (unordered IN (0, 1)),   -- "Unordered output:": lines in any order
  whole_file INTEGER NOT NULL CHECK (whole_file IN (0, 1)),  -- code is a whole file, package clause and all
  CHECK (unordered = 0 OR output IS NOT NULL),
  UNIQUE (package_id, position),
  FOREIGN KEY (symbol_id, package_id) REFERENCES symbol(id, package_id)
) STRICT;
CREATE INDEX example_symbol ON example(symbol_id);

------------------------------------------------------------------------------------------------
-- Articles
------------------------------------------------------------------------------------------------

-- A document from go.dev: spec, mem, effective-go, faq and code.
CREATE TABLE article (
  slug     TEXT PRIMARY KEY CHECK (slug <> '' AND slug NOT GLOB '*[^a-z0-9-]*'),
  title    TEXT NOT NULL CHECK (title <> ''),
  subtitle TEXT CHECK (subtitle <> ''),             -- such as the spec's language version
  lead     TEXT CHECK (lead <> ''),                 -- Markdown before the first heading
  position INTEGER NOT NULL UNIQUE CHECK (position >= 0)
) STRICT;

-- A heading of an article and what follows it up to the next heading, whatever its level.
CREATE TABLE article_section (
  article_slug TEXT NOT NULL REFERENCES article(slug),
  position     INTEGER NOT NULL CHECK (position >= 0),
  anchor       TEXT NOT NULL CHECK (anchor <> ''),        -- the heading's id on go.dev
  level        INTEGER NOT NULL CHECK (level BETWEEN 1 AND 6),
  title        TEXT NOT NULL CHECK (title <> ''),         -- plain text
  body         TEXT CHECK (body <> ''),                   -- Markdown; NULL when the next heading follows at once
  PRIMARY KEY (article_slug, position),
  UNIQUE (article_slug, anchor)
) STRICT;

------------------------------------------------------------------------------------------------
-- Search
------------------------------------------------------------------------------------------------

-- What search finds: a package by its import path and synopsis, a symbol by its qualified name
-- (crypto/sha256.New, strings.Builder.Grow) and synopsis, a section by its title and text. One
-- row per package, symbol and section.
CREATE TABLE search_entry (
  id               INTEGER PRIMARY KEY,
  kind             TEXT NOT NULL CHECK (kind IN ('package', 'symbol', 'section')),
  package_id       INTEGER UNIQUE REFERENCES package(id),
  symbol_id        INTEGER UNIQUE REFERENCES symbol(id),
  article_slug     TEXT,
  section_position INTEGER,
  title            TEXT NOT NULL CHECK (title <> ''),
  body             TEXT CHECK (body <> ''),               -- plain text
  CHECK ((kind = 'package') = (package_id IS NOT NULL)),
  CHECK ((kind = 'symbol') = (symbol_id IS NOT NULL)),
  CHECK ((kind = 'section') = (article_slug IS NOT NULL)),
  CHECK ((article_slug IS NULL) = (section_position IS NULL)),
  UNIQUE (article_slug, section_position),
  FOREIGN KEY (article_slug, section_position) REFERENCES article_section(article_slug, position)
) STRICT;

-- The full-text index over search_entry, kept in step by the triggers below.
CREATE VIRTUAL TABLE search USING fts5(
  title, body, content = 'search_entry', content_rowid = 'id', tokenize = 'unicode61 remove_diacritics 2'
);
CREATE TRIGGER search_insert AFTER INSERT ON search_entry BEGIN
  INSERT INTO search(rowid, title, body) VALUES (new.id, new.title, new.body);
END;
CREATE TRIGGER search_delete AFTER DELETE ON search_entry BEGIN
  INSERT INTO search(search, rowid, title, body) VALUES ('delete', old.id, old.title, old.body);
END;
CREATE TRIGGER search_update AFTER UPDATE OF title, body ON search_entry BEGIN
  INSERT INTO search(search, rowid, title, body) VALUES ('delete', old.id, old.title, old.body);
  INSERT INTO search(rowid, title, body) VALUES (new.id, new.title, new.body);
END;

-- 'GSMR' marks the file as the Gosmos Go reference; user_version is its schema version.
PRAGMA application_id = 1196641618;
PRAGMA user_version = 1;
