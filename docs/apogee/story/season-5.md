# Apogee story bible: Season 5

### Season 5: keeping secrets

After Season 4, the learner starts the next assignment on Monday 10 July 2028. Their new badge
says **Ground Security Engineer**, still No. 039, based at Port Ellery headquarters, since
Monday 10 July 2028. On Friday 7 July Cirrine Atmospherics signed a bigger contract with three
conditions: from Friday 25 August nobody but Cirrine and Apogee may read Cirrine-1's commands,
not even someone with a dish listening to the uplink; every file delivered to Cirrine comes
with a fingerprint Cirrine's archive can check; and every key Apogee uses for Cirrine is
replaced by Tuesday 31 October, then once a year. Mirela leads the work and reviews every
change ("LGTM once it fails closed"); she plans **Rotation Day** for Tuesday 17 October: every
key in the network replaced, at every station, without missing a single pass. Yusra carries
the customers' wishes, Tomás the radio, Ingrid the passes, Sunniva and Paz their stations, and
Nkem asks the whiteboard questions.

Story time runs from Monday 10 July to Friday 3 November 2028: thirteen chapters: seven of about
a week, four (chapters 6, 9, 10 and 12) of about two, chapter 11 of an afternoon and a night
(Friday 6 to Saturday 7 October), and chapter 13, the epilogue after the overlap ends, of a week
and a half (Wednesday 25 October to Friday 3 November). No new people, places, satellites or customers
are introduced. Customers speak through Yusra ("TL;DR: ...") and Furrowlight through Hollis's
capitals. The Cape Wren contractor (the bible's "a contractor visits monthly") stays unnamed.

#### The season, chapter by chapter

| Ch | Week of | Where | Who needs what | What is at stake |
|---|---|---|---|---|
| 1 | Mon 10 Jul 2028 | Mirela's desk, Port Ellery | Mirela, before anything is sealed: "Every secret this season starts as random bytes. Show me where ours come from." The draft nonce helper for Cirrine's commands and the ops console's session tokens both come from `math/rand/v2` seeded with the clock | Whoever knows roughly when a token or nonce was made can make it again. On Friday 14 July at 16:00 Mirela makes the season's master key, `k-2028-07`, with your helper |
| 2 | Mon 17 Jul | the backend desks, with dlinkd's delivery folder | Yusra: "TL;DR: Cirrine wants proof that every file is exactly what left the satellite." Cirrine's archive checks SHA3-256, Apogee's own deliveries SHA-256 | Season 1's CRC-32 catches radio noise, not a person: Mirela changes "0.71" to "0.17" in a Furrowlight field report and the CRC still matches. Cirrine's first fingerprinted delivery is on Friday 21 July |
| 3 | Mon 24 Jul | Mirela's desk | Mirela reviews the signing path before sealing builds on it, and finds three weak spots: Pipit-3's flight software checks a homemade signature, `sha256(key + command)`; a `\|` inside an argument signs like two arguments; the replay guard counts a command before checking its signature | Pipit-3 would obey a command nobody signed, and one forged sequence number could lock a satellite out of every later command. Mirela signs off on Friday 28 July after her twelve planted lines |
| 4 | Mon 31 Jul | the Fishbowl | On Sunday 30 July the Fishbowl's spare laptop goes missing from the night-shift desk. It held a backup of the ops console's user file: 26 accounts, every password an unsalted SHA-256. Ingrid logs it; Mirela: "Assume whoever has it has a fast computer and a list of common passwords" | Every password is changed by Friday 4 August, and the console's login reopens at 16:00 that day with records that make every guess at a stolen copy cost 600,000 rounds |
| 5 | Mon 7 Aug | the lab bench at Port Ellery | Mirela's key inventory: `k-2027-03`, the one key that signs commands for every satellite, sits in uplinkd, in the lab bench's spare Pipit-2 flight computer and in two old test scripts | Whoever reads it out of any one of them can command all seven satellites. Each satellite gets keys of its own, derived from `k-2028-07`, which only uplinkd holds; they go up on the passes of Thursday 10 August, and `k-2027-03` retires at 00:00 UTC on Friday 11 August |
| 6 | Mon 14 Aug, to Fri 25 Aug | Port Ellery's roof and the test van; the Fishbowl on Friday 25 August | Tomás points the test van's 1.2 m dish at Port Ellery's uplink and reads Pipit-2's commands (S-band, 2065 MHz) off the air in plain text. Cirrine's condition starts on Friday 25 August | Cirrine-1's first sealed command, `PING`, goes up at 07:30 UTC on Friday 25 August, on pass `cirrine-1/YST/15287` (Ystmark, LOS 07:41:12), and Cirrine-1 answers |
| 7 | Mon 28 Aug | Yusra's desk | Hollis, Monday 09:02: "FURROW-A 0827 1402 NEVER ARRIVED". dlinkd says `furrow-a-20280827-1402.tif` (booking `bk-5274`) landed at 14:13:05 on Sunday; Furrowlight's receiver says nothing. Season 3's webhook signature can't settle it: Furrowlight holds the same secret, so either side could have written it | Yusra: "TL;DR: receipts nobody but Apogee can sign and anyone can check." Receipts start on Friday 1 September; that day Furrowlight checks Sunday's receipt against the file already in its own storage: it matches. The file had arrived; their receiver missed the event |
| 8 | Mon 4 Sep | Ystmark, over the backhaul | Sunday 3 September a storm cuts Ystmark's fiber, and the station falls back on the slow satellite backhaul, a link another company runs. Sunniva: "Fiber cut. Backhaul up. Everything on it readable." | This week's command file for Ystmark, and next month Rotation Day's keys, must cross a link whose operator can read and change anything. On Thursday 7 September the file crosses sealed under a key the two ends agreed over the open link |
| 9 | Mon 11 Sep, to Fri 22 Sep | Mirela's desk; Salar Quieto; the safe | Mirela's audit: anything that reaches a station's network can talk to its uplinkd and its dish controllers. On Tuesday 12 September Paz finds an old test script on a spare PC still asking Chiquita's controller for its position every minute: harmless, and nobody knew | Every station and server gets a certificate from Apogee's own certificate authority. The root is made on Monday 18 September and locked in the safe; every certificate is issued by Friday 22 September |
| 10 | Mon 25 Sep, to Fri 6 Oct | Port Ellery and every station | Tomás shows that uplinkd's traffic to Lantern Atoll crosses another company's network as plain HTTP. Mirela: TLS 1.3 everywhere, both ends proven by Apogee's certificates | Rotation Day's keys will travel this path. By Friday 6 October every station speaks TLS 1.3 with both ends proven, and a sealed Cirrine-1 command crosses it to Ystmark untouched |
| 11 | Fri 6 Oct, to Sat 7 Oct | Yusra's desk; Brinewatch's maintenance window | Brinewatch can install Rotation Day's new API key and webhook secret only in its monthly window (the first Saturday of the month, 06:00 to 08:00 UTC) and only the way its older systems allow. Its first two secrets went out XORed with a 256-byte pad its courier brought in November 2027, both from the pad's first byte | Miss the window, and on Rotation Day Brinewatch is locked out of the booking service or keeps an old secret. The pad is retired; Brinewatch's 2021 RSA key (1,024 bits) is refused; its new 3,072-bit keys arrive by 17:30; at 06:00 on Saturday both secrets go out locked to its receiving key, and Brinewatch opens them |
| 12 | Mon 9 Oct; Rotation Day Tue 17 Oct; overlap ends Tue 24 Oct | the Fishbowl | Ingrid runs Rotation Day like a pass, from the Cirrine-1 pass that crosses midnight to Furrow-A over Salar Quieto at 14:02 | Every key in the network replaced, at every station, without missing a single pass. At 10:12 the Cape Wren contractor reports that the old station computer she took away for recycling still holds its disk, and with it the key of certificate 1005 |
| 13 | Wed 25 Oct, to Fri 3 Nov | Yusra's desk, Mirela's desk, the Fishbowl | Cirrine's yearly review passes Rotation Day with one question, which Yusra forwards: could someone holding one engineer's console password stow a dish at Ystmark during Cirrine-1's pass? It could | From Friday 3 November at 16:00 the console asks for a password and a code from the engineer's phone, or one recovery code; reset links keep only a digest, last 30 minutes and work once; ten wrong codes in a row lock an account's codes until a reset |

The opening (Monday 10 July) hands over the new badge, shows Cirrine's three conditions and
Rotation Day's date, introduces Mirela (this season's reviewer), Yusra and Tomás, and sends the
learner to Mirela's question about random numbers. Words of the first week: key, nonce,
fingerprint, the ops console.

The finale (chapter 12, Tuesday 17 October): ten passes flown, none missed; the Cirrine-1
command sealed with `k-2028-07` at 23:59:58 reaches the satellite at 00:00:03 and is accepted;
every station on its new certificate by 01:15:40; certificate 1005 refused from 10:12; on 24
October `k-2028-07` retires and the second revocation list names every old certificate. Mirela
signs the log: "Every key replaced. Every pass flown. It failed closed every time we asked it
to." The debrief points at Season 6: in November the learner carries the pager again, and
Lantern Atoll's generator has started to sound different.

