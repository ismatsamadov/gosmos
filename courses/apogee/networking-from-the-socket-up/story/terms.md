# azimuth

The compass direction from a dish to a satellite, in degrees clockwise from north: 90° is east, 180° south, 270° west. The rise azimuth is where a pass begins.

# backhaul

The link that carries a station's data back to headquarters. When storms cut Ystmark's fiber cable, it falls back to a slow backhaul through a satellite.

# banner

The one line of text a dish controller sends the moment a connection opens, before anything else: its name, the word READY and its protocol version, as in `YST-1 READY PROTO 1`.

# booking service

Apogee's web service, open to every customer since 1 February 2028, through which customers' programs book passes and hear when their data has landed. Its program is called bookingd.

# calibration

Checking a dish against a known reference and correcting it, so it points exactly where it should and measures signals truly. The test van carries a small dish to stations for this.

# checksum

A short number computed from a piece of data and sent or kept with it. Compute it again later, and a different answer means the data was damaged or changed.

# chunk

One numbered piece of a file coming down from a satellite. Chunks are numbered from 0, so the missing ones can be listed and asked for again; an HTTP chunk, on the station server, is instead one piece of an upload or an answer, sent with its size first.

# Cirrine

Short for Cirrine Atmospherics, Apogee's weather-data customer. Its satellite is Cirrine-1.

# Cirrine-1

The weather satellite operated by Cirrine Atmospherics. Its commands must be checked before the ground sends them up.

# collector

Headquarters' program at Port Ellery that fetches each recording from the station server in the hut, as soon as it hears one has arrived.

# CRC

A standard kind of checksum, good at catching the damage radio noise does: a chunk whose CRC no longer matches its data was damaged on the way. Short for cyclic redundancy check.

# deploy

Putting a new version of a program into service: the running one stops, and the new one starts in its place.

# dish

A station's bowl-shaped antenna, which turns to follow a satellite across the sky. Apogee's dishes range from 1.2 m on the test van to 5.4 m at Salar Quieto.

# dish controller

The computer in a dish's radome that turns the dish. Ystmark's new ones, YST-1 for Dish 1 and YST-2 for Dish 2, talk in raw bytes over a network connection.

# dishd

The dish service you write this season: the program that talks to Ystmark's dish controllers and shows both dishes on its watch port. It runs at Port Ellery.

# downlink

The radio link from a satellite down to the ground, the data that comes down it, and the part of Apogee's software that receives and delivers that data.

# elevation

How high a satellite stands above the horizon as seen from a station, as an angle: 0° at the horizon, 90° straight overhead.

# feed

The station server's stream of new arrivals: the collector at Port Ellery keeps one request open on it and hears of each recording the moment it is stored.

# fiber

The fiber-optic cable that links Ystmark to the mainland, about 40 ms there and back. Storms cut it, and then only the slow backhaul is left.

# firmware

The program built into a device, such as a dish controller, that makes it work. YST-1 and YST-2 run firmware 2.4.

# Furrow-A

One of Furrowlight's two crop-imaging satellites. Its photographs arrive in numbered pieces that the ground must put back together.

# Furrow-B

Furrowlight's second crop-imaging satellite, following Furrow-A in the same orbit about half a trip around Earth behind it.

# Furrowlight

Apogee's customer whose two satellites, Furrow-A and Furrow-B, photograph farmland. Yusra handles its requests and delivery deadlines.

# go/no-go

The poll ops runs before every pass: each question, from whether the station is up to whether the battery is healthy, gets a GO or a NO-GO, and a single NO-GO calls the pass off.

# Hollis

Hollis Vey, who runs operations at Furrowlight and writes every email in capitals, subject lines included.

# jog

A small move of a dish from where it already points, along one axis, in hundredths of a degree: Tomás's calibration jog was 0.35° of azimuth. Run twice, a jog moves the dish twice.

# live port

The recorder's own network port in the hut. During each pass it sends the satellite's telemetry from there to the live view at Port Ellery, and it answers there when asked whether it is recording.

# live view

