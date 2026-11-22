# alert

A message the pager sends about a problem: FIRING when it starts and RESOLVED when it's over. Who receives it depends on its rule's severity.

# AOS

The moment a satellite comes into a station's view and the link can start: the beginning of a pass. Short for acquisition of signal.

# API key

The secret a customer's program sends the booking service to prove which customer it is: the customer's ID, a dot, then 32 secret letters and digits. Whoever holds it can book passes as that customer.

# availability

The share of a station's booked passes that were really flown: Lantern Atoll's 1,177 flown of 1,199 booked in the last quarter of 2028 is 0.98165, or 98.17%.

# availability target

The availability a station promises its customers, written in basis points: 9950, or 99.5%, at every station, and 9900 at Lantern Atoll, whose palms and downpours alone cost about one pass in a hundred.

# basis point

A hundredth of a percent. Apogee writes its availability targets in basis points, so 9950 means 99.50% of booked passes flown.

# benchmark

A measurement of how long a piece of code takes, made by running it many times over and dividing the time by the number of runs.

# booking

A pass a customer has reserved through the booking service. Each has an ID, such as `bk-4988`, and may carry a short note from the customer.

# booking service

Apogee's web service, open to every customer since 1 February 2028, through which customers' programs book passes and hear when their data has landed. Its program is called bookingd.

# bookingd

The program that is the booking service: it books passes for customers and tells them when their data lands. Like Apogee's other services, its name ends in d, for a program that runs on its own in the background.

# Brine-3

Brinewatch's working satellite, which tracks ships. Its orbit stays too far south for Ystmark to see it.

# Brinewatch

Apogee's customer that tracks ships at sea with its satellite Brine-3. Brine-1 was lost at launch, and Brine-2 never left its deployer.

# burn rate

How fast lost passes are spending an error budget: at 1 the budget lasts exactly as long as planned, and at 50 a month's budget is gone in fourteen hours.

# Cape Wren

Apogee's ground station in the far south: unstaffed, and the windiest site it owns, so its dish stows itself when the wind reaches 72 km/h. Its code is CWR.

# caretaker

The person on site at Lantern Atoll who looks after the station: when the night shift asks, the caretaker checks the dish, the radio, the antennas and the generator.

# Cirrine

Short for Cirrine Atmospherics, Apogee's weather-data customer. Its satellite is Cirrine-1.

# Cirrine-1

The weather satellite operated by Cirrine Atmospherics. Its commands must be checked before the ground sends them up.

# counter

A metric that only ever goes up, such as the passes a station has failed. It starts again from 0 when its program restarts.

# debug port

A private address that each of Apogee's programs answers on, reachable only from the station's or headquarters' own network, where engineers change a running program's log level and read its numbers without restarting it.

# deploy

Putting a new version of a program into service: the running copy stops and the new one starts in its place, so anything the old one kept only in memory is gone.

# dish

A station's bowl-shaped antenna, which turns to follow a satellite across the sky. Apogee's dishes range from 1.2 m on the test van to 5.4 m at Salar Quieto.

# downpour

A sudden, heavy tropical rain. From November to February, Lantern Atoll gets one most evenings and nights, and it can drown out a satellite's signal for minutes at a time.

# elevation

How high a satellite stands above the horizon as seen from a station, as an angle: 0° at the horizon, 90° straight overhead.

# error budget

The passes a station may lose while still keeping its availability target, counted in whole passes: 11 of Lantern Atoll's 1,199 in the last quarter of 2028. What is left of it is room for change, such as a deploy or a dish's recalibration.

# fails closed

Said of a program that refuses whenever something is wrong or unclear, instead of letting it through. Mirela approves nothing until it does: "LGTM once it fails closed."

# flapping

What a rule does when its value hovers around a single line: it fires and resolves over and over, such as for a disk bobbing around 90% full.

# flight recorder

A trace that a program keeps running in its own memory, holding only the last stretch of what it did, so that when something fails those last seconds can be saved to a file.

