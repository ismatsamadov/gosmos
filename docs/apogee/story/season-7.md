# Apogee story bible: Season 7

### Season 7: scale

On Monday 5 March 2029 Tomás tells the building that investors have put money into Apogee,
and draws the network it pays for on the Fishbowl's glass: fifteen new stations by June, twenty
in all, and the satellites of four customers' launches, eighty in all by the end of June. The
learner's new badge says **Senior Planning Engineer**, still No. 039, based at Port Ellery
headquarters, since Monday 5 March 2029. Their job is passplan's rewrite. The passplan Apogee
has flown since 2023, which Nkem wrote in Apogee's first summer, plans one day of passes at a
time, in one goroutine: 6.1 seconds for Season 6's seven satellites at five stations, 7 seconds
for eight. For eighty satellites at twenty stations it needs 4 minutes 39 seconds a day, nearly
eight hours for a hundred days, and the new one must plan a hundred days, correctly and fast. You lead the rewrite,
**passplan 2**; Nkem reviews every change and asks, one question at a time, what happens when
two planners touch the same pass. The season ends on take-over day, Thursday 28 June, when
passplan 2 takes over from the old planner on a live day, under the heaviest schedule Apogee
has ever flown.

Story time runs from Monday 5 March to Thursday 28 June 2029: twelve chapters, nine of about
a week and three (chapters 8, 9 and 11) of about two; chapter 12 runs from Monday 11 June to
take-over day. Nkem carries the season; Tomás draws the network and makes one claim too many
about speed; Ingrid owns take-over day; Yusra brings the customers' wishes and Hollis his
capitals; Mirela reviews the telemetry parser; Joss Imbert and Professor Marr of Kessick
Polytechnic return with Thistle-2. No new people are introduced. The new places are fifteen
stations, each an Apogee dish on a host's land; the new satellites are Thistle-2 and three
customers' constellations, named by rule from names already in the bible. The hosts stay
unnamed, as the Cape Wren contractor does.

In plain words, for the opening: Apogee is growing from five dishes to twenty and from a
handful of satellites to eighty, and the program that decides which dish listens to which
satellite, and when, has to be rebuilt for that size without missing a pass.

#### The season, chapter by chapter

Everything happens at Port Ellery headquarters; the new stations speak through their hosts'
forms, and Sunniva and Paz only in passing. Times are UTC.

| Ch | Week of (2029) | Where | Who needs what | What is at stake |
|---|---|---|---|---|
| 1 | Mon 5 Mar | the Fishbowl's glass, then the backend desks | Nkem: before passplan 2, a smaller pipeline built the same way. tlmd's ingest reads the station telemetry feeds in one loop; with twenty feeds it falls behind, and when two stations hear the same satellite its frames arrive out of order, so "the newest frame" goes backward. Wednesday 7 March, from 21:00, is Nkem's drill: twenty feeds, the five real stations and fifteen made from the new stations' positions, replaying the day's lines. Thistle-2 leaves its rocket at 10:20 on Thursday 8 March | The Board and the pager read every satellite's newest frame from this feed. Thistle-2's first beacon reaches Port Ellery at 17:08 on Thursday (`thistle-2/PEL/00005`), and Joss and Professor Marr are in the Fishbowl to hear it |
| 2 | Mon 12 Mar | the backend desks; passplan's new machine in Port Ellery's rack, twelve cores | Nkem: passplan 2's first draft predicts each satellite in its own goroutine, and each day in another: eight thousand at once. On Friday night, 9 March, the test machine ran out of memory. "Twelve planner slots. Never more." The station dish loaders and the Board take what passplan sends at their own pace | A planner that starves its machine stops the Board that shares it; a queue nobody bounds grows until the machine falls over |
| 3 | Mon 19 Mar | Ingrid's desk in the Fishbowl, then the backend desks | Ingrid: on Monday the Board froze for six minutes while a hundred callers each loaded the station table. On Wednesday 21 March two of passplan 2's planning runs finished out of order, and on the test Board beside Ingrid's the older plan replaced the newer one | The test Board showed a plan twenty minutes old; on take-over day that would be the real one. Nkem's whiteboard question: "Why does a program without a race behave?" |
| 4 | Mon 26 Mar | Yusra's desk | Yusra: "TL;DR: customers want to book a whole season ahead." Hollis: "WE PLAN A SEASON NOT A WEEK." bookingd's catalog offers fourteen days today; the third quarter (1 July to 30 September) holds 672,292 passes at stations that can serve them | The old passplan hands every service a slice of every pass: about eighty-four megabytes for each one that asks |
| 5 | Mon 2 Apr | the Fishbowl | The first wave of new stations goes live on Monday 2 April: Tirrin Rock, Varsk Island, Calidra Atoll, Ilha Brumosa, Ostrel Island. From the same day passplan 2 runs in shadow, beside the old planner, and the two schedules are compared contact by contact every hour. The Board wants one station's contacts, bookingd one customer's, Ingrid the next twelve | Two planners that disagree about one pass would fly it twice or not at all; the first wave's schedules must be right from their first pass. Sunniva: "Tirrin up. Two dishes. Ice on both." |
| 6 | Mon 9 Apr | the backend desks | Nkem: the shadow run sorts every candidate pass of a hundred days into one slice per station. A priority queue merges the eighty satellites' passes one at a time | Cirrine's twelve new satellites leave their rocket on Tuesday 17 April; the planner must be ready for them |
| 7 | Mon 16 Apr | Mirela's desk; Kessick Polytechnic by message | On Saturday 14 April Thistle-2's student flight software sent `batt_v=nan` when a sensor glitched; tlmd's parser accepted it, and the pager's battery rule never fired, since nothing compares below a value that isn't a number. At 05:10 a line with a missing field crashed the parser and stopped every station's telemetry for forty minutes. Joss apologizes; Mirela: "Your parser accepts the line. Now show me the lines it refuses" | Cirrine's twelve send their first telemetry through this parser on Tuesday 17 April |
| 8 | Mon 23 Apr, to Fri 4 May | the backend desks | Nkem: passplan 2's timers (a replan every ten minutes, a deadline on every planning run, a wait until each contact's AOS, a pause for news to settle) are tested with real sleeps: the suite takes eleven minutes and fails one run in five | The second wave goes live on Monday 7 May; a test that lies about time would let a broken replan through |
| 9 | Mon 7 May, to Fri 18 May | the Fishbowl's glass, the backend desks | The second wave goes live on Monday 7 May: Ilvane Spit, Cape Solvane, Saltreach Island, Corith Point, Nimrie Point. Tomás writes on the glass: "passplan 2: forty times faster!" Nkem: "Faster at what?" Brinewatch's thirty-six leave their rocket on Wednesday 16 May | A planner that only looks faster would be trusted on take-over day |
| 10 | Mon 21 May | the backend desks | With fifty-six satellites, every replan of a hundred days leaves the machine's memory full of garbage, and the test Board's updates stutter while the collector cleans up. Nkem: "Halve the planner's garbage before the third wave" | The third wave adds five stations and the last launch adds twenty-four satellites |
| 11 | Mon 28 May, to Fri 8 Jun | the backend desks | The third wave's hosts send their forms: Vorrel Island's says 0 dishes, Ilha Calmaria's lists a band "K". On Friday 1 June a test deploy of passplan 2 read `slots: 0` and planned nothing. Nkem: one validator for every settings struct, its rules in tags | The third wave goes live on Monday 4 June: Hasku Point, Ilha Calmaria, Quindle Island, Punta Ventisca, Vorrel Island. A wrong form takes a station's passes away |
| 12 | Mon 11 Jun, to Thu 28 Jun | the Fishbowl | Ingrid runs take-over day like a pass: a full rehearsal in shadow on Monday 25 June, Furrowlight's launch on Wednesday 27 June, and the take-over at 00:00 UTC on Thursday 28 June | The heaviest schedule Apogee has ever flown: 2,011 contacts in one day at twenty stations, with twenty-four new Furrows in their first day of life |

The opening (Monday 5 March) shows Tomás's drawing on the glass (five dots becoming twenty,
eight satellites becoming eighty), hands over the badge, introduces Nkem (reviewer, author of
the old passplan), Tomás, Ingrid (take-over day is hers) and Yusra, and sends the learner to
Nkem's first question about tlmd's ingest. Words of the first week: passplan, planner slot,
telemetry feed, the drill.

The finale (chapter 12, Thursday 28 June): passplan 2 takes over at midnight without touching
one of the old planner's twelve contacts in progress; flies 2,011 contacts through Cape Wren's
wind; and every one of Furrowlight's twenty-four new satellites, first heard the evening
before, is heard all day. Ingrid writes one line in the log at 00:00 on Friday: "passplan 2:
GO for every day." Hollis, through Yusra: "TWENTY-FOUR FURROWS. EVERY ONE HEARD." The debrief
points at July: the whole ground segment, rehearsed live for a week.

#### The network grows

Apogee's network at the end of June 2029: the five stations of the bible's table and fifteen
new ones, each an Apogee dish on a **host's** land (a lighthouse authority, a weather
observatory, a fisheries harbor, a research camp): the host gives power, a line to the internet
and a pair of hands; nobody from Apogee lives there. Each new station arrives as a dish and a
rack in a shipping container and goes live in one of three **waves**, at 00:00 UTC on a
Monday. Coordinates are invented and land in open ocean, as the bible's do. A new station has
one 3.7 m dish for S and X unless the table says otherwise; X-band is left out where tropical
rain would fade it, as at Lantern Atoll. Only Port Ellery and Lantern Atoll send on UHF.

| Code | Name | Lat, lon | Mask | Dishes | Bands | Live | Character |
|---|---|---|---|---|---|---|---|
| YST | Ystmark | 77.8° N, 2.0° E | 5° | 2 | S X | flying | the bible's |
| PEL | Port Ellery | 46.9° N, 38.1° W | 5° | 1 | UHF S | flying | the bible's |
| SQT | Salar Quieto | 23.4° S, 104.6° W | 3° | 2 | S X | flying | the bible's |
| LAN | Lantern Atoll | 1.8° S, 82.5° E | 8° | 1 | UHF S | flying | the bible's |
| CWR | Cape Wren | 53.9° S, 95.0° W | 5° | 1 | S X | flying | the bible's |
| TIR | Tirrin Rock | 69.5° N, 12.0° W | 5° | 2 | S X | Mon 2 Apr | A bare rock in the Norwegian Sea with a lighthouse and two radomes; spray ices them in winter. Ystmark's partner for polar passes |
| VSK | Varsk Island | 51.0° N, 160.0° W | 7° | 1 | S X | Mon 2 Apr | A North Pacific fishing harbor; fog most mornings, and a hill to the north |
| CAL | Calidra Atoll | 0.5° N, 150.0° W | 8° | 1 | S | Mon 2 Apr | An equatorial atoll; palms on the horizon and rain fade, as at Lantern Atoll |
| BRM | Ilha Brumosa | 30.0° S, 15.0° W | 5° | 1 | S X | Mon 2 Apr | A misty South Atlantic island ("brumosa" means misty); its host is a weather observatory |
| OST | Ostrel Island | 42.5° S, 70.0° E | 5° | 1 | S X | Mon 2 Apr | In the westerly gales of the southern Indian Ocean; the dish stows at 72 km/h, like Cape Wren's |
| ILV | Ilvane Spit | 72.0° N, 165.0° W | 3° | 2 | S X | Mon 7 May | A gravel spit in the Chukchi Sea, flat to every horizon; sea ice until June |
| CSL | Cape Solvane | 35.0° N, 45.0° W | 5° | 1 | S X | Mon 7 May | A volcanic cape in the mid-Atlantic; its host keeps vines on the slope below the dish |
| SRH | Saltreach Island | 17.0° N, 135.0° E | 6° | 1 | S | Mon 7 May | Typhoon country in the western Pacific |
| CRP | Corith Point | 10.0° N, 65.0° E | 6° | 1 | S | Mon 7 May | A point of an Arabian Sea island; the monsoon arrives in June |
| NIM | Nimrie Point | 45.0° S, 150.0° W | 5° | 1 | S X | Mon 7 May | A lonely point in the South Pacific; the supply ship calls twice a year |
| HSK | Hasku Point | 32.0° N, 170.0° E | 4° | 1 | S X | Mon 4 Jun | An old airstrip on a North Pacific island |
| ICM | Ilha Calmaria | 8.0° N, 30.0° W | 7° | 1 | S | Mon 4 Jun | In the calm belt of the tropical Atlantic ("calmaria" means calm), with sudden squalls |
| QND | Quindle Island | 25.0° S, 95.0° E | 5° | 1 | S X | Mon 4 Jun | An island of the eastern Indian Ocean, in the trade winds |
| PVT | Punta Ventisca | 60.5° S, 40.0° E | 5° | 1 | S X | Mon 4 Jun | A cape in the Southern Ocean ("ventisca" means blizzard); the radome heater never stops |
| VRL | Vorrel Island | 64.0° S, 140.0° W | 4° | 2 | S X | Mon 4 Jun | An icebound island in the far south: the south's Ystmark |

