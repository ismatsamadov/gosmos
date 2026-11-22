# ACK

A frame that says yes: the satellite ran the command whose number it carries. Short for acknowledgment.

# alert

A message the pager sends about a problem: FIRING when it starts and RESOLVED when it's over. Who receives it depends on its rule's severity.

# AOS

The moment a satellite comes into a station's view and the link can start: the beginning of a pass. Short for acquisition of signal.

# associated data

The readable part sent beside sealed bytes, such as a sealed command's header: anyone can read it, but changing it makes the whole command fail to open.

# backhaul

The link that carries a station's data back to headquarters. When storms cut Ystmark's fiber cable, it falls back to a slow backhaul through a satellite.

# banner

The one line of text a station's radio sends the moment a call opens, before anything else: its station code, the word READY and its protocol version, as in `PEL READY PROTO 1`.

# booking

A pass a customer has reserved at one of Apogee's stations, with an ID such as `bk-4494`: which satellite, which station, and when the pass starts and ends.

# booking service

The program that lets customers' programs find passes, book them and hear when their data has landed, day and night. It has been open to every customer since 1 February 2028.

# bookingd

The program that is the booking service: it books passes for customers and tells them when their data lands. Like Apogee's other services, its name ends in d, for a program that runs on its own in the background.

# Calidra Atoll

An Apogee station on an atoll near the equator in the Pacific, live since Monday 2 April 2029. Palms line its horizon and tropical rain fades its signals, as at Lantern Atoll. Its code is CAL.

# candidate

A predicted pass a station could fly: the station offers one of the satellite's bands, and the satellite climbs at least 10° above the horizon.

# Cape Wren

Apogee's ground station in the far south: unstaffed, and the windiest site it owns, so its dish stows itself when the wind reaches 72 km/h. Its code is CWR.

# chunk

One numbered piece of a file coming down from a satellite. Chunks are numbered from 0, so the missing ones can be listed and asked for again.

# Cirrine

Short for Cirrine Atmospherics, Apogee's weather-data customer. Its first satellite is Cirrine-1, and twelve more launched on Tuesday 17 April 2029.

# Cirrine-1

Cirrine Atmospherics' first weather satellite. Since August 2028 every command to it goes up sealed, as Cirrine's contract asks.

# Cirrine-5

One of the twelve weather satellites Cirrine Atmospherics launched on Tuesday 17 April 2029. In the investors' week, its file of Tuesday 24 July lands eighteen seconds after Cirrine's receiver goes down for a restart.

# Cirrine-7

One of the twelve weather satellites Cirrine Atmospherics launched on Tuesday 17 April 2029. Like all thirteen, it takes only sealed commands.

# contact

A pass that has been given a dish: what a station will fly. The plan gives a dish to each candidate it can, in order of AOS.

# Corith Point

An Apogee station on a point of an island in the Arabian Sea, live since Monday 7 May 2029; the monsoon arrives there in June. Its code is CRP.

# counter

A metric that only ever goes up, such as the passes a station has failed. It starts again from 0 when its program restarts.

# CRC

A standard kind of check number, good at catching the damage a noisy line does: a frame whose CRC no longer matches the bytes it covers was damaged on the way. Short for cyclic redundancy check.

# derived key

A key worked out from a master key and a label, such as `sign|cirrine-1`, by a one-way recipe: each label gives an unrelated key, and no derived key leads back to the master or to any other.

# dish

A station's bowl-shaped antenna, which turns to follow a satellite across the sky. Apogee's dishes range from 1.2 m on the test van to 5.4 m at Salar Quieto.

# dish loader

The program at each station that hands its dishes their contacts, a day's plan at a time.

# dlinkd

The part of Apogee's software that puts each file from a satellite back together and delivers it to the customer.

# dropped pass

A candidate that finds no dish free at its AOS, because the station's dishes are all still busy with earlier contacts; the plan leaves it out.

# DUMP_FILE

The command that asks a satellite to send a file down again: the whole file, or only the chunks it lists by number. A satellite keeps what it sends, so a later pass can ask for what an earlier one missed.

# error budget

The passes a station may lose while still keeping its availability target, counted in whole passes: 11 of Lantern Atoll's 1,199 in the last quarter of 2028. What is left of it is room for change, such as a deploy or a dish's recalibration.

# fail closed

Said of a program that refuses whenever something is wrong or unclear, instead of letting it through. Mirela approves nothing until it does: "LGTM once it fails closed."

# fiber

The fiber-optic cable that links Ystmark to the mainland, about 40 ms there and back. Storms cut it, and then only the slow backhaul is left.

# frame

linked: no

One message on the pass line: its length, its type, a number, a body and a CRC. A TLM frame carries one unit of telemetry, a satellite's readings from one moment, such as battery, temperature and mode.

# Furrow-A

The first of Furrowlight's crop-imaging satellites, flying since before you joined Apogee. Its photographs arrive in numbered pieces that the ground puts back together.

# Furrow-K