# frame

linked: no

One unit of telemetry as it comes down: a satellite's readings from one moment, such as battery, temperature and mode, numbered by its sequence number.

# Furrowlight

Apogee's customer whose two satellites, Furrow-A and Furrow-B, photograph farmland. Yusra handles its requests and delivery deadlines.

# gauge

A metric that goes up and down, such as the passes under way at a station right now, or how full its disk is.

# generator

The diesel engine that makes Lantern Atoll's electricity, far from any power line. In late October 2028 it started to sound different.

# Ground Reliability Engineer

Your role from Monday 6 November 2028, still badge No. 039: you keep Apogee's programs working day and night, take your turn on call, and find out why something failed.

# ground station

A site on Earth with dishes and radios for talking to satellites as they pass overhead. Apogee has five, from Ystmark in the Arctic to Cape Wren in the far south.

# heap

The memory a running program takes, as it goes, to keep its data in. Whatever nothing can reach any more is freed by the program's garbage collector; the rest stays in use.

# histogram

A metric that sorts measurements into buckets by size, such as the satellites' answers that took up to half a second, half a second to 1 second, 1 to 2 seconds, and so on, and keeps a count and a total of them all.

# Hollis

Hollis Vey, who runs operations at Furrowlight and writes every email in capitals, subject lines included.

# hot path

The code a program spends most of its time in, and the first place to look when the program is slow.

# hysteresis

Using two lines instead of one: a rule fires at one value and is over only below a lower one, so a value hovering between them can't make it flap.

# incident review

A meeting after something has gone wrong, where the team works out what happened, why, and what to change, without blaming anyone.

# lab bench

The test setup at Port Ellery headquarters where every change to Apogee's programs runs first, against one of Apogee's own test satellites, before it reaches a station.

# label

One of the things a metric is split by, such as satellite or outcome, so that each pair of values gets a number of its own. A label's values must be few: satellites, never pass IDs.

# Lantern Atoll

Apogee's ground station on an island near the equator, behind a horizon of palm trees, with one dish, a diesel generator and a caretaker on site. Its code is LAN.

# LGTM

Short for looks good to me: what a reviewer writes to approve a change. Mirela adds a condition to hers: "LGTM once it fails closed."

# log level

How much a log line matters: DEBUG for fine detail, INFO for what normally happens, WARN for trouble that may pass, and ERROR for a failure.

# log line

One entry in a program's log, the record it keeps of what it does, written as it happens. At three in the morning, the log lines are often all anyone has to go on.

# log rotation

Starting a fresh log file when the current one reaches its size limit or a new day begins, and renaming the old one under its day and a number, so that no file grows without end.

# LOS

The moment a satellite sets out of a station's view and the link ends: the end of a pass. Short for loss of signal.

# master key

The secret that every satellite's uplink keys are made from, which no program but uplinkd at headquarters ever uses: 32 random bytes, replaced in versions named by month, such as `k-2028-07`.

# metric

A number a program keeps about its own work and publishes for others to read, such as the passes it has flown or how full a disk is.

# Monday restart

The line Ingrid added to her checklist after the page of 8 November, "Monday 08:00: restart LAN uplinkd": every Monday morning Lantern Atoll's uplinkd starts afresh, until anyone knows why a restart helps.

# morning list

Ingrid's list of things that can wait for daylight, read at the 08:00 meeting in the Fishbowl. Warning alerts go on it instead of waking anyone.

# night shift

The team on duty overnight in the Fishbowl, who run the passes between midnight and morning.

# on call

Being the engineer the pager wakes, at any hour, when something breaks. Engineers take turns, a week at a time, and the Board shows whose week it is.

# ops console

Apogee's internal web tool, running since long before you joined, where engineers log in to stow a dish, run a pass by hand or read a station's log.

# page

linked: no

An alert that wakes the on-call engineer, whatever the hour; to page someone is to send one.

# pager