The epilogue (chapter 13, Wednesday 25 October to Friday 3 November) closes the season after its
climax. Wednesday 25 October: the reset file holds every token as it was sent. Thursday 26 October: its oldest link
(March 2027) still works, a link Nkem used in May works again, and the reset form answers "no account
named" to strangers. Tuesday 31 October: Ingrid's phone runs 20 seconds slow. Wednesday 1 November:
Mirela reads a code over Nkem's shoulder, and Ingrid asks why six digits are fine now. Thursday 2
November: Sunniva's phone falls into a snowdrift at Ystmark (snow 40 cm; "Found in April, maybe"),
and Sunniva signs in from the sheet. Friday 3 November, 16:00: the console asks every engineer for a
password and a code; Mirela's password without her phone is refused ("LGTM. It failed closed.");
Yusra to Cirrine: "TL;DR: a stolen password alone now opens nothing." The debrief ends where Season 6
begins: on Monday at 08:00, Mirela hands the learner the pager.

#### New systems and things

- **The ops console**, Apogee's internal web tool where engineers log in to stow a dish, run a
  pass by hand or read a station's log. It is older than the learner. Its user file holds one
  password record per account (chapter 4).
- **uplinkd's station side.** uplinkd runs at headquarters, where it signs, seals and schedules
  commands, and at every station, where it queues command lines for the next pass and hands
  them to the dish. In readmes: "Ystmark's uplinkd". No new service name.
- **Apogee's certificate authority**: a root certificate kept offline in the safe at Port
  Ellery and an intermediate on a machine in the Fishbowl (chapter 9).