One of the twenty-four crop-imaging satellites Furrowlight launched on Wednesday 27 June 2029. Like every Furrow, it sends its photographs down in numbered pieces.

# Furrowlight

Apogee's customer whose satellites, Furrow-A, Furrow-B and the twenty-four launched on Wednesday 27 June 2029, photograph farmland. Hollis runs its operations, and Yusra handles its requests and delivery deadlines.

# go/no-go

The poll ops runs before every pass: each question gets a GO or a NO-GO, and a single NO-GO calls the pass off. For the investors' week, it is Ingrid's call on the whole week, made from the rehearsal's report.

# Ground Segment Lead

Your role from Monday 2 July 2029: you lead the rehearsal of Apogee's ground segment, every program on the ground, from planning a pass to delivering its data.

# ground station

A site on Earth with dishes and radios for talking to satellites as they pass overhead. Apogee has twenty, from Ystmark in the Arctic to Vorrel Island in the far south.

# histogram

A metric that sorts measurements into buckets by size, such as the satellites' answers that took up to half a second, half a second to 1 second, 1 to 2 seconds, and so on, and keeps a count and a total of them all.

# HMAC

A signature computed from a message and a secret key, so nobody without the key can forge one. Short for hash-based message authentication code; Apogee uses HMAC-SHA256.

# Hollis

Hollis Vey, who runs operations at Furrowlight and writes every email in capitals, subject lines included.

# host

The owner of the land a new Apogee station stands on, such as a lighthouse authority or a weather observatory. The host gives power, a line to the internet and a pair of hands; nobody from Apogee lives there.

# Ilha Brumosa

An Apogee station on a misty island in the South Atlantic, live since Monday 2 April 2029; "brumosa" means misty, and its host is a weather observatory. Its code is BRM.

# Joss

Joss Imbert of Kessick Polytechnic, once a student on Professor Marr's Thistle team, who now leads the flight software of Thistle-2.

# Kessick's receiver

The program Joss wrote at Kessick Polytechnic to hear Apogee say that a Thistle-2 file has landed: the university's first webhook receiver.

# key ID

The name of one version of a key, such as `k-2029-07` for the master key made in July 2029. A sealed command carries it, so the satellite knows which of its keys opens it.

# landing

The moment a booking's file, such as one of Furrow-K's images, has come down from the satellite and arrived at Apogee in full, ready for the customer, who is told about each one by webhook.

# Lantern Atoll

Apogee's ground station on an island near the equator, behind a horizon of palm trees, with one dish, a diesel generator and a caretaker on site. Its code is LAN.

# LOAD_KEY

The command that gives a satellite its own two keys from a new master key, always sealed under a master it already holds. Once it has run, the satellite can check commands made with the new master.

# LOS

The moment a satellite sets out of a station's view and the link ends: the end of a pass. Short for loss of signal.

# master key

The secret that every satellite's uplink keys are made from, kept only at headquarters: 32 random bytes, replaced in versions named by month, such as `k-2028-10` and `k-2029-07`.

# maximum elevation

The highest a satellite climbs above the horizon during a pass, often written max elevation. Low passes are short, far away and weak.

# Mirela

Mirela Dunn, Apogee's security engineer and guardian of the uplink. She signs off a review with "LGTM once it fails closed": it looks good to her once it refuses whatever is wrong.

# morning list

Ingrid's list of things that can wait for daylight, read at the 08:00 meeting in the Fishbowl. Warning alerts go on it instead of waking anyone.

# NAK

A frame that says no: a station refusing a pass, or passing on a satellite's refusal of a command, with a code and a reason. Short for negative acknowledgment.

# night shift

The team on duty overnight in the Fishbowl, who run the passes between midnight and morning.

# nonce

A number used only once: 12 random bytes sent in the open beside each sealed command. It may be public, but it must never repeat under the same key.

# on call

Being the engineer the pager wakes, at any hour, when something breaks. Engineers take turns, a week at a time, and the Board shows whose week it is.

# page

linked: no

An alert that wakes the on-call engineer, whatever the hour; to page someone is to send one.

# pager

Apogee's alarm service: it checks incoming data against its rules and pages the on-call engineer when something is wrong. Carrying the pager means being the one it wakes.

# pass

linked: no

One flyover: the five to fifteen minutes when one satellite is in view of one station, from the moment it rises into view to the moment it sets. Once the plan gives it a dish, it is a contact.

# pass archive

passplan's file of predicted passes: one line per pass, with its pass ID, AOS, LOS and maximum elevation.

# pass ID

A pass's name: its satellite, its station and the satellite's orbit number in five digits, as in `pipit-1/PEL/04182`.

# pass line

The connection the rehearsal holds with a station's radio for one pass, and the messages it carries both ways.

# passplan

The part of Apogee's software that predicts every pass, builds the day's schedule and feeds the Board.

# passplan 2

The new passplan, whose rewrite you led in 2029: it plans a hundred days ahead for twenty stations and eighty satellites, and took over from the old one on Thursday 28 June.

# PING

The simplest command: it changes nothing on the satellite and only asks it to answer, so a pass can show that its commands get through.

