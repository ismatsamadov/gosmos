# AOS

The moment a satellite comes into a station's view and the link can start: the beginning of a pass. Short for acquisition of signal.

# band

linked: no

A range of radio frequencies, such as UHF, S-band or X-band. Ops also sorts batteries into bands, from FULL down to CRITICAL. In 2029 each host's form lists the bands its station uses, and Apogee books passes on three: UHF, S and X.

# beacon

A short status message a satellite broadcasts by itself, without being asked, with readings such as its battery, temperature and mode. Pipit-1 sends one every 30 seconds.

# benchmark

A measurement of how long a piece of code takes, made by running it many times over and dividing the time by the number of runs. In May 2029, Tomás's "forty times faster" for passplan 2 came from benchmarks that measured different work.

# bookingd

The program that is the booking service: it books passes for customers and tells them when their data lands. Like Apogee's other services, its name ends in d, for a program that runs on its own in the background.

# Brine-3

Brinewatch's working satellite, which tracks ships. Its orbit stays too far south for Ystmark to see it.

# Brinewatch

Apogee's customer that tracks ships at sea with its satellite Brine-3. Thirty-six more Brinewatch satellites launch on Wednesday 16 May 2029.

# Calidra Atoll

A new Apogee station on an atoll near the equator in the Pacific, live from Monday 2 April 2029. Palms line its horizon and tropical rain fades its signals, as at Lantern Atoll. Its code is CAL.

# candidate

A pass worth planning: one over a station that shares a band with the satellite, climbing to a max elevation of 10° or more. passplan 2 tries to give each candidate a dish.

# Cape Wren

Apogee's ground station in the far south: unstaffed, and the windiest site it owns, so its dish stows itself when the wind reaches 72 km/h. Its code is CWR.

# catalog

The passes Apogee predicts for the next 14 days, which customers can book through the booking service. A pass that isn't in it can't be booked. In March 2029 it grows to a whole season ahead.

# Cirrine

Short for Cirrine Atmospherics, Apogee's weather-data customer. Its satellite is Cirrine-1, and twelve more launch on Tuesday 17 April 2029.

# Cirrine-1

The weather satellite operated by Cirrine Atmospherics. Its commands must be checked before the ground sends them up.

# contact

A pass that passplan has given one of a station's dishes: what that station will fly. A pass no dish is free for is dropped, and the station's contacts in time order are its schedule.

# cut-over

The instant one planner takes over from another. passplan 2's is at 00:00 UTC on take-over day, Thursday 28 June 2029: from then on it plans every contact, and the old planner flies only those already in progress.

# deploy

Putting a new version of a program into service: the running one stops, and the new one starts in its place. On Friday 1 June 2029, a test deploy of passplan 2 read 0 planner slots and planned nothing.

# dish

A station's bowl-shaped antenna, which turns to follow a satellite across the sky. Apogee's dishes range from 1.2 m on the test van to 5.4 m at Salar Quieto.

# dish loader

The program at each station that takes the contacts passplan sends and loads them for its dish, so the dish knows which satellite to follow and when. Over a host's thin line it takes two uploads at a time.

# elements

The few numbers that place a satellite on its orbit, such as where it last crossed the equator and how far along its orbit it was at a given moment. passplan predicts every pass from them; a new satellite's first ones come from the company that launched it.

# fails closed

Said of a program that refuses whenever something is wrong or unclear, instead of letting it through. Mirela approves nothing until it does: "LGTM once it fails closed."

# frame

linked: no

One unit of telemetry as it comes down: a satellite's readings from one moment, such as battery, temperature and mode, numbered by its sequence number.

# Furrow-A

One of Furrowlight's crop-imaging satellites. Its photographs arrive in numbered pieces that the ground must put back together. Twenty-four more Furrows join it on Wednesday 27 June 2029.

# Furrowlight

Apogee's customer whose satellites, Furrow-A and Furrow-B, photograph farmland; twenty-four more launch on Wednesday 27 June 2029. Yusra handles its requests and delivery deadlines.

# ground station

