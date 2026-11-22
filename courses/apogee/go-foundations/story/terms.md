# anemometer

An instrument that measures wind speed. Cape Wren's software uses its reading to decide when to park the dish.

# anomaly

Anything in a satellite's data or a station's state that is outside its limits or unexpected, such as a battery under 6.6 V. The pager's rules decide which ones wake someone.

# AOS

The moment a satellite comes into a station's view and the link can start: the beginning of a pass. Short for acquisition of signal.

# azimuth

The compass direction from a dish to a satellite, in degrees clockwise from north: 90° is east, 180° south, 270° west. The rise azimuth is where a pass begins.

# backhaul

The link that carries a station's data back to headquarters. When storms cut Ystmark's fiber cable, it falls back to a slow backhaul through a satellite.

# band

linked: no

A range of radio frequencies, such as UHF, S-band or X-band. Ops also sorts batteries into bands, from FULL down to CRITICAL.

# beacon

A short status message a satellite broadcasts by itself, without being asked, with readings such as its battery, temperature and mode. Pipit-1 sends one every 30 seconds.

# Brine-2

Brinewatch's satellite that never left its release device after launch. It cannot take part in a pass and should not be in the working fleet registry.

# Brine-3

Brinewatch's working satellite, which tracks ships. Its orbit stays too far south for Ystmark to see it.

# Brinewatch

Apogee's customer that tracks ships at sea with its satellite Brine-3. Brine-1 was lost at launch, and Brine-2 never left its deployer.

# calibration

Checking a dish against a known reference and correcting it, so it points exactly where it should and measures signals truly. The test van carries a small dish to stations for this.

# call sign

The short name a radio station uses to identify itself on the air, written only with capital letters, digits, hyphens and underscores.

# checksum

A short number computed from a piece of data and sent or kept with it. Compute it again later, and a different answer means the data was damaged or changed.

# Chiquita

The nickname of Salar Quieto's smaller dish, 2.4 m across. Its big sister is La Grande.

# chunk

One numbered piece of a file coming down from a satellite. Chunks are numbered from 0, so the missing ones can be listed and asked for again.

# Cirrine

Short for Cirrine Atmospherics, Apogee's weather-data customer. Its satellite is Cirrine-1.

# Cirrine-1

The weather satellite operated by Cirrine Atmospherics. Its commands must be checked before the ground sends them up.

# CRC

A standard kind of checksum, good at catching the damage radio noise does: a chunk whose CRC no longer matches its data was damaged on the way. Short for cyclic redundancy check.

# cubesat

A small satellite built from standard units, cubes 10 cm on a side. Pipit-1 is a 3U, three units long, about the size of a loaf of bread.

# dB

A scale for how much stronger or weaker one signal is than another: every 3 dB doubles or halves the power, and every 10 dB is a factor of ten. Short for decibel.

# dBm

A signal's power on the decibel scale, counted from one milliwatt: 0 dBm is 1 milliwatt and 30 dBm is 1 watt. A signal from orbit often arrives near −100 dBm.

# deployer

The spring-loaded box on a rocket that carries small satellites to orbit and pushes them out. Brine-2 never left its deployer.

# DETUMBLE

The mode of a satellite fresh off its rocket and still spinning, while it brings itself steady. Until then, ops listens only to its beacon.

# dish

A station's bowl-shaped antenna, which turns to follow a satellite across the sky. Apogee's dishes range from 1.2 m on the test van to 5.4 m at Salar Quieto.

# Doppler shift

The change in a signal's frequency caused by the satellite's motion: higher as it comes toward you, lower as it moves away, like a passing siren. Pipit-1's beacon shifts by up to about 10 kHz.

# downlink

The radio link from a satellite down to the ground, the data that comes down it, and the part of Apogee's software that receives and delivers that data.

# elevation

How high a satellite stands above the horizon as seen from a station, as an angle: 0° at the horizon, 90° straight overhead.

# footprint

The patch of Earth from which a satellite can be seen above a given elevation at one moment. Any station inside it can talk to the satellite, and higher orbits have bigger footprints.

# frame

linked: no

One unit of telemetry as it comes down: a satellite's readings from one moment, such as battery, temperature and mode, numbered by its sequence number.

# Furrow-A

One of Furrowlight's two crop-imaging satellites. Its photographs arrive in numbered pieces that the ground must put back together.

# Furrow-B

Furrowlight's second crop-imaging satellite, following Furrow-A in the same orbit about half a trip around Earth behind it.

# Furrowlight

Apogee's customer whose two satellites, Furrow-A and Furrow-B, photograph farmland. Yusra handles its requests and delivery deadlines.

# go/no-go

The poll ops runs before every pass: each question, from whether the station is up to whether the battery is healthy, gets a GO or a NO-GO, and a single NO-GO calls the pass off.

# ground station

