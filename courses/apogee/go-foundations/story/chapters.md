# first-contact

## Briefing

Day one. At Apogee's headquarters in Port Ellery, Ingrid runs the ops room, where every pass is planned to the minute. This week she needs two things from you: the message the room sends up to Pipit-1, one of Apogee's own satellites, and a roll call of all five stations for Friday. By then, the whole room will be listening for your words to come back from orbit.

## Debrief

Friday morning, Pipit-1's beacon brings all 29 bytes of your message back down: *Apogee ground calling Pipit-1*. Ingrid puts your roll call on the Board, four stations up and Cape Wren stowed for wind, and writes one line in the log: "Echo received. New engineer is GO."

# numbers-from-orbit

## Briefing

Week two. The onboarding sheet holds the numbers every tool in the ops room is about to share: altitudes, radio frequencies, battery voltages, orbit times. All of it was entered by hand, and no two copies of a number agree, down to the size of the Earth. The sheet goes live this week, and Ingrid and Tomás need it right, because a wrong number there becomes wrong everywhere.

## Debrief

The onboarding sheet is live, with one Earth radius, 6371 km, and orbit periods it works out for itself: 94.99 minutes for Pipit-1, from its altitude and one constant. Tomás checks it against his notebook: "Altitude in, period out. Now the sheet can't disagree with itself!"

# the-pass-board

## Briefing

Week three. The Board, the big screen in the Fishbowl, Apogee's glass-walled ops room, lists the next passes in a jumble: AOS times like 4:7, frequencies with random decimals, columns that wander. Ingrid wants it readable at a glance before Furrowlight, the customer whose satellites photograph farmland, walks in on Friday at nine.

## Debrief

At nine on Friday, Furrowlight's team walks into the Fishbowl and reads the Board without help: four passes in straight columns, AOS times like 04:12, maximum elevations to one decimal. Ingrid allows herself one line in the log: "09:00 visit. Board read at a glance."

# orbit-math

## Briefing

Week four. Every tool at Apogee has copied the orbit formulas from Tomás's battered notebook, and no two copies agree. One read a nine-minute pass as nine hours and kept Port Ellery booked for all nine. Before Pipit-2's next pass over Salar Quieto, Tomás wants the formulas written once and every number for the pass worked out from them.

## Debrief

The notebook is a toolbox now. As Pipit-2 rises over Salar Quieto, Tomás holds one card with every number your code worked out: a 95.50-minute orbit, 1,815 km to the satellite at 10° up, 164.64 dB lost on the way down. He taps it: "My notebook does its own arithmetic!"

# go-or-no-go

## Briefing

Week five. Before every pass, the ops room runs a go/no-go poll: is the station up, is its dish free to move, is the wind safe, will the satellite climb high enough above the horizon, is its battery healthy? Ingrid asks each question by hand, and a skipped one can cost a pass. She wants the poll automated before Furrowlight's Furrow-A makes its evening pass over Salar Quieto.

## Debrief

Twenty minutes before Furrow-A rises over Salar Quieto, your poll checks the station, the wind at 18 km/h, a highest point of 41.5° and a 7.9 V battery, and calls GO. Ingrid reads it out to the Fishbowl: "Poll says GO. AOS in twenty."

# counting-down-to-aos

## Briefing

Launch week. On Monday, Pipit-3, Apogee's newest satellite, rides to orbit on a rocket crowded with other small satellites, and nobody knows yet which new dot in the sky is theirs. Its first passes are only predictions, its battery is low and its beacon is faint. Tomás, Sunniva and Ingrid will listen on every pass until Pipit-3 answers.

## Debrief

Pipit-3 answers. Your sweep finds its faint beacon 3.5 kHz above the predicted frequency, and its battery will be full in ten orbits. Tomás pulls off his headphones, grinning: "Our newest pipit, singing a little sharp!"

# rules-you-can-hand-around

## Briefing