A site on Earth with dishes and radios for talking to satellites as they pass overhead. Apogee has five, from Ystmark in the Arctic to Cape Wren in the far south, and will have twenty by the end of June 2029.

# hand-over

What passes from the old planner to the new one at the cut-over: when each dish is free. A dish is free once the old planner's last contact on it has ended and the turnaround after it has passed.

# Hollis

Hollis Vey, who runs operations at Furrowlight and writes every email in capitals, subject lines included.

# host

The owner of the land a new Apogee station stands on, such as a lighthouse authority or a weather observatory. The host gives power, a line to the internet and a pair of hands; nobody from Apogee lives there.

# Ilha Brumosa

A new Apogee station on a misty island in the South Atlantic, live from Monday 2 April 2029; "brumosa" means misty, and its host is a weather observatory. Its code is BRM.

# Ilha Calmaria

A new Apogee station on an island in the calm belt of the tropical Atlantic, with sudden squalls, live from Monday 4 June 2029; "calmaria" means calm. Its code is ICM, and its one dish works on S-band alone.

# Ilvane Spit

A new Apogee station on a gravel spit in the Chukchi Sea, flat to every horizon and locked in sea ice until June, with two dishes, live from Monday 7 May 2029. Its code is ILV.

# Joss

Joss Imbert of Kessick Polytechnic, once a student on Professor Marr's Thistle team, who now leads the flight software of Thistle-2.

# Kessick Polytechnic

The university whose students built Thistle and now Thistle-2. Professor Marr leads the team, and Joss leads Thistle-2's flight software.

# LGTM

Short for looks good to me: what a reviewer writes to approve a change. Mirela adds a condition to hers: "LGTM once it fails closed."

# LOS

The moment a satellite sets out of a station's view and the link ends: the end of a pass. Short for loss of signal.

# max elevation

The highest a satellite climbs above the horizon during a pass, in full the maximum elevation. Low passes are short, far away and weak.

# mutant

A copy of the code with one bug planted on purpose. Good tests fail on every mutant; one that passes them shows a bug the tests would miss.

# Ostrel Island

A new Apogee station on an island in the westerly gales of the southern Indian Ocean, live from Monday 2 April 2029. Its dish stows at 72 km/h, like Cape Wren's. Its code is OST.

# pager

Apogee's alarm service: it checks incoming data against its rules and pages the on-call engineer when something is wrong. Carrying the pager means being the one it wakes.

# pass

linked: no

One flyover: the five to fifteen minutes when one satellite is in view of one station, from the moment it rises into view to the moment it sets. Once the plan gives it a dish, it is a contact.

# pass archive

Apogee's file of every past pass, one per line: its pass ID, AOS, LOS and max elevation. bookingd answers customers' questions about earlier passes from it.

# pass ID

A pass's name: its satellite, its station and the satellite's orbit number in five digits, as in `pipit-1/PEL/04182`. From Monday 9 April 2029, when two passes rise in the same second, passplan 2 puts them in the order of their pass IDs.

# passplan

The part of Apogee's software that predicts every pass, builds the day's schedule and feeds the Board.

# passplan 2

The new passplan, whose rewrite you lead in 2029: it plans a hundred days ahead for twenty stations and eighty satellites, and takes over from the old one on Thursday 28 June.

# Pipit-1

Apogee's small test satellite. A new engineer's first message goes up to it and is echoed back in its status broadcast.

# Pipit-3

Apogee's newest test satellite, launched on 5 April 2027. Finding its first faint broadcasts took a search across radio frequencies.

# planner slot

One of the twelve jobs passplan 2 may run at once, such as predicting one satellite or planning one station: one for each core of its machine.

# planning run

One go of passplan 2 at the whole plan: it predicts the passes, gives them dishes and publishes the schedule. Each run takes the next run number, and its schedule carries that number as its version.

# Port Ellery

Apogee's headquarters, home of the Fishbowl, with a dish on the roof where every new piece of software runs first. Its station code is PEL.

# Professor Marr