The table order (the five, then the waves in this order) is the **table order** every list
of stations uses. In readmes a new station is its name first and its code after: "Tirrin Rock
(TIR)".

**The fleet** grows from eight to eighty: the seven the bible's fleet still has (Thistle
re-entered in July 2027), Thistle-2 on Thursday 8 March, then three launches. Satellites of a
launch leave the rocket ("deploy") one after another and fly in a line around the same orbit;
a new satellite's first days are its LEOP. IDs and names follow the bible's rule: `cirrine-2`
is Cirrine-2, `brine-17` is Brine-17, `furrow-q` is Furrow-Q.

| Satellites | Operator | Orbit | Bands | Deployed (Epoch below) | Story |
|---|---|---|---|---|---|
| `pipit-1`, `pipit-2`, `pipit-3`, `furrow-a`, `furrow-b`, `brine-3`, `cirrine-1` | as in the bible | as in the bible | Pipit-1 UHF; Pipit-2 S; Pipit-3 S X; the Furrows X S; Brine-3 S; Cirrine-1 S | flying | the bible's |
| `thistle-2` | Kessick Polytechnic | 505 km, 97.4° | UHF | Thu 8 Mar, 10:20:00 | Kessick's second student cubesat; Joss now leads its flight software. First beacon at Port Ellery, 17:08:39 (`thistle-2/PEL/00005`) |
| `cirrine-2` to `cirrine-13` (12) | Cirrine Atmospherics | 600 km, 97.8° | S | Tue 17 Apr, 09:40:00 | the bigger contract of Season 5; first contact `cirrine-3/VSK/00001` at 09:41:01 |
| `brine-4` to `brine-39` (36) | Brinewatch | 480 km, 53.0° | S | Wed 16 May, 13:30:00 | ship tracking; like Brine-3, never over Ystmark |
| `furrow-c` to `furrow-z` (24) | Furrowlight | 510 km, 97.4° | X S | Wed 27 Jun, 20:40:00 | lift-off 18:05, deployed one about every four minutes from 19:10; first contact `furrow-g/YST/00001` at 20:40:16, the last of the 24 first heard at 22:28:50 (`furrow-c/TIR/00002`) |

The fleet order (the seven, Thistle-2, then the three launches in this order, each by
number or letter) is the **fleet order**.

#### The pass model

passplan predicts passes from each satellite's **elements**, the few numbers that place it on
its orbit (`Elements`, below: Season 6's type, to which this season adds `SatID` and `Orbit`). The model is written here once and copied unchanged into every
lesson that predicts a pass, as `network.go` (the twenty stations and eighty satellites) and
`model.go` (`predict`); `go run ./cmd/passes` runs the same code. It is the model Season 3's
reviewers settled, ported to Go: a circular orbit around a round Earth, the orbit's plane turning with the mean
Sun for a sun-synchronous satellite (one tilted past 90°) and by Earth's bulge for any other,
Earth turning beneath it, and a station's mask. It watches the satellite's elevation every 20
seconds (less often while the satellite is far below the horizon), finds each instant it
crosses the mask to a hundred-millionth of a second and rounds it to the second, and takes the
highest of the elevations a second apart from AOS as the maximum. A pass shorter than two
minutes doesn't count. Near the top of its track a satellite crosses a high station's latitude
twice on one orbit; the station sees it rise once and set once: one pass. Pass IDs follow the
bible's rule: the orbit is the one on which the satellite crosses the station's latitude during
the pass (the top or bottom of its track where it never reaches that latitude), counted from
`Rev0` at `Epoch`; a new satellite's orbits count from 1 at its deployment. (Season 3's first
version numbered a pass by the orbit at its AOS, which gave two passes one ID at stations just
north of the equator; the rest agree.) Elements of the eighty:

| Satellites | Epoch | Rev0 | NodeDeg | PhaseDeg |
|---|---|---|---|---|
| `pipit-1` | 2029-01-01T00:00:00Z | 14094 | 201.4926025817 | 176.9565679403 |
| `pipit-2` | 2029-01-01T00:00:00Z | 19752 | 234.8439482122 | 339.4065995135 |
| `pipit-3` | 2029-01-01T00:00:00Z | 9689 | 150 | 40 |
| `furrow-a` | 2029-01-01T00:00:00Z | 21818 | 294.5055051072 | 52.6136621172 |
| `furrow-b` | 2029-01-01T00:00:00Z | 21815 | 297.5851231763 | 231.0891932221 |
| `brine-3` | 2029-01-01T00:00:00Z | 15209 | 179.2270109042 | 188.3141532885 |
| `cirrine-1` | 2029-01-01T00:00:00Z | 17206 | 282.5236543824 | 268.6744811484 |
| `thistle-2` | 2029-03-08T10:20:00Z | 1 | −120 | 0 |
| `cirrine-n`, n = 2 to 13 | 2029-04-17T09:40:00Z | 1 | −150 | 30 × (n − 2) |
| `brine-n`, n = 4 to 39 | 2029-05-16T13:30:00Z | 1 | 120 | 10 × (n − 4) |
| `furrow-c` … `furrow-z`, k = 0 to 23 | 2029-06-27T20:40:00Z | 1 | 10 | 15 × k |

Pipit-1, Furrow-A, Furrow-B, Brine-3 and Cirrine-1 are Season 3's settled fits, and Pipit-2 is
Season 4's, fitted to the Build gctl night's pass over Ystmark (`pipit-2/YST/13287`, 04:31:00 to
04:40:33 on 30 October 2027, 47.2°) with its orbits counted from 19752 at 1 January 2029; all six
are carried to 1 January 2029 (Furrow-A and Furrow-B with `RevsPerDay` 867/57: 867 orbits in 57
days, a little faster than their altitude gives); every other satellite has `RevsPerDay` 0. So the
orbit numbers continue the earlier seasons', and Furrow-A still rises over Salar Quieto at
14:02:10 on 28 June (`furrow-a/SQT/24534`, to 14:12:04, 47.9°), as on launch day in 2028, though
that day both of the station's dishes are taken by then. Pipit-1 crosses Port Ellery at 03:39:12
and 05:13:13 on 28 June (`pipit-1/PEL/16795`, `/16796`), each time a few minutes after Thistle-2,
which gets the station's one dish.

```go
// Site is a station as passplan sees it.
type Site struct {
	Station
	Dishes int       // the dishes passplan can give passes to
	Live   time.Time // the first moment it flies a pass; zero means always
}

// Elements place a satellite on its orbit for the pass model.
type Elements struct {
	SatID    string
	Orbit    Orbit     // altitude and inclination
	Epoch    time.Time // the moment the numbers below describe; no pass before it
	Rev0     int       // the orbit number at Epoch
	NodeDeg  float64   // east longitude of the northbound equator crossing, at Epoch
	PhaseDeg float64   // how far past that crossing the satellite is at Epoch, 0 to 360
	// RevsPerDay is the satellite's orbits a day, where tracking measured them; zero means
	// the period its altitude gives.
	RevsPerDay float64
}

// Network is everything passplan plans for.
type Network struct {
	Sites      []Site      // the twenty, in table order
	Satellites []Satellite // the eighty, in fleet order
	Elements   []Elements  // the same eighty, in the same order
}
```

`network.go`, exactly (it needs `fmt` and `time`):

```go
// newNetwork returns the network as it stands at the end of June 2029.
func newNetwork() *Network {
	sx, s, us := []Band{SBand, XBand}, []Band{SBand}, []Band{UHF, SBand}
	w1 := time.Date(2029, 4, 2, 0, 0, 0, 0, time.UTC)
	w2 := time.Date(2029, 5, 7, 0, 0, 0, 0, time.UTC)
	w3 := time.Date(2029, 6, 4, 0, 0, 0, 0, time.UTC)
	site := func(code, name string, lat, lon, mask float64, bands []Band, dishes int, live time.Time) Site {
		st := Station{Code: code, Name: name, LatDeg: lat, LonDeg: lon, MaskDeg: mask, Bands: bands}
		return Site{Station: st, Dishes: dishes, Live: live}
	}
	var always time.Time
	n := &Network{Sites: []Site{
		site("YST", "Ystmark", 77.8, 2.0, 5, sx, 2, always),
		site("PEL", "Port Ellery", 46.9, -38.1, 5, us, 1, always),
		site("SQT", "Salar Quieto", -23.4, -104.6, 3, sx, 2, always),
		site("LAN", "Lantern Atoll", -1.8, 82.5, 8, us, 1, always),
		site("CWR", "Cape Wren", -53.9, -95.0, 5, sx, 1, always),
		site("TIR", "Tirrin Rock", 69.5, -12.0, 5, sx, 2, w1),
		site("VSK", "Varsk Island", 51.0, -160.0, 7, sx, 1, w1),
		site("CAL", "Calidra Atoll", 0.5, -150.0, 8, s, 1, w1),
		site("BRM", "Ilha Brumosa", -30.0, -15.0, 5, sx, 1, w1),
		site("OST", "Ostrel Island", -42.5, 70.0, 5, sx, 1, w1),
		site("ILV", "Ilvane Spit", 72.0, -165.0, 3, sx, 2, w2),
		site("CSL", "Cape Solvane", 35.0, -45.0, 5, sx, 1, w2),
		site("SRH", "Saltreach Island", 17.0, 135.0, 6, s, 1, w2),
		site("CRP", "Corith Point", 10.0, 65.0, 6, s, 1, w2),
		site("NIM", "Nimrie Point", -45.0, -150.0, 5, sx, 1, w2),
		site("HSK", "Hasku Point", 32.0, 170.0, 4, sx, 1, w3),
		site("ICM", "Ilha Calmaria", 8.0, -30.0, 7, s, 1, w3),
		site("QND", "Quindle Island", -25.0, 95.0, 5, sx, 1, w3),
		site("PVT", "Punta Ventisca", -60.5, 40.0, 5, sx, 1, w3),
		site("VRL", "Vorrel Island", -64.0, -140.0, 4, sx, 2, w3),
	}}
	add := func(id, name, op string, o Orbit, bands []Band, epoch time.Time, rev0 int, node, phase, revs float64) {
		n.Satellites = append(n.Satellites, Satellite{ID: id, Name: name, Operator: op, Orbit: o, Bands: bands, Mode: "NOMINAL"})
		n.Elements = append(n.Elements, Elements{SatID: id, Orbit: o, Epoch: epoch, Rev0: rev0, NodeDeg: node, PhaseDeg: phase, RevsPerDay: revs})
	}
	jan := time.Date(2029, 1, 1, 0, 0, 0, 0, time.UTC)
	xs := []Band{XBand, SBand}
	// Pipit-1, Pipit-2, Furrow-A, Furrow-B, Brine-3 and Cirrine-1 fitted to passes the
	// stations recorded (Furrow-A and Furrow-B make 867 orbits in 57 days); Pipit-3 from
	// Apogee's own tracking.
	add("pipit-1", "Pipit-1", "Apogee", Orbit{525, 97.5}, []Band{UHF}, jan, 14094, 201.4926025817, 176.9565679403, 0)
	add("pipit-2", "Pipit-2", "Apogee", Orbit{550, 97.6}, s, jan, 19752, 234.8439482122, 339.4065995135, 0)
	add("pipit-3", "Pipit-3", "Apogee", Orbit{510, 97.4}, sx, jan, 9689, 150, 40, 0)
	add("furrow-a", "Furrow-A", "Furrowlight", Orbit{510, 97.4}, xs, jan, 21818, 294.5055051072, 52.6136621172, 867.0/57)
	add("furrow-b", "Furrow-B", "Furrowlight", Orbit{510, 97.4}, xs, jan, 21815, 297.5851231763, 231.0891932221, 867.0/57)
	add("brine-3", "Brine-3", "Brinewatch", Orbit{480, 53.0}, s, jan, 15209, 179.2270109042, 188.3141532885, 0)
	add("cirrine-1", "Cirrine-1", "Cirrine Atmospherics", Orbit{600, 97.8}, s, jan, 17206, 282.5236543824, 268.6744811484, 0)
	add("thistle-2", "Thistle-2", "Kessick Polytechnic", Orbit{505, 97.4}, []Band{UHF},
		time.Date(2029, 3, 8, 10, 20, 0, 0, time.UTC), 1, -120, 0, 0)
	for k := 2; k <= 13; k++ {
		add(fmt.Sprintf("cirrine-%d", k), fmt.Sprintf("Cirrine-%d", k), "Cirrine Atmospherics", Orbit{600, 97.8}, s,
			time.Date(2029, 4, 17, 9, 40, 0, 0, time.UTC), 1, -150, float64(30*(k-2)), 0)
	}
	for k := 4; k <= 39; k++ {
		add(fmt.Sprintf("brine-%d", k), fmt.Sprintf("Brine-%d", k), "Brinewatch", Orbit{480, 53.0}, s,
			time.Date(2029, 5, 16, 13, 30, 0, 0, time.UTC), 1, 120, float64(10*(k-4)), 0)
	}
	for k := 0; k < 24; k++ {
		add(fmt.Sprintf("furrow-%c", 'c'+k), fmt.Sprintf("Furrow-%c", 'C'+k), "Furrowlight", Orbit{510, 97.4}, xs,
			time.Date(2029, 6, 27, 20, 40, 0, 0, time.UTC), 1, 10, float64(15*k), 0)
	}
	return n
}
```