Apogee's alarm service: it checks incoming data against its rules and pages the on-call engineer when something is wrong. Carrying the pager means being the one it wakes.

# pass

linked: no

One flyover: the five to fifteen minutes when one satellite is in view of one station, from the moment it rises into view to the moment it sets. Also called a contact.

# pass ID

A pass's name: its satellite, its station and the satellite's orbit number in five digits, as in `pipit-1/PEL/04182`.

# passplan

The part of Apogee's software that predicts every pass, builds the day's schedule and feeds the Board.

# personal data

Anything that tells who a person is or how to reach them, such as a name or an email address. Apogee keeps customers' contact addresses only to write to them, and never copies one whole into a log.

# Port Ellery

Apogee's headquarters, home of the Fishbowl, with a dish on the roof where every new piece of software runs first. Its station code is PEL.

# profile

A record of where a running program spends its time or keeps its memory, line by line, saved to a file that tools such as `go tool pprof` can read.

# recorder

linked: no

The program at each station that records each pass as it comes down. At Lantern Atoll, the recordings wait on the station computer's disk until headquarters collects them.

# redaction

Replacing a secret, or anything else a log must not keep, with a placeholder such as `[redacted]` before the line is written, so the line keeps its shape and loses the secret.

# rota

The on-call timetable: which engineer is on call each week. A week starts on Monday at 08:00 UTC, when the name on the Board changes.

# Rotation Day

Tuesday 17 October 2028, when Mirela replaced every key in Apogee's network, at every station, without missing a single pass.

# runbook

The written steps the night shift follows when something goes wrong at a station, one step at a time, before anyone has to be woken.

# Salar Quieto

Apogee's ground station on a salt flat 3,840 m up, with the lowest horizon and the thinnest air in the network. Its code is SQT.

# sampling

Keeping only some of many identical log lines, such as the first ten a minute and then one in a hundred, and noting how many were left out, so that a flood of one warning can't fill a disk.

# severity

How much an alert matters, set by its rule: Critical pages the person on call, Warning goes on the morning list, and Info is only logged.

# silence

A stretch of time when the pager holds back a station's alerts, such as while its generator is being repaired. A silence always has an end.

# span

One timed step of a piece of work, such as a pass or its wait for a worker, recorded with an ID, its start and end, and the ID of the step it belongs to. All the spans of one pass share a trace ID.

# stow

To park a dish in a safe position, usually pointing straight up, so wind can't damage it. Ystmark's and Cape Wren's dishes stow when the wind reaches 72 km/h.

# structured log

A log whose lines are records rather than free text: each holds a fixed message and named values, so a program can search, filter and count them.

# the Board

The big screen in the Fishbowl that lists the next passes and the state of every station. Ingrid owns it.

# the Fishbowl

Apogee's ops room at headquarters in Port Ellery, named for its glass walls on three sides. Every pass is run from here, under the Board, the big screen that lists the next passes.

# threshold

The value at which a rule's problem starts, such as a disk 90% full. A rule may also ask the value to stay there a while before anyone is told.

# TL;DR

A one-line summary at the top of a message, for readers in a hurry; Yusra starts every message with one. Short for too long; didn't read.

# token

The short-lived secret the booking service gives a customer's program in exchange for its API key: good for one hour, it carries the customer's ID, the time it expires and a signature.

# trace

A record of what a program did and when, step by step, saved to a file that tools such as `go tool trace` show on a timeline.

# uplinkd

The program that signs and sends Apogee's commands up to the satellites. It runs at headquarters, and its station side at every ground station queues the commands for the next pass.

# UTC

The world's reference clock, the same everywhere on Earth, with no time zones or daylight saving; everything at Apogee runs on it. Short for Coordinated Universal Time.

# worker

One of the four slots on which a station's uplinkd flies its passes: a pass takes a free worker before it starts and gives it back when it ends. With all four held, the next pass has to wait.

# Ystmark

Apogee's ground station on an Arctic island, the farthest north in the network. Its code is YST.
