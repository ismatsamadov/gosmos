# Cryptography in Go: syllabus

Season 5 of Apogee, July to October 2028. Cirrine Atmospherics signs a bigger contract on one
condition: nobody but Cirrine and Apogee may read its commands. Mirela leads the work, and the
learner learns each tool of Go's standard library by protecting something real, up to Rotation
Day, when every key in the network is replaced without missing a pass.

**Authoring status:** all eleven chapters are implemented as 88 lessons. Chapter 1 was read in
the app in light and dark before the rest was written; chapters 2–11 were written from this
plan two chapters at a time and every chapter independently reviewed. The world, the
characters and the writing rules are in `story.md`, whose section "Season 5: keeping secrets"
fixes every key, format, vector and message. The lessons in `courses/apogee.sqlite` are the
source of truth, and the larger changes from this plan are listed under "Changes in writing
and review" at the end.

- **Title:** Cryptography in Go
- **Slug:** `cryptography-in-go` · **Kind:** Course
- **Short description:** Protect Apogee's commands, files, passwords and links with Go's
  crypto packages, up to the day every key is replaced.
- **Description:** Cirrine Atmospherics has signed a bigger contract on one condition above
  all: nobody but Cirrine and Apogee may read its commands, not even someone with a dish
  listening to the uplink. With Mirela leading and reviewing, you learn each tool by protecting
  something real: numbers nobody can guess, fingerprints for delivered files, signed commands
  nobody can extend or replay, passwords stored so every guess at a stolen copy costs 600,000
  rounds, a key for every satellite from one master key, commands sealed for one reader,
  receipts anyone can check, a key agreed over an open link, Apogee's own certificate
  authority, and TLS with both ends proven. Every chapter shows the mistake first, then the way
  that refuses it. The season ends on Rotation Day: every key in the network replaced, at every
  station, without missing a single pass. Standard library only, every check on your own Mac.
- **Story:** an opening, each chapter's briefing and debrief, and the field guide, told as
  `docs/storytelling.md` says, following the bible's spine for the season.

## At a glance

| # | Chapter | Slug | Go topics | Lessons |
|---|---|---|---|---|
| 1 | Numbers Nobody Can Guess | `numbers-nobody-can-guess` | Make keys, nonces, salts and tokens with `crypto/rand` and `rand.Text`, measure their strength in bits, and see why `math/rand` never makes secrets. | 8 |
| 2 | Fingerprints for Files | `fingerprints-for-files` | Hash bytes and streams with `crypto/sha256` and `crypto/sha3`, name digests by their hash, and check files against a manifest. | 8 |
| 3 | Signed by the Ground | `signed-by-the-ground` | Use `crypto/hmac` in depth: why a hash of key and message isn't a MAC, unambiguous canonical forms, verify-then-act, `hmac.Equal`, `crypto/subtle` and one hash per message. | 8 |
| 4 | Passwords at Rest | `passwords-at-rest` | Store passwords with salts and `crypto/pbkdf2`, parse self-describing records, compare in constant time and upgrade old records at login. | 8 |
| 5 | Keys from Keys | `keys-from-keys` | Derive a key per satellite and per purpose with `crypto/hkdf`, keep secrets out of logs, and load a versioned key ring. | 7 |
| 6 | Sealed Messages | `sealed-messages` | Encrypt and authenticate with `crypto/aes` and `crypto/cipher`'s GCM: tags, nonces, associated data, a strict wire format and checks in the right order. | 8 |
| 7 | Signatures Anyone Can Check | `signatures-anyone-can-check` | Sign and verify with `crypto/ed25519`, publish public keys as PEM, sign exact bytes with `json.RawMessage`, and choose keys by ID. | 8 |
| 8 | Agreeing on a Secret | `agreeing-on-a-secret` | Agree on keys with `crypto/ecdh` X25519, authenticate the shares, derive session keys, keep forward secrecy, and encapsulate with `crypto/mlkem`. | 8 |
| 9 | Who Are You, Server? | `who-are-you-server` | Build a private certificate authority with `crypto/x509` and `encoding/pem`, issue certificates, and verify chains, names, usages and validity. | 8 |
| 10 | Every Connection Encrypted | `every-connection-encrypted` | Serve and dial with `crypto/tls`: your own roots, TLS 1.3 only, mutual TLS, authorization by certificate, and TLS on a raw socket. | 9 |
| 11 | Rotation Day | `rotation-day` | Rotate keys with versions and overlap, revoke at once with revocation lists and `VerifyConnection`, and run a whole day of passes through the change. | 8 |

88 lessons: 6 code, 68 code_tests and 14 quizzes (84.1% hands-on). Difficulty rises from 6 to
10. Every chapter shows the mistake it guards against before the right way: the starter holds
the mistake and the read-only test plays whoever exploits it. No test depends on random
output: tests check sizes, round trips and refusals, or fixed vectors from the bible. Every
server listens on 127.0.0.1 on a port the system picks; clocks are injected; no run comes near
the app's 10 seconds.

## Chapter 1 · Numbers Nobody Can Guess

`numbers-nobody-can-guess` · 8 lessons

Fill bytes with `crypto/rand.Read`, make text secrets with `rand.Text`, count bits of strength,
and see why seeded `math/rand/v2` and in-memory counters never make secrets or nonces.

#### 1.1 Predictable by the Clock
`predictable-by-the-clock` · code_tests · difficulty 6
- **Objective:** Fill a nonce from `crypto/rand.Read` instead of a `math/rand/v2` generator.
- **Scenario:** Monday 10 July: the draft nonce helper for Cirrine's sealed commands seeds `math/rand/v2` with the current second, and Mirela shows that knowing the second is enough to make the same nonce again.
- **Exercise:** `newNonce() []byte` builds 12 bytes from `rand.New(rand.NewPCG(uint64(time.Now().Unix()), 0))`, one `Uint32` per byte; you fill a 12-byte slice with `crypto/rand.Read` and drop the generator.
- **Check:** `TestAttackerPredicts`: the test notes the time around the call and rebuilds the nonce from every second from 2 s before to 2 s after with the starter's own recipe; no rebuild may match (the starter's always does). `TestLength`: 12 bytes. `TestTwoDiffer`: two calls differ. `TestNoMathRand`: `go/parser` finds no import of `math/rand` or `math/rand/v2` in `main.go`.

#### 1.2 Thirty-Two Random Bytes
`thirty-two-random-bytes` · code_tests · difficulty 6
- **Objective:** Make a 256-bit key with `crypto/rand` and read it back from 64 hex digits.
- **Scenario:** On Friday Mirela will make the season's master key, `k-2028-07`, and she wants the helper that makes it, and the parser that reads the key file, first.
- **Exercise:** `newKey(id string) Key` returns `Key{ID: id}` with a nil `Secret`, and `parseSecret(text string) ([]byte, error)` returns `hex.DecodeString`'s result unchecked; you fill 32 random bytes, and make `parseSecret` refuse anything but exactly 64 hex digits.
- **Check:** `newKey("k-2028-07")` keeps the ID and has 32 bytes; two keys differ; the test master's 64 digits (`75706c69…6573`) parse to its 32 bytes; 62 digits (`hex.DecodeString` alone accepts them), 66 digits, a `g` and an empty string → errors.

