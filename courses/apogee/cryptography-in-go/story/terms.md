# AOS

The moment a satellite comes into a station's view and the link can start: the beginning of a pass. Short for acquisition of signal.

# API key

The secret a customer's programs send with every request to the booking service, to prove which customer is asking: the customer's name, a dot and 32 hex digits, such as `brinewatch.0f1e…e1f0`.

# associated data

The readable part sent beside sealed bytes, such as a sealed command's header, which the tag covers too: anyone can read it, but changing it makes the whole command fail to open.

# authenticator app

An app on an engineer's phone that shows a new six-digit code every 30 seconds, made from a secret only the phone and the ops console hold. From 3 November 2028 the console asks for one after the password.

# backhaul

The link that carries a station's data back to headquarters. When storms cut Ystmark's fiber cable, it falls back to a slow backhaul through a satellite.

# birthday bound

How soon some two of many tries share a value: after about the square root of the number of values possible. Codes of six hex digits have 16.7 million values, yet two probably match after about 4,800 tries, while matching one given code takes millions.

# booking

A pass a customer has reserved through the booking service. Each has an ID, such as `bk-4988`, and may carry a short note from the customer.

# booking service

Apogee's web service, open to every customer since 1 February 2028, through which customers' programs book passes and hear when their data has landed. Its program is called bookingd.

# Brinewatch

Apogee's ship-tracking customer, whose satellite, Brine-3, hears the radio signals that ships send at sea. Its older systems lock and sign with RSA keys.

# CA

Short for certificate authority.

# canonical form

The one exact text of a command that gets signed: its fields joined by bars, arguments last, as in `pipit-2|118|SET_MODE|SAFE`. Its fields are checked so that no two commands share a text.

# Cape Wren

Apogee's ground station in the far south: unstaffed, and the windiest site it owns, so its dish stows itself when the wind reaches 72 km/h. Its code is CWR.

# certificate

A signed file that ties a public key to a name, such as a station's host name, for a set of dates and uses. It counts only if a chain of signatures leads from it to a root you trust.

# certificate authority

Whoever signs certificates, vouching that each public key belongs to its name. Apogee's own has a root kept in the safe and an intermediate, the stations CA, that signs the stations' certificates.

# chain

The line of signatures from a certificate up to a root: Ystmark's certificate, signed by the stations CA, which the root signed.

# Chiquita

The nickname of Salar Quieto's smaller dish, 2.4 m across. Its big sister is La Grande.

# Cirrine

Short for Cirrine Atmospherics, Apogee's weather-data customer. Its satellite is Cirrine-1.

# Cirrine's contract

The bigger contract Cirrine Atmospherics signed on Friday 7 July 2028, with three conditions: from 25 August nobody but Cirrine and Apogee may read Cirrine-1's commands; every file delivered to Cirrine comes with a fingerprint; and every key Apogee uses for Cirrine is replaced by 31 October, then once a year.

# Cirrine-1

The weather satellite operated by Cirrine Atmospherics. Its commands must be checked before the ground sends them up.

# CRC-32

A standard kind of checksum, 4 bytes long, good at catching the damage radio noise does: a chunk whose CRC-32 no longer matches its data was damaged on the way. Short for 32-bit cyclic redundancy check.

# customer guide

Yusra's pages for customers' engineers: how their programs should talk to the booking service, and what each of its answers means.

# derived key

A key worked out from a master key and a label, such as `sign|cirrine-1`, by a one-way recipe: each label gives an unrelated key, and no derived key leads back to the master or to any other.

# digest

The value of fixed size that a hash works out from any bytes, such as the 32 bytes of a SHA-256 digest. Apogee writes it with its hash's name, as in `sha256:fe5e…bd66`. Another word for a fingerprint.

# dish

A station's bowl-shaped antenna, which turns to follow a satellite across the sky. Apogee's dishes range from 1.2 m on the test van to 5.4 m at Salar Quieto.

# dish controller

The computer that turns a dish, one for each dish. Programs on the station's network talk to it to point, stow or move the dish.

# dishd