`model.go`, exactly (it needs `fmt`, `math`, `strconv` and `time`):

```go
// The pass model Season 3's reviewers settled: a circular orbit around a round Earth, the
// orbit's plane turning with the mean Sun for a sun-synchronous satellite (one tilted past
// 90°) and by Earth's bulge for any other, and Earth turning beneath it. A pass is the time the
// satellite's elevation over a station stays at or above the station's mask: its AOS and LOS
// are the instants the elevation crosses the mask, rounded to the second, and its maximum
// elevation is the highest of its elevations a second apart from AOS, to a tenth of a degree.

const (
	earthRadiusKm = 6371.0
	muKm3PerS2    = 398600.4418
	j2            = 1.08263e-3                       // how strongly Earth's bulge turns an orbit
	earthRate     = 2 * math.Pi / 86164.0905         // Earth's turn, radians a second
	sunRate       = 2 * math.Pi / (365.2422 * 86400) // the mean Sun's, radians a second
)

func rad(d float64) float64 { return d * math.Pi / 180 }
func deg(r float64) float64 { return r * 180 / math.Pi }

// periodS returns an orbit's period in seconds, as its altitude gives it.
func periodS(o Orbit) float64 {
	r := earthRadiusKm + o.AltKm
	return 2 * math.Pi * math.Sqrt(r*r*r/muKm3PerS2)
}

// skyTrack is one satellite over one station: everything predict needs, worked out once.
// Its times are seconds after the satellite's Epoch.
type skyTrack struct {
	a, cosI, sinI   float64    // the orbit's radius in km, and its tilt
	n, u0           float64    // radians a second along the orbit, and how far along at Epoch
	node0, nodeRate float64    // the node's longitude at Epoch, and how fast it moves over Earth
	st, up          [3]float64 // the station, km from Earth's center, and straight up there
	mask, reach     float64    // the mask, and the widest angle at Earth's center it sees out to
	speed           float64    // no satellite and station draw apart or together faster, radians a second
}

func newSkyTrack(e Elements, s Site) *skyTrack {
	k := &skyTrack{a: earthRadiusKm + e.Orbit.AltKm, mask: s.MaskDeg}
	i := rad(e.Orbit.IncDeg)
	k.cosI, k.sinI = math.Cos(i), math.Sin(i)
	k.n = math.Sqrt(muKm3PerS2 / (k.a * k.a * k.a))
	if e.RevsPerDay > 0 {
		k.n = e.RevsPerDay * 2 * math.Pi / 86400
	}
	drift := -1.5 * k.n * j2 * (earthRadiusKm / k.a) * (earthRadiusKm / k.a) * k.cosI
	if e.Orbit.IncDeg > 90 {
		drift = sunRate // sun-synchronous: the plane keeps its angle to the Sun
	}
	k.node0, k.nodeRate, k.u0 = rad(e.NodeDeg), drift-earthRate, rad(e.PhaseDeg)
	la, lo := rad(s.LatDeg), rad(s.LonDeg)
	k.up = [3]float64{math.Cos(la) * math.Cos(lo), math.Cos(la) * math.Sin(lo), math.Sin(la)}
	for j := range 3 {
		k.st[j] = earthRadiusKm * k.up[j]
	}
	m := rad(s.MaskDeg)
	k.reach = math.Acos(earthRadiusKm*math.Cos(m)/k.a) - m
	k.speed = k.n + earthRate + math.Abs(drift)
	return k
}

// pos is the satellite's place over the turning Earth, km from its center.
func (k *skyTrack) pos(t float64) (x, y, z float64) {
	su, cu := math.Sincos(k.u0 + k.n*t)
	sn, cn := math.Sincos(k.node0 + k.nodeRate*t)
	return k.a * (cn*cu - sn*su*k.cosI), k.a * (sn*cu + cn*su*k.cosI), k.a * su * k.sinI
}

// elev is the satellite's elevation over the station at t, in degrees.
func (k *skyTrack) elev(t float64) float64 {
	x, y, z := k.pos(t)
	dx, dy, dz := x-k.st[0], y-k.st[1], z-k.st[2]
	d := math.Sqrt(dx*dx + dy*dy + dz*dz)
	return deg(math.Asin((dx*k.up[0] + dy*k.up[1] + dz*k.up[2]) / d))
}

// apart is the angle at Earth's center between the satellite and the station.
func (k *skyTrack) apart(t float64) float64 {
	x, y, z := k.pos(t)
	return math.Acos(max(-1, min(1, (x*k.up[0]+y*k.up[1]+z*k.up[2])/k.a)))
}

// above is the elevation over the mask: negative while the satellite is below it.
func (k *skyTrack) above(t float64) float64 { return k.elev(t) - k.mask }

// orbit is the number of the orbit on which the satellite, near t, crosses the station's
// latitude: Rev0 at Epoch, one more at each northbound equator crossing. A satellite whose
// track never reaches the station's latitude counts from the top or bottom of its track.
func (k *skyTrack) orbit(rev0 int, t float64) int {
	u := k.u0 + k.n*t
	var marks []float64 // where along an orbit the satellite is at the station's latitude
	if s := k.up[2] / k.sinI; s >= 1 {
		marks = []float64{math.Pi / 2}
	} else if s <= -1 {
		marks = []float64{3 * math.Pi / 2}
	} else {
		marks = []float64{math.Asin(s), math.Pi - math.Asin(s)}
	}
	best := math.Inf(1)
	var at float64
	for _, m := range marks {
		c := m + 2*math.Pi*math.Round((u-m)/(2*math.Pi)) // the nearest such point to u
		if d := math.Abs(c - u); d < best {
			best, at = d, c
		}
	}
	return rev0 + int(math.Floor(at/(2*math.Pi)))
}

// predict returns e's passes over s whose AOS is in [from, to), by AOS.
// No pass is before Epoch or before the station goes live, and a pass shorter than two
// minutes doesn't count. Near the top of its track a satellite crosses a high station's
// latitude twice on one orbit; the station sees it rise once and set once: one pass.
func predict(e Elements, s Site, from, to time.Time) []PassView {
	from = later(from, later(e.Epoch, s.Live))
	if !from.Before(to) {
		return nil
	}
	k := newSkyTrack(e, s)
	start, end := from.Sub(e.Epoch).Seconds(), to.Sub(e.Epoch).Seconds()
	at := func(t float64) time.Time { // to the nearest second
		return e.Epoch.Add(time.Duration(math.RoundToEven(t)) * time.Second)
	}
	tenths := func(v float64) float64 { // to a tenth, as the settled model prints it
		f, _ := strconv.ParseFloat(strconv.FormatFloat(v, 'f', 1, 64), 64)
		return f
	}
	var out []PassView
	const step = 20.0
	t := start - 1800 // half an hour early: a pass in progress at from is seen rising
	prev := k.above(t)
	for t < end+60 {
		if prev < 0 {
			// Far below the mask the satellite can't come up for a while: skip what it can't.
			if gap := k.apart(t) - k.reach - 0.01; gap > 0 {
				if jump := math.Floor(gap/k.speed/step) - 1; jump >= 1 {
					t += jump * step
				}
			}
		}
		t += step
		cur := k.above(t)
		if prev < 0 && cur >= 0 {
			aos := k.crossing(t-step, t)
			peak, top := k.highest(aos, int(2*k.reach/(k.n-earthRate))+60)
			end := aos + float64(2*peak) + 120
			for k.above(end) >= 0 {
				end += 60
			}
			los := k.crossing(aos+float64(peak), end)
			p := Pass{SatID: e.SatID, Station: s.Code, AOS: at(aos), LOS: at(los), MaxElevDeg: tenths(top)}
			id := fmt.Sprintf("%s/%s/%05d", e.SatID, s.Code, k.orbit(e.Rev0, aos+float64(peak)))
			out = append(out, PassView{ID: id, Pass: p})
			t = los
			cur = k.above(t)
		}
		prev = cur
	}
	kept := out[:0]
	for _, p := range out {
		if !p.AOS.Before(from) && p.AOS.Before(to) && p.LOS.Sub(p.AOS) >= 2*time.Minute {
			kept = append(kept, p)
		}
	}
	return kept
}

// crossing returns the instant between lo and hi where the elevation crosses the mask, below
// it at one end and not at the other: the Illinois method, to a hundred-millionth of a second.
func (k *skyTrack) crossing(lo, hi float64) float64 {
	f := k.above
	flo, fhi := f(lo), f(hi)
	side := 0
	for range 100 {
		m := (lo*fhi - hi*flo) / (fhi - flo)
		if !(m > lo && m < hi) {
			m = (lo + hi) / 2
		}
		fm := f(m)
		if (fm >= 0) == (fhi >= 0) {
			hi, fhi = m, fm
			if side == 1 {
				flo /= 2
			}
			side = 1
		} else {
			lo, flo = m, fm
			if side == -1 {
				fhi /= 2
			}
			side = -1
		}
		if hi-lo < 1e-8 || fm == 0 {
			return m
		}
	}
	return (lo + hi) / 2
}

// highest returns the whole second i in [0, n] after aos when the elevation is largest, and
// that elevation. The elevation rises to one peak, then falls, so a binary search on which
// way it is heading finds the peak.
func (k *skyTrack) highest(aos float64, n int) (int, float64) {
	f := func(i int) float64 { return k.elev(aos + float64(i)) }
	lo, hi := 0, n
	for lo < hi {
		m := (lo + hi) / 2
		if f(m+1) > f(m) {
			lo = m + 1
		} else {
			hi = m
		}
	}
	return lo, f(lo)
}

// later returns the later of two times.
func later(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// earlier returns the earlier of two times.
func earlier(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
```