The week after launch. Ingrid's checklist has become one long chain of customer exceptions, and only Nkem dares to edit it. Furrowlight needs passes high above the horizon; Cirrine needs a different minimum. Before tonight's bookings, the team needs rules it can keep separately and run together without skipping a refusal.

## Debrief

Tonight's poll checks every customer's rule in order and stops at the first refusal. Each rule keeps its own minimum, and each checklist its own count. Ingrid closes the booking sheet: "Every rule asked. Nobody skipped. GO."

# modeling-the-network

## Briefing

Week eight. Everything Apogee knows about its network, five stations and eight satellites with their orbits, lives in a spreadsheet nobody trusts anymore, and every tool copies its mistakes. A wrong mask there means booking passes a dish can't use. Nkem wants the network written down in code instead: each fact once, in shapes every later tool can build on.

## Debrief

Apogee's network now lives in code: every station as one record, every satellite with its orbit inside it, every pass with an AOS and an LOS, and a test van that shares the stations' dish record. Nkem approves the review: "Good. Everything we build next stands on this."

# the-days-schedule

## Briefing

Week nine. Three Pipits and five customer satellites add up to dozens of passes a day. Predictions arrive one at a time, in no particular order, yet the Board must list the next few by AOS, and no dish may ever be booked for two satellites at once. Ingrid needs the whole day in one schedule she can trust.

## Debrief

The day's schedule now grows as predictions arrive, sorts itself by AOS and never lets a draft move a live pass, and tonight your check catches Furrow-B and Pipit-3 booked on Ystmark's one working dish. Sunniva: "Wind calm. Pipit-3 moved, one dish, one satellite."

# whos-who-in-the-fleet

## Briefing

Week ten. Furrowlight is angry: their dashboard credited one of their passes to Brine-2, a satellite that never made it off its rocket. Apogee keeps its fleet in a list someone searches by hand. Yusra needs answers Furrowlight can trust: which satellites exist, whose they are, and exactly which pieces of a Furrow-A image never came down.

## Debrief

Before nine, Furrowlight has the number of every missing chunk of last night's Furrow-A image, and Brine-2 is gone from a registry that answers by ID. Yusra forwards their reply: "TL;DR: they'll ask Furrow-A for only those chunks tonight. And they said thank you."

# changing-things-in-place

## Briefing

Week eleven. A storm is heading for Ystmark, Apogee's Arctic station. Its two dishes must stow when the wind reaches 72 km/h and go back to following satellites when it drops, hour by hour through the night, and the Board must show the truth about each one. Sunniva will be outside checking the heaters, so the ground software has to keep the record straight.

## Debrief

Through the night, the wind at Ystmark reads 64, 75, 58 and 81 km/h, and both dishes stow and return to work on cue, each stow logged with its reason. Sunniva, at 05:00: "Wind 81, both dishes stowed, Board agrees. Coffee holding."

# every-radio-speaks

## Briefing

Week twelve. Apogee's stations speak three kinds of radio: UHF for small satellites like Pipit-1, S-band for commands and status, X-band for Furrowlight's big images. Every new radio has meant another copy of the pass software, and Paz has fitted a fourth at Salar Quieto. Tomás wants one routine for every radio; Mirela wants alarm rules she can add without touching the pager.

## Debrief

One pass routine now drives every radio, and Paz's new radio joins it without any change: 9,000 MB in 240 seconds. Mirela adds her first two alarm rules, battery low and too hot, without opening the pager. Tomás beams: "Four radios, one routine, like one plug for every socket!"

# when-the-frame-is-wrong

## Briefing

Week thirteen. Thistle, Kessick Polytechnic's student satellite, has started sending trouble: a battery at 6.2 V, temperatures of 900 °C, readings that arrive twice. Apogee's ground software stores them all without a word, so the students plot garbage and nobody can tell a sick satellite from a bad reading. Ingrid wants every bad reading caught, named and logged.