The dish service at Port Ellery, the one program that sends commands to Ystmark's dish controllers. It proves who it is to them with a client certificate of its own.

# dlinkd

The part of Apogee's software that puts each file from a satellite back together and delivers it to the customer, with its fingerprint.

# dynamic truncation

How a 20-byte HMAC becomes a six-digit code: the low 4 bits of its last byte say where to read 4 bytes, and those, with the top bit cleared, make a number whose last six digits are the code.

# fails closed

Said of a program that refuses whenever something is wrong or unclear, instead of letting it through. Mirela approves nothing until it does: "LGTM once it fails closed."

# fiber

The fiber-optic cable that links Ystmark to the mainland, about 40 ms there and back. Storms cut it, and then only the slow backhaul is left.

# fingerprint

A short value worked out from every byte of a file. Change any byte and the fingerprint changes, and nobody can make a different file with the same one. Also called a digest.

# flight software

The programs that run on a satellite's own computer. They check each command that comes up and carry it out.

# forward secrecy

The promise that a key stolen later opens nothing recorded earlier: every session uses a new key, which nobody keeps once the session ends.

# Furrow-A

One of Furrowlight's two crop-imaging satellites. Its photographs arrive in numbered pieces that the ground must put back together.

# Furrowlight

Apogee's customer whose two satellites, Furrow-A and Furrow-B, photograph farmland. Yusra handles its requests and delivery deadlines.

# Ground Security Engineer

Your role from Monday 10 July 2028, still badge No. 039: you protect Apogee's commands, files, passwords and links, and Mirela reviews your work.

# hash

A recipe that turns any bytes, from none to many gigabytes, into a digest of fixed size. A good one, such as SHA-256, gives a digest that looks unrelated when even one byte changes, and nobody can find two inputs with the same digest.

# HMAC

A signature computed from a message and a secret key that only the two ends share, such as Apogee and a satellite, so nobody without the key can forge one. Short for hash-based message authentication code; Apogee signs commands with HMAC-SHA256, and authenticator apps use HMAC-SHA-1.

# Hollis

Hollis Vey, the ops manager at Furrowlight, who writes subject lines and booking notes in capitals.

# identity key

A station's lasting signing key on the backhaul. The other end holds its public half in advance, so a share signed with it can only have come from whoever holds the key: the station itself.

# intermediate

A certificate authority signed by the root, which signs everyday certificates so that the root can stay offline. Apogee's is the stations CA, `Apogee Stations CA 2028`.

# iteration count

How many rounds the slow recipe that turns a password into its stored key repeats: 600,000 for every new record at Apogee. Anyone guessing at a stolen record pays every round for every guess.

# key

linked: no

A secret that locks, unlocks or signs commands and data: whoever holds it can do what it does. Apogee's keys are 32 random bytes, and each version has an ID such as `k-2028-07`.

# key agreement

Two ends arriving at the same secret key over an open link: each sends only the public half of a key it made, and with today's computers nobody who sees only those halves can work out the secret.

# key ID

The name of one version of a key, such as `k-2028-07` for the master key made in July 2028, or `r-2028-08` for a receipt key. A sealed command or a receipt carries it, so the other end knows which key to use.

# landing

dlinkd's word to the booking service that a booking's file has been delivered: which file, how many bytes, and when it was complete.

# Lantern Atoll

Apogee's ground station on an atoll near the equator, run on diesel power with a caretaker on site. Its code is LAN.

# length extension

A trick against a tag made as the hash of a key followed by a message: from the tag and the key's length alone, anyone can make a valid tag for that message with more bytes added to its end. An HMAC can't be extended this way.

# LGTM

Short for looks good to me: what a reviewer writes to approve a change. Mirela adds a condition to hers: "LGTM once it fails closed."

# LOS

The moment a satellite sets out of a station's view and the link ends: the end of a pass. Short for loss of signal.

# machine in the middle

Someone on the link between two ends who can pass on, read or change whatever crosses it. On the backhaul, the company that runs the link could be one.

# manifest