#### 1.3 Quick Check: Randomness
`quick-check-randomness` · quiz · difficulty 6
- **Objective:** Check which generators may make secrets, and what `crypto/rand.Read` promises.
- **Scenario:** Mirela's three questions before the practice key goes in the safe.
- **Exercise:** Three questions.
- **Check:** `math/rand/v2` started from a known seed gives the same numbers every run, so it never makes secrets; `crypto/rand.Read` never returns an error and always fills the slice (if the system can't supply randomness, the program stops); a 32-byte key has 256 bits, 2^256 possible values.

#### 1.4 Tokens People Can Type
`tokens-people-can-type` · code_tests · difficulty 6
- **Objective:** Make a secret text with `rand.Text`.
- **Scenario:** The ops console's session tokens are the login time in nanoseconds, written in decimal; Mirela wants ones nobody can guess, still short enough to paste.
- **Exercise:** `newSessionToken() string` returns `strconv.FormatInt(time.Now().UnixNano(), 10)`; you return `rand.Text()`.
- **Check:** at least 26 characters, every one from `A`–`Z` and `2`–`7`; two tokens differ; the test's attacker, who knows the time around the call, reads the starter's token as a time inside that window, and finds nothing of the kind in yours.

#### 1.5 How Many Bits
`how-many-bits` · code · difficulty 6
- **Objective:** Measure a secret's strength in bits: its length times log2 of its alphabet's size.
- **Scenario:** Ingrid asks whether the console's old six-digit codes were "basically as good" as the new tokens.
- **Exercise:** `bits(length, alphabet int) float64` returns `float64(length * alphabet)`; you return `float64(length) * math.Log2(float64(alphabet))`. `main` prints four schemes with `%.1f`.
- **Check:** `6 digits: 19.9 bits` / `4 bytes in hex: 32.0 bits` / `rand.Text: 130.0 bits` / `32-byte key: 256.0 bits` (the starter prints 60.0, 128.0, 832.0 and 8192.0). The readme's callout: Apogee's floor for any secret is 128 bits.

#### 1.6 A Counter Forgets
`a-counter-forgets` · code_tests · difficulty 7
- **Objective:** Choose random nonces when a counter can't be trusted to survive a restart.
- **Scenario:** Nkem's sketch numbers Cirrine's nonces 0, 1, 2…, and the count starts at 0 again whenever uplinkd restarts.
- **Exercise:** `type Noncer struct{ next uint64 }` and `(n *Noncer) Nonce() []byte`, which writes the counter big-endian into 12 bytes; you return 12 bytes from `crypto/rand` instead and remove the counter.
- **Check:** the first nonces of two `Noncer` values, two runs of uplinkd, differ (the starter's are both zero); 1,000 nonces from one value are all 12 bytes and all distinct.

#### 1.7 Quick Check: Nonces
`quick-check-nonces` · quiz · difficulty 6
- **Objective:** Check what a nonce must be and how far random ones go.
- **Scenario:** Mirela's two questions after the counter goes.
- **Exercise:** Two questions.
- **Check:** A nonce may be public but must never repeat under the same key, and random 12-byte nonces stay safe for 2^32 messages per key; a counter kept only in memory repeats its first values after a restart.

#### 1.8 Mirela's Random Kit
`mirelas-random-kit` · code_tests · difficulty 7
- **Objective:** Keep every secret-making helper in one file that uses only `crypto/rand`.
- **Scenario:** Friday 14 July, 16:00: Mirela makes `k-2028-07` with your helpers, and wants one file, `random.go`, that every later tool takes its secrets from.
- **Exercise:** `random.go` holds `newKey`, `newNonce`, `newSalt` and `newSessionToken`, returning nil or `""`; you write all four: 32, 12 and 16 bytes from `crypto/rand.Read`, and `rand.Text()`.
- **Check:** sizes 32, 12 and 16, a key keeping its ID, a token of 26 or more base32 characters; two calls of each differ; `TestOnlyCryptoRand` parses every `.go` file in the lesson and finds `crypto/rand` imported in `random.go` and `math/rand` and `math/rand/v2` nowhere.

## Chapter 2 · Fingerprints for Files

`fingerprints-for-files` · 8 lessons

Hash bytes with `sha256.Sum256`, streams through `hash.Hash` with `io.Copy` and
`io.MultiWriter`, pick SHA-256 or SHA3-256 by name, and check delivered files against a
manifest.

#### 2.1 A File's Fingerprint
`a-files-fingerprint` · code · difficulty 6
- **Objective:** Compute a SHA-256 digest with `sha256.Sum256` and print it in hex.
- **Scenario:** Monday 17 July: Tomás tries a fingerprint on the first words you ever sent up, `Apogee ground calling Pipit-1`, and on the same words for Pipit-2.
- **Exercise:** `main` prints each text's CRC-32 with `%08x`; you print the digest size and the first 8 bytes of each `sha256.Sum256` with `%x`.
- **Check:** `digest size: 32 bytes` / `Pipit-1: fe5e0f73525ffc9f` / `Pipit-2: 8cbc29efabfcd2fb`.

#### 2.2 Noise Is Not an Attacker
`noise-is-not-an-attacker` · code_tests · difficulty 6
- **Objective:** Detect deliberate changes with a SHA-256 digest, where a CRC-32 catches only noise.
- **Scenario:** Mirela changes "0.71" to "0.17" in a Furrowlight field report, adds four bytes, and the CRC-32 the delivery note carries still matches.
- **Exercise:** `fingerprint(file []byte) string` returns the CRC-32 as 8 hex digits; you return the 64 hex digits of `sha256.Sum256`.
- **Check:** the original report (bible fixture) → `72198855a7d55dd3…0aae`; Mirela's forgery (`…0a1a15bcd1`) → a different fingerprint (the starter gives both `f6ced6c0`); one byte changed anywhere in a 1 KiB file → a different fingerprint; an empty file → `e3b0c442…b855`.

#### 2.3 Hashing a Stream
`hashing-a-stream` · code_tests · difficulty 6
- **Objective:** Hash a stream of any length by copying it into a `hash.Hash` with `io.Copy`.
- **Scenario:** A Furrow-A image is 10.9 GB, and dlinkd's fingerprint must cover every byte of it without holding it in memory.
- **Exercise:** `fingerprintStream(r io.Reader) (string, error)` hashes only what one `Read` call returns; you `io.Copy` the whole stream into `sha256.New()` and return the hex of `Sum(nil)`, or `""` and the copy's error.
- **Check:** a 1 MiB stream delivered 4 KiB per `Read` → the same digest as `sha256.Sum256` of all of it (the starter's covers 4 KiB); an empty stream → `e3b0c442…`; a stream that fails after 64 KiB → `""` and that error.

#### 2.4 Quick Check: Hashes
`quick-check-hashes` · quiz · difficulty 6
- **Objective:** Check what a digest proves and what it doesn't.
- **Scenario:** Yusra's three questions before she promises Cirrine anything.
- **Exercise:** Three questions, about the forged field report shown in the readme.
- **Check:** Anyone can compute a SHA-256 digest, so it proves a file unchanged only if the digest itself arrives safely; four chosen bytes can make any file match a given CRC-32; a SHA-256 digest is 32 bytes for any file, empty or 10.9 GB.

#### 2.5 Digests That Name Their Hash
`digests-that-name-their-hash` · code_tests · difficulty 7
- **Objective:** Read digests written as `<hash>:<hex>` strictly, so two hashes are never confused.
- **Scenario:** Cirrine's archive uses SHA3-256 and Apogee's deliveries SHA-256, and a bare 64-digit string could be either.
- **Exercise:** `parseDigest(s string) (alg string, sum []byte, err error)` splits at the first `:` and decodes whatever follows; you accept only `sha256` and `sha3-256` and exactly 64 lowercase hex digits, with the bible's error texts.
- **Check:** `sha256:fe5e…bd66` → `sha256` and 32 bytes; `sha3-256:7f7e…84ea` → `sha3-256`; a bare digest, `SHA256:` in capitals, `md5:`, 63 digits and uppercase hex → errors; `sha256:` with 64 digits and a trailing space → error.

#### 2.6 Cirrine's Hash
`cirrines-hash` · code_tests · difficulty 7
- **Objective:** Choose a hash by name through the `hash.Hash` interface, with `sha3.New256`.
- **Scenario:** Cirrine checks every delivery against its archive, which has used SHA3-256 since its first satellite.
- **Exercise:** `newHash(alg string) (hash.Hash, error)` knows only `sha256`; you add `sha3-256`. `fingerprintAs(r io.Reader, alg string) (string, error)`, given, returns a named digest through it.
- **Check:** `fingerprintAs` of `Apogee ground calling Pipit-1` as `sha3-256` → `sha3-256:7f7e6c3e…84ea`, as `sha256` → `sha256:fe5e0f73…bd66`; an unknown name → an error before anything is read (the test's reader counts calls).

#### 2.7 Fingerprint While You Save
`fingerprint-while-you-save` · code_tests · difficulty 7
- **Objective:** Hash bytes as they are written, with `io.MultiWriter`.
- **Scenario:** dlinkd writes each image once into the delivery folder, and reading 10.9 GB back to fingerprint it would double the work.
- **Exercise:** `save(dst io.Writer, src io.Reader) (string, error)` copies, then returns `""`; you copy into `io.MultiWriter(dst, h)` and return `sha256:` and the hex digest, or `""` and the error.
- **Check:** the saved bytes equal the source and the digest equals `sha256:` of the source; the source is read once (counted); a destination that fails after 1 KiB → `""` and that error.

#### 2.8 The Delivery Manifest
`the-delivery-manifest` · code_tests · difficulty 7
- **Objective:** Check a folder of delivered files against a manifest of named digests and sizes.
- **Scenario:** Friday 21 July: Cirrine's first fingerprinted delivery, and Yusra wants every file's state in one report.
- **Exercise:** `checkManifest(fsys fs.FS, manifest io.Reader) ([]string, error)` returns nil; you read each line `<digest> <bytes> <file>`, check the size, then hash the file from `fsys` with the digest's own hash, and report `ok`, `wrong size`, `changed` or `missing` and the file, in manifest order.
- **Check:** an `fstest.MapFS` holding one whole SHA-256 file, one whole SHA3-256 file, one changed file and one of the wrong size, with a fifth listed but absent → five lines in order; a line with two fields → `manifest line 3: want <digest> <bytes> <file>`; an empty manifest → no lines and no error.

## Chapter 3 · Signed by the Ground

`signed-by-the-ground` · 8 lessons

Authenticate with `crypto/hmac` rather than a hash of key and message, make canonical forms
unambiguous, verify before changing state, compare with `hmac.Equal` and
`subtle.ConstantTimeCompare`, and give every message its own HMAC.

#### 3.1 A Homemade Signature
`a-homemade-signature` · code_tests · difficulty 7
- **Objective:** Authenticate a message with `hmac.New(sha256.New, key)`, never with `sha256(key + message)`.
- **Scenario:** Monday 24 July: Pipit-3's flight software checks commands with `sha256(key + command)`, and Mirela makes a valid tag for a command with `|REBOOT` on the end without knowing the key.
- **Exercise:** `tag(key []byte, canonical string) []byte` returns `sha256.Sum256` of the key followed by the command; you return its HMAC-SHA256. `valid(key []byte, canonical string, t []byte) bool`, given, recomputes and compares with `hmac.Equal`.
- **Check:** `TestLengthExtension` (the bible's recipe): from the tag of `pipit-3|212|SET_MODE|SAFE` and the key's length alone, the test builds a tag for that command plus SHA-256's padding plus `|REBOOT`; `valid` must refuse it (the starter accepts). `TestTag`: `tag(benchKey, "pipit-2|118|SET_MODE|SAFE")` → Season 1's `512545bc…3b4d`.

#### 3.2 One Bar Too Many
`one-bar-too-many` · code_tests · difficulty 7
- **Objective:** Make a canonical form unambiguous by refusing separators and empty fields inside it.
- **Scenario:** `SET_MODE` with the one argument `SAFE|REBOOT` signs exactly like `SET_MODE` with two arguments, `SAFE` and `REBOOT`.
- **Exercise:** `canonical(c Command) (string, error)` joins the fields with `|`; you refuse fields as the bible's canonical rules say, with an error naming the field.
- **Check:** `pipit-2|118|SET_MODE|SAFE` unchanged; an argument `SAFE|REBOOT` → error, so the two commands can no longer share a form; an empty argument, `SAFE REBOOT` with a space, a tab, `Seq: 0`, the satellite `Pipit-2` and the operation `set_mode` → errors.

#### 3.3 Quick Check: Canonical Forms
`quick-check-canonical-forms` · quiz · difficulty 6
- **Objective:** Check why signed bytes must have exactly one meaning.
- **Scenario:** Nkem's two questions on Mirela's review.
- **Exercise:** Two questions.
- **Check:** Two commands with the same canonical text share every signature, so a signature for one works for the other; refusing `|` inside fields (or writing each field's length before it) gives every text one meaning.

#### 3.4 Verify Before You Count
`verify-before-you-count` · code_tests · difficulty 7
- **Objective:** Change state only after a message's signature checks out.
- **Scenario:** The replay guard records a command's sequence number before checking its signature, so one forged command numbered 18,446,744,073,709,551,615 makes every real command afterwards look like a replay.
- **Exercise:** `(g *Guard) Accept(c Command, key []byte) error` records `c.Seq` when it is above the last, then verifies; you verify first (`canonical`, then `verify` from Season 1), then compare with the last, then record.
- **Check:** the forged `math.MaxUint64` command → `ErrForged`, and a signed `Seq 119` after it is still accepted (the starter refuses it as a replay); `118` twice → `ErrReplay` the second time; `119` then `120` → both accepted; a zero `Guard` works; errors are checked with `errors.Is`.

#### 3.5 Same Length, Same Time
`same-length-same-time` · code_tests · difficulty 7
- **Objective:** Compare secrets of any length in constant time by comparing fixed-size digests with `subtle.ConstantTimeCompare`.
- **Scenario:** The ops console checks the night shift's maintenance phrase with `==`, which stops at the first wrong character; `ConstantTimeCompare` stops at once when lengths differ.
- **Exercise:** `samePhrase(typed, stored string) bool` uses `==`; you hash both with `sha256.Sum256` and compare the digests with `subtle.ConstantTimeCompare`.
- **Check:** equal phrases → true; one wrong character, a prefix, a longer phrase and an empty phrase → false; `TestNoEquals` parses `main.go`: `samePhrase` calls `subtle.ConstantTimeCompare` and has no `==` or `!=`.

#### 3.6 Quick Check: Comparing Secrets
`quick-check-comparing-secrets` · quiz · difficulty 7
- **Objective:** Check when a comparison must take constant time.
- **Scenario:** Mirela's three questions.
- **Exercise:** Three questions.
- **Check:** An HMAC check needs `hmac.Equal`, because how fast a wrong tag is refused tells an attacker how much of it was right; comparing two public file digests needs no constant time; `subtle.ConstantTimeCompare` answers at once for different lengths, so variable-length secrets are compared through their digests.

#### 3.7 One Signer, One Message
`one-signer-one-message` · code_tests · difficulty 7
- **Objective:** Start a new HMAC for every message, since a `hash.Hash` keeps everything written to it.
- **Scenario:** uplinkd's `Signer` keeps one HMAC for speed, and the second command of the night was signed together with the first.
- **Exercise:** `type Signer struct{ key []byte; mac hash.Hash }`; `NewSigner(key)` makes one HMAC and `(s *Signer) Sign(canonical string) string` writes into it and sums; you make `Sign` start its own `hmac.New(sha256.New, s.key)` and remove `mac`.
- **Check:** the second of two signatures equals a fresh HMAC of the second command alone (the starter's covers both; the test stops there with `Fatalf`); then 8 goroutines × 50 commands sign at once and every signature verifies.

#### 3.8 Mirela's Planted Lines
`mirelas-planted-lines` · code_tests · difficulty 8
- **Objective:** Check a signed command line from end to end, failing closed on every bad line.
- **Scenario:** Friday 28 July: Mirela feeds twelve lines she wrote herself to the uplink's checker; it must accept the two good ones and refuse the other ten.
- **Exercise:** `checkLine(line string, key []byte, g *Guard) (Command, error)` refuses everything; you parse Season 1's line `CMD pipit-2 118 SET_MODE SAFE sig=<64 hex>` with `strings.Fields`, build it with `canonical`, verify, and pass it to the guard, in the bible's order.
- **Check:** lines 1 and 2 accepted, in order; refused: a changed argument, a changed sequence number, a bar inside an argument, no `sig=`, 63 hex digits, a signature by another key, `sig=` not last, a forged maximum sequence number, a replay of line 1, an empty line; after every refusal the guard is unchanged.

## Chapter 4 · Passwords at Rest

`passwords-at-rest` · 8 lessons

Salt passwords, slow every guess with `pbkdf2.Key`, store self-describing records, compare with
`subtle.ConstantTimeCompare`, upgrade old records at login, and refuse unknown users and wrong
passwords alike.

#### 4.1 One Table Cracks Them All
`one-table-cracks-them-all` · code_tests · difficulty 7
- **Objective:** Salt every stored password with 16 random bytes, so no precomputed table matches it.
- **Scenario:** Monday 31 July: the Fishbowl's spare laptop is gone with a copy of the ops console's user file, and the test plays whoever has it, with a table of 1,000 common passwords and their SHA-256 digests.
- **Exercise:** `hashPassword(pw string) string` returns `sha256$<hex>`; you return the lesson's step `salted$<32 hex salt>$<64 hex of sha256(salt + password)>` with a fresh salt from `newSalt`.
- **Check:** none of the test accounts' records is found in the table (the starter's `kiosk` record, `fishbowl`, is found at once); the same password hashed twice gives two different records (the starter's are equal); each salt is 16 bytes and its digest recomputes.

#### 4.2 Make Every Guess Cost
`make-every-guess-cost` · code_tests · difficulty 7
- **Objective:** Derive a password's stored key with `pbkdf2.Key`, SHA-256 and 600,000 iterations.
- **Scenario:** A salt stops the table, not the guessing: one laptop tries a billion salted SHA-256 guesses a second.
- **Exercise:** `derive(pw string, salt []byte) ([]byte, error)` returns one SHA-256 of salt and password; you return `pbkdf2.Key(sha256.New, pw, salt, iterations, 32)` with the constant `iterations = 600000`.
- **Check:** Mirela's password and the fixture salt → `6d91ceab…d938`, 32 bytes; another salt → another key; `iterations` is 600000; at most two derivations in the test.

#### 4.3 Quick Check: Salt and Cost
`quick-check-salt-and-cost` · quiz · difficulty 7
- **Objective:** Check what salts and iteration counts each defend against.
- **Scenario:** Ingrid's three questions for the incident log.
- **Exercise:** Three questions.
- **Check:** A salt makes equal passwords store differently and defeats precomputed tables, but doesn't slow a single guess; 600,000 iterations make every guess 600,000 times the work for anyone holding the file; the salt isn't secret and is stored beside the key.

#### 4.4 The Record Says How
`the-record-says-how` · code_tests · difficulty 7
- **Objective:** Store the algorithm, count and salt with every password key, and parse such records strictly.
- **Scenario:** Next year's count will be higher, and today's records must still check out then.
- **Exercise:** `formatRecord(r Record) string` and `parseRecord(s string) (Record, error)` are stubs; you write `pbkdf2-sha256$<count>$<32 hex>$<64 hex>` and read it back as the bible's record rules say.
- **Check:** a round trip; Mirela's record → 600000, a 16-byte salt and a 32-byte key; `sha256$…`, a count of 0, of 999 and of 10,000,001, a 30-digit salt, a 62-digit key, a missing field and a fifth field → errors.

#### 4.5 Checking a Login
`checking-a-login` · code_tests · difficulty 8
- **Objective:** Check a password by deriving again with the record's own count and salt, comparing in constant time.
- **Scenario:** Mirela logs in with her new password, and the record she made on Monday must accept it.
- **Exercise:** `checkPassword(r Record, pw string) (bool, error)` returns `false, nil`; you derive with `pbkdf2.Key` and `r.Iterations` and `r.Salt`, and compare with `subtle.ConstantTimeCompare`.
- **Check:** Mirela's 600,000 record accepts `fails closed, always` and refuses `fails closed, always.` and `""`; the `bench` record at 1,000 iterations (`5b893367…74b3`) accepts the same password; `TestConstantTime` parses `main.go`: `subtle.ConstantTimeCompare` or `hmac.Equal`, and no `bytes.Equal` or `==` on the keys.

#### 4.6 Upgrade at Login
`upgrade-at-login` · code_tests · difficulty 8
- **Objective:** Replace a weak record with a new one at the moment its password is proven.
- **Scenario:** Twelve engineers are on leave until after Friday; their old records must improve the next time they log in, and not before.
- **Exercise:** `login(stored, pw string) (ok bool, upgraded string, err error)` handles only `pbkdf2-sha256` records; you also check old `sha256$` records (in constant time), and return a new 600,000-iteration record whenever the password matched an old record or one under 600,000 iterations.
- **Check:** `kiosk`'s old record with `fishbowl` → ok and a new record that itself accepts `fishbowl`; a wrong password → not ok and no new record; Mirela's current record → ok and `""`; the `bench` 1,000-iteration record → ok and upgraded; a record that doesn't parse → an error.

#### 4.7 The Same Answer for Everyone
`the-same-answer-for-everyone` · code_tests · difficulty 8
- **Objective:** Refuse an unknown user and a wrong password with the same error after the same work.
- **Scenario:** The console answered "no such user" at once and "wrong password" a moment later, so anyone could list its 26 accounts.
- **Exercise:** `(c *Console) Check(user, pw string) error` returns `no such user` before deriving anything; you return `errInvalidLogin` for both, checking an unknown user against the fixed dummy record so the work is the same.
- **Check:** an unknown user and a wrong password → the same error value (`errors.Is`, message `invalid user or password`); each costs exactly one `kdf` call (the test's `kdf` counts); the right password → nil.

#### 4.8 The Console Signs In
`the-console-signs-in` · code_tests · difficulty 8
- **Objective:** Combine strict records, constant-time checks, upgrades and random session tokens in one login.
- **Scenario:** Friday 4 August: every engineer has a new password, and the console's login opens again at 16:00.
- **Exercise:** `(c *Console) Login(user, pw string) (token string, err error)` is a stub; you check through `Check`, save an upgraded record through `c.save`, then issue `newSessionToken()` and keep only its SHA-256 in hex in `c.sessions`.
- **Check:** a good login returns a token of 26 or more characters whose SHA-256 is in `sessions` and which is itself nowhere in it; `kiosk`'s old record is saved upgraded exactly once; a failed login returns `errInvalidLogin`, saves nothing and adds no session; two logins get different tokens.

## Chapter 5 · Keys from Keys

`keys-from-keys` · 7 lessons

Derive a key per satellite and per purpose with `hkdf.Key`, split HKDF into `Extract` and
`Expand`, hide keys from `fmt` with `String` and `GoString`, and load versioned master keys
strictly.

#### 5.1 One Key Opens Everything
`one-key-opens-everything` · code_tests · difficulty 7
- **Objective:** Derive a separate key for each satellite from one master key with `hkdf.Key`.
- **Scenario:** Monday 7 August: the key read out of the lab bench's spare Pipit-2 flight computer signs a command that Cirrine-1 accepts.
- **Exercise:** `signKey(master Key, satID string) ([]byte, error)` returns `master.Secret` for every satellite; you return `hkdf.Key(sha256.New, master.Secret, salt, "sign|"+satID, 32)` with the given `salt`.
- **Check:** with the test master, pipit-2 → `82f94563…da7c` and cirrine-1 → `18668bfe…0213`, 32 bytes each; a command signed with pipit-2's key fails `verify` under cirrine-1's (the starter's passes); an empty satellite ID → error.

#### 5.2 Salt and Info
`salt-and-info` · code · difficulty 7
- **Objective:** See that HKDF's info string decides the key, so one master gives unrelated keys.
- **Scenario:** Tomás wants the fleet's new keys side by side, the first four bytes of each.
- **Exercise:** `main` derives every line with the same info; you pass each line's label as the info.
- **Check:** `sign|pipit-2   82f94563` / `sign|furrow-a  00a1988a` / `sign|cirrine-1 18668bfe` / `seal|cirrine-1 427a8ae7`.

#### 5.3 Quick Check: Derived Keys
`quick-check-derived-keys` · quiz · difficulty 7
- **Objective:** Check what HKDF is for and what it is not.
- **Scenario:** Nkem's three questions before the satellites get their keys.
- **Exercise:** Three questions.
- **Check:** One derived key reveals neither the master nor another satellite's key; HKDF is for keys that are already random, while a password needs PBKDF2's slow work; changing the salt or the info changes every derived key.

#### 5.4 One Key, One Job
`one-key-one-job` · code_tests · difficulty 7
- **Objective:** Derive a different key for each purpose, so no key serves two algorithms.
- **Scenario:** Cirrine-1 will hold one key that checks signed commands and another that opens sealed ones.
- **Exercise:** `satKey(master Key, purpose, satID string) ([]byte, error)` always derives `sign|<sat>`; you derive `<purpose>|<sat>`, refuse purposes other than `sign` and `seal`, and refuse satellites outside the fleet.
- **Check:** cirrine-1's `sign` → `18668bfe…` and `seal` → `427a8ae7…0171`, different; `wrap` → error; `thistle` and `""` → errors.

#### 5.5 Keys Never Print
`keys-never-print` · code_tests · difficulty 8
- **Objective:** Keep secrets out of logs with `String` and `GoString` methods.
- **Scenario:** Last night's uplinkd log holds all 32 bytes of `k-2028-07`, printed with `%v`.
- **Exercise:** `Key` and `KeyRing` have no methods; you add `Key.String`, `Key.GoString` and `(*KeyRing).String` as the bible writes them; the ring's needs its own, because fmt can't call methods on the keys in its unexported field.
- **Check:** `%v`, `%+v`, `%s` and `%#v` of a `Key`, a `[]Key`, a struct holding a `Key` and a `*KeyRing` never show the secret (as text, decimal bytes or hex) and do show the ID; `TestRing` fails until the ring has its own `String`, since a ring without one prints the bytes even when `Key` has its methods.

#### 5.6 Extract Once, Expand Many
`extract-once-expand-many` · code_tests · difficulty 8
- **Objective:** Split HKDF into one `hkdf.Extract` and an `hkdf.Expand` per key.
- **Scenario:** Every satellite needs its two keys at once: fourteen derivations from one master.
- **Exercise:** `fleetKeys(master Key, satIDs []string) (map[string][]byte, error)` returns an empty map; you extract once with the salt and expand `sign|<sat>` and `seal|<sat>` for each satellite.
- **Check:** the seven satellites → 14 keys, each equal to `hkdf.Key` with the same salt and info; an empty list → an empty map; `TestExtractOnce` parses `main.go`: one `hkdf.Extract`, and no `hkdf.Key`.

#### 5.7 The Key Ring
`the-key-ring` · code_tests · difficulty 8
- **Objective:** Load versioned master keys from a strict file and derive any satellite's key by key ID.
- **Scenario:** Friday 11 August, 00:00 UTC: `k-2027-03` retires, and every satellite answers only to its own keys.
- **Exercise:** `loadRing(r io.Reader) (*KeyRing, error)` and `(r *KeyRing) SatKey(keyID, purpose, satID string) ([]byte, error)` are stubs; you read `uplink-keys.txt` as the bible says and derive with the named version.
- **Check:** August's file → one version, `k-2028-07`, activating 2028-08-11T00:00:00Z; `SatKey("k-2028-07", "seal", "cirrine-1")` → `427a8ae7…`; `k-2027-03` → `ErrUnknownKey`; comments and blank lines skipped; a duplicate ID, a 63-digit secret and a time `2028-13-01T00:00:00Z` → errors naming the line.

## Chapter 6 · Sealed Messages

`sealed-messages` · 8 lessons

Encrypt and authenticate with `aes.NewCipher` and `cipher.NewGCM`, give every message a fresh
nonce, bind a readable header as associated data, and write, parse and open the sealed command
line in a fixed order.

#### 6.1 Secret but Not Safe
`secret-but-not-safe` · code_tests · difficulty 8
- **Objective:** Encrypt with AES-GCM, which refuses any changed ciphertext, instead of counter mode alone.
- **Scenario:** Monday 14 August: Mirela's first draft encrypts with AES in counter mode, and without the key she turns `SET_MODE|SAFE` into `SET_MODE|LEOP` by flipping bits.
- **Exercise:** `seal(key, nonce, plaintext []byte) []byte` and `open(key, nonce, sealed []byte) ([]byte, error)` use `cipher.NewCTR` with the 12-byte nonce padded to 16; you use `cipher.NewGCM`'s `Seal` and `Open`, no associated data yet.
- **Check:** a round trip; the test's bit flip (it XORs `SAFE` with `LEOP` at bytes 9 to 12) → `open` returns an error (the starter returns `SET_MODE|LEOP`); one changed byte in the tag, a cut-off tag and another key → errors; the sealed text is 16 bytes longer than the plaintext.

#### 6.2 What the Tag Adds
`what-the-tag-adds` · code · difficulty 7
- **Objective:** See GCM's sizes and its refusal: a 12-byte nonce, a 16-byte tag, an error for a changed byte.
- **Scenario:** Yusra wants to tell Cirrine how much longer a sealed command is than a plain one.
- **Exercise:** `main` seals `SET_MODE|SAFE` with Cirrine-1's test seal key and the fixed nonce and prints zeros; you print the lengths, `NonceSize()`, `Overhead()` and the error from opening a copy with its last byte changed.
- **Check:** `plaintext: 13 bytes` / `sealed: 29 bytes` / `nonce: 12 bytes, tag: 16 bytes` / `tampered: cipher: message authentication failed`.

#### 6.3 Never the Same Nonce Twice
`never-the-same-nonce-twice` · code_tests · difficulty 8
- **Objective:** Give every sealed message a fresh random nonce and send it with the ciphertext.
- **Scenario:** The draft seals with one fixed nonce; from two commands sealed under it, Mirela reads the second by XORing the ciphertexts with the first, known command.
- **Exercise:** `sealCmd(aead cipher.AEAD, plaintext []byte) (nonce, sealed []byte)` uses a constant nonce; you take a new one from `newNonce()` (chapter 1) for every call.
- **Check:** two seals of the same plaintext differ in nonce and ciphertext; sealing `SET_MODE|SAFE` (known) and `DUMP_FILE|occ-0825.bin`, the test's XOR attack recovers the second command's first 13 bytes from the starter and nothing from yours; each sealed text opens with its own nonce.

#### 6.4 Quick Check: Sealing
`quick-check-sealing` · quiz · difficulty 7
- **Objective:** Check what GCM promises and what it needs from you.
- **Scenario:** Nkem's three questions before the first sealed pass.
- **Exercise:** Three questions.
- **Check:** A changed byte anywhere makes `Open` fail, so nothing tampered is ever read; a nonce may travel in the clear but must never repeat under one key; `Open` with the right key and the wrong nonce fails.

#### 6.5 Readable but Bound
`readable-but-bound` · code_tests · difficulty 8
- **Objective:** Authenticate a readable header as associated data, so nobody can change it.
- **Scenario:** Ystmark must read a sealed command's satellite and number to schedule it, and Mirela moves an old sealed command to a new number to slip it past the replay guard.
- **Exercise:** `sealWithHeader(aead cipher.AEAD, s Sealed, plaintext []byte) Sealed` and `openWithHeader(aead cipher.AEAD, s Sealed) ([]byte, error)` pass nil associated data; you pass `s.AD()` (given: `SEAL|2|cirrine-1|4471|k-2028-07`) to both.
- **Check:** a round trip; the same ciphertext with the header's sequence number, satellite, key ID or version changed → `Open` fails (the starter opens all four); the header fields stay readable.

#### 6.6 The Sealed Line
`the-sealed-line` · code_tests · difficulty 8
- **Objective:** Write and parse a binary message's text line strictly, with hex fields.
- **Scenario:** Ystmark's uplinkd queues lines of text, and a sealed command must survive the trip as one line.
- **Exercise:** `formatSealed(s Sealed) string` and `parseSealed(line string) (Sealed, error)` are stubs; you write the bible's seven-field `SEAL 2 …` line and parse it back, wrapping `ErrMalformed`, or `ErrVersion` for any version but 2.
- **Check:** the bible's sealed `PING` parses to its fields and formats back byte for byte; 6 or 8 fields, a sequence number of 0 or `44a1`, a 22-digit nonce, an odd-length or uppercase ciphertext, a ciphertext under 17 bytes and one over 272 → `ErrMalformed`; `SEAL 3 …` → `ErrVersion`.

#### 6.7 Open in the Right Order
`open-in-the-right-order` · code_tests · difficulty 8
- **Objective:** Check a sealed command in a fixed order, and record nothing until it opens.
- **Scenario:** Cirrine-1's flight software follows Mirela's order, and its stand-in on the ground must behave exactly the same before Friday.
- **Exercise:** `(s *Satellite) Accept(line string) (Command, error)` opens first and records `Seq` before opening; you follow the bible's order of checks and record `Seq` only after `Open` succeeds.
- **Check:** the bible's sealed `PING` → `cirrine-1`, 4471, `PING`; then, with `errors.Is`: the same line again → `ErrReplay`; a line for pipit-2 → `ErrWrongSat`; key `k-2027-03` → `ErrUnknownKey`; one changed ciphertext byte → `ErrForged`; a forged line numbered `math.MaxUint64` → `ErrForged`, and a good line numbered 4472 is still accepted after it.

#### 6.8 Cirrine's First Sealed Command
`cirrines-first-sealed-command` · code_tests · difficulty 9
- **Objective:** Seal commands from the key ring end to end, so only the satellite can read them.
- **Scenario:** Friday 25 August, 07:30 UTC: Cirrine-1 rises over Ystmark, and its first sealed command, `PING`, goes up.
- **Exercise:** `sealCommand(ring *KeyRing, keyID string, c Command) (string, error)` returns `""`; you check the command with `canonical`, derive the satellite's seal key from the ring, take a fresh nonce, seal `OP|ARGS` with the header as associated data, and format the line.
- **Check:** the line opens in the satellite stand-in as the same command; the line holds neither `PING` nor its hex; two seals of one command differ; a command with a bar in an argument, or an unknown key ID → error and no line.

## Chapter 7 · Signatures Anyone Can Check

`signatures-anyone-can-check` · 8 lessons

Sign and verify with `crypto/ed25519`, publish public keys as PEM through
`x509.MarshalPKIXPublicKey`, sign the exact bytes you send with `json.RawMessage`, choose
verifying keys by ID, and check a delivered file against its receipt.

#### 7.1 Receipts Both Sides Can Forge
`receipts-both-sides-can-forge` · code_tests · difficulty 8
- **Objective:** Sign with an Ed25519 private key and let others verify with only the public key.
- **Scenario:** Monday 28 August: with the secret both sides hold, the test, playing Furrowlight, writes a receipt saying the file landed a day late, and Apogee's own check accepts it.
- **Exercise:** `Notary` signs with HMAC: `newNotary`, `Sign`, `Published` (which hands customers the secret) and `checkReceipt(published, msg, sig []byte) bool`; you make the notary an Ed25519 key from its seed, sign with `ed25519.Sign`, publish the 32-byte public key, and verify with `ed25519.Verify` after checking the key's length.
- **Check:** Apogee's signature verifies with `Published()` alone; the test, holding everything a customer holds, can't make a signature `checkReceipt` accepts for a changed receipt (the starter accepts its forgery); a changed receipt byte → false; a 31-byte public key → false, not a panic; a 31-byte seed → `newNotary` error.

#### 7.2 Keys From a Seed
`keys-from-a-seed` · code · difficulty 7
- **Objective:** Make an Ed25519 key pair from a 32-byte seed with `ed25519.NewKeyFromSeed`.
- **Scenario:** Mirela makes the test receipt key from a seed, so every test signs the same way.
- **Exercise:** `main` prints placeholders; you derive the key from the SHA-256 of `receipts r-2028-08`, sign a short receipt, and print the sizes, the public key's first 8 bytes and two checks.
- **Check:** `public key: 32 bytes, ba6c6226d2c3c056` / `signature: 64 bytes` / `verifies: true` / `one byte changed: false`.

#### 7.3 Quick Check: Signatures
`quick-check-signatures` · quiz · difficulty 7
- **Objective:** Check what a public key can and can't do.
- **Scenario:** Yusra's three questions before she writes the customer guide's page.
- **Exercise:** Three questions.
- **Check:** Anyone with the public key can verify, and nobody can sign with it; Ed25519 signs a message the same way every time, so a test may expect a fixed signature; an HMAC can't settle a dispute, because both sides hold its key.

#### 7.4 The Guide's Public Key
`the-guides-public-key` · code_tests · difficulty 8
- **Objective:** Publish a public key as a PEM block with `x509.MarshalPKIXPublicKey`, and read one back safely.
- **Scenario:** Yusra's customer guide must give every customer Apogee's receipt key in a form their own tools read.
- **Exercise:** `publishKey(pub ed25519.PublicKey) ([]byte, error)` and `readKey(pemText []byte) (ed25519.PublicKey, error)` are stubs; you encode a `PUBLIC KEY` block, and decode it with `pem.Decode`, `x509.ParsePKIXPublicKey` and a type assertion.
- **Check:** a round trip of the test key; a `PRIVATE KEY` block, text after the block, empty input → errors; an ECDSA P-256 public key (the test makes one) → `not an Ed25519 key`.

#### 7.5 Sign the Bytes You Send
`sign-the-bytes-you-send` · code_tests · difficulty 8
- **Objective:** Verify a signature over the exact bytes received before decoding them, kept with `json.RawMessage`.
- **Scenario:** A customer's tool decoded a receipt and re-encoded it before checking, so it accepted a receipt carrying a field Apogee never signed.
- **Exercise:** `wrapReceipt(n *Notary, r Receipt) ([]byte, error)` is given; `openEnvelope(pub ed25519.PublicKey, env []byte) (Receipt, error)` decodes the receipt, marshals it again and verifies that; you verify the envelope's raw `receipt` bytes first, then decode them strictly.
- **Check:** Apogee's envelope opens; the same receipt with `"refund":true` added inside → refused (the starter accepts it); the same receipt spaced differently → refused; a `sig` of 127 hex digits → `ErrBadEnvelope`; an unknown field in a receipt the test signs itself → refused by strict decoding.

#### 7.6 Never the Key It Brings
`never-the-key-it-brings` · code_tests · difficulty 8
- **Objective:** Choose the verifying key by its ID from your own list, never from the message.
- **Scenario:** A draft envelope carried its public key inside, so anyone could sign with a key of their own and send it along.
- **Exercise:** `verifyEnvelope(keys map[string]ed25519.PublicKey, env []byte) (Receipt, error)` verifies with the envelope's `pub` field when there is one; you look the key up by `key_id` in `keys`, refuse envelopes with fields outside `Envelope`, and check the receipt's own `key_id`.
- **Check:** an envelope signed with `r-2028-08` → the Receipt; one signed with the test's own key and carrying it → refused (the starter accepts it); `key_id` `r-2027-01` → `ErrUnknownSigner`; an inner `key_id` that differs from the envelope's → `ErrBadEnvelope`.

#### 7.7 Every Landing Gets a Receipt
`every-landing-gets-a-receipt` · code_tests · difficulty 8
- **Objective:** Issue a signed receipt from a landing, with the file's digest and the signing time.
- **Scenario:** From Friday 1 September, bookingd signs a receipt for every file dlinkd delivers.
- **Exercise:** `issue(n *Notary, l Landing, b Booking, digest string, now func() time.Time) ([]byte, error)` returns nil; you fill a `Receipt` (version 1, the notary's key ID, the booking's satellite and station, the landing's file, bytes and time, `signed_at` from `now` truncated to the second) and wrap it.
- **Check:** the envelope opens with the published key and every field matches the landing and the booking; `signed_at` is the fake clock's time without its milliseconds; a digest that `parseDigest` refuses → error and nothing signed; a landing for another booking → error.

#### 7.8 Proof for Furrowlight
`proof-for-furrowlight` · code_tests · difficulty 9
- **Objective:** Check a delivery from the customer's side: the signature, the key's window and the file against the receipt's digest.
- **Scenario:** Friday 1 September: Furrowlight checks the receipt for Sunday's Furrow-A file against the file already in its own storage, with the checker from Yusra's guide.
- **Exercise:** `checkDelivery(keys []ReceiptKey, env []byte, file io.Reader) (Receipt, error)` returns an error; you follow the bible's order of checks, then hash the file with the receipt's own algorithm and compare its size and digest.
- **Check:** the fixture's receipt and its 64 KiB file → the Receipt; the file with one byte changed or one byte short → `ErrWrongFile`; a receipt signed after its key's retire time → `ErrOutsideWindow`; a revoked key → `ErrRevoked`; a receipt for a SHA3-256 file checks too.

## Chapter 8 · Agreeing on a Secret

`agreeing-on-a-secret` · 8 lessons

Agree on a shared secret with `ecdh.X25519`, check public shares, sign them with Ed25519
against known keys, derive session keys with HKDF, make fresh keys per session, and
encapsulate a key with ML-KEM-768.

#### 8.1 Two Halves of a Secret
`two-halves-of-a-secret` · code_tests · difficulty 8
- **Objective:** Agree on a shared secret with X25519: `GenerateKey`, `PublicKey` and `ECDH`.
- **Scenario:** Monday 4 September: Ystmark's fiber is cut; Port Ellery and Ystmark each make a key, send only its public half over the backhaul, and arrive at the same 32 bytes.
- **Exercise:** `agree(mine *ecdh.PrivateKey, theirs []byte) ([]byte, error)` returns nil; you parse `theirs` with `ecdh.X25519().NewPublicKey` and call `mine.ECDH`.
- **Check:** two fresh key pairs agree in both directions, 32 bytes; the fixed keys `ephemeral PEL` and `ephemeral YST` agree on `c50d64c8…cc7a`; a third key gets a different secret.

#### 8.2 Shares on the Wire
`shares-on-the-wire` · code_tests · difficulty 8
- **Objective:** Send public shares as hex and refuse malformed or degenerate ones.
- **Scenario:** Ystmark's share arrives as 64 hex digits over a link that sometimes garbles a line.
- **Exercise:** `parseShare(text string) (*ecdh.PublicKey, error)` and `agreeHex(mine *ecdh.PrivateKey, text string) ([]byte, error)` decode without checks; you require 64 lowercase hex digits and pass on `ECDH`'s error.
- **Check:** PEL's share `a10c608e…a93e` parses; 62 digits, a `z`, uppercase → errors; 64 zeros parse, but `agreeHex` returns an error, never 32 zero bytes.

#### 8.3 Quick Check: Key Agreement
`quick-check-key-agreement` · quiz · difficulty 8
- **Objective:** Check what key agreement hides and what it doesn't prove.
- **Scenario:** Sunniva's telegraphic questions over the backhaul.
- **Exercise:** Three questions.
- **Check:** Whoever records both public shares still can't compute the secret; the shares alone prove nothing about who sent them; a share is public, and its private key never leaves its station.

#### 8.4 Someone in the Middle
`someone-in-the-middle` · code_tests · difficulty 9
- **Objective:** Authenticate key-agreement shares with Ed25519 signatures from keys known in advance.
- **Scenario:** The test plays whoever runs the backhaul: it swaps both shares for its own and ends up with a key for each side, while both sides believe they agreed.
- **Exercise:** `hello(id Identity, share *ecdh.PrivateKey) string` writes `HELLO <code> <share hex>`, and `checkHello(line string, known map[string]ed25519.PublicKey) (string, *ecdh.PublicKey, error)` accepts any share; you sign `backhaul share|<code>|<share hex>` and append the signature, and check every field as the bible says.
- **Check:** PEL's and YST's HELLOs verify; the middle's HELLO claiming to be YST, signed with the `identity MID` key, → `ErrBadHello`, so its swap fails (the starter lets both swaps through); a changed share, an unknown code and a missing signature → `ErrBadHello`.

#### 8.5 Never the Raw Secret
`never-the-raw-secret` · code_tests · difficulty 8
- **Objective:** Turn a shared secret into a session key with HKDF, binding both shares into the info.
- **Scenario:** The draft used the X25519 output itself as the AES key.
- **Exercise:** `sessionKey(shared, dialerShare, answererShare []byte) ([]byte, error)` returns `shared`; you derive 32 bytes with `hkdf.Key`, no salt, and the info `apogee backhaul v1|<dialer hex>|<answerer hex>`.
- **Check:** the fixed vector `24e57dc8…958e`; swapping the two shares → another key; the key differs from the raw secret; a secret that isn't 32 bytes → error.

#### 8.6 Fresh Keys Every Session
`fresh-keys-every-session` · code_tests · difficulty 8
- **Objective:** Make a new X25519 key for every session, so a key stolen later opens no recording.
- **Scenario:** The draft keeps one X25519 key in Ystmark's config, and whoever steals it later can open every backhaul session recorded before.
- **Exercise:** `type Backhaul struct{ static *ecdh.PrivateKey }` and `(b *Backhaul) NewSession() (*ecdh.PrivateKey, error)` returns `static`; you generate a fresh key for every call and remove `static`.
- **Check:** two sessions' shares differ; the test, holding everything the `Backhaul` keeps after a session, can't recompute that session's secret from the recorded shares (with the starter it can); each share is 32 bytes.

#### 8.7 Ready for Bigger Computers
`ready-for-bigger-computers` · code_tests · difficulty 8
- **Objective:** Agree on a key with ML-KEM-768: `mlkem.GenerateKey768`, `Encapsulate` and `Decapsulate`.
- **Scenario:** Tomás has read that recordings made today may be opened by machines built years from now, and wants a key exchange meant for them.
- **Exercise:** `wrapFor(ek []byte) (key, ciphertext []byte, err error)` and `unwrap(dk *mlkem.DecapsulationKey768, ciphertext []byte) ([]byte, error)` are stubs; you parse the encapsulation key with `mlkem.NewEncapsulationKey768`, encapsulate, and decapsulate.
- **Check:** both sides get the same 32 bytes; the encapsulation key is 1,184 bytes and the ciphertext 1,088; a 100-byte ciphertext → error; one changed ciphertext byte → no error but a different key (the readme says why: the session then fails at its first sealed message).

#### 8.8 Ystmark's Storm Key
`ystmarks-storm-key` · code_tests · difficulty 9
- **Objective:** Run an authenticated key agreement over a connection and seal messages with the session key.
- **Scenario:** Thursday 7 September: the fiber is still cut, and this week's command file for Ystmark must cross the backhaul sealed.
- **Exercise:** `dial(conn io.ReadWriter, me Identity, peer string, known map[string]ed25519.PublicKey) ([]byte, error)` and `answer(conn io.ReadWriter, me Identity, known map[string]ed25519.PublicKey) (string, []byte, error)` are stubs; you exchange HELLOs with fresh shares, check them, derive the session key, and the given `sendData` and `readData` carry the file as `DATA` lines.
- **Check:** over `net.Pipe()`, PEL and YST end with the same key and YST reads the file's three lines; a middle relaying with its own shares → refused before any DATA line; a DATA line changed in transit → error; the whole exchange ends within 1 s.

## Chapter 9 · Who Are You, Server?

`who-are-you-server` · 8 lessons

Build a private certificate authority with `x509.CreateCertificate`, read PEM certificates,
issue server and client certificates, and verify chains, names, usages and validity at an
injected time.

#### 9.1 A Name Is Not Proof
`a-name-is-not-proof` · code_tests · difficulty 8
- **Objective:** Trust a certificate only when `Verify` chains it to your own roots.
- **Scenario:** Monday 11 September: the draft checker accepts any certificate whose name says `yst.apogee.internal`, and the test makes such a certificate in a millisecond.
- **Exercise:** `trusted(cert *x509.Certificate, host string, roots, inters *x509.CertPool, now time.Time) error` compares `cert.Subject.CommonName` with `host`; you call `cert.Verify` with `Roots`, `Intermediates`, `DNSName` and `CurrentTime`.
- **Check:** Ystmark's certificate → nil; the impostor (same names, from a root also named `Apogee Root CA 2028`) → error (the starter accepts it); Ystmark's certificate for `pel.apogee.internal` → error; after its `NotAfter` → error.

#### 9.2 What a Certificate Says
`what-a-certificate-says` · code · difficulty 8
- **Objective:** Read a PEM certificate with `pem.Decode` and `x509.ParseCertificate`.
- **Scenario:** Paz wants to know what Ystmark's new certificate actually claims, before it goes anywhere near Salar Quieto.
- **Exercise:** `main` prints the file's size; you decode `yst.pem` (read-only), parse its first block and print six fields.
- **Check:** `subject: yst.apogee.internal` / `issuer: Apogee Stations CA 2028` / `serial: 1001` / `valid: 2028-09-22 to 2028-12-21` / `names: [yst.apogee.internal]` / `CA: false`.

#### 9.3 Apogee's Own Root
`apogees-own-root` · code_tests · difficulty 9
- **Objective:** Make a self-signed root certificate with `x509.CreateCertificate` and an Ed25519 key.
- **Scenario:** Monday 18 September, 10:00: Mirela makes the root on a laptop that never touches a network, then locks it in the safe.
- **Exercise:** `newRoot(key ed25519.PrivateKey, now time.Time) (*x509.Certificate, error)` returns nil; you fill the template as the bible's table says (serial 1, `Apogee Root CA 2028`, organization Apogee, 10 years from `now`, a CA with path length 1, cert and CRL signing), create it with the root as its own parent, and parse the result.
- **Check:** subject and issuer `Apogee Root CA 2028`; `IsCA`, `MaxPathLen` 1, `NotAfter` exactly 10 years after `now`; cert-sign key usage; it verifies as the only root in its own pool; a nil key → error, never a panic.

#### 9.4 Quick Check: Certificates
`quick-check-certificates` · quiz · difficulty 8
- **Objective:** Check what a certificate binds together, and why the root stays offline.
- **Scenario:** Nkem's three questions after the root goes in the safe.
- **Exercise:** Three questions.
- **Check:** A certificate binds a public key to names, signed by its issuer; anyone can make a certificate with any name, but only the root's key makes one that chains to Apogee's root; whoever holds the root's private key can certify anything, so it never leaves the safe.

#### 9.5 A Certificate for Ystmark
`a-certificate-for-ystmark` · code_tests · difficulty 9
- **Objective:** Issue a server certificate from an intermediate, with its name, usage and 90 days of validity.
- **Scenario:** Sunniva's station gets the first certificate the stations CA ever issues.
- **Exercise:** `issueServer(ca Issuer, host string, serial int64, pub ed25519.PublicKey, now time.Time) (*x509.Certificate, error)` returns nil; you set the common name and DNS name to `host`, server-auth usage, digital-signature key usage, `now` to `now` plus 90 days, and sign with the CA's key.
- **Check:** it chains through `Apogee Stations CA 2028` to the root for `yst.apogee.internal`; it is refused when client auth is asked for; serial 1001, `NotAfter` 2028-12-21T00:00:00Z for a `now` of 2028-09-22; it isn't a CA.

#### 9.6 The Middle of the Chain
`the-middle-of-the-chain` · code_tests · difficulty 9
- **Objective:** Verify through an intermediate, and refuse a leaf that tries to act as a CA.
- **Scenario:** The root stays in the safe, so stations present their certificate with the intermediate, and Mirela checks that a station's certificate can't sign another.
- **Exercise:** `verifyStation(leaf *x509.Certificate, chain []*x509.Certificate, roots *x509.CertPool, host string, now time.Time) error` passes no intermediates; you put the chain into an `Intermediates` pool and ask for server-auth usage.
- **Check:** Ystmark's leaf with its chain → nil; without the chain → error; a certificate the test signs with Ystmark's leaf key → error; an intermediate from the impostor root → error.

#### 9.7 Too Early, Too Late
`too-early-too-late` · code_tests · difficulty 9
- **Objective:** Check validity at an injected time, knowing that both ends of a certificate's validity count.
- **Scenario:** After a diesel restart, uplinkd at Lantern Atoll ran 40 minutes slow, and a new certificate looked "not yet valid" to it.
- **Exercise:** `checkAt(leaf *x509.Certificate, chain []*x509.Certificate, roots *x509.CertPool, host string, now func() time.Time) error` checks at `time.Now()`; you pass `now()` as `CurrentTime`.
- **Check:** at `NotBefore` exactly → nil; one second before → error; at `NotAfter` exactly → nil; one second after → error (the starter, on today's clock, fails all four).

#### 9.8 Certificates for Every Station
`certificates-for-every-station` · code_tests · difficulty 9
- **Objective:** Issue a set of certificates and save them as PEM files, private keys with mode 0600.
- **Scenario:** Friday 22 September: every station, uplinkd and passplan get their certificates before the weekend's passes.
- **Exercise:** `issueAll(ca Issuer, dir string, now time.Time) error` is a stub; you issue the five station certificates (1001 to 1005) and the client certificates of uplinkd (2001) and passplan (2002), and write each `<name>.pem` (leaf, then intermediate) and `<name>.key` (PKCS #8) through a temporary file and a rename.
- **Check:** seven pairs under `t.TempDir()`, each loading with `tls.X509KeyPair`; each certificate verifies for its own name and usage only; every `.key` has mode 0600; an issuer without a key → error and no files written.

## Chapter 10 · Every Connection Encrypted

`every-connection-encrypted` · 9 lessons

Dial and serve with `crypto/tls`: your own roots and server name, your own certificate chain,
TLS 1.3 only, client certificates, authorization from `VerifiedChains`, and TLS over a raw TCP
connection.

#### 10.1 Trusting Anyone
`trusting-anyone` · code_tests · difficulty 9
- **Objective:** Verify a TLS server against your own roots and its name, never with `InsecureSkipVerify`.
- **Scenario:** Monday 25 September: the draft client set `InsecureSkipVerify` to get past an error, and the test's impostor, an `httptest.NewTLSServer` with a certificate for `example.com`, receives a command meant for Ystmark.
- **Exercise:** `newClient(roots *x509.CertPool, host string) *http.Client` sets `InsecureSkipVerify: true`; you set `RootCAs`, `ServerName` and a 2-second timeout.
- **Check:** Ystmark's test server, with Apogee's certificate, answers `ok YST`; the impostor → the request fails (the starter reaches it); the right root asking for `pel.apogee.internal` → fails.

#### 10.2 Ystmark's Own Certificate
`ystmarks-own-certificate` · code_tests · difficulty 9
- **Objective:** Serve TLS with your own certificate chain through `tls.Config.Certificates`, on an injected clock.
- **Scenario:** Ystmark's uplinkd starts using the certificate it was issued on 22 September.
- **Exercise:** `serverConfig(leaf, inter *x509.Certificate, key ed25519.PrivateKey, now func() time.Time) *tls.Config` returns an empty config, so httptest serves its own certificate; you build a `tls.Certificate` with the leaf first, then the intermediate, and set `Time: now`.
- **Check:** a client trusting only Apogee's root and asking for `yst.apogee.internal` completes a request at the fake time; the server presents serial 1001 with its intermediate; a client whose clock is past `NotAfter` fails the handshake.

#### 10.3 Nothing Older Than 1.3
`nothing-older-than-one-three` · code_tests · difficulty 9
- **Objective:** Require TLS 1.3 with `MinVersion`, and read what was agreed from `ConnectionState`.
- **Scenario:** An old tool at Salar Quieto still speaks TLS 1.2, and Mirela wants it refused, not quietly accepted.
- **Exercise:** `harden(cfg *tls.Config)` does nothing; you set `MinVersion: tls.VersionTLS13`.
- **Check:** a client limited to TLS 1.2 is refused; an ordinary client gets version 1.3 and the key exchange `X25519MLKEM768`, the hybrid with chapter 8's ML-KEM; `harden` leaves the config's other fields as they were.

#### 10.4 Quick Check: Secure Connections
`quick-check-secure-connections` · quiz · difficulty 8
- **Objective:** Check what a TLS client checks and what each setting does.
- **Scenario:** Nkem's three questions before the first station switches.
- **Exercise:** Three questions.
- **Check:** `ServerName` is the name the server's certificate is checked against, even when you dial an IP address; with `RootCAs` nil the system's roots are used, and they don't include Apogee's; `InsecureSkipVerify` turns off every check of the server's certificate.

#### 10.5 Both Ends Proven
`both-ends-proven` · code_tests · difficulty 9
- **Objective:** Require and verify client certificates with `ClientAuth` and `ClientCAs`.
- **Scenario:** Ystmark must refuse any client that can't prove it is one of Apogee's servers.
- **Exercise:** `requireClients(cfg *tls.Config, roots *x509.CertPool)` does nothing and `clientConfig(roots, host, cert, now)` sends no certificate; you set `ClientAuth: tls.RequireAndVerifyClientCert` and `ClientCAs`, and add the client's chain.
- **Check:** uplinkd's client → 200; no client certificate → refused; a client certificate from the impostor root → refused; Ystmark's own server certificate used as a client's → refused (wrong usage).

#### 10.6 Who Is Calling
`who-is-calling` · code_tests · difficulty 9
- **Objective:** Authorize a request by the verified client certificate in `r.TLS.VerifiedChains`.
- **Scenario:** Every Apogee server may ask Ystmark for its status, but only uplinkd may queue a command.
- **Exercise:** `caller(r *http.Request) (string, bool)` returns `"", false` and `stationHandler` queues any POST; you return the common name of `r.TLS.VerifiedChains[0][0]`, and answer 403 `only uplinkd may queue commands` for any other caller.
- **Check:** uplinkd → 202 `queued cirrine-1 4471`; passplan → 403 and nothing queued; `GET /status` from passplan → 200 `ok YST`; a request without TLS (a recorder, `r.TLS` nil) → 403, never a panic.

#### 10.7 Quick Check: Mutual Proof
`quick-check-mutual-proof` · quiz · difficulty 9
- **Objective:** Check how mutual TLS identifies a caller and who decides what it may do.
- **Scenario:** Mirela's two questions.
- **Exercise:** Two questions.
- **Check:** `VerifiedChains` holds only chains the handshake verified, while `PeerCertificates` is whatever the client sent; a verified certificate says who is calling, and the handler still decides what that caller may do.

#### 10.8 Encrypting a Bare Socket
`encrypting-a-bare-socket` · code_tests · difficulty 9
- **Objective:** Wrap a raw TCP connection in TLS with `tls.Dialer`, handshaking within a deadline.
- **Scenario:** Ystmark's dish controllers speak their own raw socket protocol, and Sunniva wants the link from uplinkd to them encrypted too.
- **Exercise:** `dialController(ctx context.Context, addr string, cfg *tls.Config) (net.Conn, error)` dials plain TCP; you dial with a `tls.Dialer` (or `tls.Client` and `HandshakeContext`), so the handshake obeys `ctx`.
- **Check:** the TLS stand-in controller echoes `PING yst`; a plain-TCP listener that never answers → an error within the context's 200 ms; a controller with the impostor's certificate → error.

#### 10.9 The Uplink, End to End
`the-uplink-end-to-end` · code_tests · difficulty 10
- **Objective:** Send a sealed command across mutual TLS 1.3 to a station that relays it without being able to read it.
- **Scenario:** Friday 6 October: every station speaks TLS 1.3 with both ends proven, and Cirrine-1's sealed commands cross it to Ystmark.
- **Exercise:** `send`, `stationHandler` and `Queue` are given; you write `uplinkClient(roots *x509.CertPool, host string, cert tls.Certificate, now func() time.Time) *http.Client` and `stationServer(leaf, inter *x509.Certificate, key ed25519.PrivateKey, roots *x509.CertPool, now func() time.Time) *tls.Config` with everything this chapter required.
- **Check:** a sealed line for cirrine-1 arrives byte for byte in Ystmark's queue and still opens in the satellite stand-in, while the station holds no seal key; TLS 1.3 agreed; refused: an impostor station, an impostor uplinkd, a TLS 1.2 client, and passplan posting a command; no request takes over 1 s.

## Chapter 11 · Rotation Day

`rotation-day` · 8 lessons

Accept old and new key versions through an overlap, sign with the newest, send new keys sealed
under old ones, revoke with `x509.CreateRevocationList` and `VerifyConnection`, schedule each
switch between passes, and run the whole day.

#### 11.1 The Midnight Switch
`the-midnight-switch` · code_tests · difficulty 9
- **Objective:** Accept a key and its successor side by side through an overlap window.
- **Scenario:** Monday 9 October, in Mirela's rehearsal: a command sealed with `k-2028-07` at 23:59:58 reaches Cirrine-1 at 00:00:03, is refused, and the pass is lost.
- **Exercise:** `accepts(ring []KeyVersion, id string, now time.Time) bool` accepts only the version that would sign at `now`, a hard switch; you follow the bible's rules: every version in the ring is accepted until its `Retire` unless revoked, the next one even before its `Activate`.
- **Check:** `k-2028-07` at 2028-10-17T00:00:03Z → true (the starter says false); `k-2028-10` on 9 October → true (the starter says false); `k-2028-07` at exactly 2028-10-24T00:00:00Z → false and one second earlier → true; a revoked version → false at any time; an ID not in the ring → false.

#### 11.2 Sign With the Newest
`sign-with-the-newest` · code_tests · difficulty 9
- **Objective:** Choose the signing key by time: the newest activated version that isn't retired or revoked.
- **Scenario:** uplinkd must switch to `k-2028-10` at exactly 00:00:00 on Rotation Day, and never seal with a revoked key.
- **Exercise:** `signingKey(ring []KeyVersion, now time.Time) (KeyVersion, error)` returns the first version in the slice; you pick by `Activate` as the bible's rules say, and return `no key to sign with at <time>` when none fits.
- **Check:** 2028-10-16T23:59:59Z → `k-2028-07`; 2028-10-17T00:00:00Z → `k-2028-10`; with `k-2028-10` revoked, on 20 October → `k-2028-07` and on 25 October → the error; the ring in any order gives the same answers.

#### 11.3 Quick Check: Overlap
`quick-check-overlap` · quiz · difficulty 8
- **Objective:** Check the rules of overlap, retirement and revocation.
- **Scenario:** Ingrid's three questions for Rotation Day's checklist.
- **Exercise:** Three questions.
- **Check:** Verifiers get the new key before signers start using it; a receipt key stays published after it stops signing, so old receipts still verify; a retired key is refused from its retire instant, a revoked one at once.

#### 11.4 New Keys Before the Day
`new-keys-before-the-day` · code_tests · difficulty 10
- **Objective:** Send a new key version to a satellite inside a command sealed under its current key.
- **Scenario:** The week of 9 October: every satellite must hold `k-2028-10` before Rotation Day, and the new keys may never cross the air readable.
- **Exercise:** `loadKeyCommand(ring *KeyRing, satID string, seq uint64, next KeyVersion) (string, error)` returns `""`; you derive the satellite's new sign and seal keys from `next` and seal `LOAD_KEY|<id>|<sign hex>|<seal hex>` with `sealCommand` under the current version.
- **Check:** the stand-in satellite opens it and then accepts a command sealed with `k-2028-10`, and still accepts `k-2028-07`; neither new key appears in the line, in hex or raw; a LOAD_KEY sealed under a key the satellite doesn't hold → refused and nothing loaded; a second LOAD_KEY with an old sequence number → `ErrReplay`.

#### 11.5 A List of the Revoked
`a-list-of-the-revoked` · code_tests · difficulty 10
- **Objective:** Publish revocations in a signed list with `x509.CreateRevocationList`, and check certificates against it.
- **Scenario:** Tuesday 17 October, 10:12 UTC: the Cape Wren contractor reports that the old station computer she took away for recycling still holds its disk, and with it the key of certificate 1005.
- **Exercise:** `revoke(ca Issuer, number int64, entries []x509.RevocationListEntry, now time.Time) ([]byte, error)` and `isRevoked(crlDER []byte, ca, cert *x509.Certificate, now time.Time) (bool, error)` are stubs; you create the list (`ThisUpdate` now, `NextUpdate` 7 days later), and when checking, parse it, verify it with `CheckSignatureFrom`, check that it is current, and compare serial numbers.
- **Check:** list 1 revokes 1005 and not 1105; a list signed by the impostor CA → error, never "not revoked"; a list one second past its `NextUpdate` → error, and at `NextUpdate` exactly → still current; list number 1.

#### 11.6 Refused at the Handshake
`refused-at-the-handshake` · code_tests · difficulty 10
- **Objective:** Refuse revoked peers during the TLS handshake with `VerifyConnection`.
- **Scenario:** uplinkd's connections to Cape Wren must fail the moment certificate 1005 is listed, though its chain still verifies.
- **Exercise:** `withRevocations(cfg *tls.Config, current func() []byte, ca *x509.Certificate, now func() time.Time)` does nothing; you set `VerifyConnection` to check the peer's leaf with `isRevoked` against `current()`, answering `certificate 1005 is revoked`.
- **Check:** a server presenting 1105 → connects; 1005 → the handshake fails with that error; a list with a bad signature → every handshake fails; the list changes between two connections and the second sees the new one.

#### 11.7 Between Two Passes
`between-two-passes` · code_tests · difficulty 9
- **Objective:** Find a station's first quiet gap in a day of passes, with sorting and `time` arithmetic.
- **Scenario:** Ingrid's rule for Rotation Day: a station changes certificates only in 10 quiet minutes, 5 clear of any pass on either side.
- **Exercise:** `switchTime(passes []Pass, station string, from time.Time) (time.Time, error)` returns `from`; you sort that station's passes by AOS and return the bible's earliest fitting start.
- **Check:** Rotation Day's passes → YST, SQT and CWR at 01:00:00, PEL at 01:12:10, LAN at 01:15:40 (the starter switches PEL and LAN during a pass); passes in any order, mixed with other stations'; a station with no passes → `from`; a station booked every 20 minutes for a day → `no 10-minute gap for … before …`.

#### 11.8 Rotation Day, Hour by Hour
`rotation-day-hour-by-hour` · code_tests · difficulty 10
- **Objective:** Run a whole day of passes through a key and certificate rotation, failing closed and missing nothing.
- **Scenario:** Tuesday 17 October: Ingrid runs Rotation Day like a pass, from the Cirrine-1 pass that crosses midnight to Furrow-A over Salar Quieto at 14:02.
- **Exercise:** `runDay(d Day) Report` is a stub; using the season's helpers (given), you walk the day in time order: seal each command with the key `signingKey` picks when it is sealed, deliver it to the satellite stand-in 5 seconds later, connect to each station with the certificate in force (switching at its `switchTime`), apply revocation list 1 at 10:12, and report each pass as the bible's report lines say.
- **Check:** the real day → `Report{Flown: 10, Missed: 0}` and its ten lines, the midnight pass's being `23:52:10 cirrine-1 YST flown k-2028-07+k-2028-10 1001`; every station's passes after its switch show its 11xx certificate; Ingrid's rehearsal day, where Lantern Atoll never switched and list 1 also revokes 1004 at 05:00 → `Flown: 9, Missed: 1` and `09:14:20 brine-3 LAN missed: certificate 1004 is revoked`; the run ends well under a second.

## Changes in writing and review

- **Security claims** are exact and never overstated: a changed byte gives a digest that looks
  unrelated (not "changes all of it"); `==` stops at the first difference it finds, in chunks;
  a fast hash of a common password falls to a precomputed list; forward secrecy doesn't protect
  against a stolen identity key; claims about future quantum computers stay modest.
- **Tests** that parse source accept every correct solution: helpers, length checks and
  `ConstantTimeCompare(...) == 1` all pass, while the named mistakes still fail.
- **Chapter 4** `matches` spends one dummy derivation on old `sha256$` records too, so an old
  record can't reveal an account name by timing; 4.6 also tests a record stronger than 600,000
  iterations.
- **Chapter 5** `(*KeyRing).GoString` hides secrets under `%#v`; 5.5 names the number verbs
  that still print a key's bytes. 5.1 carries chapter 3's `tag` and `verify` unchanged.
- **Chapter 6** sealed commands check the sequence number before opening, with the reason in
  6.7: a number already seen is refused before any work, and a number is recorded only after
  `Open` succeeds. Only a SEAL line can be refused for its version.
- **Chapters 8–9** a stolen identity key can sign shares until it is replaced; 9.5 refuses an
  issuer without its certificate or key, and its certificate 1001 is signed on Friday.
- **Chapters 10–11** `newClient` takes the injected clock, since Ystmark's certificate is valid
  only from 22 September; 11.4 loads new keys through `Satellite.Receive`, so `Accept` keeps
  chapter 6's text; 11.6 also checks resumed sessions.
- **Given code** keeps the text the learner last saw, including `KeyVersion`'s comments.