The model is hidden where a lesson only needs passes and read-only where a lesson calls it by
name. Its type is `Predictor` when code takes it as a value:
`type Predictor func(e Elements, s Site, from, to time.Time) []PassView`.

**Reference numbers.** A copy of the shared files is right only if it gives these, all
computed with Go 1.27.3 (the calibration writer checks them before chapter 4 is written):

| What | Number |
|---|---|
| Passes on Thursday 28 June (every satellite over every station) | 7,485 (Furrow-A: 105) |
| ... at stations that can serve the satellite (`canServe`, below) | 7,289 |
| ... of those, candidates: max elevation at least 10° | 6,102 |
| Servable passes, 1 to 7 July | 51,174 (the first: `brine-27/VSK/00697` 00:00:04, `furrow-s/SRH/00049` 00:00:15, `brine-16/CRP/00697` 00:00:18) |
| Third quarter, 1 July to 30 September: all passes / servable | 690,261 / 672,292 |
| Furrowlight's 26 satellites, third quarter, servable | 249,045 (1 to 14 July: 37,917, the first `furrow-s/SRH/00049` 00:00:15) |
| Contacts planned for Monday 5 March (5 stations, 7 satellites) | 128, from 142 candidates |
| Contacts for Monday 2 April (10 stations, 8 satellites) | 260, from 286 candidates |
| Tirrin Rock alone on Tuesday 17 April | 120 contacts, from 148 candidates |
| Contacts for Wednesday 27 June (the old planner's last day) | 1,653 |
| Contacts for 28 June planned from scratch (no hand-over, no stow) | 2,020, from 6,102 candidates, 4,082 dropped; Cape Wren 76 of its 360 |
| From the hand-over, 28 June to 5 July (lesson 12.8's week) | 14,282 contacts, from 42,853 candidates, 28,571 dropped |

#### Planning: the rule

A **candidate** is a predicted pass at a station that can serve the satellite (they share a
band), whose max elevation is at least `Config.MinElevDeg` (10°), and that no stow overlaps.
Each station is planned on its own, in **AOS order, then pass ID**: a pass takes the
lowest-numbered dish that is free at its AOS, a dish being free once its last contact's LOS
plus the **turnaround** (2 minutes, to turn and settle) has passed; a pass no dish can take
is **dropped**. A pass given a dish is a **contact**. Stations never share a dish, so
stations can be planned at once, in any order, and the plan is the same. `plan.go`, exactly
(it needs `cmp`, `iter`, `slices`, `strings` and `time`):

```go
// byAOS orders passes by AOS, then by pass ID.
func byAOS(a, b PassView) int {
	return cmp.Or(a.AOS.Compare(b.AOS), strings.Compare(a.ID, b.ID))
}

// canServe reports whether the station offers one of the satellite's bands.
func canServe(s Site, sat Satellite) bool {
	return slices.ContainsFunc(sat.Bands, func(b Band) bool { return slices.Contains(s.Bands, b) })
}

// stowedDuring reports whether a stow at the pass's station overlaps the pass.
func stowedDuring(p Pass, stows []Stow) bool {
	for _, st := range stows {
		if st.Station == p.Station && (st.To.IsZero() || p.AOS.Before(st.To)) && st.From.Before(p.LOS) {
			return true
		}
	}
	return false
}

// assign gives one station's candidate passes dishes. The passes come in AOS
// order, then ID; free[d] is when dish d+1 is next free (zero: free now).
// It returns the contacts in AOS order and how many passes found no dish.
func assign(s Site, passes iter.Seq[PassView], free []time.Time, turnaround time.Duration) ([]Contact, int) {
	next := make([]time.Time, s.Dishes)
	copy(next, free)
	var out []Contact
	dropped := 0
	for p := range passes {
		d := slices.IndexFunc(next, func(t time.Time) bool { return !p.AOS.Before(t) })
		if d < 0 {
			dropped++
			continue
		}
		next[d] = p.LOS.Add(turnaround)
		out = append(out, Contact{PassView: p, Dish: d + 1})
	}
	return out, dropped
}
```

`byAOS` is the comparator everywhere: schedules, merges, catalogs. Chapter 6 teaches it;
`plan.go` is given from chapter 6 on. Before that, a lesson that needs it gives `byAOS` read-only or hidden as it stands until 6.1, comparing the AOS alone (`// byAOS orders passes by AOS.` above `return a.AOS.Compare(b.AOS)`), and no check depends on the order of passes that rise in the same second; 6.1 adds the pass ID after Port Ellery's tie of Sunday 8 April, the first two candidates at one station to rise in the same second since 2 April.

**The window.** passplan 2 plans a **window** of a hundred days: from 00:00 UTC of the current
day to 00:00 a hundred days later (28 June to 6 October on take-over day). A schedule holds
every contact whose AOS is in its window.

**The hand-over** (two planners, one dish). At a **cut-over** instant, contacts of the current
schedule whose AOS is before it stay as they are, flown by whoever planned them; planning
starts at the cut-over with each dish free at the latest LOS plus turnaround of its contacts
with AOS before the cut-over. The same rule makes every replan of passplan 2 a **continuation**:
replanning at any instant with nothing new gives exactly the contacts the schedule already
had from that instant on. At 00:00 on 28 June the old planner's schedule has twelve contacts in
progress and two that ended less than two minutes before midnight; planned without the
hand-over, ten of passplan 2's first contacts would take a dish still busy with one of the
old planner's (the first: `furrow-j/SRH/00003` at 00:00:40, on the dish `furrow-k/SRH/00003`
holds until 00:05:00). Run 1 is itself a replan, of the old planner's schedule from the cut-over, so passplan 2's schedules on take-over day also carry the old planner's contacts, which rose before their window, and every later replan's hand-over still waits for their dishes. Planned from the hand-over alone, version 1 would forget them, and the 00:10 tick would give Ilvane Spit's dish 2 to `cirrine-4/ILV/01069` (AOS 00:10:23) while `furrow-f/ILV/00003`'s turnaround lasts until 00:10:26: a fourth version, and one dish touched by both planners.

**Publishing.** Every planning run takes the next **run number** when it starts (counting from
1 when passplan 2 starts) and its schedule carries it as `Version`. A run whose contacts equal
the current schedule's publishes nothing. The publisher keeps the current schedule in an
`atomic.Pointer` and accepts a schedule only from a newer run (compare-and-swap); a schedule is
never changed after it is published. passplan 2 never runs two plans at once on purpose: a new
trigger cancels the run in progress and starts a fresh one, and the publisher still refuses
anything older than what it holds.

```go
// Window is the span a plan covers: From included, To excluded.
type Window struct {
	From time.Time
	To   time.Time
}

// Stow is a span when a station's dishes are parked. A zero To lasts until further notice.
type Stow struct {
	Station string    // "CWR"
	From    time.Time // 2029-06-28T09:47:00Z
	To      time.Time // 2029-06-28T13:23:00Z
	Reason  string    // "wind 74 km/h"
}

// Contact is a pass given a dish: what a station will fly.
type Contact struct {
	PassView     // the pass and its ID
	Dish     int `json:"dish"` // 1 or 2
}

// Schedule is one published plan. Nobody changes it after it is published.
type Schedule struct {
	Version  uint64    // the run that made it
	Window             // what it plans
	Made     time.Time // when its run started
	Contacts []Contact // every contact with AOS in the window, by AOS then ID
	Dropped  int       // candidates no dish could take
}
```

#### passplan 2's settings, bounds and stages

```go
// Config is passplan 2's settings.
type Config struct {
	Slots      int           `check:"min=1,max=64" default:"12"`
	Days       int           `check:"min=1,max=120" default:"100"`
	MinElevDeg float64       `check:"min=0,max=45" default:"10"`
	Turnaround time.Duration `check:"min=0s,max=10m" default:"2m"`
	Replan     time.Duration `check:"min=1m,max=1h" default:"10m"`
	RunLimit   time.Duration `check:"min=1s,max=5m" default:"30s"`
	Queue      int           `check:"min=1,max=4096" default:"256"`
	Quiet      time.Duration `check:"min=0s,max=1m" default:"5s"`
}
```

The bounds, everywhere:

| Bound | Value | Where |
|---|---|---|
| Planner slots | 12: jobs (a satellite to predict, a station to plan) running at once, one per core of passplan's machine | chapters 2, 12 |
| Window | 100 days | chapters 4, 12 |
| Turnaround | 2 minutes | the rule |
| Minimum max elevation for a contact | 10° | the rule |
| Replan | every 10 minutes, and at once on news | chapters 8, 12 |
| Run limit | 30 seconds per planning run; past it the run is canceled | chapter 8 |
| Quiet | 5 seconds: news arriving closer together than this makes one replan | chapter 8 |
| Dish loader queue | 256 contacts per station; when full, passplan waits | chapter 2 |
| Dish loader uploads | 2 at once per station (the hosts' lines are thin) | chapter 2 |
| The Board's queue | 4 schedules; when full, the oldest waiting one is dropped, never the newest | chapter 2 |
| tlmd's stage channels | room for 16 values each | chapter 1 |
| tlmd's parse workers | 4, at most 4 lines in flight | chapter 1 |

The **stages** of a planning run, in order:

1. **Load** the network once for every caller (`sync.OnceValues`, chapter 3).
2. **Predict**: for each station, every satellite it can serve, lazily, a day at a time
   (`Passes`, chapter 4), at most twelve stations at once (chapter 2's group, chapter 12).
3. **Merge**: each station's eighty sorted streams into one stream in AOS order (`MergeBy`,
   a priority queue, chapter 6), filtered to candidates (chapter 5's adapters).
4. **Assign**: the rule above, per station, from the dishes' free times (chapter 12).
5. **Collect**: the twenty stations' contacts merged into one schedule by AOS (`MergeBy`).
6. **Publish**: compare-and-swap into the publisher, newer runs only (chapter 3).
7. **Deliver**: each station's contacts to its dish loader through a queue of 256 (wait when
   full), the schedule to the Board through a queue of 4 (drop the oldest), chapter 2.

A run is cancellable at every stage: every goroutine selects on its context, the predictor
is asked again only between days, and a canceled run publishes nothing and leaves no
goroutine behind.

#### tlmd's ingest (chapter 1)

The ingest is a pipeline: one `read` stage per telemetry feed, `merge`, `parse`, then the
`newest` sink. A **telemetry feed** is the stream of telemetry lines one station hears; lines
are the bible's telemetry lines. A satellite heard by two stations appears in both feeds, so its
frames can arrive twice and out of order; the newest frame is the one with the highest
sequence number, whatever order it arrived in.

```go
// Heard is one telemetry line as a station heard it.
type Heard struct {
	Station string // "TIR"
	Line    string // "cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=NOMINAL"
}

// Health is each satellite's newest frame (its highest Seq), by satellite ID.
type Health map[string]TelemetryFrame
```

- Each stage owns the channel it returns: it alone sends on it and it closes it when it is done.
  Channels have room for 16 values. Every send and every receive also selects on `ctx.Done()`.
- `parseLine` is given read-only in chapter 1: Season 1's reading of a line, by `strings.Fields`,
  `strconv.ParseUint`, `time.Parse(time.RFC3339, ...)` and `strconv.ParseFloat`. The `parse` stage
  drops a line `parseLine` refuses and goes on. Chapter 7 makes `parseLine` strict.
- `readFeed` reports a broken feed by cancelling the run with `context.WithCancelCause`'s
  cancel, the cause wrapped as `feed YST: <err>`; `ingest` returns `context.Cause(ctx)` after
  every stage has stopped, or nil and the Health when every feed ended.
- `parseInOrder` parses with four workers and sends frames in the order their lines arrived:
  each line gets a channel of its own for its result, queued in order (room for `workers`), and
  the sender reads those channels in turn. At most `workers` lines are parsed at once.

#### Telemetry lines, strictly (chapter 7)

A line tlmd accepts is exactly what `formatLine` writes: six fields, one space apart, nothing
before or after.

```go
// formatLine writes f as a telemetry line, the one form tlmd accepts.
func formatLine(f TelemetryFrame) string {
	return fmt.Sprintf("%s %d %s batt_v=%.2f temp_c=%.1f mode=%s",
		f.SatID, f.Seq, f.At.UTC().Format("2006-01-02T15:04:05Z"), f.BattV, f.TempC, f.Mode)
}
```

- The satellite ID is 1 to 24 characters of lowercase letters, digits and `-`, starting with a
  letter. The sequence number is digits only, no sign, no leading zero, at most 4294967295. The
  time is `2006-01-02T15:04:05Z` exactly. `batt_v` has exactly two decimals and is from 0.00 to
  20.00; `temp_c` exactly one and is from −100.0 to 150.0, and never `NaN` or an infinity (a
  NaN compares false with every number, which is how `batt_v=nan` slipped past the range check).
  The mode is `NOMINAL`, `SAFE`, `DETUMBLE` or `LEOP`.
- The strict `parseLine` checks fields in order and returns the first problem, wrapped as
  `line: <problem>`: `want 6 fields, got 5`; `satellite "Pipit-1": want lowercase letters,
  digits and -`; `sequence "-1": want a whole number from 0 to 4294967295`; `time
  "2029-04-14T03:12:00+00:00": want 2006-01-02T15:04:05Z`; `field 4: want batt_v=`; `batt_v
  "nan": want volts like 7.91`; `batt_v 31.20: want 0.00 to 20.00`; `temp_c "1e1": want
  degrees like -4.5`; `mode "nominal": want NOMINAL, SAFE, DETUMBLE or LEOP`. Its last check
  is the round trip: a line that parses but isn't `formatLine` of its own frame (such as
  `batt_v=+7.91` or `batt_v=07.91`) is `line: not canonical`.
- Thistle-2's lines of Saturday 14 April, as Joss forwarded them, include
  `thistle-2 2214 2029-04-14T03:36:12Z batt_v=nan temp_c=1.0 mode=SAFE` (heard at Port Ellery,
  pass `thistle-2/PEL/00560`) and the line that crashed the old parser on the next pass, `thistle-2 2290 2029-04-14T05:10:41Z batt_v=7.40` (four fields). Thistle-2 numbers a frame about every 75 seconds. Mirela keeps four of Saturday's lines as corpus files: those two (`saturday-0336`, `saturday-0510`), `thistle-2 2216 2029-04-14T03:38:42Z batt_v=7.42 temp_c=-Inf mode=SAFE` (`saturday-0338`) and the 05:13:11 line that lost its satellite ID and starts with a space (`saturday-0513`). Joss's fourteen odd lines (7.8) are frame 2214 with the battery at 7.91 V written fourteen wrong ways: `batt_v=+7.91`, `batt_v=07.91`, a time ending `+00:00`, the sequence `02214`, `batt_v=7.9`, `temp_c=1`, `temp_c=1e0`, `batt_v=NaN`, `mode=safe`, `Thistle-2`, the sequence `-1`, a time ending `.000Z`, two spaces after the satellite ID, and a seventh field `crc=6e`.
- The **Pipit beacon** (chapter 7.6) is the bible's 12 bytes, decoded into
  `Beacon{Sat uint8; Seq uint16; Mode uint8; BattMV uint16; TempDeciC int16; Flags uint8}`:
  `decodeBeacon(b []byte) (Beacon, error)` refuses a length other than 12 (`beacon: want 12
  bytes, got 7`), wrong sync bytes (`beacon: no sync`) and a bad XOR (`beacon: bad check
  byte`); `encodeBeacon(bc Beacon) []byte` writes it back. Every Pipit sends this beacon, and Pipit-1 is satellite number 1 and Pipit-2 number 2, as in Season 1.

How fuzz lessons run (the app's `go test` runs a fuzz test's seed corpus only, never the
fuzzer): the read-only `main_test.go` holds `FuzzParseLine` (and `FuzzBeacon`). It adds every
seed in the learner's `seeds.go` (`var lineSeeds []string`) with `f.Add`, and `go test` also
runs every file in `testdata/fuzz/FuzzParseLine/` (format: the line `go test fuzz v1`, then
`string("...")`). Its target calls the learner's property, `checkLine(t, parse, line)` in
`checks.go`, through `safely`, which turns a panic into the failure `parseLine("...") panicked:
<value>`, so a crash reads as a failed check rather than a stack dump. Mirela's reference parser, which accepts exactly the lines `formatLine` writes with a valid satellite ID, readings in range and a known mode, is hidden in `reference_test.go`; mutant parsers built from it live in the test file; tests that prove the seeds or the property catch them re-run the test binary
(`os.Args[0] -test.run=^FuzzParseLine$`, `APOGEE_MUTANT=<name>`) and expect it to fail, as
Season 1's chapter 15 did. Readmes say once how to fuzz for real in a terminal:
`go test -fuzz=FuzzParseLine -fuzztime=30s`, and that a failure it finds is saved under
`testdata/fuzz/FuzzParseLine/` and runs as a seed from then on.

#### Iterators and generic types (chapters 4 to 6)

```go
type Predictor func(e Elements, s Site, from, to time.Time) []PassView

func Days(from time.Time, n int) iter.Seq[time.Time]                               // 4.1: n midnights
func Passes(predict Predictor, e Elements, s Site, w Window) iter.Seq[PassView]     // 4.2: a day at a time
func (sch *Schedule) At(station string) iter.Seq2[int, Contact]                     // 4.4: (dish, contact) by AOS
func Catalog(n *Network, w Window, predict Predictor) iter.Seq[PassView]            // 4.5: servable, in byAOS order
func ReadPasses(r io.Reader) iter.Seq2[PassView, error]                             // 4.6: the pass archive
func parsePass(line string) (PassView, error)                                       // 4.6: given, one archive line
func Archive(open func() (io.ReadCloser, error)) iter.Seq2[PassView, error]         // 4.7: closes on every way out
func Bookable(n *Network, w Window, predict Predictor, operator string) iter.Seq[PassView] // 4.8

func Filter[V any](seq iter.Seq[V], keep func(V) bool) iter.Seq[V]   // 5.2
func Map[V, W any](seq iter.Seq[V], f func(V) W) iter.Seq[W]         // 5.3
func Take[V any](seq iter.Seq[V], n int) iter.Seq[V]                 // 5.3
func Zip[V any](a, b iter.Seq[V]) iter.Seq2[V, V]                    // 5.6: pairs, until either ends
func Merge[V any](a, b iter.Seq[V], cmp func(V, V) int) iter.Seq[V] // 5.7: two sorted into one
func Compare(old, new iter.Seq[Contact]) iter.Seq[Change]            // 5.9: the shadow run
func changes(olds, news []Contact) []Change                          // 5.9: given, one second's contacts matched by pass ID

// Change is one difference between two schedules.
type Change struct {
	Kind    string  // "old only", "new only" or "dish"
	Contact Contact // the old planner's for "old only", the new planner's otherwise
}

type passHeap []PassView // 6.2: heap.Interface by byAOS

// Queue is a priority queue: Pop returns the smallest value by cmp.
type Queue[T any] struct{ h *heapOf[T] } // heapOf implements heap.Interface for a slice and cmp

func NewQueue[T any](cmp func(a, b T) int) *Queue[T] // 6.4
func (q *Queue[T]) Push(v T)
func (q *Queue[T]) Pop() (T, bool)  // the zero value and false when empty
func (q *Queue[T]) Len() int
func (q *Queue[T]) Drain() iter.Seq[T] // 6.5: pops until empty or the loop stops
func (q *Queue[T]) Peek() (T, bool)    // 6.6

func MergeBy[T any](cmp func(a, b T) int, seqs ...iter.Seq[T]) iter.Seq[T] // 6.7: k sorted into one
func stationPasses(n *Network, s Site, w Window, predict Predictor, minElevDeg float64) iter.Seq[PassView] // 6.8
```

- An iterator stops the moment `yield` returns false and never calls it again (Go panics with `runtime error: range function continued iteration after function for loop body returned false`).
  An iterator that opened something closes it on every way out, with `defer` inside the
  iterator function.
- Generated sequences are made a day at a time: `Passes` asks the predictor for one day, yields
  it, and only then asks for the next; `Catalog` predicts every pair for one day, sorts that
  day by `byAOS`, yields it, and moves on. Breaking out of either after the first pass costs
  one day of prediction, never the quarter.
- `Merge` and `MergeBy` keep ties stable: on equal keys the earlier sequence's value comes
  first. They use `iter.Pull` and call every `stop` they got, whether the loop ends or breaks.
  `MergeBy` holds one value per input sequence at a time, in a `Queue`.
- `Compare` walks both schedules a second at a time: it gathers each schedule's contacts with the earliest AOS either has left and yields what `changes` finds among them (the old planner's `old only` and `dish` changes in its order, then `new only` in passplan 2's), so contacts at different stations that rise in the same second (5 April, 05:00:11: `brine-3/SQT/16651` and `furrow-a/VSK/23251`) match in whatever order each planner lists them. On Thursday 5 April Tirrin Rock's host parked both dishes from 07:30 to 08:24 to clear their ice, a stow only the old planner heard of: its 273 contacts and passplan 2's 276 differ in four, all at Tirrin Rock (`new only` `furrow-b/TIR/23250` 07:38:06, `pipit-2/TIR/21175` 07:52:22 and `pipit-3/TIR/11124` 08:23:30, then `dish` `furrow-a/TIR/23253` 08:24:16, on dish 2 for passplan 2 and dish 1 for the old planner). passplan 2 hears of stows from Friday 6 April, and planned again with the stow, Thursday has 273 contacts in both.
- The **pass archive** is one pass per line, `<pass ID> <AOS> <LOS> <max elev>`, times in RFC
  3339 UTC, elevation with one decimal:
  `furrow-a/SQT/24534 2029-06-28T14:02:10Z 2029-06-28T14:12:04Z 47.9`. A bad line yields its
  error, `archive line 3: <problem>`, and ends the sequence. The archive is Apogee's file of every past pass; lessons quote it from 1 January 2029 on, where the model's passes begin. `parsePass` reads a line; its problems read `3 fields, want 4`, `bad pass ID "…"`, `bad AOS "…"`, `bad LOS "…"` and `bad max elevation "…"`. A failed read yields the reader's error, unwrapped, after the passes before it, and `Archive` yields `open`'s error unwrapped.

#### Time in a bubble (chapter 8)

Lessons test timers with `testing/synctest`. A bubble's clock starts at 2000-01-01 00:00 UTC,
so code takes `now func() time.Time` for instants and uses the `time` package (timers,
tickers, `Sleep`, context deadlines) for waiting; a test shifts the bubble's clock to 2029:

```go
synctest.Test(t, func(t *testing.T) {
	shift := time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC).Sub(time.Now())
	now := func() time.Time { return time.Now().Add(shift).UTC() }
	// ... a day of passes, in a millisecond of real time
})
```

```go
func waitUntil(ctx context.Context, now func() time.Time, at time.Time) error             // 8.1
func every(ctx context.Context, d time.Duration, f func(at time.Time))                     // 8.4
func planWithin(ctx context.Context, limit time.Duration, plan func(context.Context) (*Schedule, error)) (*Schedule, error) // 8.5
func newDebouncer(quiet time.Duration, fire func()) *Debouncer                             // 8.6
func (d *Debouncer) Notify()
func (d *Debouncer) Stop()
func flyDay(ctx context.Context, now func() time.Time, contacts []Contact, fly func(Contact)) error // 8.8
```

- `waitUntil` returns nil at `at` exactly (at once when `at` has passed) and `ctx.Err()` if the
  context ends first; it stops its timer either way.
- `every` calls `f` at each tick until the context ends, then stops its ticker and returns: a day
  has 144 ticks of 10 minutes.
- `planWithin` gives the run a context that ends after `limit` and cancels it when the run returns; a run past its limit returns `plan: took longer than 30s: context deadline exceeded` (the message names `limit`), which wraps `context.DeadlineExceeded`; any other error, a caller's cancel included, comes back unchanged.
- A `Debouncer` calls `fire` once, `quiet` after the last of a burst of `Notify` calls; `Stop` cancels a pending fire, and a later `Notify` arms it again. It holds one `time.AfterFunc` timer, made by the first `Notify` and `Reset` by later ones under a mutex; passplan 2's quiet is 5 seconds.
- `flyDay` waits until each contact's AOS in order and calls `fly`; it returns `ctx.Err()` when
  the context ends, having flown exactly the contacts whose AOS came first.
- When a timer and `ctx.Done()` are ready together, a `select` takes either, so no test ends a context at the instant a timer is due: 8.4 cancels at 24:00:01 and 09:47, between ticks, and 8.8 at noon, between contacts at 11:59:51 and 12:00:01.
- Learners write synctest checks in `checks.go` (8.2 and 8.7), proved against mutants as above, re-run with `APOGEE_MUTANT` as in chapter 7. The test runs the learner's check twenty times, stopping at the first failure, so a check that passes only by luck (one that looks without `synctest.Wait`) fails; and it fails a check that calls the code on the real clock, which it tells by a real-clock reading taken just before (a call within a minute of it is outside any bubble; a bubble's clock starts in 2000 and may be moved to any year). A count that a timer's function changes is an atomic, or under a mutex: the race detector sees no order between a test's read and a timer's later write, even in fake time. Inside a bubble a real network read keeps the clock still (in-memory fakes such as `net.Pipe` don't), `T.Run` isn't allowed, and every goroutine the bubble starts must end before the test function returns.