The program at Port Ellery that shows Ingrid each pass at Ystmark as it comes down: every frame the recorder sends, in order and each once, and what the recorder says when the frames stop.

# LOS

The moment a satellite sets out of a station's view and the link ends: the end of a pass. Short for loss of signal.

# midnight sun

The weeks when the sun never sets, so far north. At Ystmark it stays up all night from 19 April.

# ops room

The room where ops runs the passes. At Apogee it is the Fishbowl, at headquarters in Port Ellery.

# pass

linked: no

One flyover: the five to fifteen minutes when one satellite is in view of one station, from the moment it rises into view to the moment it sets. Also called a contact.

# pass ID

A pass's name: its satellite, its station and the satellite's orbit number in five digits, as in `pipit-1/PEL/04182`.

# Pipit-1

Apogee's small test satellite. A new engineer's first message goes up to it and is echoed back in its status broadcast.

# Pipit-2

Apogee's radio test satellite. The team uses it to check commands protected with a shared secret key.

# polar orbit

An orbit tilted close to 90° to the equator, so it crosses near both poles and, as the Earth turns beneath it, passes over every place. Ystmark sees such satellites on almost every revolution.

# Port Ellery

Apogee's headquarters, home of the Fishbowl, with a 3.0 m dish on its roof: the station where every new piece of software is tried first. Its code is PEL.

# radome

A dome over a dish that keeps out wind, snow and ice but lets radio signals through. Ystmark's two dishes sit under heated radomes.

# recorder

The program in the hut's rack that records each pass at Ystmark as it comes down. It is a Go program, and from June it sends every recording to the station server.

# recording

Everything one pass brought down, kept as one file by the recorder. A pass of Furrow-A makes a recording of about 11 GB.

# replay

Sending an old, properly signed command a second time to make a satellite repeat it. Sequence numbers stop it, because an old number is refused.

# revolution

One trip around the Earth; another word for one orbit.

# sequence number

A counter that goes up by one with every frame or command. A gap means something was lost, and a repeat means a duplicate or a replay.

# signature

A string of digits attached to a command that proves who sent it and that nothing in it changed on the way. A satellite refuses any command whose signature doesn't check out.

# simulator

A program that behaves exactly like something else. The maker of Ystmark's dish controllers sent one on a stick, so dishd can be tried on the bench before the real controllers are installed.

# station server

The server in Ystmark's hut that keeps every finished recording until headquarters collects it, so a long outage of the fiber loses nothing. You build it on a bare socket.

# stow

To park a dish in a safe position, usually pointing straight up, so wind can't damage it. Ystmark's and Cape Wren's dishes stow when the wind reaches 72 km/h.

# telemetry

A satellite's readings of its own health, such as its battery voltage, its temperature and its mode, sent down in numbered frames during each pass.

# the Fishbowl

Apogee's ops room at headquarters in Port Ellery, named for its glass walls on three sides. Every pass is run from here, under the Board.

# the hut

Ystmark's one building, below the ridge with the two radomes: the bench computer, Sunniva's terminal, the recorders' rack and the coffee machine.

# the island's line

Ystmark's link to the mainland: normally the fiber, about 40 ms there and back; in storms only the backhaul, about 0.7 s, enough for commands but never for big files.

# the manual

The maker's book for the new dish controllers: how to connect to them and what every byte they send means. Its page one is wrong.

# TL;DR

A one-line summary at the top of a message, for readers in a hurry; Yusra starts every message with one. Short for too long; didn't read.

# UTC

The world's reference clock, the same everywhere on Earth, with no time zones or daylight saving; everything at Apogee runs on it. Short for Coordinated Universal Time.

# watch port

dishd's listening port, where anyone allowed can connect and watch both dishes: a greeting line, then one line per dish every second.

# wind lock

A dish controller's own guard against storms: at a wind of 72 km/h or more it stows its dish and refuses every command to unstow or move it. The dish stays stowed until it is told to unstow.

# Ystmark

Apogee's Arctic ground station, on an island at 77.8° N, with two 3.7 m dishes under radomes. It sees every polar orbit on almost every revolution. Its code is YST.