# Pipit-2

Apogee's radio test satellite. The team uses it to check commands protected with a shared secret key.

# Port Ellery

Apogee's headquarters, home of the Fishbowl, with a dish on the roof where every new piece of software runs first. Its code is PEL.

# radio

linked: no

A station's radio: the equipment behind its dishes that talks to satellites. The rehearsal calls it at the start of each pass.

# replay

Sending an old, properly signed command a second time to make a satellite repeat it. Sequence numbers stop it, because an old number is refused.

# retire

To stop accepting a key version from a set instant, as planned: `k-2028-10` retires at 00:00 UTC on 1 August 2029, a week after `k-2029-07` takes over.

# revoke

To withdraw a key or a certificate for good, at once, before its time is up, because it may be in the wrong hands. Whatever its dates say, a revoked one is refused, and it is never brought back.

# rota

The on-call timetable: which engineer is on call each week. A week starts on Monday at 08:00 UTC, when the name on the Board changes.

# science file

A file of what a satellite's instrument measured, not of its own health. Thistle-2's first comes from its radiation counter.

# sealed command

A command locked for one satellite under its own key, so that nobody on the way can read it, and any change makes it fail to open. Its header, which names the satellite, the number and the key ID, stays readable.

# sequence number

The number each command to a satellite carries: it must be higher than the last one the satellite ran, or the command is refused. A repeat means a duplicate or a replay.

# signature

A string of digits attached to a command or a webhook that proves who made it and that nothing in it changed on the way. A satellite refuses any command whose signature doesn't check out.

# simulated time

The week's own clock in the rehearsal's tests: every pass starts and ends on time, yet a whole week passes in moments.

# stow

To park a dish in a safe position, usually pointing straight up, so wind can't damage it. Ystmark's and Cape Wren's dishes stow when the wind reaches 72 km/h.

# telemetry

The data a satellite sends down about itself, such as battery voltage, temperatures and mode, so the ground can tell how it is doing. Often written TLM.

# the Board

The big screen in the Fishbowl that lists the next passes and the state of every station. Ingrid owns it.

# the drill

Nkem's rehearsal of tlmd's new pipeline on Wednesday 7 March 2029: twenty telemetry feeds at once, five from Apogee's stations and fifteen made from where the new stations would stand.

# the Fishbowl

Apogee's ops room at headquarters in Port Ellery, named for its glass walls on three sides. Every pass is run from here, under the Board, the big screen that lists the next passes.

# the investors' week

Monday 23 July 2029, 00:00 UTC, to Monday 30 July, 00:00 UTC: the week the investors who put money into Apogee watch the whole network, wanting every booked pass flown.

# the rehearsal

The program you build for the investors' week: Apogee's ground segment, its own code, run against the simulated world instead of the real one, so the whole week can be flown in seconds.

# the simulated world

Nkem's program that plays every station's radio, every satellite and every customer's receiver for the rehearsal, with whatever trouble Ingrid writes into it. She grew it from her drill of March.

# Thistle

Kessick Polytechnic's first student satellite, which re-entered the atmosphere in July 2027. Thistle-2 followed it in March 2029.

# Thistle-2

Kessick Polytechnic's second student satellite, which left its rocket on Thursday 8 March 2029. Its first beacon reached Port Ellery at 17:08 that day.

# Tirrin Rock

An Apogee station on a bare rock in the Norwegian Sea, with a lighthouse and two dishes, live since Monday 2 April 2029. Its code is TIR, and it shares polar passes with Ystmark.

# TL;DR

A one-line summary at the top of a message, for readers in a hurry; Yusra starts every message with one. Short for too long; didn't read.

# uplinkd

The program that signs and sends Apogee's commands up to the satellites. It runs at headquarters, and its station side at every ground station queues the commands for the next pass.

# UTC

The world's reference clock, the same everywhere on Earth, with no time zones or daylight saving; everything at Apogee runs on it. Short for Coordinated Universal Time.

# Varsk Island

A new Apogee station at a fishing harbor in the North Pacific, with fog most mornings and a hill to the north, live since Monday 2 April 2029. Its code is VSK.

# Vorrel Island

A new Apogee station on an icebound island in the far south, the south's Ystmark, live since Monday 4 June 2029 with two dishes. Its code is VRL.

# webhook

A message Apogee sends to a customer's own program, its receiver, the moment something happens, such as a file landing, instead of waiting to be asked. Each one is signed with a secret Apogee and the customer share, so the receiver can tell it came from Apogee.

# wind lock

A dish controller's own guard against storms: at a wind of 72 km/h or more it stows its dish and refuses every command to unstow or move it. The dish stays stowed until it is told to unstow.

# window

The span of time a plan covers, from 00:00 UTC on its first day; the rehearsal's window is the days it flies.

# Ystmark

Apogee's ground station on an Arctic island, the farthest north in the network. Its code is YST.

# Yusra

Yusra Belkadi, who looks after Apogee's customers and turns their worries into lists of what must be done. Every message she writes starts with a one-line summary.