A site on Earth with dishes and radios for talking to satellites as they pass overhead. Apogee has five, from Ystmark in the Arctic to Cape Wren in the far south.

# groundctl

The Go code behind Apogee's ground stations, which your own code grows into. Its parts predict passes, check telemetry, sign commands, bring files down and wake the on-call engineer.

# HMAC

A signature computed from a message and a secret key that only Apogee and the satellite share, so nobody without the key can forge one. Short for hash-based message authentication code; Apogee uses HMAC-SHA256.

# housekeeping

A satellite's data about its own health, such as battery, temperatures and mode, as opposed to the images or measurements it flies to make.

# Hz

Cycles per second, the unit of a radio signal's frequency. Short for hertz; a kHz is a thousand of them and an MHz a million.

# incident review

A meeting after something has gone wrong, where the team works out what happened, why, and what to change, without blaming anyone.

# inclination

The tilt of an orbit to the equator, in degrees. It caps how far north or south a satellite flies, which is why Thistle, at 51.6°, never rises over Ystmark.

# Ka-band

Radio frequencies higher than X-band, which carry data faster still and fade even more in rain. Paz is testing a Ka-band radio on La Grande.

# Kessick Polytechnic

The university whose students built Thistle. Professor Marr leads the team, and Joss works with the readings it sends down.

# kHz

A thousand cycles per second: Doppler shift moves Pipit-1's beacon by up to about 10 kHz as it passes. Short for kilohertz.

# La Grande

The nickname of Salar Quieto's big dish, 5.4 m across, the largest in Apogee's network. Its small sister is Chiquita.

# LEOP

A satellite's first days in orbit after launch, while ops finds it, checks it over and brings it into normal work. Short for launch and early orbit phase.

# link budget

The sum, in dB, of every gain and loss between a transmitter and a receiver, which tells you whether a signal will arrive strong enough to be heard.

# LOS

The moment a satellite sets out of a station's view and the link ends: the end of a pass. Short for loss of signal.

# margin

How much stronger a received signal is than the weakest the receiver can still use, in dB. A positive margin means the link works, and more is safer.

# mask

The lowest elevation a station can use, set by trees, hills or radio noise on its horizon. Salar Quieto's is 3°, and Lantern Atoll's palm trees push its to 8°.

# maximum elevation

The highest a satellite climbs above the horizon during a pass, often written max elevation. Low passes are short, far away and weak.

# Mbps

How fast a link carries data, in millions of bits per second: Furrow-A's X-band link runs at 150 Mbps. Short for megabits per second.

# MHz

A million cycles per second: Pipit-1's beacon is at 437.250 MHz, and Pipit-2 talks at 2245 MHz. Short for megahertz.

# mode

The state a satellite reports it is in, which decides what ops does on a pass: NOMINAL, SAFE, DETUMBLE or LEOP.

# mutant

A copy of the code with one bug planted on purpose. Good tests fail on every mutant; one that passes them shows a bug the tests would miss.

# mV

A thousandth of a volt: 8200 mV is 8.2 V, a full battery. Short for millivolt.

# NOMINAL

A satellite's normal mode: everything works as planned, and its passes go ahead as booked.

# on-call

Being the engineer the pager wakes, at any hour, when something breaks. Engineers take turns, a week at a time.

# operator

linked: no

The company or team that owns and flies a satellite: Kessick Polytechnic is Thistle's operator. Operators book passes from Apogee.

# ops

The team that runs every pass, from the pre-pass checklist to LOS, led by Ingrid in the Fishbowl. Short for operations.

# ops room

The room where ops runs the passes. At Apogee it is the Fishbowl, at headquarters in Port Ellery.

# orbit elements

Numbers describing a satellite's orbit at a stated time. A planner uses them to predict when a satellite will come into a station's view.

# page

linked: no

An alert that wakes the on-call engineer, whatever the hour; to page someone is to send one.

# pager

Apogee's alarm service: it checks incoming data against its rules and pages the on-call engineer when something is wrong. Carrying the pager means being the one it wakes.

# pass

linked: no

One flyover: the five to fifteen minutes when one satellite is in view of one station, from AOS to LOS. Also called a contact.

# pass ID

A pass's name: its satellite, its station and the satellite's orbit number in five digits, as in `pipit-1/PEL/04182`.

# passplan

The part of Apogee's software that predicts every pass, builds the day's schedule and feeds the Board.

# path loss

How much a radio signal weakens on its way from satellite to dish, in dB. It grows with the distance and with the frequency.

# payload

The equipment that does a satellite's main job, such as a camera or a weather instrument. It may be switched off while the satellite keeps its basic systems running.

# period

linked: no

The time a satellite takes to go once around the Earth: about 95 minutes for Pipit-1, at 525 km.

# Pipit-1

Apogee's small test satellite. A new engineer's first message goes up to it and is echoed back in its status broadcast.