- **The backhaul** (already in the bible: Ystmark's satellite link when the fiber is cut) is
  run by another company, unnamed; anything on it may be read or changed.
- **Receipts**: bookingd signs one for every delivered file; customers fetch it with
  `GET /bookings/{id}/receipt` (token, owner only; 404 `no receipt for bk-5274 yet` before the
  file lands) and check it with the public keys in the customer guide (chapter 7). No lesson
  runs this route; lessons sign and check envelopes as bytes.
- **Brinewatch's older systems** (chapter 11): its security office (unnamed; it speaks through Yusra) takes secrets only on its pad or locked to an RSA key, signs with RSA, and changes its systems only in its monthly maintenance window, the first Saturday of the month, 06:00 to 08:00 UTC. The pad, 256 random bytes on a USB stick its courier brought in November 2027 (one copy at each end), carried its API key that November and its webhook secret in January 2028, both XORed from the pad's first byte; it is retired on Friday 6 October 2028, and Mirela wipes Apogee's copy on the 7th. Its 2021 RSA key has 1,024 bits and is refused; its new receiving key (OAEP) and signing key (PSS), 3,072 bits each, are made that Friday, and Yusra reads the signing key's whole fingerprint back to Brinewatch by phone.

#### Every key, its size and its encoding

Hex is always lowercase. "Random" always means `crypto/rand`; nothing secret ever comes from
`math/rand` or `math/rand/v2`. Secret key files have mode 0600 and are written whole through a
temporary sibling and a rename (Season 2's way).

| Key | Made from | Size | Written as | Lives in |
|---|---|---|---|---|
| Uplink master `k-YYYY-MM` | random | 32 bytes | 64 hex digits in `uplink-keys.txt` | uplinkd at headquarters only |
| Satellite sign key | HKDF of a master | 32 bytes | never written; LOAD_KEY carries it sealed | the satellite; uplinkd derives it when needed |
| Satellite seal key | HKDF of a master | 32 bytes (AES-256) | as above | as above |
| `k-2027-03` (Season 1) | as made in 2027 | as made in 2027 | — | uplinkd; retired 2028-08-11T00:00:00Z; never in `uplink-keys.txt` |
| Command nonce | random, fresh for every seal | 12 bytes | 24 hex digits in the SEAL line | the line itself |
| Password salt | random, per record | 16 bytes | 32 hex digits in the record | the ops console's user file |
| Password key | PBKDF2-HMAC-SHA256, 600,000 iterations | 32 bytes | 64 hex digits in the record | the user file |
| Session token | `rand.Text()` | 26 characters `A`–`Z`, `2`–`7` (at least 128 bits) | the token; the server keeps only its SHA-256 in hex | the engineer's browser |
| Receipt key `r-YYYY-MM` | random 32-byte Ed25519 seed | public 32 bytes, signature 64 | seed: 64 hex in bookingd's key file; public: 64 hex in `receipt-keys.txt` and a PEM `PUBLIC KEY` block in the customer guide; signature: 128 hex | bookingd (seed); everyone (public) |
| Station identity (backhaul) | random Ed25519 seed | public 32 bytes | 64 hex | each station; the other end knows the public key in advance |
| Backhaul share | `ecdh.X25519().GenerateKey`, fresh per session | 32 bytes | 64 hex in HELLO | never stored |
| Backhaul session key | HKDF of the X25519 secret | 32 bytes | never written | memory, one session |
| CA and certificate keys | random Ed25519 seed | public 32 bytes | PEM: `CERTIFICATE` blocks; private keys as PKCS #8 `PRIVATE KEY` blocks, mode 0600 | root: offline in the safe; intermediate: the CA machine; each leaf: its own server |
| Brinewatch's pad (to 6 Oct 2028) | random | 256 bytes | on a USB stick | Brinewatch and Apogee, one copy each; Apogee's wiped 7 Oct |
| Brinewatch's RSA keys (receiving, signing) | Brinewatch's own, 6 Oct 2028 | 3,072 bits; ciphertexts and signatures 384 bytes | PEM `PUBLIC KEY` blocks (PKIX); fingerprint `sha256:` + 64 hex of the PKIX bytes | Brinewatch (private); Apogee (public) |
| bookingd token key, webhook secrets, API keys (Season 3) | random | as Season 3 | as Season 3 | replaced on Rotation Day; Brinewatch's new API key and webhook secret reach it on 7 October locked to its RSA receiving key (chapter 11); no lesson covers the others |
| Reset token | `rand.Text()` | 26 characters (130 bits) | in the emailed link; the console keeps its SHA-256 in hex and an expiry 30 minutes on | the link; one pending reset per account |
| Code secret (TOTP) | random | 20 bytes (160 bits) | set up once in the engineer's authenticator app | the phone and the console, whole (it can't be a digest) |
| Recovery code | `rand.Text()` | 26 characters, eight a sheet | printed once; the console keeps each one's SHA-256 in hex | the engineer's sheet |

Ed25519 keys are always Ed25519 (`crypto/ed25519`). The only RSA keys are Brinewatch's (chapter 11), and no lesson generates one: the tests' four were made once with `rsa.GenerateKey` (Go 1.27.3) and are embedded as PKCS #8 PEM. ML-KEM-768
sizes (chapter 8): encapsulation key 1,184 bytes, ciphertext 1,088 bytes, shared key 32 bytes,
decapsulation key seed 64 bytes.

#### Test keys and fixed vectors

Lessons never contain flight keys. Test keys are readable on purpose. Symmetric test keys are
ASCII text of the right length; every test seed or private key of Ed25519 or X25519 is the
SHA-256 of a label (`sha256.Sum256([]byte(label))`). Every value below was computed with
Go 1.27.3.

| Test key or value | Exactly |
|---|---|
| Season 1's lab bench HMAC key | ASCII `lab bench key, never flown`; `pipit-2\|118\|SET_MODE\|SAFE` signs as `512545bcb43a5bb6ccb9f9c8511f4cc8c5fedeb38c226376af5b5aaae60e3b4d` |
| Master `k-2028-07` | ASCII `uplink test master key, 32 bytes` (hex `75706c696e6b2074657374206d6173746572206b65792c203332206279746573`) |
| Master `k-2028-10` | ASCII `uplink next master key, 32 bytes` (hex `75706c696e6b206e657874206d6173746572206b65792c203332206279746573`) |
| HKDF salt (production too: it is public) | ASCII `apogee uplink keys v1` (21 bytes) |
| Derived from `k-2028-07` | `sign\|pipit-2` `82f945638cbdd9a034dfab2bc578fd7a720d73e5701d3abcf561cd825847da7c`; `sign\|furrow-a` `00a1988ac17a102e03bf1e93099ecfeb710b18ec76e69dac6e30251a674a8dd3`; `sign\|cirrine-1` `18668bfe29024937ff84a61df34e8a0310a4702a74834b969f56e2c6ccd00213`; `seal\|cirrine-1` `427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171`; `seal\|pipit-2` `53fdd2ccce4f186b0bdeb36626b6e8591ac1c5f59a00a80956c348dbdbaa2c3a` |
| Derived from `k-2028-10` | `sign\|pipit-2` `514e80107e64a13250a9bf9f8c5b1370995bea12a047f057a351c5c11f7522ad`; `sign\|cirrine-1` `3c280cbe585cfec33b8ac6693e6e0802fc0c532efc2af548cb1b3a85054418ba`; `seal\|cirrine-1` `a70fbff3ac2b05cbf0a66c3fc0eb76ceb1ec55269766703f715f4ec43c78663d` |
| Fixed nonce (vectors only; production nonces are always random) | hex `a1b2c3d4e5f60718293a4b5c` |
| Sealed `PING`, `cirrine-1`, seq 4471, `k-2028-07`, fixed nonce | `SEAL 2 cirrine-1 4471 k-2028-07 a1b2c3d4e5f60718293a4b5c c1f2a9798d38081d764bbcfb5e35e2598d2584a4` |
| Password fixture | password `fails closed, always`, salt ASCII `mirela test salt` (hex `6d6972656c6120746573742073616c74`); 600,000 iterations → key `6d91ceab779ff574b23820b133a740fa76d8b6a0deb5d6f0da2f94a6508dd938`; 1,000 iterations → `5b893367890cac1547d6c07d906a4b46b8a89733bf1f7fb3d7ba24f8e16174b3` |
| Old console record | `kiosk`: `sha256$1c4bb463c52d458ae6e9225430c45c753c5ae417bb4e7e4530773aaa0cda201b` (the password `fishbowl`) |
| Receipt keys | seed = SHA-256 of `receipts r-2028-08` → public `ba6c6226d2c3c056684e9d9a5e489db44a148e27e3fcb8f2ba16a6bddb0a86b2`; of `receipts r-2028-10` → public `3c3ae60d0f3d9f2f20ae8d68ee3b8996b233ed694c73b20f9091e3d14a494f12` |
| Station identities | seeds = SHA-256 of `identity PEL`, `identity YST`, `identity MID` (the one in the middle) → public `1497e6fc165f9bcd9fa8540b6a74089da8e02b669244eb2c9386fdd80588bcae`, `cacb8869a23066fe505e739fc60a7246f839d6a25e9994c1cd9a073de7cdd355`, `90aa0565e70d224291d92e937c4227ed71c6806cbda864e28bc7a6f9e8713b95` |
| Backhaul shares (vectors) | X25519 private keys = SHA-256 of `ephemeral PEL`, `ephemeral YST`, `ephemeral MID` → shares `a10c608ed294ad1d437b121ce63f146ee081f200b8b71126bc8f5642e69aa93e`, `b6f8a5a974e292828807a52f32334ac3925ea872ca82de318f53e971fafc5a66`, `1c85634a8ffb99861b81aa64cdc22c9f386010f5795f8212d89551372b360021`; PEL with YST agree on `c50d64c87ddd744d6c24842274ac98d74c33be2bdadd89463eae753cae91cc7a`; session key `24e57dc80a28df831da79d5b65636ba94ed13a4274908279576a15428a06958e` |
| CA and certificate keys | seeds = SHA-256 of `ca root 2028`, `ca stations 2028`, `cert <name> <serial>` (such as `cert yst.apogee.internal 1001`, `cert uplinkd 2001`); impostors: `ca someone else`, `cert impostor` |
| Digests | SHA-256 of `Apogee ground calling Pipit-1` `fe5e0f73525ffc9f862167a80cb444fbc55de2a916041901a86964737b2fbd66`, of `...Pipit-2` `8cbc29efabfcd2fb143763fd987194d89474355d64403f55e0def74fab81e72f`; SHA3-256 of `...Pipit-1` `7f7e6c3e2e3a3f1120a59350ba37cf6b382ad4644e840a4e15a35a179d9584ea`; SHA-256 of nothing `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| Brinewatch's RSA test keys (PEM in each RSA lesson's `keys_test.go`, byte for byte the same) | receiving, 3,072 bits, `sha256:d2499efc1976897aa14ea5d3948c4524bb826ad42b6ac1b43c5772f6b2ab8c1e`; signing, 3,072 bits, `sha256:e47dfcb6d8f2b0e68d9cc9ebfc169c6dbe786e4930e5f12e74238ff24d65f508`; the 2021 key, 1,024 bits, `sha256:6998de4f989149257401c3037508a903a2a573483f60b3eeb980c9abe5d2b2fe`; the impostor's, 2,048 bits, `sha256:66eb315a02a52f08680f0a365829b8e5108739e09f78673b27c6bfee0ea5e87e` (SHA-256 of each PKIX public key) |
| Brinewatch's test pad | 256 bytes: the SHA-256 of `brinewatch pad 0` to `brinewatch pad 7`, one after another |
| Brinewatch's new secrets (test values) | API key `brinewatch.20281017a1b2c3d4e5f60718293a4b5c`; webhook secret ASCII `brinewatch webhook secret 2028-10` (the old ones are Season 3's) |
| Short codes (lesson 11.2) | texts `brinewatch <n>`: 2 bytes, 48 and 307 share `bc8c` (308 tries); 3 bytes, 20 and 4774 share `f9f338` (4,775 tries); matching text 0's code: 2 bytes at text 96636, 3 bytes none in 100,000 |
| TOTP vectors (13.4, 13.5) | key ASCII `12345678901234567890` (the RFCs'); RFC 4226 counters 0–9 → `755224 287082 359152 969429 338314 254676 287922 162583 399871 520489`; RFC 6238 SHA-1 at 59, 1111111109, 1111111111, 1234567890, 2000000000 and 20000000000 s → `287082 081804 050471 005924 279037 353130` (the RFC's eight digits, cut to six) |
| Phone test secrets (chapter 13) | 20 ASCII bytes each: `mirela's phone, 2028`, `ingrid's phone, 2028`, `nkem's phone in 2028`, `sunniva's phone 2028`; Mirela's codes for the steps around 2028-10-30T10:00:00Z (step 61883760): `462810`, `687056`, `173810` |
| Chapter 13's passwords | `fails closed, always` (Mirela), `wind 64, gusts 81` (Sunniva), `fishbowl` (kiosk); after the test reset `two proofs, always`; 13.9's tests derive with a 1,000-round stand-in for PBKDF2 |

Since Go 1.26, `ecdh` key generation ignores a reader passed to it, so a fixed "random"
reader never makes fixed keys: fixed keys come from `NewKeyFromSeed` (Ed25519) and
`ecdh.X25519().NewPrivateKey` (X25519). No test uses `testing/cryptotest.SetGlobalRandom` to
expect fixed output.

#### Random values (chapter 1)

- The helpers in `random.go`, used by every later chapter: `newKey(id string) Key` (32 random
  bytes), `newNonce() []byte` (12), `newSalt() []byte` (16), `newSessionToken() string`
  (`rand.Text()`). `crypto/rand.Read` never returns an error and always fills its slice; if the
  system can't supply randomness, the program stops. Code calls it without checking a result.
- A secret has at least 128 bits: length × log2(alphabet size), when every symbol is picked
  at random. Six digits are 19.9 bits;
  `rand.Text()` is 26 base32 characters, 130 bits.
- Random 12-byte nonces stay safe for 2^32 seals under one key; Apogee's busiest key seals a
  few thousand a year. A nonce counter kept in memory starts again after a restart, so Apogee
  never uses counters for nonces.

#### Digests (chapter 2)

- A **digest string** names its hash: `sha256:<64 hex>` or `sha3-256:<64 hex>`. Nothing else
  is accepted: no capitals, no other hash, no bare hex. `parseDigest` errors:
  `unknown hash "md5"`, `digest must be 64 lowercase hex digits`, `digest must be <hash>:<hex>`.
- Files are hashed as streams (`io.Copy` into the `hash.Hash`), never read whole: a Furrow-A
  image is 10.9 GB. Digests of public files are compared with `==` or `bytes.Equal`;
  constant-time comparison is for secrets.
- The **delivery manifest** (`manifest.txt` beside the files): one line per file,
  `<digest> <bytes> <file>`, single spaces; file names never contain spaces. The check reports
  one line per manifest line, in order: `ok <file>`, `wrong size <file>` (checked before
  hashing), `changed <file>`, `missing <file>`. A line without exactly three fields, a bad
  digest or a bad size is an error, `manifest line 3: want <digest> <bytes> <file>`.
- **The forged field report** (lessons 2.2 and 2.4): the original is the 20 bytes
  `NDVI field 12: 0.71\n` (hex `4e445649206669656c642031323a20302e37310a`, SHA-256
  `72198855a7d55dd3cbd5fff4d06168339fa8d32086a6cd18899400210b610aae`). Mirela's forgery is
  `NDVI field 12: 0.17\n` followed by the four bytes `1a 15 bc d1` (hex
  `4e445649206669656c642031323a20302e31370a1a15bcd1`, SHA-256
  `4b402892e69f9ddac861d7107795c4a44f68499ae7f4bc0d1ee1acea71093b8f`). Both have CRC-32
  `f6ced6c0`.

#### Signed commands, hardened (chapter 3)

Season 1's formats are unchanged: the canonical command `pipit-2|118|SET_MODE|SAFE`, the wire
line `CMD pipit-2 118 SET_MODE SAFE sig=<64 hex digits>`, HMAC-SHA256, sequence numbers
strictly increasing per satellite. Season 5 adds rules, never a new layout:

- **Canonical fields** (`canonical(c Command) (string, error)`): the satellite ID is 1 or more
  lowercase letters, digits and hyphens; `Seq` is at least 1; the operation is 1 or more
  capital letters and `_`;
  each argument is 1 or more printable ASCII characters from `!` to `~` other than `|`. Any
  other field is refused (`argument 1 contains "|"`), so no two commands share a canonical form.
- **Order of checks** for a CMD line: parse the line (exactly one `sig=` field, last, 64 hex
  digits); build the canonical form; verify the HMAC with `hmac.Equal`; only then compare
  `Seq` with the last accepted; only then record it. A refusal changes nothing.
- **The homemade signature** Mirela finds in Pipit-3's flight software is `sha256(key ||
  canonical)`. Lesson 3.1's read-only test forges it without the key by SHA-256 length
  extension: it rebuilds a hash state from a tag (the 4 bytes `sha\x03`, the 32 tag bytes, 64
  zero bytes, then the big-endian count of bytes hashed so far including padding) with
  `encoding.BinaryUnmarshaler` on `sha256.New()`, writes `|REBOOT`, and sends the command with
  SHA-256's padding and `|REBOOT` appended. HMAC refuses that forgery.
- **One HMAC per message.** A `hash.Hash` keeps everything written to it; signers start a new
  `hmac.New` for every message, which also makes them safe for many goroutines at once.
- **Comparing secrets.** MACs: `hmac.Equal` on bytes. Secrets of different lengths (a typed
  phrase): compare their SHA-256 digests with `subtle.ConstantTimeCompare`, because it returns
  at once when lengths differ.
- Errors, checked with `errors.Is`: `ErrForged` (`command failed authentication`, never more
  detail) and `ErrReplay` (`replayed sequence number: 118, last 118`), shared with sealed
  commands below.

#### The ops console's passwords (chapter 4)

- **Record**: `pbkdf2-sha256$<iterations>$<32 hex salt>$<64 hex key>`, from
  `pbkdf2.Key(sha256.New, password, salt, iterations, 32)` with a fresh 16-byte salt. New
  records use 600,000 iterations. A record with a count from 1,000 to 10,000,000 is checked
  with its own count; anything else, any other algorithm, or a salt or key of the wrong length
  doesn't parse. Mirela's: `pbkdf2-sha256$600000$6d6972656c6120746573742073616c74$6d91ceab779ff574b23820b133a740fa76d8b6a0deb5d6f0da2f94a6508dd938`.
- **Old records** (before 4 August): `sha256$<64 hex>`, the unsalted SHA-256 of the password.
  Lesson 4.1's intermediate step, `salted$<32 hex salt>$<64 hex of sha256(salt + password)>`,
  is lesson-only and never stored by the console.
- **The user file** `console-users.txt`: one line per account, `<user> <record>`; user names
  are lowercase, such as `mirela`, `kiosk`, `bench`.
- **Login rules**: derive again with the record's own count and salt and compare with
  `subtle.ConstantTimeCompare`; an unknown user and a wrong password give the same error,
  `invalid user or password` (`errInvalidLogin`), after the same work: one derivation against a
  fixed dummy record for an unknown user, and one for an old `sha256$` record too, which would
  otherwise answer at once and show which accounts still have one (`matches`, from 4.8). A
  record under 600,000 iterations still answers faster until it is upgraded. A successful login
  against an old record or one under 600,000 iterations saves a new record (upgrade at login);
  a session token is `rand.Text()`, and the console keeps only its SHA-256 in hex.
- **Cost in tests**: one derivation at 600,000 iterations takes about 60 ms on an Apple-silicon
  Mac. A lesson's tests run at most six of them; other cases use 1,000-iteration records.

```go
// Record is a parsed password record.
type Record struct {
	Iterations int    // 600000
	Salt       []byte // 16 bytes
	Key        []byte // 32 bytes
}

// Console is the ops console's login (lessons 4.7 and 4.8).
type Console struct {
	records  map[string]string                                    // user → record text
	kdf      func(pw string, salt []byte, iter int) ([]byte, error) // PBKDF2-HMAC-SHA256; tests count calls
	save     func(user, record string) error                      // stores an upgraded record
	sessions map[string]string                                    // hex SHA-256 of a token → user
}
```

The dummy record an unknown user is checked against is a fixed, valid 600,000-iteration record
whose password nobody knows; the check always fails after one `kdf` call, like a wrong password.

#### Keys from keys (chapter 5)

- A satellite key is `hkdf.Key(sha256.New, master.Secret, []byte("apogee uplink keys v1"),
  purpose+"|"+satID, 32)`, with purpose `sign` (its HMAC key for CMD lines) or `seal` (its
  AES-256-GCM key for SEAL lines). One master, many unrelated keys: a satellite holds only its
  own two keys for each master version it knows. HKDF is for keys that are already random; a
  password is never an HKDF input (that is PBKDF2's job).
- **`uplink-keys.txt`** (uplinkd only, mode 0600): `#` lines and blank lines are ignored; each
  other line is `<id> <activate> <retire or -> <64 hex secret>`, optionally followed by the word
  `revoked`; times are RFC 3339 UTC. A duplicate ID, a bad time or a secret that isn't 64 hex
  digits is an error naming the line. In August the file holds one line:
  `k-2028-07 2028-08-11T00:00:00Z - 75706c69...6573`.
- **The fleet** for derivations: `pipit-1`, `pipit-2`, `pipit-3`, `furrow-a`, `furrow-b`,
  `brine-3`, `cirrine-1` (Thistle re-entered in 2027).
- **Printing keys.** `Key` has `String()` returning `k-2028-07 (secret hidden)` and
  `GoString()` returning `Key{ID: "k-2028-07", Secret: hidden}`. fmt can't call methods on a
  value in an unexported field and prints its bytes, so `KeyRing` has its own `String()`:
  `key ring: k-2028-07, k-2028-10`, IDs in activation order, and a `GoString()` with the same
  text, since `%#v` calls GoString. The number verbs (`%d`, `%o`, `%b`, `%c`, `%U`) skip both
  methods and still print a key's bytes, so keys are logged with `%v` (lesson 5.5).

#### Sealed commands (chapter 6)

Uplink format 2 (format 1 is Season 1's signed CMD line). Every satellite's flight software
opens SEAL lines since August; Cirrine-1's commands are always sealed from 25 August; other
operators' routine commands may stay CMD lines; LOAD_KEY is always sealed.

```go
// Sealed is a command sealed for one satellite: only its header is readable.
type Sealed struct {
	Version    int    // 2
	SatID      string // "cirrine-1"
	Seq        uint64 // as Command.Seq: strictly increasing per satellite
	KeyID      string // "k-2028-07": the master version its key comes from
	Nonce      []byte // 12 bytes from crypto/rand, new for every seal
	Ciphertext []byte // AES-256-GCM of "OP|ARG...", its 16-byte tag included
}
```

- **The line**: `SEAL 2 <sat> <seq> <key id> <nonce: 24 hex> <ciphertext: hex>`, seven fields
  separated by single spaces. The ciphertext is 17 to 272 bytes (34 to 544 hex digits): a
  plaintext of 1 to 256 bytes plus the tag. Readmes show it shortened: `SEAL 2 cirrine-1 4471
  k-2028-07 a1b2…4b5c c1f2…84a4`.
- **The key**: the satellite's `seal` key from the named master (chapter 5).
- **Associated data**: the header as text, `SEAL|2|cirrine-1|4471|k-2028-07`, authenticated
  but not hidden, so the station can schedule the command without reading it, and nobody can
  move a sealed command to another satellite, number or key.
- **Plaintext**: the operation and arguments joined by `|`, checked by chapter 3's canonical
  rules: `PING`, `SET_MODE|SAFE`. After opening, the full canonical command is
  `cirrine-1|4471|PING`.
- **Order of checks** (the satellite and its stand-in in tests): parse (`ErrMalformed`, also
  for a line that doesn't start with `SEAL`); version 2 (`ErrVersion`, only for a SEAL line);
  satellite is its own (`ErrWrongSat`); key ID held (`ErrUnknownKey`) and not revoked
  (`ErrRevoked`) or retired (`ErrRetired`); `Seq` above the last accepted (`ErrReplay`); `Open`
  with the associated data (`ErrForged`); only then record `Seq`. A
  refusal changes nothing, so a forged line numbered 18,446,744,073,709,551,615 can't lock the
  satellite out. (Chapter 3 verifies before it compares `Seq`; here a number already seen is
  refused before any work is spent opening it, and 6.7's readme says so.)
- **Errors**, all `errors.New` sentinels checked with `errors.Is`, wrapped with context:

  ```go
  var (
  	ErrMalformed = errors.New("malformed command line") // "malformed command line: want 7 fields"
  	ErrVersion   = errors.New("unsupported version")    // "unsupported version: 3"
  	ErrWrongSat  = errors.New("not for this satellite") // `not for this satellite: "pipit-2"`
  	// ErrUnknownKey means the ring holds no version with the key ID asked for.
  	ErrUnknownKey = errors.New("unknown key")
  	ErrRevoked    = errors.New("key revoked") // `key revoked: "k-2028-10"`
  	ErrRetired    = errors.New("key retired") // `key retired: "k-2028-07"`
  	// ErrReplay refuses a command whose sequence number isn't above the last one accepted.
  	ErrReplay = errors.New("replayed sequence number")
  	// ErrForged refuses a command whose signature doesn't check out, and never says more.
  	ErrForged = errors.New("command failed authentication")
  )
  ```

  `ErrUnknownKey`, `ErrReplay` and `ErrForged` keep the text of their first declarations
  (chapters 5 and 3). Wrapped, they read `unknown key: "k-2027-03"` and
  `replayed sequence number: 4471, last 4471`.

- **LOAD_KEY** (chapter 12): the plaintext `LOAD_KEY|<new key id>|<64 hex sign key>|<64 hex
  seal key>` (148 bytes for `k-2028-10`), sealed under the satellite's current seal key. After
  opening it, the satellite accepts the new version at once; it signs nothing, so "activate"
  doesn't apply to it.
- The station never holds a satellite's seal key: Ystmark's uplinkd queues Cirrine-1's SEAL
  lines without being able to read them.

#### Receipts (chapter 7)

```go
// Receipt is Apogee's signed word that a booking's file was delivered.
type Receipt struct {
	Version   int       `json:"v"`          // 1
	KeyID     string    `json:"key_id"`     // "r-2028-08": the key that signed it
	BookingID string    `json:"booking_id"` // "bk-5274"
	SatID     string    `json:"sat_id"`
	Station   string    `json:"station"`
	File      string    `json:"file"`       // "furrow-a-20280827-1402.tif"
	Bytes     int64     `json:"bytes"`
	Digest    string    `json:"digest"`     // "sha256:<64 hex>", or "sha3-256:" for Cirrine
	LandedAt  time.Time `json:"landed_at"`  // from the Landing (Season 3)
	SignedAt  time.Time `json:"signed_at"`  // bookingd's clock, truncated to the second
}

// Envelope carries a receipt and its signature.
type Envelope struct {
	KeyID   string          `json:"key_id"`  // the same as the receipt's
	Receipt json.RawMessage `json:"receipt"` // the exact bytes that were signed
	Sig     string          `json:"sig"`     // 128 hex digits: Ed25519 over Receipt's bytes
}

// ReceiptKey is one published receipt key and its signing window.
type ReceiptKey struct {
	ID       string            // "r-2028-08"
	Public   ed25519.PublicKey // 32 bytes
	Activate time.Time         // receipts signed from here...
	Retire   time.Time         // ...until here (zero: still signing)
	Revoked  bool
}

// Notary is bookingd's receipt signer.
type Notary struct {
	KeyID string
	priv  ed25519.PrivateKey
}
```

- **Signing**: `json.Marshal(receipt)` (no newline, no indentation) is signed with
  `ed25519.Sign`; the envelope is `json.Marshal(Envelope{...})` with those exact bytes as
  `receipt`. Example receipt bytes:
  `{"v":1,"key_id":"r-2028-08","booking_id":"bk-5274","sat_id":"furrow-a","station":"SQT","file":"furrow-a-20280827-1402.tif","bytes":10874220000,"digest":"sha256:<64 hex>","landed_at":"2028-08-27T14:13:05Z","signed_at":"2028-09-01T09:00:00Z"}`.
  The archived file is `<file>.receipt.json`, the envelope as sent.
- **Checking**, in this order: decode the envelope strictly (unknown fields refused; 64 KiB at
  most) else `ErrBadEnvelope`; find the key by the envelope's `key_id` in your own published list,
  never a key the envelope brings, else `ErrUnknownSigner`; refuse a revoked key
  (`ErrRevoked`); `sig` must be 128 hex digits and verify over the raw `receipt` bytes, else
  `ErrBadSignature`; only then decode the receipt strictly; its `key_id` must equal the
  envelope's (`ErrBadEnvelope`); `signed_at` must fall in `[Activate, Retire)` of that key,
  else `ErrOutsideWindow`. Checking a file: the same size, and the digest of the file, hashed
  with the receipt's own algorithm, equal to the receipt's, else `ErrWrongFile`.
- `ed25519.Verify` panics on a public key that isn't 32 bytes: check the length first.
- **`receipt-keys.txt`** (published in the customer guide, beside each key's PEM block):
  `<id> <activate> <retire or -> <64 hex public key>`, `revoked` appended when revoked:

  ```
  r-2028-08 2028-09-01T00:00:00Z 2028-10-17T00:00:00Z ba6c6226...86b2
  r-2028-10 2028-10-17T00:00:00Z - 3c3ae60d...4f12
  ```

- Errors: `ErrBadEnvelope` (`malformed receipt`), `ErrUnknownSigner` (`unknown receipt key`),
  `ErrBadSignature` (`receipt signature does not verify`), `ErrOutsideWindow` (`receipt signed
  outside its key's window`), `ErrWrongFile` (`file does not match receipt`), and `ErrRevoked`
  from the sealed commands.
- `newNotary(keyID string, seed []byte) (*Notary, error)` refuses a seed that isn't 32 bytes
  (`ed25519.NewKeyFromSeed` would panic).

#### The backhaul handshake (chapter 8)

Port Ellery headquarters (`PEL`) dials; the station (`YST`) answers. Each side holds its own
Ed25519 identity (`Identity{Code string; Key ed25519.PrivateKey}`) and knows the other's
public identity key in advance (`map[string]ed25519.PublicKey`, by station code).

1. Each side makes a fresh X25519 key for this session: `ecdh.X25519().GenerateKey(rand.Reader)`.
2. Each sends one line, PEL first: `HELLO <code> <64 hex share> <128 hex signature>`. The
   signature is the side's identity key over the bytes `backhaul share|<code>|<64 hex share>`.
   A line is checked before its share is used: exactly four fields, a known code that isn't
   your own, 64 lowercase hex digits, a signature that verifies, else `ErrBadHello` (`bad hello`,
   wrapped: `bad hello: unknown station "MID"`, `bad hello: signature does not verify`).
3. Both compute the X25519 secret (an all-zero result is an error, never a key) and derive
   `hkdf.Key(sha256.New, secret, nil, "apogee backhaul v1|<PEL share hex>|<YST share hex>", 32)`:
   the dialer's share first.
4. Data lines: `DATA <24 hex nonce> <hex ciphertext>`, AES-256-GCM under the session key with
   associated data `DATA|<from code>|<to code>`, a fresh random nonce each. A line that doesn't
   open ends the session with an error. DATA lines carry no sequence number, so a machine in
   the middle can drop, repeat or reorder them within a session, though never read or change
   one; the CMD lines inside carry their own numbers. The HKDF info binds the two shares, not
   the station codes, which is sound because each share's signature covers its code.
5. Nothing of the session is kept afterwards: a stolen station later can't open a recording.

Tests run both sides over `net.Pipe()` with a deadline of 1 s; the side in the middle is the
test, holding the `identity MID` key.

#### Apogee's certificate authority (chapter 9)

All keys Ed25519; all times UTC; serial numbers fixed in tests (in production, random 128-bit
from `crypto/rand`, with `rand.Int`).

| Certificate | Subject (CN, O=Apogee) | Issuer | Serial | Valid | Shape |
|---|---|---|---|---|---|
| Root | `Apogee Root CA 2028` | itself | 1 | 2028-09-18T00:00:00Z to 2038-09-18T00:00:00Z (10 years) | CA, `MaxPathLen` 1, key usage cert sign and CRL sign |
| Intermediate | `Apogee Stations CA 2028` | the root | 2 | 2028-09-18T00:00:00Z to 2030-09-18T00:00:00Z (2 years) | CA, `MaxPathLen` 0 with `MaxPathLenZero`, cert sign and CRL sign |
| Stations | `yst.apogee.internal`, `pel.`, `sqt.`, `lan.`, `cwr.` | the intermediate | 1001 YST, 1002 PEL, 1003 SQT, 1004 LAN, 1005 CWR | 2028-09-22T00:00:00Z for 90 days, to 2028-12-21T00:00:00Z | DNS name = CN; server auth; digital signature |
| uplinkd (headquarters) | `uplinkd`, DNS `uplinkd.apogee.internal` | the intermediate | 2001 | as stations | client auth; digital signature |
| passplan | `passplan`, DNS `passplan.apogee.internal` | the intermediate | 2002 | as stations | client auth (it may read status, never queue) |
| dishd (Port Ellery) | `dishd`, DNS `dishd.apogee.internal` | the intermediate | 2003 | as stations | client auth (it dials Ystmark's dish controllers) |
| Rotation Day's | the same names | the intermediate | 1101–1105, 2101–2103 | 2028-10-09T00:00:00Z to 2029-01-07T00:00:00Z | as above |

- Host names end in `.apogee.internal`; `.internal` is reserved for private networks, so no
  name can belong to anyone outside.
- A certificate is valid from `NotBefore` to `NotAfter`, **both instants included** (crypto/x509
  refuses only a time before `NotBefore` or after `NotAfter`).
- Files: `<name>.pem` holds the leaf then the intermediate (`CERTIFICATE` blocks);
  `<name>.key` the PKCS #8 private key (`PRIVATE KEY`), mode 0600; `roots.pem` the root.
  Names: `yst`, `pel`, `sqt`, `lan`, `cwr`, `uplinkd`, `passplan`, `dishd`. `tls.X509KeyPair`
  loads a pair.
- `Verify` always gets `Roots` (a pool with the root only), `Intermediates` (the chain the peer
  sent), `DNSName` where a name is checked, `KeyUsages` (server or client auth) and
  `CurrentTime` from the injected clock. Never leave `Roots` nil: on a Mac that hands the
  check to the system verifier, which doesn't know Apogee's root and refuses Ed25519
  certificates.
- The impostor in tests is a root also named `Apogee Root CA 2028` (key `ca someone else`)
  issuing `yst.apogee.internal`: same names, wrong key.
- `Issuer{Cert *x509.Certificate; Key ed25519.PrivateKey}` is a CA's certificate and key.

#### Every connection encrypted (chapter 10)

- **Both ends**: TLS 1.3 only (`MinVersion: tls.VersionTLS13`); the server presents its leaf
  and the intermediate; the client sets `RootCAs` (Apogee's root) and `ServerName`
  (`yst.apogee.internal`), since tests dial 127.0.0.1; the server sets `ClientAuth:
  tls.RequireAndVerifyClientCert` and `ClientCAs` (Apogee's root); both set `Time` to the
  injected clock. Go's default key exchange is the hybrid `X25519MLKEM768`; nobody sets
  `CurvePreferences`. `InsecureSkipVerify` never appears outside lesson 10.1's starter.
- **Who may do what**: the caller is the common name of `r.TLS.VerifiedChains[0][0]` (never
  `PeerCertificates`, which is only what the client sent). Any verified Apogee client may read
  status; only `uplinkd` may queue commands. A request without `r.TLS` is refused, never a panic.
- **Ystmark's uplinkd over HTTPS** (the station side; all answers `text/plain`):

  | Route | Answer |
  |---|---|
  | `GET /status` | 200 `ok YST` and a newline |
  | `POST /uplink` | body: one CMD or SEAL line, at most 1024 bytes, syntax checked only; 202 `queued cirrine-1 4471`; 403 `only uplinkd may queue commands`; 400 `bad command line` |

- **The dish controllers** speak their own raw socket protocol (Season 4), and dishd at Port
  Ellery is still the only program that talks to Ystmark's. This season wraps its connection to
  a controller in TLS (`tls.Dialer`, or `tls.Client` with `HandshakeContext`): dishd proves
  itself with its client certificate, 2003, and a controller presents its station's
  certificate, 1001 for Ystmark. In tests the controller is a stand-in on
  `net.Listen("tcp", "127.0.0.1:0")` behind `tls.Server` that echoes each line it reads,
  `PING yst` in, `PING yst` out; it never claims to be the real protocol.

#### Brinewatch's older keys (chapter 11)

- **The pad:** `Pad{key []byte; used int}`; `(*Pad).Seal` XORs with the next unused bytes and refuses with `pad has 7 bytes left, need 25`, spending nothing. A refusal changes nothing.
- **Short codes:** `shortCode(text, n)` is the hex of the first n bytes of SHA-256; Brinewatch's tool shows 3 bytes (six hex digits, 16.7 million codes).
- **`readRSA(pemText []byte) (*rsa.PublicKey, error)`:** one `PUBLIC KEY` block with only white space after it (`want one PUBLIC KEY block`), then `x509.ParsePKIXPublicKey` and a type switch. Errors: `want at least 2048 bits, got 1024`, `want an RSA key, got Ed25519`, `want an RSA key, got X25519`, `want an RSA key, got %T`.
- **OAEP:** `encryptTo` and `decryptWith` use `sha256.New()` and labels `apogee api key` and `apogee webhook secret`; a 3,072-bit key gives 384-byte ciphertexts and locks at most 318 bytes.
- **PSS:** `signRequest` and `checkRequest` work on the SHA-256 digest with `crypto.SHA256` and `nil` options; Brinewatch's tool uses `PSSSaltLengthEqualsHash`, which `nil` verification accepts.
- **The request:** `Request{Text, Sig, Key []byte}`, with Text `brinewatch handover sha256:<64 hex>` (`keyFingerprint`: `sha256:` + SHA-256 of `x509.MarshalPKIXPublicKey`).
- **The handover:** check the signature, read the key, compare the line exactly, then two lines `api-key <hex>` and `webhook-secret <hex>` (768 hex digits each). Any failure gives an error and no lines.

#### Rotation Day (chapter 12)

```go
// KeyVersion is one version of a key and the times it may be used.
type KeyVersion struct {
	Key                // Season 1's Key: ID and Secret
	Activate time.Time // signers and sealers use it from this instant
	Retire   time.Time // nothing it made is accepted from this instant; zero: not scheduled
	Revoked  bool      // refused at once, whatever the times; never undone
}
```

**The rules**, the same for every kind of key:

1. Verifiers accept a version from the moment it is installed, before its `Activate`, so
   nothing sealed after the switch arrives at a verifier that doesn't know the key yet.
2. A signer uses the newest version (by `Activate`) that has activated, isn't revoked and
   hasn't retired. With none, it refuses to sign: `no key to sign with at <RFC 3339 time>`.
3. Windows are half open: a version signs from `Activate` included and is refused from
   `Retire` included.
4. **Overlap**: an uplink master's `Retire` is its successor's `Activate` plus 7 days, so a
   command sealed with the old key just before the switch is still accepted after it.
   Receipt keys don't overlap in signing (the old one's `Retire` is the new one's `Activate`),
   and stay published for as long as their receipts are kept: a receipt is judged by its
   `signed_at`, a command by when it arrives.
5. **Revocation** wins over every window, at once, and is never undone; a revoked receipt key
   makes all its receipts unacceptable. Certificates are revoked through the intermediate's
   revocation list.
6. A CMD line names no key: during an overlap a satellite tries its keys newest first. SEAL
   lines name theirs.

**Revocation lists** (`x509.CreateRevocationList`, signed by `Apogee Stations CA 2028`):
`Number` counts from 1; `NextUpdate` is `ThisUpdate` plus 7 days; a list is current from
`ThisUpdate` to `NextUpdate`, both included. A checker refuses a list whose signature doesn't
verify against the intermediate (`CheckSignatureFrom`) or that isn't current, and then refuses
every certificate rather than trusting it (fail closed). The TLS check runs in
`VerifyConnection` (it also runs on resumed sessions) and answers
`certificate 1005 is revoked`.

**The timeline** (UTC):

| When | What |
|---|---|
| Fri 2028-07-14 16:00 | `k-2028-07` made (Activate 2028-08-11T00:00:00Z) |
| Thu 2028-08-10 | every satellite's own `k-2028-07` keys go up inside its operator's August flight-software update, which also teaches it to open SEAL lines |
| Fri 2028-08-11 00:00 | `k-2027-03` retires; satellites answer only to their own keys |
| Mon 2028-08-28 | `r-2028-08` made; published; Activate 2028-09-01T00:00:00Z |
| Fri 2028-10-06 | Brinewatch's pad retired; its new RSA keys made, and its new API key and webhook secret |
| Sat 2028-10-07 06:00 | Brinewatch's maintenance window: both new secrets handed over locked to its receiving key; they take effect on Rotation Day |
| Mon 2028-10-09 | `k-2028-10` made (Activate 2028-10-17T00:00:00Z); `k-2028-07` gets Retire 2028-10-24T00:00:00Z; `r-2028-10` published (Activate 2028-10-17T00:00:00Z) and `r-2028-08` gets Retire 2028-10-17T00:00:00Z; certificates 1101–1105 and 2101–2103 issued |
| Mon 9 to Fri 13 Oct | each satellite gets LOAD_KEY for `k-2028-10` on its first pass, sealed under its `k-2028-07` seal key |
| Tue 2028-10-17 00:00:00 | Rotation Day: uplinkd seals and signs with `k-2028-10`; bookingd signs receipts with `r-2028-10` |
| 00:30:00 | uplinkd's client certificate 2101 replaces 2001, and dishd's 2103 replaces 2003 (stations and controllers accept both until 24 October) |
| 01:00:00 to 01:15:40 | each station switches to its new certificate in its first quiet gap (below) |
| 10:12:00 | revocation list 1: certificate 1005 (Cape Wren's old station computer left with its disk) |
| Tue 2028-10-24 00:00:00 | `k-2028-07` retires; revocation list 2: 1001 to 1005 and 2001 to 2003 |

**Rotation Day's passes** (Season 1's `Pass`; lessons 12.7 and 12.8). Each pass gets one sealed
`PING`; the Cirrine-1 pass across midnight gets two, sealed at 23:59:58 and 00:01:00. A command
reaches its satellite 5 seconds after it is sealed.

| Satellite | Station | AOS | LOS |
|---|---|---|---|
| cirrine-1 | YST | 2028-10-16T23:52:10Z | 2028-10-17T00:03:40Z |
| pipit-3 | PEL | 00:58:30 | 01:07:10 |
| brine-3 | LAN | 01:02:00 | 01:10:40 |
| furrow-b | CWR | 02:51:20 | 03:00:44 |
| pipit-1 | PEL | 04:12:05 | 04:22:41 |
| pipit-2 | YST | 04:31:10 | 04:41:52 |
| furrow-a | CWR | 04:58:30 | 05:08:02 |
| cirrine-1 | YST | 07:30:00 | 07:41:12 |
| brine-3 | LAN | 09:14:20 | 09:22:31 |
| furrow-a | SQT | 14:02:10 | 14:11:52 |

**The switch rule** (`switchTime`): switches start from 01:00:00; a switch takes 10 minutes and
must keep 5 minutes clear of every pass at that station on each side, so it can't overlap
`[AOS − 5 min, LOS + 5 min)`. The first such start is `from` itself or 5 minutes after some
pass's LOS. Results: YST 01:00:00, SQT 01:00:00, CWR 01:00:00, PEL 01:12:10, LAN 01:15:40. No
fit within 24 hours of `from` is an error: `no 10-minute gap for CWR before
2028-10-18T01:00:00Z`.

**The finale's report** (`Report{Flown, Missed int; Lines []string}`), one line per pass in AOS
order: its AOS as `15:04:05`, the satellite, the station, `flown`, the key IDs its commands
were sealed with (joined by `+` when there are two) and the serial of the station certificate
in force at AOS: `04:12:05 pipit-1 PEL flown k-2028-10 1102`,
`23:52:10 cirrine-1 YST flown k-2028-07+k-2028-10 1001`. A miss names its first refusal:
`09:14:20 brine-3 LAN missed: certificate 1004 is revoked`. The real day: 10 flown, 0 missed.

#### More than a password (chapter 13)

- **Reset links.** The console emails a link holding a `rand.Text()` token and keeps one pending
  `Reset{Hash, Expires}` per account: `Hash` is the token's hex SHA-256 (`tokenDigest`), `Expires` 30
  minutes after the start (`resetFor`), on an injected clock; a new request replaces the pending one.
  `Redeem` accepts a token before `Expires` (refused from it), compares digests with
  `subtle.ConstantTimeCompare`, and deletes the reset only on success. Every refusal is
  `errInvalidReset` (`invalid or expired reset link`), unwrapped, and changes nothing. The form answers
  every name `If that account exists, a reset link is on its way.` (`resetSent`) and queues mail
  (`send`) only for an account that exists (`known`).
- **Codes (TOTP).** RFC 6238 over RFC 4226: HMAC-SHA-1 of the count of 30-second steps since 1970
  (`timeStep`), as 8 big-endian bytes; dynamic truncation; six digits with leading zeros (`hotp`,
  `totp`). Each engineer's authenticator app and the console share 20 random bytes, set up once; the
  console keeps the secret whole and guards it like a key. `(*Codes).Check` accepts the codes of the
  step before, the current step and the step after, in constant time, only for a step above the last
  one it accepted for that user, and then records it. An account without a secret has no code.
- **Recovery codes.** A sheet of eight (`sheetSize`) `rand.Text()` codes, printed once; the console keeps
  only their hex SHA-256 (`Recovery.unused`). Each code works once, and a new sheet replaces the old one.
- **The gate.** From Friday 3 November 2028, 16:00 UTC, `SignIn` checks the password with
  `Console.Check`, then refuses an account with ten wrong codes in a row (`maxMisses`), then takes a code
  from the phone or else one recovery code; every refusal is `errInvalidSignIn` (`invalid user, password
  or code`). A wrong password spends nothing and counts no miss; a right code clears the count.
  `ResetPassword` redeems the link, saves a new record, ends the account's sessions and clears its
  misses; it signs nobody in. The kiosk's shared account has neither a phone nor a sheet, so from 16:00
  it signs nobody in.
- **Declarations.** `Reset`, `Resets`, `Codes`, `Recovery` and `Gate` are as `reset.go`, `codes.go`,
  `recovery.go` and `gate.go` give them in 13.9.

#### The code, by chapter

Each lesson is its own `package main` (module `lesson`). A lesson gives the earlier pieces it
needs as read-only files under these names, and uses only the fields it needs. `Key`,
`Command`, `Pass` and Season 3's `Booking` and `Landing` are unchanged.

| Chapter | Names it adds (and files) |
|---|---|
| 1 | `newKey`, `newNonce`, `newSalt`, `newSessionToken` (`random.go`); `bits` and `Noncer` are lesson-only |
| 2 | `parseDigest`, `newHash`, `fingerprintAs` (`digest.go`), `checkManifest` (`manifest.go`); `fingerprint`, `fingerprintStream` and `save` are lesson-only |
| 3 | `Command` (`command.go`), `tag`, `verify`, `canonical` and its helpers `checkField`, `satIDRune`, `opRune`, `argRune` (`mac.go`), `Guard`, `(*Guard).Accept`, `ErrForged`, `ErrReplay` (`guard.go`), `checkLine` (`line.go`); `samePhrase`, `Signer` and `NewSigner` stay in their lessons' `main.go` |
| 4 | `Record`, `derive`, `formatRecord`, `parseRecord`, `checkPassword`, `newRecord`, `login`, `Console`, `(*Console).Check`, `(*Console).Login`, `(*Console).matches`, `(*Console).newRecord`, `needsUpgrade`, `dummyRecord`, `errInvalidLogin` (`password.go`, `console.go`); `checkStored` stays in 4.6's `main.go` |
| 5 | `signKey`, `satKey`, `fleetKeys`, `KeyVersion`, `KeyRing`, `loadRing`, `(*KeyRing).SatKey`, `(*KeyRing).Versions`, `(*KeyRing).String`, `(*KeyRing).GoString`, `Key.String`, `Key.GoString`, `parseVersion`, `(*KeyRing).add`, `ErrUnknownKey` (`derive.go`, `ring.go`; `%#v` calls GoString, so both types need one) |
| 6 | `seal`, `open`, `sealCmd`, `Sealed`, `(Sealed).AD`, `formatSealed`, `parseSealed`, `Satellite`, `SatKeys`, `(*Satellite).Accept`, `sealCommand`, the `Err…` sentinels (`seal.go`, `satellite.go`) |
| 7 | `Notary`, `newNotary`, `(*Notary).Sign`, `(*Notary).Published`, `checkReceipt`, `publishKey`, `readKey`, `Receipt`, `Envelope`, `ReceiptKey`, `wrapReceipt`, `openEnvelope`, `verifyEnvelope`, `issue`, `checkDelivery` (`receipt.go`) |
| 8 | `agree`, `parseShare`, `agreeHex`, `Identity`, `hello`, `checkHello`, `ErrBadHello`, `sessionKey`, `Backhaul`, `(*Backhaul).NewSession`, `wrapFor`, `unwrap`, `dial`, `answer` (`agree.go`, `backhaul.go`). Never a type named `Station`: that name is the bible's ground station |
| 9 | `trusted`, `Issuer`, `newRoot`, `newIntermediate`, `issueServer`, `issueClient`, `verifyStation`, `checkAt`, `issueAll` (`ca.go`) |
| 10 | `newClient`, `serverConfig`, `harden`, `requireClients`, `clientConfig`, `caller`, `Queue`, `stationHandler`, `dialController`, `uplinkClient`, `stationServer`, `send` (`tlsconf.go`, `station.go`) |
| 11 | `Pad`, `(*Pad).Seal`, `shortCode`, `firstMatch` (lesson-only, `main.go`); `readRSA`, `encryptTo`, `decryptWith`, `signRequest`, `checkRequest`, `keyFingerprint`, `Request`, `apiKeyLabel`, `webhookLabel` (`brinewatch.go`); `handOver` (`main.go`) |
| 12 | `accepts`, `signingKey`, `loadKeyCommand`, `revoke`, `isRevoked`, `withRevocations`, `switchTime` (`rotate.go`), `(*Satellite).Receive` and `load` (`satellite.go`), `Slot`, `Day`, `(Day).Passes`, `(Day).Connect`, `(Day).Ping`, `handshake`, `Report`, `runDay` (`day.go`, `main.go`) |
| 13 | `resetFor`, `resetSent`, `errInvalidReset`, `Reset`, `Resets`, `tokenDigest`, `(*Resets).Start`, `(*Resets).Redeem`, `(*Resets).Request` (`reset.go`); `stepSeconds`, `hotp`, `timeStep`, `totp` (`totp.go`); `Codes`, `(*Codes).Check` (`codes.go`); `sheetSize`, `Recovery`, `(*Recovery).Issue`, `(*Recovery).Use` (`recovery.go`); `maxMisses`, `errInvalidSignIn`, `Gate`, `(*Console).admit`, `(*Console).setPassword`, `(*Console).signOut` (`gate.go`); `(*Gate).SignIn`, `(*Gate).ResetPassword` (`main.go`). Chapter 4's `Console`, its methods and the record helpers are given unchanged |

Signatures every chapter after the one that defines them uses exactly:

```go
func newKey(id string) Key
func newNonce() []byte
func newSalt() []byte
func newSessionToken() string
func parseDigest(s string) (alg string, sum []byte, err error)
func newHash(alg string) (hash.Hash, error)
func fingerprintAs(r io.Reader, alg string) (string, error) // "sha256:<hex>"
func canonical(c Command) (string, error)
func verify(canonical, sigHex string, key []byte) bool // Season 1's, unchanged
func (g *Guard) Accept(c Command, key []byte) error
func checkLine(line string, key []byte, g *Guard) (Command, error)
func derive(pw string, salt []byte) ([]byte, error)
func parseRecord(s string) (Record, error)
func formatRecord(r Record) string
func checkPassword(r Record, pw string) (bool, error)
func signKey(master Key, satID string) ([]byte, error)          // 5.1: the "sign" key only
func satKey(master Key, purpose, satID string) ([]byte, error) // from 5.4: purpose "sign" or "seal"
func loadRing(r io.Reader) (*KeyRing, error)
func (r *KeyRing) SatKey(keyID, purpose, satID string) ([]byte, error)
func (r *KeyRing) Versions() []KeyVersion // by Activate
func (s Sealed) AD() []byte
func formatSealed(s Sealed) string
func parseSealed(line string) (Sealed, error)
func (s *Satellite) Accept(line string) (Command, error)
func sealCommand(ring *KeyRing, keyID string, c Command) (string, error)
func newNotary(keyID string, seed []byte) (*Notary, error)
func (n *Notary) Sign(msg []byte) []byte
func (n *Notary) Published() []byte // what customers are given
func wrapReceipt(n *Notary, r Receipt) ([]byte, error)
func verifyEnvelope(keys map[string]ed25519.PublicKey, env []byte) (Receipt, error)
func checkDelivery(keys []ReceiptKey, env []byte, file io.Reader) (Receipt, error)
func hello(id Identity, share *ecdh.PrivateKey) string
func checkHello(line string, known map[string]ed25519.PublicKey) (code string, share *ecdh.PublicKey, err error)
func sessionKey(shared, dialerShare, answererShare []byte) ([]byte, error)
func dial(conn io.ReadWriter, me Identity, peer string, known map[string]ed25519.PublicKey) ([]byte, error)
func answer(conn io.ReadWriter, me Identity, known map[string]ed25519.PublicKey) (peer string, key []byte, err error)
func newRoot(key ed25519.PrivateKey, now time.Time) (*x509.Certificate, error)
func newIntermediate(root Issuer, key ed25519.PrivateKey, now time.Time) (*x509.Certificate, error)
func issueServer(ca Issuer, host string, serial int64, pub ed25519.PublicKey, now time.Time) (*x509.Certificate, error)
func issueClient(ca Issuer, name string, serial int64, pub ed25519.PublicKey, now time.Time) (*x509.Certificate, error)
func verifyStation(leaf *x509.Certificate, chain []*x509.Certificate, roots *x509.CertPool, host string, now time.Time) error
func serverConfig(leaf, inter *x509.Certificate, key ed25519.PrivateKey, now func() time.Time) *tls.Config
func clientConfig(roots *x509.CertPool, host string, cert tls.Certificate, now func() time.Time) *tls.Config
func caller(r *http.Request) (string, bool)
func stationHandler(code string, q *Queue) http.Handler
func dialController(ctx context.Context, addr string, cfg *tls.Config) (net.Conn, error)
func send(ctx context.Context, c *http.Client, url, line string) error
func readRSA(pemText []byte) (*rsa.PublicKey, error)
func encryptTo(pub *rsa.PublicKey, secret []byte, label string) ([]byte, error)
func checkRequest(pub *rsa.PublicKey, req, sig []byte) error
func keyFingerprint(pub *rsa.PublicKey) (string, error) // "sha256:<hex>"
func accepts(ring []KeyVersion, id string, now time.Time) bool
func signingKey(ring []KeyVersion, now time.Time) (KeyVersion, error)
func loadKeyCommand(ring *KeyRing, satID string, seq uint64, next KeyVersion) (string, error)
func revoke(ca Issuer, number int64, entries []x509.RevocationListEntry, now time.Time) ([]byte, error)
func isRevoked(crlDER []byte, ca, cert *x509.Certificate, now time.Time) (bool, error)
func withRevocations(cfg *tls.Config, current func() []byte, ca *x509.Certificate, now func() time.Time)
func switchTime(passes []Pass, station string, from time.Time) (time.Time, error)
func tokenDigest(token string) string
func (r *Resets) Redeem(user, token string) error
func (r *Resets) Request(user string) string
func hotp(secret []byte, counter uint64) string
func totp(secret []byte, t time.Time) string
func (c *Codes) Check(user, code string) bool
func (r *Recovery) Use(user, code string) bool
func (g *Gate) SignIn(user, pw, code string) (string, error)
func (g *Gate) ResetPassword(user, token, pw string) error
```

`Satellite` (the flight software's stand-in) holds `ID`, `Keys map[string]SatKeys` (by key ID;
`SatKeys{Sign, Seal []byte; Retire time.Time; Revoked bool}`), `Last uint64` and
`Now func() time.Time`; from chapter 12 its `Receive` checks a line as `Accept` does, then
installs the keys a LOAD_KEY carries (`load`). `Guard` holds `last map[string]uint64` (a zero
`Guard` works). `Queue` holds a mutex and the queued lines.

#### How lessons test it

- **Never on random output.** Tests check properties (sizes, alphabets, two draws differing,
  round trips, a changed byte refused) or fixed vectors from the table above, computed in the
  test with the standard library where the algorithm is deterministic (SHA-2, SHA-3, HMAC,
  PBKDF2, HKDF, AES-GCM with a fixed nonce, Ed25519, X25519 from fixed private keys, ML-KEM
  round trips). A property that random output could break by chance (two 12-byte nonces equal)
  has odds below 2^-90. RSA-OAEP and RSA-PSS are randomized, so tests open or verify their output and never compare its bytes.
- **The attacker is the test.** One lesson in every chapter shows the mistake the chapter
  guards against, in Mirela's voice, before the right way: lessons 1.1, 2.2, 3.1, 4.1, 5.1,
  6.1, 7.1, 8.4 (after two lessons of basics), 9.1, 10.1, 11.1, 12.1 and 13.1. The starter holds the
  mistake, and the read-only test plays the person who exploits it (predicts the nonce, forges
  the CRC, extends the hash, cracks the table, signs for another satellite, flips the
  counter-mode bytes, forges the shared-secret receipt, sits in the middle, presents a
  same-named certificate, impersonates a server, reads one pad-sealed secret from another, arrives a second after midnight, or pastes a digest from a copy of the reset file). The
  solution makes the attack fail.
- **Reading the code.** Where right and wrong code give the same answers (constant-time
  comparison, no `math/rand`, extract once), the test parses the learner's file with
  `go/parser`, as Season 1's lesson 14.10 did, and says what it looked for.
- **PBKDF2 budget**: at most six 600,000-iteration derivations per lesson run.
- **Servers and connections**: `httptest.NewUnstartedServer` with `TLS` set, then
  `StartTLS()`, or `tls.Server` on `net.Listen("tcp", "127.0.0.1:0")`; never a fixed port. A
  server's `Config.ErrorLog` is `log.New(io.Discard, "", 0)`, so refused handshakes print
  nothing. Clients have `RootCAs`, `ServerName` and a timeout of at most 2 s; clocks are
  injected through `tls.Config.Time` and `VerifyOptions.CurrentTime`; every wait has a guard
  of about a second; sleeps are 50 ms at most.
- **Files** live under `t.TempDir()`; filesystems for manifests are `fstest.MapFS`.
- Messages read `call = got, want want`, with hex shortened to its first 8 digits and `…` when
  the full value is longer than 50 characters.
- Readme examples and outputs stay within 50 characters a line: hex is shown as its first and
  last four digits around `…`, and PEM blocks (60-character lines) are never shown whole.

**Topics added this season** (front matter): `crypto`. Lessons keep using `bytes`, `errors`,
`files`, `http`, `json`, `security`, `streams`, `strings`, `testing`, `time` and the rest.

**New field guide words** (world words, in plain language): Ground Security Engineer, the ops
console, Cirrine's contract, Rotation Day, random bytes, nonce, fingerprint (digest), salt,
password record, iteration count, master key, key ID, derived key, sealed command, associated
data, tag, receipt, public key, private key, key agreement, share, the middle (a machine in the
middle), forward secrecy, certificate, certificate authority, root, intermediate, chain,
revocation list, overlap, retire, revoke, TLS, mutual TLS, post-quantum, pad, birthday bound, RSA, OAEP, PSS, API key, webhook secret, Brinewatch, reset link, authenticator app, TOTP, dynamic truncation, recovery code.

**Name checks** to run before the season ships: `Apogee Root CA 2028`, `Apogee Stations CA
2028`, `apogee.internal` (a reserved top-level domain, so only the company name matters), and
"Rotation Day" as a phrase. No new people, places, satellites or customers.