Professor Elspeth Marr, who leads the student satellite team at Kessick Polytechnic: first Thistle, which fell out of orbit in 2027, and now Thistle-2.

# safe mode

The minimal state a satellite falls back to after a fault, or when the ground sends it there, with its beacon on and its main instruments off until the ground brings it back. Its frames give its mode as SAFE.

# Salar Quieto

Apogee's ground station on a salt flat 3,840 m up, with the lowest horizon and the thinnest air in the network. Its code is SQT.

# Senior Planning Engineer

Your role from Monday 5 March 2029, still badge No. 039: you lead the rewrite of passplan, and Nkem reviews your work.

# sequence number

A counter that goes up by one with every frame or command. A gap means something was lost, and a repeat means a duplicate or a replay.

# shadow run

passplan 2 planning beside the old planner from Monday 2 April 2029, flying nothing, so that the two schedules can be compared contact by contact every hour.

# station runner

The part of passplan 2 that flies the plan: it waits for each contact's AOS, then starts that contact at its station.

# stow

To park a dish in a safe position, usually pointing straight up, so wind can't damage it. Ystmark's and Cape Wren's dishes stow when the wind reaches 72 km/h, and from April 2029 so does Ostrel Island's.

# take-over day

Thursday 28 June 2029, when passplan 2 takes over from the old planner at 00:00 UTC, on a live day of passes.

# telemetry

The data a satellite sends down about itself, such as battery voltage, temperatures and mode, so the ground can tell how it is doing. Often written TLM.

# telemetry feed

The stream of telemetry lines one ground station hears, one line per frame, on its way to tlmd. A satellite heard by two stations is in both of their feeds.

# test Board

A second Board beside the real one in the Fishbowl, which shows passplan 2's plans while the old passplan still plans the passes Apogee flies. On take-over day, passplan 2's plans go to the real Board.

# the Board

The big screen in the Fishbowl that lists the next passes and the state of every station. Ingrid owns it, and the building's other screens show it too.

# the drill

Nkem's rehearsal of tlmd's new pipeline on Wednesday 7 March 2029: twenty telemetry feeds at once, five from Apogee's stations and fifteen made from where the new stations will stand.

# the Fishbowl

Apogee's ops room at headquarters in Port Ellery, named for its glass walls on three sides. Every pass is run from here, under the Board, the big screen that lists the next passes.

# Thistle-2

Kessick Polytechnic's second student satellite, which leaves its rocket on Thursday 8 March 2029. Its first beacon reaches Port Ellery at 17:08 that day.

# Tirrin Rock

A new Apogee station on a bare rock in the Norwegian Sea, with a lighthouse and two dishes, live from Monday 2 April 2029. Its code is TIR, and it shares polar passes with Ystmark.

# TL;DR

A one-line summary at the top of a message, for readers in a hurry; Yusra starts every message with one. Short for too long; didn't read.

# tlmd

The program that takes in telemetry: it reads every station's telemetry feed, turns each line into a frame and keeps every satellite's newest frame for the Board and the pager. Like Apogee's other services, its name ends in d, for a program that runs on its own in the background.

# turnaround

The two minutes a dish needs after one contact's LOS to turn and settle before it can take the next pass.

# UTC

The world's reference clock, the same everywhere on Earth, with no time zones or daylight saving; everything at Apogee runs on it. Short for Coordinated Universal Time.

# Varsk Island

A new Apogee station at a fishing harbor in the North Pacific, with fog most mornings and a hill to the north, live from Monday 2 April 2029. Its code is VSK.

# Vorrel Island

A new Apogee station on an icebound island in the far south, the south's Ystmark, live from Monday 4 June 2029 with two dishes. Its code is VRL.

# wave

One of the three Mondays in 2029 when five new stations go live at once, at 00:00 UTC: 2 April, 7 May and 4 June.

# window

The stretch of time a plan covers. passplan 2's runs from 00:00 UTC of the current day to 00:00 a hundred days later.

# Ystmark

Apogee's ground station on an Arctic island, the farthest north in the network. Its code is YST.