#### Benchmarks and allocations (chapters 9 and 10)

- A learner's benchmark is a function in `bench.go`, such as `benchPlanDay(b *testing.B)`, with
  `for b.Loop()`. Read-only tests run it with `testing.Benchmark` after
  `flag.Set("test.benchtime", "20x")`, so it runs exactly 20 iterations in milliseconds.
  Tests never check nanoseconds: they check iterations, allocations, bytes and the learner's
  own metrics (`b.ReportMetric`). Readmes give the terminal command for the timings: `go test -bench=. -benchmem -count=10`. Tests that set `-test.benchtime` put it back when they end, so a terminal's benchmarks keep their time. A sub-benchmark doesn't inherit `b.ReportAllocs` from its parent; lesson 9.8's test reads the six lines of a re-run of its test binary (`-test.run=^$ -test.bench=^BenchmarkWindows$ -test.benchtime=2x -test.cpu=1`).
- Allocation budgets use `testing.AllocsPerRun` and `BenchmarkResult.AllocsPerOp`, computed on
  Go 1.27.3 by the writer, never guessed. Counts drift under `-race`, so every budget test skips
  there: the test files include `race_on_test.go` (`//go:build race`, `const raceOn = true`)
  and `race_off_test.go` (`//go:build !race`, `const raceOn = false`), and a budget test
  begins `if raceOn { t.Skip("allocation counts differ under -race") }`. Without `-race`, code that formats with `fmt` or pulls with `iter.Pull` still allocates a little more when the collector runs (`planNew` up to about 0.4%; exact with `GOGC=off`), so a budget on such code allows a margin (9.8: 1.5%, its benchmark re-run on one CPU); `planDay` and `assign` count exactly. A range-over-func loop counts the same in a test only when the compiler inlines its function there: out of line, the loop's closure, its state and the variables it writes go on the heap (3 allocations for 10.1's `collect` written with `slices.Grow`), so 10.1 subtracts what `collect` makes for no passes.
