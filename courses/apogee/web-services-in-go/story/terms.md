# account page

The small signed-in page beside the booking service, at `/account/`, where a customer's people see their own bookings in a browser, and can cancel one, without asking an engineer.

# AOS

The moment a satellite comes into a station's view and the link can start: the beginning of a pass. Short for acquisition of signal.

# booking

A pass a customer has reserved at one of Apogee's stations, with an ID such as `bk-4494`: which satellite, which station, and when the pass starts and ends.

# booking service

The program Apogee is building so that customers' programs can find passes, book them and hear when their data has landed, day and night. It opens to every customer on 1 February 2028.

# bookingd

The program that is the booking service: it books passes for customers and tells them when their data lands. Like Apogee's other services, its name ends in d, for a program that runs on its own in the background.

# Brine-3

Brinewatch's working satellite, which tracks ships. Its orbit stays too far south for Ystmark to see it.

# Brinewatch

Apogee's ship-tracking customer. Its satellite, Brine-3, picks up the radio signals that ships send, all over the oceans.

# Cape Wren

Apogee's ground station in the far south: unstaffed, and the windiest site it owns, so its dish stows itself when the wind reaches 72 km/h. Its code is CWR.

# catalog

The passes Apogee predicts for the next 14 days, which customers can book through the booking service. A pass that isn't in it can't be booked.

# Cirrine

Short for Cirrine Atmospherics, Apogee's weather-data customer. Its satellite is Cirrine-1.

# command file

The commands a customer wants sent up to its satellite during one booked pass, one to a line, in a plain text file it attaches to the booking. The station sends them up while the satellite is in view.

# customer guide

Yusra's pages for customers' engineers: how their programs should talk to the booking service, and what each of its answers means.

# deploy

Putting a new version of a program into service: the running copy stops and the new one starts in its place, so anything the old one kept only in memory is gone.

# Furrow-A

One of Furrowlight's two crop-imaging satellites. Its photographs arrive in numbered pieces that the ground must put back together.

# Furrow-B

Furrowlight's second crop-imaging satellite, following Furrow-A in the same orbit about half a trip around Earth behind it.

# Furrowlight

Apogee's customer whose two satellites, Furrow-A and Furrow-B, photograph farmland. Yusra handles its requests and delivery deadlines.

# go/no-go

The poll ops runs before every pass: each question, from whether the station is up to whether the battery is healthy, gets a GO or a NO-GO, and a single NO-GO calls the pass off.

# HMAC

A signature computed from a message and a secret key, so nobody without the key can forge one. Short for hash-based message authentication code; Apogee uses HMAC-SHA256.

# Hollis

Hollis Vey, who runs operations at Furrowlight and writes every email in capitals, subject lines included.

# landing

The moment a booking's file, such as one of Furrow-A's images, has come down from the satellite and arrived at Apogee in full, ready for the customer. The booking service tells the customer about each one.

# LGTM

Short for looks good to me: what a reviewer writes to approve a change. Mirela signs off her reviews with "LGTM once it fails closed."

# LOS

The moment a satellite sets out of a station's view and the link ends: the end of a pass. Short for loss of signal.

# mutant

A copy of the code with one bug planted on purpose. Good tests fail on every mutant; one that passes them shows a bug the tests would miss.

# on-call

Being the engineer the pager wakes, at any hour, when something breaks. Engineers take turns, a week at a time.

# ops

The team that runs every pass, from the pre-pass checklist to LOS, led by Ingrid in the Fishbowl. Short for operations.

# pager

Apogee's alarm service: it checks incoming data against its rules and pages the on-call engineer when something is wrong. Carrying the pager means being the one it wakes.

# pass

linked: no

One flyover: the five to fifteen minutes when one satellite is in view of one station, from the moment it rises into view to the moment it sets. Also called a contact.

# pass ID

A pass's name: its satellite, its station and the satellite's orbit number in five digits, as in `pipit-1/PEL/04182`.

# Pipit-1

Apogee's small test satellite. A new engineer's first message goes up to it and is echoed back in its status broadcast.

# plan

linked: no

How many bookings a customer may make in one calendar month under its contract with Apogee, counted in UTC and started again at midnight on the 1st. Apogee's own bookings have no plan.

# planning page

Furrowlight's own web page, where its engineers choose which passes to book. Its script asks the booking service for passes straight from the browser.

# Salar Quieto

Apogee's ground station on a salt flat 3,840 m up, with the lowest horizon and the thinnest air in the network. Its code is SQT.

# the beta

The weeks before the booking service opens to every customer, when a few customers try it with their own programs and tell Yusra what goes wrong.

# the Fishbowl

Apogee's ops room at headquarters in Port Ellery, named for its glass walls on three sides. Every pass is run from here, under the Board, the big screen that lists the next passes.

# Thistle

The student satellite built at Kessick Polytechnic. Its orbit sank lower all spring, and in July 2027 it fell out of orbit and burned up.

# TL;DR

A one-line summary at the top of a message, for readers in a hurry; Yusra starts every message with one. Short for too long; didn't read.

# UHF

Radio frequencies from 300 to 3,000 MHz; Apogee uses the stretch near 435 MHz, where small satellites such as Pipit-1 send beacons and slow data. Port Ellery and Lantern Atoll send on it, Salar Quieto and Cape Wren only listen, and Ystmark has none. Short for ultra high frequency.

# UTC

The world's reference clock, the same everywhere on Earth, with no time zones or daylight saving; everything at Apogee runs on it. Short for Coordinated Universal Time.

# weather mast

The pole at a station that carries the instruments measuring its wind and temperature. When a mast falls silent, the weather service has no reading to give for that station.

# weather service

The program Tomás wrote in 2024 that reads every station's weather mast once a minute and gives other programs, such as the go/no-go poll, the latest wind and temperature.

# X-band

Radio frequencies around 8,000 MHz that carry big files fast, such as Furrow-A's images at 150 Mbps, but fade in heavy rain.

# Ystmark

Apogee's ground station on an Arctic island, the farthest north in the network. Its code is YST.