## Debrief

Every bad reading now gets a name, a reason and a line in the log: Thistle's reading 4182 shows up as a repeat and as 900 °C at once, and a repeat alone wakes nobody. Ingrid scrolls the new log: "Copy. Now Thistle can tell us what's wrong."

# reading-the-raw-feed

## Briefing

Week fourteen. Data reaches the Fishbowl raw: Thistle's readings as lines of text, Pipit beacons as a dozen bare numbers, place names from Salar Quieto that the Board cuts short. A single stray letter has turned a battery reading into 0 V. Pipit-2 is about to get its first signed commands, carrying proof that they came from Apogee, and Mirela will review every line.

## Debrief

Thistle's lines, Pipit beacons and Paz's place names now read true, a stray letter raises an error instead of becoming 0 V, and Pipit-2 accepts its first signed command. Mirela isn't done: "Your code accepts the command. Now show me the one it refuses."

# trust-but-verify

## Briefing

Week fifteen. An old command, sent a second time, almost reached Cirrine-1, a customer's weather satellite, because of one wrong character in the uplink. Mirela is changing the rules: whoever writes uplink code now also writes the checks that prove it works. She will plant broken copies of the code, and your checks have to catch every one.

## Debrief

Your checks catch all six of Mirela's planted bugs, from a copy that accepts any signature to one that refuses even a properly signed command, and the uplink now carries checks its own author wrote. Mirela signs off: "Six bugs planted, six caught. The uplink is yours to touch."

# splitting-groundctl

## Briefing

Week sixteen. groundctl, the software behind Apogee's ground stations, has become one big tangle, and three teams keep tripping over each other's changes in it. Nkem draws boxes on the whiteboard: orbit, telemetry, uplink, alert. Each team will own a box and show the others only what they need, the signing keys stay locked inside uplink, and the first release ships Friday.

## Debrief

On Friday, groundctl v1.0.0 ships in clean boxes: orbit, telemetry, uplink, alert, and a small severity box that telemetry and alert share. The signing keys are reachable only through uplink. Nkem: "Now each team can change its own box without breaking anyone else's."

# one-helper-every-type

## Briefing

Week seventeen. groundctl is split into tidy boxes, but each box keeps its own copies of the same small tools: one picks the passes worth keeping, another the readings, one finds the highest elevation, another the warmest battery. Every copy is one more place for a bug to hide. Nkem wants each tool written once, for passes, readings or anything else.

## Debrief

One filter, one maximum, one ring of recent readings and one set now serve all of groundctl, and Ingrid has her answer: Lantern Atoll and Port Ellery both heard Thistle and Pipit-1 tonight. Nkem deletes the old copies one by one: "Smaller. Good."

# all-stations-at-once

## Briefing

Week eighteen. Thistle sinks a little lower every day, and before the week is out it will burn up in the atmosphere. Its student team at Kessick Polytechnic wants every station that can hear it listening on every pass, but Apogee's software checks the stations one after another, too slowly to catch a pass at all of them. These are their last chances to hear it.

## Debrief

On Thistle's last pass, four stations listen at once and your code collects every beacon, 9001 to 9004; the last comes down at Salar Quieto before Thistle burns up. Joss sends one line from Kessick Polytechnic: "Thank you for listening to the end." Next week, the pager is yours.

# locks-and-deadlines

## Briefing

On-call week. Mirela hands you the pager: when something breaks, it wakes you, and by Friday you carry it alone. The station status every service reads goes wrong when they all use it at once, a download keeps writing after LOS, and a stuck uplink holds a dish for an hour. Friday is a full day of passes, with a storm building at Ystmark.

## Debrief

Friday's pass day closes your first week on call: six passes done, one cut short when Ystmark stows in 81 km/h winds, its last two skipped, not stuck. Nkem: "Works, and it says what it means." Later in July, the test van heads up the Quebrada Añil, beyond any signal.