- The **old planner** (`planOld` in `planners.go`, given read-only in chapters 9 and 10; a copy in chapter 10 is 9's file byte for byte) is the 2023 passplan's planning brought to the rule and run on passplan 2's pass model: it predicts every pair into one slice, sorts that slice by station, then AOS, then ID, and assigns station by station, growing its slices as it goes and keeping each contact's dish in a map keyed by pass ID. The **new planner** (`planNew`) is per station: chapter 6's `stationPasses` (candidates merged by `MergeBy`), assigned, then the stations' contacts merged by `MergeBy` into a slice sized ahead. Both return the same contacts for the same input, and the new one sorts nothing in place. On one core it is no faster, and heavier: predicting passes is nearly all of both (91 to 96% of `planOld`'s time, 96 to 98% of its allocations), and `MergeBy` starts an `iter.Pull` for every satellite at every station (about six allocations each) and boxes every value it queues (two allocations, through `container/heap`'s `any`). For 6, 12 and 24 hours from midnight on Friday 18 May (303, 619 and 1,237 contacts), `planNew` takes up to about 15% longer and makes 2.1 to 3.3 times the allocations: about 23,600, 30,500 and 43,100 against 7,150, 12,040 and 20,970. No planner that predicts its own passes can make a fraction of `planOld`'s allocations: predicting the 1,564 pairs that share a band is 93 to 94% of them. That is chapter 9's truth on the glass ("Not faster. Heavier. That we can fix."): no later lesson or story line says passplan 2 plans faster than the old planner on one core, or allocates less. Where a story line sets passplan 2 against the old passplan's 4 minutes 39 seconds a day, the difference is the old passplan's own, older pass prediction and passplan 2's twelve planner slots, never passplan 2's planning. Chapter 10 cuts garbage around the merge: `planStationFat` is its own function (`planNew`, readable in 9.6 and 9.8, has no string keys), and chapter 6's `queue.go` and `merge.go` stay as they are, since chapter 12 gives them read-only.
- The **dish loader line** (chapter 10), one contact per line, a day per file: AOS and LOS as `15:04:05`, the dish, the pass ID: `05:29:26 05:37:52 1 furrow-a/VSK/24529`.
- Chapter 10's pieces: a host's report (10.2) is `CWR contacts: 76`, then one line a contact in AOS order, `cirrine-1/CWR/19862 dish 1`; the export (10.5) writes each station's dish loader lines in one write from an 8 KB buffer (Ystmark's 199 lines on take-over day, 7,811 bytes, are the longest day); a planner slot's scratch holds 512 passes (Ystmark's 439 candidates are the most), and `putScratch` keeps none over 65,536 (6.8 MB). The per-station step `planStationFat` (10.10) copies the candidates into a slice grown from nil, skips repeats with a map keyed by `fmt.Sprintf("%s %s %d", SatID, Station, AOS.Unix())`, sorts the copy and assigns it: on take-over day 1,471 allocations and 213 KB at Cape Wren, 1,789 and 251 KB at Ystmark; `planStationLean` 11 and 60 KB, 13 and 93 KB. Planning take-over day at twenty stations from a network loaded once (each station's candidates predicted, the step, the host report and the export of its dish loader lines) makes 81,657 allocations and 12.0 MB with Monday 21 May's code and 42,705 and 6.0 MB with Friday 25 May's (the debrief's "6 MB of garbage instead of 12"); of what is left, predicting passes is 99% of the allocations and 77% of the bytes. No line compares these with `planOld`.
- A `sync.Pool` may lose its values to any collection, keeps a part for each P (a CPU the program runs on), and after a collection allocates that part again on its next use (2 allocations, `New` not called). No check needs a value back from a pool through a collection: 10.6 counts its rounds with the collector off (`debug.SetGCPercent(-1)`) and checks that a 65,536-pass scratch is kept on one CPU (`runtime.GOMAXPROCS(1)`) with the collector off, skipped under `-race`, where `Put` drops one value in four; 10.10 measures on one CPU. With a forced collection before every pool call, every check of 10.6 and 10.10 still passes but 10.6's count, which then reads the pool's own 4 allocations a round.
- **Tomás's strip** (10.8 and 10.9, Thursday 24 May, 17:10 and 17:40): a preview of take-over
  day on the test Board, one `glance` for each of the window's 613,496 candidates. As Tomás
  declared it (`Sat` and `Station`, uint8 places in fleet and table order; `AOS` and `LOS`,
  int64 Unix seconds; `MaxElevDeg` float32; `Dish` uint8, 0 when dropped; `Booked` bool) a
  `glance` is 32 bytes, aligned to 8: 6 bytes of padding after `Station` and 2 at the end,
  drawn `12......3333333344444444555567..`, 19,631,872 bytes for the strip. His notebook's sum
  is 24 bytes, 14,723,904 bytes; ordered by alignment, largest first, a `glance` is exactly
  that, 4.9 MB less. Neither the compiler nor `go vet` reorders fields or reports padding.
  The strip is the test Board's memory, not passplan 2's garbage, so no other number of
  chapter 10 counts it.

#### Struct tags (chapter 11)

One validator reads the `check` tag of every exported field; `applyDefaults` reads `default`.
A `check` tag is rules separated by commas: `required` (not the zero value; a slice not empty),
`min=<n>` and `max=<n>` (ints, floats, and durations written as duration text), `len=<n>` (a
string's length in characters, or each string of a `[]string`), `oneof=<a> <b> ...` (a string,
or each string of a `[]string`). A field whose type has no use for a rule is a mistake in the tag, and says so: `Config.Slots: rule len=3 does not fit int`; so does a `min` or `max` bound that doesn't parse for the field's type (`Config.Slots: rule max=6A does not fit int`). An unknown rule is `HostConfig.Name: unknown rule "sometimes"`. A `time.Duration` is told by its type (its kind is `int64`); `oneof` takes words separated by spaces and `len` counts characters, both on a `string` or each string of a `[]string`. A field with no `check` key, or an empty one, has no rules.

```go
// HostConfig is the form a station's host fills in when its dish is installed.
type HostConfig struct {
	Code    string   `check:"required,len=3"`
	Name    string   `check:"required"`
	LatDeg  float64  `check:"min=-90,max=90"`
	LonDeg  float64  `check:"min=-180,max=180"`
	MaskDeg float64  `check:"min=0,max=15"`
	Dishes  int      `check:"min=1,max=2"`
	Bands   []string `check:"required,oneof=UHF S X"`
}

// PipelineConfig is tlmd's ingest settings.
type PipelineConfig struct {
	Workers int      `check:"min=1,max=32" default:"4"`
	Buffer  int      `check:"min=0,max=1024" default:"16"`
	Feeds   []string `check:"required,len=3"` // station codes
}
```

- `Validate(v any) error` takes a struct or a pointer to one and checks every exported field in order and walks into each struct field, embedded or not, right after that field's own rules (pointer fields and slices are not followed), returning every problem joined with `errors.Join`, or nil; a nested struct's problems name its own type (`HostConfig.Dishes` inside a struct that embeds `HostConfig`). Each
  problem is `<Type>.<Field>: <problem>`: `HostConfig.Dishes: 0 is below min 1`,
  `Config.Slots: 70 is above max 64`, `Config.Turnaround: 15m0s is above max 10m0s`,
  `HostConfig.Bands: "K" is not one of UHF S X`, `HostConfig.Code: "VRLX" is not 3 characters`,
  `HostConfig.Name: required`. Ints print with `%d`, floats with `%g`, durations with
  `String()`, strings with `%q`. A value that isn't a struct is `validate: want a struct, got %T`: `got int`, `got []main.HostConfig`.