# Pipit-2

Apogee's radio test satellite. The team uses it to check commands protected with a shared secret key.

# Pipit-3

Apogee's newest test satellite, which launches on 5 April 2027. Its first broadcasts are faint, so the team must search the radio frequencies to find it.

# polar orbit

An orbit tilted close to 90° to the equator, so it crosses near both poles and, as the Earth turns beneath it, passes over every place. Ystmark sees such satellites on almost every revolution.

# prediction

A forecast, worked out from a satellite's orbit, of when it will pass over a station and how high it will climb, minute by minute.

# Professor Marr

Professor Elspeth Marr, who leads Thistle's student team at Kessick Polytechnic. She wants its last broadcasts heard before it returns to Earth's atmosphere.

# Quebrada Añil

The ravine road that is the only way up to Salar Quieto, 3,840 m above the sea.

# radome

A dome over a dish that keeps out wind, snow and ice but lets radio signals through. Ystmark's two dishes sit under heated radomes.

# rain fade

The weakening of a high-frequency signal, such as X-band, as rain absorbs and scatters it. It is why Lantern Atoll, with its tropical downpours, has no X-band.

# range rate

How fast the distance between dish and satellite is changing, in km/s: negative while the satellite approaches, positive as it moves away.

# re-enter

To fall out of orbit into the air, where the heat of the fall burns a small satellite up. Low orbits sink a little every day until the satellite re-enters.

# recovery pass

A pass spent bringing a satellite out of safe mode: ops reads what went wrong and sends the commands that return it to normal work.

# replay

Sending an old, properly signed command a second time to make a satellite repeat it. Sequence numbers stop it, because an old number is refused.

# replay guard

The part of the uplink that refuses replayed commands: it accepts a command only if its sequence number is higher than the last one accepted.

# revolution

One trip around the Earth; another word for one orbit.

# S-band

Radio frequencies around 2,000 to 4,000 MHz, used for commands going up and housekeeping coming down. Pipit-2 talks at 2245 MHz.

# SAFE

Safe mode: the minimal state a satellite falls back to after a fault, with its beacon on and its main instruments off, until ops brings it back.

# sequence number

A counter that goes up by one with every frame or command. A gap means something was lost, and a repeat means a duplicate or a replay.

# signature

A string of digits attached to a command that proves who sent it and that nothing in it changed on the way. A satellite refuses any command whose signature doesn't check out.

# slant range

The straight-line distance from a dish to a satellite: for Pipit-1, 525 km when it is straight overhead and about 2,640 km at the horizon.

# SSO

A near-polar orbit that passes over each place at about the same local time every day, so the light below looks alike from pass to pass. Short for sun-synchronous orbit.

# stand-up

A short daily team meeting, held standing so it stays short, where each person says what they are working on.

# stow

To park a dish in a safe position, usually pointing straight up, so wind can't damage it. Ystmark's and Cape Wren's dishes stow when the wind reaches 72 km/h.

# sync bytes

Fixed values at the very start of a frame that tell the receiver where the frame begins. Every Pipit beacon starts with the two letters AP.

# telecommand

A command sent up to a satellite, such as a change of mode. At Apogee every telecommand is signed, so the satellite knows it truly came from Apogee.

# telemetry

The data a satellite sends down about itself, such as battery voltage, temperatures and mode, so the ground can tell how it is doing. Often written TLM.

# test van

Apogee's station on wheels: a van based at Port Ellery with one 1.2 m S-band dish, driven out to other sites for calibration.

# the Board

The big screen in the Fishbowl that lists the next passes and the state of every station. Ingrid owns it.

# the Fishbowl

Apogee's ops room at headquarters in Port Ellery, named for its glass walls on three sides. Every pass is run from here, under the Board.

# Thistle

The student satellite built at Kessick Polytechnic. Its battery readings grow unreliable as its orbit gets lower, and the team listens through its final pass.

# TL;DR

A one-line summary at the top of a message, for readers in a hurry; Yusra starts every message with one. Short for too long; didn't read.

# track

To keep a dish pointed at a satellite as it crosses the sky. A radio that tracks also follows the satellite's Doppler shift by itself.

# UHF

At Apogee, the slice of radio frequencies around 435 to 438 MHz, where small satellites such as Pipit-1 and Thistle send their beacons. Every station but Ystmark hears them on yagis, antennas made of rows of short metal rods. Short for ultra high frequency.

# uplink

The radio link from the ground up to a satellite, the commands sent up it, and the part of Apogee's software that signs and sends them.

# UTC

The world's reference clock, the same everywhere on Earth, with no time zones or daylight saving; everything at Apogee runs on it. Short for Coordinated Universal Time.

# X-band

Radio frequencies around 8,000 MHz that carry big files fast, such as Furrow-A's images at 150 Mbps, but fade in heavy rain.