The list delivered beside a customer's files: one line per file, with its digest, its size in bytes and its name, so anyone can check that every file arrived whole and unchanged.

# master key

The secret that every satellite's uplink keys are made from, which no program but uplinkd at headquarters ever uses: 32 random bytes, replaced in versions named by month, such as `k-2028-07`.

# mode

linked: no

The state a satellite reports it is in, which decides what ops does on a pass: NOMINAL, SAFE, DETUMBLE or LEOP.

# mutual TLS

TLS in which both ends present a certificate, so a server knows which program is calling as surely as the program knows which server answered. From October 2028, every link between Apogee's programs uses it.

# night shift

The team on duty overnight in the Fishbowl, who run the passes between midnight and morning.

# nonce

A number used only once: 12 random bytes sent in the open beside each locked command. It may be public, but it must never repeat under the same key.

# OAEP

The safe way to lock a message to an RSA public key: random bytes are mixed in, so the same message never locks the same way twice, and a label is checked, so a locked message opens only for the purpose it was locked for. Short for optimal asymmetric encryption padding.

# ops console

Apogee's internal web tool, running since long before you joined, where engineers log in to stow a dish, run a pass by hand or read a station's log.

# overlap

The days when a key and its successor are both accepted, so a command sealed with the old key just before the switch still counts when it arrives after it. Rotation Day's overlap runs from 17 to 24 October 2028.

# pad

Random bytes that two ends each hold a copy of, carried by hand in advance, which hide a message XORed with them. Used once, byte for byte, a pad gives away nothing but the message's length; used twice, it leaks. Brinewatch's came on a USB stick in November 2027.

# pager

Apogee's alarm service: it checks incoming data against its rules and pages the on-call engineer when something is wrong. Carrying the pager means being the one it wakes.

# pass

linked: no

One flyover: the five to fifteen minutes when one satellite is in view of one station, from the moment it rises into view to the moment it sets. Also called a contact.

# passplan

The part of Apogee's software that predicts every pass, builds the day's schedule and feeds the Board.

# password record

The line the ops console keeps for each account instead of its password: the recipe's name, the iteration count, the salt and the stored key, as in `pbkdf2-sha256$600000$…`. It can check a password but doesn't hold it.

# Pipit-1

Apogee's small test satellite. A new engineer's first message goes up to it and is echoed back in its status broadcast.

# Pipit-2

Apogee's radio test satellite. The team uses it to check commands protected with a shared secret key.

# Pipit-3

Apogee's newest test satellite, launched on 5 April 2027. Finding its first faint broadcasts took a search across radio frequencies.

# Port Ellery

Apogee's headquarters, home of the Fishbowl, with a dish on the roof where every new piece of software runs first. Its station code is PEL.

# post-quantum

Designed to stay safe even against a quantum computer, a machine that, if one is ever built big enough, could break today's key agreements. ML-KEM is post-quantum.

# private key

The half of a key pair its owner keeps secret, to sign with or to agree on keys. It never leaves the machine that holds it.

# PSS

The safe way to sign with an RSA private key: fresh random bytes go into every signature, so the same message signed twice gives two different signatures, and a check accepts each. Short for probabilistic signature scheme.

# public key

The half of a key pair its owner can hand to anyone. Others use it to check the owner's signatures or to agree on a key with the owner; nobody can sign with it.

# random bytes

Bytes nobody can predict or make again, taken from the computer's own secure source, never from a clock or a seed. Every key, nonce, salt and token at Apogee starts as random bytes.

# receipt

Apogee's signed word that a booking's file was delivered: which file, its size, its fingerprint and when it landed. Only Apogee can sign one, and anyone with Apogee's public key can check it.

# recovery code

One of eight codes printed once for each engineer, each good for a single sign-in in place of a code from the phone, for the day the phone is lost. The console keeps only their digests.

# replay

Sending an old, properly signed command a second time to make a satellite repeat it. Sequence numbers stop it, because an old number is refused.

# replay guard

The part of the uplink that refuses replayed commands: it accepts a command only if its sequence number is higher than the last one accepted.