- `applyDefaults(ptr any) error` sets each zero exported field that has a `default` tag, through
  a pointer (`defaults: want a pointer to a struct, got main.Config`); it leaves set fields, unexported fields and fields without a `default` key alone. Zero means "use the default", so a `Buffer: 0` becomes 16; a default the field can't read returns `setText`'s error, wrapped as `defaults: <Field>: …`.
- The forms of the third wave: Vorrel Island's says `Dishes: 0` (the host meant two); Ilha Calmaria's first form (Thursday 31 May) left the station's name blank and listed bands `S` and `K`, and its corrected form (Friday 1 June) still listed `S` and `K`; the other three are right, and all five were right by go-live. The test deploy of Friday 1 June
  had a `Config` with `Slots: 0`.

#### Take-over day

Wednesday 27 June and Thursday 28 June 2029, UTC:

| When | What |
|---|---|
| Wed 18:05 | Furrowlight's rideshare lifts off with twenty-four Furrows |
| 19:10–20:40 | They leave the rocket, one about every four minutes; elements valid from 20:40:00 |
| 19:56 | The launcher's elements arrive; the old passplan starts replanning the day for eighty satellites |
| 20:00:39 | Its plan is ready, after 4 minutes 39 seconds |
| 20:40:16 | First contact of the twenty-four, sixteen seconds after their elements begin: `furrow-g/YST/00001`; by 22:28:50 (`furrow-c/TIR/00002`) all have been heard |
| 23:30 | Ingrid's go/no-go poll for the take-over: GO |
| Thu 00:00:00 | Cut-over. passplan 2's run 1 replans the old planner's schedule from the cut-over, over 28 June to 6 October, and publishes version 1: the old planner's contacts that rose before 00:00 stay in it, flown by the old planner, and the rest is planned from the hand-over: 2,016 contacts on the day, 201,774 in the window (613,496 candidates, 411,722 dropped). The old planner keeps its twelve contacts in progress |
| 00:01:04 | passplan 2's first contact: `cirrine-4/TIR/01069` |
| 00:08:26 | The old planner's last contact ends (`furrow-f/ILV/00003`); at 00:10 it is switched off |
| 05:29:26 | `furrow-a/VSK/24529`, the day's 417th contact: Furrow-A over Varsk Island at 74.9°, Hollis's pass |
| 07:48:33 | `pipit-1/LAN/16797`, the 612th: Pipit-1's first contact of the day. Its two passes over Port Ellery went to Thistle-2, which rose a few minutes earlier each time |
| 09:47:00 | Cape Wren's wind reaches 74 km/h and its dish stows itself; the stow is recorded with no end. Run 60 publishes version 60: Cape Wren dark until further notice, 1,982 contacts on the day (193,239 in the window) |
| 13:23:00 | The wind has stayed under 72 km/h since 13:15, and Cape Wren unstows; the stow now ends at 13:23. Run 83 publishes version 83: 2,011 contacts on the day (201,769 in the window). Six Cape Wren passes were lost to the wind (`pipit-2/CWR/22443` 10:44:34, `cirrine-1/CWR/19868` 10:54:48, `furrow-b/CWR/24530` 11:51:46, `cirrine-1/CWR/19869` 12:28:44, `furrow-g/CWR/00011` 13:11:36, and `furrow-d/CWR/00011` 13:22:45, still above the horizon at 13:23); the afternoon's 21 Cape Wren contacts became 22 others, from `furrow-b/CWR/24531` at 13:23:30. Cape Wren ends the day with 71 contacts instead of 76 |
| 23:59:49 | The day's last AOS: `cirrine-2/BRM/01083` |
| Fri 00:00 | Ingrid writes in the log: "passplan 2: GO for every day." |

The day in numbers: 146 runs (run 1 at the cut-over, 143 ten-minute ticks from 00:10 to
23:50, and the two at 09:47 and 13:23); three published versions (1, 60, 83), the other runs
finding nothing new; 2,011 contacts flown at twenty stations, 819 of them for the twenty-four
new Furrows; no pass flown twice, none touched by both planners. Lessons that run the whole
day test it with their own `Config` (a replan every hour, a one-day window), so it ends in under a second. 12.9's test also puts a `dayModel` in `passModel`'s place: each pair's passes of the day predicted once and every replan's question answered from them, exactly what `predict` gives (checked for every pair from every ten-minute instant and from 09:47 and 13:23), so the day costs one day of prediction: about 0.35 s and 0.65 s of CPU, against 0.7 s and 2.3 s.

