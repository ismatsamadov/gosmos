## Every station's carrier losses must turn up in one search. Where does a station's code belong?

- [ ] In the message, as in `carrier lost at LAN`
  > Then each station has a message of its own, and a search for one message finds one station's losses.
- [x] In an attribute such as `site`, beside a message that never changes
  > The fixed message finds every carrier loss at once, and the attribute says where each one happened.
- [ ] In both the message and an attribute, to be safe
  > Any value in the message makes it differ from station to station, which is exactly what breaks the search.

## What does the last line write after `"site":"LAN"`?

- [ ] Nothing more: slog drops a value that has no key
  > slog keeps every argument; a value with no key before it is written under a key that flags it.
- [ ] `"margin_db":-1.5`, the key the radio uses for it
  > slog can't know which key you meant, and it never guesses one.
- [x] `"!BADKEY":-1.5`, since `-1.5` has no key
  > An argument left without a partner is written under `!BADKEY`, and `go vet`, which `go test` runs too, flags the call: it needs a key before the value, here `"margin_db"`.

## Without `-1.5`, the two calls make the same record. What differs between their lines?

- [x] Only the layout: `key=value` pairs, or one JSON object
  > Both lines hold the same level, message and site; a handler decides only how a record is written.
- [ ] The JSON line has no level, since JSON has no place for one
  > The JSON handler writes the level under `"level"`, as the text handler writes `level=`.
- [ ] The text line drops the site's key, since text has no keys
  > The text handler writes `site=LAN`: every attribute keeps its key in both layouts.