# reset link

The link the ops console emails to someone who has forgotten a password. Its token works once and for 30 minutes, and the console keeps only the token's digest.

# retire

To stop accepting a key version from a set instant, as planned: `k-2027-03` retires at 00:00 UTC on 11 August 2028, once every satellite holds keys of its own.

# revocation list

A dated list of the certificates a certificate authority has revoked, signed by that authority and replaced within 7 days. A list that doesn't check out, or is out of date, is refused, and with it every connection it should vouch for.

# revoke

To withdraw a key or a certificate for good, at once, before its time is up, because it may be in the wrong hands. Whatever its dates say, a revoked one is refused, and it is never brought back.

# root

The certificate every other one is checked against: it signs itself, and you trust it because you hold it yourself. Apogee's root key stays offline in the safe at Port Ellery.

# Rotation Day

Tuesday 17 October 2028, when Mirela plans to replace every key in Apogee's network, at every station, without missing a single pass.

# RSA

An older kind of key pair, whose public key can lock a message for its owner as well as check its owner's signatures. Its strength grows with its size in bits: Apogee accepts 2,048 bits and up, and Brinewatch's new keys have 3,072.

# Salar Quieto

Apogee's ground station on a salt flat 3,840 m up, with the lowest horizon and the thinnest air in the network. Its code is SQT.

# salt

Random bytes, 16 at Apogee, kept beside each stored password, so that two people with the same password never get the same record.

# sealed command

A command locked for one satellite under its own key, so that nobody on the way can read it, and any change makes it fail to open. Its header, which names the satellite, the number and the key ID, stays readable.

# sequence number

A counter that goes up by one with every frame or command. A gap means something was lost, and a repeat means a duplicate or a replay.

# session key

The key that seals everything sent in one backhaul session. Both ends derive it from the secret they agreed on, and nobody keeps it once the session ends.

# session token

The secret text the ops console gives your browser when you log in. The browser shows it with every request in place of your password, until you log out.

# SHA-256

The hash Apogee uses for its own deliveries: any bytes in, a 32-byte digest out. Short for Secure Hash Algorithm, 256 bits.

# SHA3-256

A newer hash of a different design, with a 32-byte digest like SHA-256's. Cirrine's archive has used it since Cirrine's first satellite.

# share

linked: no

The public half of the key a station makes for one backhaul session, sent as 64 hex digits. A share may be public; its private key never leaves the station.

# signature

A string of digits attached to a command, a receipt, a key share or a certificate that proves who made it and that nothing in it changed on the way. A satellite refuses any command whose signature doesn't check out.

# tag

The short value sent with a message to prove that someone holding the key made it and that nothing changed since: the 32 bytes of an HMAC, or the last 16 bytes of a sealed command.

# the Fishbowl

Apogee's ops room at headquarters in Port Ellery, named for its glass walls on three sides. Every pass is run from here, under the Board, the big screen that lists the next passes.

# TL;DR

A one-line summary at the top of a message, for readers in a hurry; Yusra starts every message with one. Short for too long; didn't read.

# TLS

The protocol that encrypts a connection between two programs and proves who is at the far end, with the certificates they present. HTTPS is HTTP over TLS. Short for Transport Layer Security.

# TOTP

Short for time-based one-time password: a code made from a shared secret and the clock, as the published standard RFC 6238 describes. Authenticator apps show them.

# uplink

The radio link from the ground up to a satellite, the commands sent up it, and the part of Apogee's software that signs and sends them.

# uplinkd

The program that signs and sends Apogee's commands up to the satellites. It runs at headquarters, and its station side at every ground station queues the commands for the next pass.

# UTC

The world's reference clock, the same everywhere on Earth, with no time zones or daylight saving; everything at Apogee runs on it. Short for Coordinated Universal Time.

# webhook secret

The secret the booking service signs its messages to a customer's program with, such as "your data has landed", so the program can check that Apogee sent them. Apogee and the customer both hold it.

# Ystmark

Apogee's ground station on an Arctic island, the farthest north in the network. Its code is YST.
