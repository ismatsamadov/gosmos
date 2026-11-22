# AOS

The moment a satellite rises into a station's view and the link can begin: the start of a pass. Short for acquisition of signal.

# azimuth

The compass direction a dish points in, or a satellite lies in, measured in degrees clockwise from north: 90° is east, 180° south and 270° west.

# booking

A pass a customer has reserved at one of Apogee's stations: which satellite, which station, and when it rises and sets. Yusra's booking export lists every one.

# booking sheet

Furrowlight's bookings as Hollis Vey exports them from his spreadsheet: a text file with one row per booking, under a header row in capitals that names the columns.

# calibration

Checking a dish against a known reference and correcting it, so it points exactly where it should and measures signals truly. The test van carries a small dish to stations for this.

# Cape Wren

Apogee's ground station in the far south: unstaffed, and the windiest site it owns, so its dish stows itself when the wind reaches 72 km/h. Its station code is CWR.

# Chiquita

The nickname of Salar Quieto's smaller dish, 2.4 m across. Its big sister is La Grande; the controller log calls them CH and LG.

# chunk

One numbered piece of a file coming down from a satellite. An image arrives as thousands of chunks, so the missing ones can be counted and asked for again.

# controller

The computer at Salar Quieto that steers La Grande and Chiquita, turning each dish to follow a satellite. It writes every pointing reading of both dishes in one log.

# controller log

The file in which Salar Quieto's dish controller, the computer that steers La Grande and Chiquita, writes one line for every pointing reading. Paz brought 200,000 of its lines down from the station.

# correction

The small angles a dish controller adds to every azimuth and elevation it is asked to point at, so that its dish points true. La Grande's controller keeps its correction in `correction.xml`, and calibration day replaces it.

# delivery report

Furrowlight's own account of what reached it: an image report for one picture, or a station report for one station's day. Each comes with a version number and a long ID.

# dish

A station's bowl-shaped antenna, which turns to follow a satellite across the sky. Apogee's dishes range from 1.2 m on the test van to 5.4 m at Salar Quieto.

# elevation

How high a dish points, or a satellite stands, above the horizon, as an angle: 0° at the horizon, 90° straight overhead.

# Furrow-B

One of Furrowlight's two crop-imaging satellites, flying in the same orbit as Furrow-A, half an orbit behind it.

# Furrowlight

A customer whose satellites photograph farmland, so farmers can plan the work in their fields. Hollis runs its operations and writes his subject lines in capitals.

# gctl

The command-line tool Apogee's engineers use to inspect passes and station records. In this season you build its offline file and reporting tools.

# ground station

A site on Earth with dishes and radios for talking to satellites as they pass overhead. Apogee has five, from Ystmark in the Arctic to Cape Wren in the far south.

# Hollis

Hollis Vey, operations manager at Furrowlight, the customer whose two satellites photograph farmland. He sends his bookings as a spreadsheet and writes in capitals.

# La Grande

The nickname of Salar Quieto's big dish, 5.4 m across, the largest in Apogee's network. Its small sister is Chiquita.

# Lantern Atoll

Apogee's ground station on a tropical island on the equator, run on diesel power with a caretaker on site. Its station code is LAN.

# LOS

The moment a satellite sets out of a station's view and the link ends: the end of a pass. Short for loss of signal.

# manifest

The list a parcel carries in `manifest.xml`: the dish it is about, when the controller made it, and every file inside with its size in bytes.

# mask

linked: no

The lowest elevation at which a station works with a satellite; below it, hills, trees or radio noise get in the way. Salar Quieto's is 3°, the lowest in the network.

# offset

How far a clock is set from UTC, in hours and minutes ahead (`+`) or behind (`-`). The van's laptop keeps the salt flat's offset, `-07:00`, all year round.

# parcel

One file that packs several others together, as Salar Quieto's dish controller sends them: a zip for a dish's settings, a gzip'd tar for its log, each with a manifest listing what is inside.

# pass

linked: no

One flyover: the five to fifteen minutes when one satellite is in view of one station, from the moment it rises into view to the moment it sets. Also called a contact.

# Port Ellery

Apogee's headquarters, home of the Fishbowl, with a dish on the roof where every new piece of software runs first. Its station code is PEL.

# Quebrada Añil

The ravine road that is the only way up to Salar Quieto, 3,840 m above the sea.

# Salar Quieto

Apogee's ground station on a high salt flat, 3,840 m up, with two dishes nicknamed La Grande and Chiquita. Its station code is SQT.

# stow

To park a dish in a safe position, usually pointing straight up, so wind can't damage it. Ystmark's and Cape Wren's dishes stow when the wind reaches 72 km/h.

# sweep

A slow scan of a dish back and forth across a target during calibration, its position logged many times a second. The dish controller writes a whole sweep on one long line.

# test van

Apogee's station on wheels: a van based at Port Ellery with one 1.2 m dish, driven out to other sites for calibration.

# the Fishbowl

Apogee's ops room at headquarters in Port Ellery, named for its glass walls on three sides. Every pass is run from here, under the Board, the big screen that lists the next passes.

# TL;DR

A one-line summary at the top of a message, for readers in a hurry; Yusra starts every message with one. Short for too long; didn't read.

# UTC

The world's reference clock, the same everywhere on Earth, with no time zones or daylight saving; everything at Apogee runs on it. Short for Coordinated Universal Time.

# Ystmark

Apogee's Arctic ground station, on an island where storms cut the fiber to the mainland and its data has to crawl over a slow satellite link instead. Its station code is YST.