**Bookings.** Where a lesson needs a booking (Season 3's `Booking`), its ID is from `bk-7001`
to `bk-7999`: Season 5 reached `bk-5274`, and the numbers below `bk-7000` are left to Season 6.
Furrowlight's booking of `furrow-a/VSK/24529` on take-over day is `bk-7412`; Kessick
Polytechnic's of Thistle-2's first pass over Port Ellery (`thistle-2/PEL/00005`) is `bk-7001`.

#### The code, by chapter

Each lesson is its own `package main` (module `lesson`). A lesson gives the pieces it needs
from earlier chapters as read-only files under these names; `network.go` and `model.go` are
hidden unless the readme sends the learner to read them. `types.go` holds the types (Season 1's and this season's: `Site`, `Elements`, `Network`, `Window`, `Contact` and the rest), cut to what the lesson needs; `network.go` is the `newNetwork` block alone.

| Chapter | The learner writes (files) | Given |
|---|---|---|
| 1 | `read`, `parse` (cancellable), `newest`, `merge`, `ingest`, `parseInOrder` (`ingest.go`) | `Heard`, `Health`, `parseLine`, `readFeed` from 1.7 (`telemetry.go`) |
| 2 | `Slots`, `newSlots`, `Acquire`, `Release`, `predictAll`, `runPool`, `feed`, `offer`, `Group`, `newGroup`, `Go`, `fail`, `Wait`, `Uploads`, `newUploads`, `(*Uploads).Do` (`slots.go`, `predict.go`, `pool.go`, `queue.go`, `group.go`, `uploads.go`) | `Elements`, `PassView`, `Contact`, `Schedule`; 2.1's `predictAll` (`predict.go`); `network.go` and `model.go` in 2.8; fake predictors in tests |
| 3 | `Tables`, `(*Tables).Network`, `networkOnce`, `Run`, `newRun`, `Stop`, `Done`, `Stats`, `(*Stats).Add`, `Snapshot`, `Publisher`, `Current`, `Update`, `Offer` (`tables.go`, `run.go`, `stats.go`, `publish.go`) | `loadNetwork() (*Network, error)` (the slow loader: in tests it counts its calls), `network.go` |
| 4 | `Days`, `Passes`, `(*Schedule).At`, `Catalog`, `ReadPasses`, `Archive`, `Bookable` (`seq.go`, `archive.go`) | `network.go` and `model.go` (hidden); `Predictor` above `Passes` in `seq.go`; `byAOS` as it stands until 6.1 (the AOS alone) and `canServe` (`plan.go`) in 4.5 and 4.8; `parsePass` (`passline.go`) in 4.6 and 4.7; in 4.4 the schedule of 2 April, planned by a hidden test file from hidden `network.go`, `model.go` and `plan.go` |
| 5 | `Filter`, `Map`, `Take`, `Zip`, `Merge`, `Compare` (`adapt.go`, `shadow.go`); `main` in 5.1 and 5.5 | `types.go`; the earlier adapters back in the editable `adapt.go`, read-only in 5.5; `Change` above `Compare` in `shadow.go`, and `changes` (`match.go`, read-only) in 5.9; the schedule of 2 April, planned from scratch by a hidden test file (a hidden `wave.go` in 5.5) from hidden `network.go`, `model.go` and `plan.go`, with `byAOS` as it stands until 6.1; in 5.9 both planners' schedules of 5 April, planned the same way by a hidden test file |
| 6 | `byAOS`, `passHeap`, `Queue`, `heapOf`, `NewQueue` and methods, `MergeBy`, `stationPasses` (`plan.go`, `queue.go`, `merge.go`) | `network.go`, `model.go`, `plan.go`, `Predictor` and `Passes` (`seq.go`), `Filter` (`adapt.go`) |
| 7 | `parseLine` (strict, `parse.go`), `lineSeeds` (`seeds.go`), `checkLine` (`checks.go`), `decodeBeacon` (`decode.go`) | `formatLine` (`telemetry.go`), `encodeBeacon` and `Beacon` (`beacon.go`), the fuzz targets, Mirela's reference parser (hidden), corpus files |
| 8 | `waitUntil`, `checkWaitUntil` (`checks.go`), `every`, `planWithin`, `Debouncer`, `checkDebouncer` (`checks.go`), `flyDay` (`timers.go`) | `Contact`, `Schedule` (`types.go`); the earlier `timers.go`, read-only, in 8.2 and 8.7; `notifier` (`checks.go`); in 8.8 the day's contacts as a fixture, planned by a hidden test file from hidden `network.go`, `model.go` and `plan.go`; Nkem's mutants in the tests |
| 9 | `benchPlanDay`, `benchFresh`, `benchStream`, `sameContacts`, `benchWindows` (`bench.go`) | `june28`, `stationDay`, `dayCounts`, `planned` and `planDay` (`day.go`, read-only in 9.1, 9.2, 9.4 and 9.5); `planOld` and `planNew` (`planners.go`, read-only in 9.6 and 9.8) with chapter 6's `seq.go`, `adapt.go`, `queue.go` and `merge.go`, hidden; `plan.go` (read-only in 9.1 to 9.5, hidden in 9.6 and 9.8); `network.go` and `model.go`, hidden; `Config` in `types.go` in 9.6 and 9.8 |
| 10 | `collect`, `stationReport`, `appendContact`, `exportStations`, `getScratch`, `putScratch`, `passKey`, `dedupe`, `planStationLean` (`lean.go`); `main`'s marks (`main.go`, 10.8); `glance`'s field order (`glance.go`, 10.9) | `planStationFat` (`fat.go`, 10.10); earlier solutions read-only: `appendContact` (`loader.go`, 10.5), the pool with `getScratch`, `putScratch` and `maxScratch` (`scratch.go`, 10.10), `passKey` (`keys.go`, 10.10); `plan.go` (read-only in 10.10, hidden in 10.1, 10.2 and 10.4 to 10.7); `types.go`; hidden `network.go`, `model.go`, chapter 9's `day.go` (10.1, 10.6, 10.7, 10.10) and chapter 8's `day_test.go` (10.2, 10.4, 10.5); in 10.8 and 10.9 only Tomás's `glance` (`glance.go`, read-only in 10.8), and in 10.8 `mark` and `candidates` (`main.go`) |
| 11 | `fieldLines`, `tagLines`, `describe`, `checkRule`, `measure`, `applyDefaults`, `Validate` (`validate.go`) | `Config`, `HostConfig`, `PipelineConfig` (`types.go`); `main`, with 11.2's `Note` (`main.go`, in 11.1 and 11.2); from 11.6, `checkRule` cutting each rule at `=` with the `min` and `max` messages and `misfit`, in the editable file; `setText` (`defaults.go`) in 11.8; in 11.10, read-only, `checkRule`, `misfit` and `measure` (`rules.go`), `applyDefaults` and `setText` (`defaults.go`) and the hosts' forms as first sent, `hostForms` and `form` (`forms.go`) |
| 12 | `planStation`, `Plan`, `handover`, `replanFrom`, `(*Replanner).Run`, `heard`, `takeOver` (`planner.go`, `replan.go`) | everything above, read-only (`network.go` and `model.go` hidden): `plan.go` (with `stowedDuring`), 4.8's `seq.go`, 5.7's `adapt.go`, chapter 6's `queue.go` and `merge.go`, `group.go`, `publish.go`, `timers.go`; `passmodel.go` (`passModel`, the model `planStation` asks, which tests replace); `Stow`, `Schedule`, `Config` and `Event` in `types.go`; the rest of the `Replanner` (`newReplanner`, `publish`) in `replan.go` from 12.6; the earlier `planner.go`, read-only, in 12.6, 12.7 and 12.9; answer keys in a hidden `ref_test.go` (all but 12.6) |

Signatures every lesson after the one that defines them uses exactly (with the ones listed in
the sections above):

```go
// chapter 1
func read(ctx context.Context, station string, r io.Reader) <-chan Heard
func readFeed(ctx context.Context, cancel context.CancelCauseFunc, station string, r io.Reader) <-chan Heard
func parse(ctx context.Context, in <-chan Heard) <-chan TelemetryFrame
func newest(in <-chan TelemetryFrame) Health
func merge(ctx context.Context, feeds ...<-chan Heard) <-chan Heard
func ingest(ctx context.Context, feeds map[string]io.Reader) (Health, error)
func parseInOrder(ctx context.Context, in <-chan Heard, workers int, parse func(string) (TelemetryFrame, error)) <-chan TelemetryFrame

// chapter 2
type Slots chan struct{}
func newSlots(n int) Slots
func (s Slots) Acquire(ctx context.Context) error // nil once a slot is held; ctx.Err() if ctx ends first
func (s Slots) Release()
func predictAll(ctx context.Context, els []Elements, slots int, predict func(context.Context, Elements) ([]PassView, error)) (map[string][]PassView, error)
func runPool(ctx context.Context, jobs <-chan Elements, workers int, do func(context.Context, Elements) error) error
func feed(ctx context.Context, contacts []Contact, q chan<- Contact) error
func offer(q chan *Schedule, s *Schedule) (dropped int)
func newGroup(ctx context.Context, limit int) (*Group, context.Context)
func (g *Group) Go(job func(ctx context.Context) error) // waits for a slot; runs nothing once the group's context has ended, and fails the group with that cause
func (g *Group) Wait() error                            // the first error, after every job has returned
func newUploads(perStation int) *Uploads
func (u *Uploads) Do(ctx context.Context, station string, upload func() error) error

// chapter 3
func loadNetwork() (*Network, error)
func (t *Tables) Network() (*Network, error)
func networkOnce(load func() (*Network, error)) func() (*Network, error)
func newRun() *Run
func (r *Run) Stop()
func (r *Run) Done() <-chan struct{}
func (s *Stats) Add(planned, dropped int) // one run's counts
func (s *Stats) Snapshot() (planned, dropped, runs int64)
func (p *Publisher) Current() *Schedule                      // nil before the first schedule
func (p *Publisher) Update(contacts []Contact) *Schedule     // a new schedule, Version one more
func (p *Publisher) Offer(s *Schedule) bool                  // only a newer Version replaces the current

// chapter 7
func parseLine(line string) (TelemetryFrame, error)
func checkLine(t *testing.T, parse func(string) (TelemetryFrame, error), line string)
func decodeBeacon(b []byte) (Beacon, error)
func encodeBeacon(bc Beacon) []byte

// chapter 8
func checkWaitUntil(t *testing.T, wait func(context.Context, func() time.Time, time.Time) error)
type notifier interface{ Notify(); Stop() }
func checkDebouncer(t *testing.T, newD func(quiet time.Duration, fire func()) notifier)

// chapter 9
func planOld(n *Network, w Window, cfg Config) []Contact
func planNew(n *Network, w Window, cfg Config) []Contact
func stationDay(code string, day time.Time) (Site, []PassView)                   // the station and its candidates that day, in fleet order
func planDay(cands []PassView, sites []Site, turnaround time.Duration) []Contact // sorts cands in place, then assigns
func benchPlanDay(b *testing.B)
func benchFresh(b *testing.B)
func benchStream(b *testing.B)
func sameContacts(old, new []Contact) error // "contact 213: old cirrine-5/ILV/00445 dish 1, new cirrine-5/ILV/00445 dish 2"
func benchWindows(b *testing.B)             // sub-benchmarks old/hours=6 ... new/hours=24

// chapter 10
func collect(seq iter.Seq[PassView], n int) []PassView
func stationReport(code string, contacts []Contact) string // "CWR contacts: 76", then "cirrine-1/CWR/19862 dish 1" a contact
func appendContact(dst []byte, c Contact) []byte
func exportStations(w io.Writer, byStation map[string][]Contact, codes []string) error
func getScratch() *[]PassView
func putScratch(p *[]PassView) // keeps none over maxScratch (65536) passes
type passKey struct {          // a pass's identity without formatting its ID
	SatID, Station string
	AOS            int64 // Unix seconds
}
func dedupe(passes []PassView) []PassView
func planStationFat(s Site, cands []PassView, turnaround time.Duration) []Contact
func planStationLean(s Site, cands []PassView, turnaround time.Duration) []Contact

// chapter 11
func fieldLines(v any) []string // "Slots int int"
func tagLines(v any) []string   // "Slots: min=1,max=64"
func describe(v any) []string   // "Slots=12"
func checkRule(name string, v reflect.Value, rule string) error
func misfit(name, rule string, t reflect.Type) error             // "Config.Slots: rule len=3 does not fit int"
func measure(v reflect.Value, text string) (int, any, any, bool) // cmp.Compare's result, v and the bound for %v; false when the bound doesn't fit v
func applyDefaults(ptr any) error
func setText(v reflect.Value, text string) error // a duration, an int, a float or a string, from text
func Validate(v any) error

// chapter 12
var passModel Predictor = predict // passmodel.go: the pass model planStation asks; tests replace it
func planStation(ctx context.Context, n *Network, s Site, w Window, free []time.Time, stows []Stow, cfg Config) ([]Contact, int, error)
func Plan(ctx context.Context, n *Network, w Window, free map[string][]time.Time, stows []Stow, cfg Config) (*Schedule, error)
func handover(cur *Schedule, at time.Time, turnaround time.Duration) map[string][]time.Time
func replanFrom(ctx context.Context, n *Network, cur *Schedule, at time.Time, stows []Stow, cfg Config) (*Schedule, error)
// Event is news for the replanner: a tick, or a stow that begins or ends.
type Event struct {
	At   time.Time
	Stow *Stow // nil for a tick; a stow with a zero To begins, one with a To ends it
}
func newReplanner(n *Network, pub *Publisher, cfg Config, now func() time.Time) *Replanner
func (r *Replanner) Run(ctx context.Context, events <-chan Event) error
func heard(stows []Stow, st Stow) []Stow // 12.7: st begun, or its station's open stow ended, in a new slice
func takeOver(ctx context.Context, n *Network, old *Schedule, at time.Time, events <-chan Event, cfg Config, now func() time.Time, fly func(Contact)) error
```

- `Plan` runs `planStation` for every station in a group of `cfg.Slots`, merges the twenty
  contact lists with `MergeBy`, and returns a schedule whose `Version` the caller sets; its
  first error cancels the other stations and is returned.
- `replanFrom` keeps `cur`'s contacts with AOS before `at`, plans from `at` with `handover(cur,
  at, cfg.Turnaround)`, and returns both together in a new slice, by AOS then ID, over `cur`'s window; its `Dropped` counts the passes from `at` on. `handover` gives a station an entry for each dish up to the highest its contacts before `at` used; a dish with no entry is free now, as `assign` reads `free`.
- A `Replanner` takes the next run number for every event (a tick or a stow), cancels the run in
  progress, plans from the event's instant with every stow it has heard of, and offers the
  result; a run whose contacts equal the current schedule's offers nothing. Its runs are numbered after the publisher's current version; `Run` also replans every `cfg.Replan` on a ticker of its own, plans each run through `planWithin(ctx, cfg.RunLimit, …)`, waits for its runs before returning (when `events` closes, the last one still finishes), and leaves word in its `news` channel (room for one) after each schedule it publishes, which the station runner reads.
- `takeOver` replans a copy of `old` given passplan 2's window (`cfg.Days` days from `at`) from `at` and publishes it as version 1; it runs a `Replanner` on `events` and flies each contact of the newest schedule at its AOS, the first after the last one flown in `byAOS` order (starting from a pass at `at` with no ID), waking on the `Replanner`'s `news` to read the schedule again; once `ctx` ends it waits for `Run` and returns `ctx.Err()`.

#### How lessons test it

- Everything stays in the process: no network, no files outside `t.TempDir()`, no real clock in
  a check. Clocks are injected or come from a synctest bubble; randomness is
  `rand.New(rand.NewPCG(1, 2))`; real sleeps are 50 ms at most.
- Concurrency checks are correct by construction: limits are measured with an `atomic.Int64`
  of jobs in flight and its highest value (a compare-and-swap loop), and parallelism with a
  barrier (fake jobs that wait until the expected number are in flight, with a one-second
  guard), never by timing.
- A **leak check** takes `runtime.NumGoroutine()` before the code runs and, after cancelling,
  polls every millisecond for up to a second until the count is back; a leak fails as
  `goroutines left after cancel = 3, want 0`. (A leaked goroutine inside a synctest bubble
  ends the whole test binary with a deadlock panic, so leak checks never use a bubble.)
- Every wait in a test has a guard of about a second, so a wrong starter fails within a second
  or two and never hangs the run. The guard is for what right code does at once: a wait on real work, such as 2.8's day of predictions (up to 0.4 s under `-race` on a busy machine), gets five seconds of its own, and the guard starts after it. A wait on a loop that can run on without yielding also caps memory: 5.9's walk fails once it has allocated 64 MB (a right walk of a day allocates 64 KB), since a loop that appends without pulling takes gigabytes within the guard's second, and the table stops at the first walk still running. 12.9's station runner test ends the day after 3,000 flights, more than a day holds, so a runner that flies one contact over and over stops. A check that can fail a wrong starter only by waiting out its guard (12.2's twelve stations at once) runs first, with a guard of its own, and the checks after it run only once it passes, so the starter fails in about a second.
- Tests plan a day or a week, lesson 12.8 included. The model watches each satellite's elevation,
  about 4 µs a pass on one core: a day of every pass (7,485) takes 35 ms, a week of servable
  passes 0.2 s, a hundred days 3 s, so a hundred days and a quarter stay story numbers.
- Where learners write checks (chapters 7 and 8) or benchmarks (chapter 9), their code goes in an
  editable non-test file and the read-only test calls it, as Season 1's chapter 15 did.
- Messages read `call = got, want want`; schedules are compared by pass ID and dish:
  `Plan(28 June)[289] = pipit-2/YST/22439 dish 2, want pipit-2/YST/22439 dish 1`.

**Topics added this season** (front matter): `iterators`, `performance`, `reflection`. Fuzzing
and synctest lessons use `testing`; atomics use `goroutines` and `mutexes`.

**New field guide words** (world words only; Go words are taught in readmes): passplan,
passplan 2, planner slot, candidate (pass), contact, dropped pass, turnaround, window, shadow
run, cut-over, hand-over, take-over day, host, wave, rideshare, deployment, elements, first
contact, telemetry feed, the drill, dish loader, pass archive, Thistle-2, planning run, test Board, station runner. Carried over and
needed again: AOS, LOS, pass, pass ID, LEOP, stow, mask, max elevation, catalog, the Board, the
Fishbowl, booking, bookingd, tlmd, UTC, benchmark (Season 6's), band (Season 1's), deploy (Season 4's).

**Name checks** to run before the season ships: the fifteen station names. A web search on
6 November 2026 found nothing real under Tirrin Rock, Varsk Island, Calidra Atoll, Ilha
Brumosa, Ostrel Island, Ilvane Spit, Cape Solvane, Saltreach Island, Corith Point, Nimrie
Point, Hasku Point, Ilha Calmaria or Punta Ventisca (near misses recorded: Cape Solander in
Australia, Kaskulu Point in Alaska, Torrin on Skye, a restaurant called Ventisca in Punta
Arenas, a Calmaria guesthouse in Portugal); Quindle Island and Vorrel Island are not yet
checked. The satellites are numbered or lettered members of names already checked (Thistle,
Cirrine, Brine, Furrow); "passplan 2" is the bible's own service with a number. No new people
or customers.
