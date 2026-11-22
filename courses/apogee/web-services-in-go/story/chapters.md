# what-the-web-speaks

## Briefing

Week one, at the backend desks beside the Fishbowl. Yusra still types forty emailed bookings a day, and soon customers' programs will book passes from Apogee's own service instead. Before anyone builds it, Nkem wants you to learn exactly what a program sends and what comes back, on a stand-in she wrote. One wrong word in a program, and it books the wrong pass.

## Debrief

Friday afternoon, your request crossed a real connection and the stand-in answered with Port Ellery. Nkem sends it to Cirrine's engineers for Monday: "Now they know what the service will say before it exists." In the Fishbowl, Ingrid still wants to know why Tuesday's go/no-go poll waited four minutes.

# asking-politely

## Briefing

Week two, in the Fishbowl. Last Tuesday the go/no-go poll waited four minutes for the weather service while Cape Wren's weather mast was silent, and Furrow-B's checklist finished after AOS. Ingrid wants the poll's wind check to answer within two seconds, or say why not, before another night of passes depends on it.

## Debrief

Friday night, the wind check answers for all five stations within two seconds: Cape Wren stows, Salar Quieto's silent mast says so, and nothing waits four minutes again. Ingrid writes one line in the log: "Poll on time. Every NO-GO has a reason."

# your-first-endpoint

## Briefing

Week three, at the backend desks. Cirrine's engineers start their client this week, and the first thing it will ask Apogee's booking service is which stations exist. Yusra wants that answer right, with her customer guide served beside it, before she sends Cirrine the address on Friday.

## Debrief

Friday 19 November, 15:00 UTC: Yusra sends Cirrine the address, and within the hour their client lists all five stations. Yusra's TL;DR of their reply: "Read the guide, asked us nothing, it worked." Hollis has already written, in capitals, to ask when it is Furrowlight's turn.

# an-address-for-every-pass

## Briefing

Monday 22 November, at the backend desks. Furrowlight's engineers start trying the booking service this week, and Hollis wants a link to every pass in his emails. Yusra adds his second wish: "TL;DR: Furrowlight wants to book a week at a time." Every station, pass and booking needs an address of its own, because a link that points at the wrong pass books the wrong pass.

## Debrief

Friday afternoon, Yusra sends Hollis the first email with a link to a pass: Furrow-A over Salar Quieto on 6 December. His answer comes back in capitals: "THIS ONE. HOW DO WE BOOK IT?" Nkem reads it over your shoulder: "Next week, his program does."

# the-first-booking

## Briefing

Monday 29 November, at the backend desks. On Friday, Hollis's program will book Furrow-A's 14:02 pass over Salar Quieto on 6 December all by itself, and Brinewatch's engineers start trying the service too. Every booking the service accepts must be a pass that exists and is still open, and Furrowlight's will be the first booking that never touches Yusra's spreadsheet.

## Debrief

Friday afternoon, `bk-4501` arrives: Furrow-A over Salar Quieto on Monday at 14:02:10, booked by a program, with no email at all. Yusra forwards it to the whole team: "TL;DR: the first booking nobody typed."

# around-every-request

## Briefing

Week of 6 December, at the backend desks. On Monday at 14:02, Furrow-A flew the first pass ever booked through the service, over Salar Quieto. That afternoon a bug met one of Cirrine's requests: their program got a dropped connection, and nobody could say which request it had been. Nkem wants every request to get an answer, an ID and a line in the log.

## Debrief

Friday afternoon, Cirrine's request from Monday got an answer instead of a dropped connection, and Yusra quoted its line in the log, `req-000002`, in her reply. Nkem's fix for the short orbit number ships with Monday's deploy.

# keeping-bookings

## Briefing

Week of 13 December, at the backend desks. On Monday a deploy restarted the beta, and Furrowlight's twelve test bookings vanished: the service had kept them only in memory. Yusra won't retire her spreadsheet until bookings survive a restart, and the last deploy before the holidays goes out on Friday.

## Debrief

Friday's deploy restarted the service at 16:00 UTC, and every booking came back, Furrowlight's January passes among them. Yusra's TL;DR: "The spreadsheet can start packing." Mirela's review is the first thing on the calendar for 3 January.

# whos-asking

## Briefing

Monday 3 January, back from the holidays, at Mirela's desk. She has spent the morning with the booking service, and her review opens with one finding: anyone who can reach it can read and cancel any booking, Cirrine's and Brinewatch's as easily as Furrowlight's. She signs off on Friday, once the service knows who is asking and answers each customer about its own.

## Debrief

Friday afternoon, each of Mirela's four broken copies failed your review while the real service passed it, and her sign-off, for once, has no condition: "LGTM, it fails closed." Yusra sends the three beta customers their keys before the weekend.

# fair-use

## Briefing

Monday 10 January, at the backend desks. Since Sunday night, Brinewatch's client for the beta, cut off from its own network, has retried without a pause: a request every millisecond. Behind the flood, Cirrine's bookings wait half a minute for an answer. Nobody at Brinewatch meant any harm, but one customer's bug must never slow every other customer down.

## Debrief

Friday's replay: Brinewatch's thousand requests got ten answers and 990 refusals, each saying when to try again, while Cirrine's bookings went through at once. Brinewatch's corrected client now waits as it is told. Yusra's next TL;DR is already in: "Furrowlight asks every five seconds whether its data has landed."

# calling-customers-back

## Briefing

Monday 17 January, at the backend desks. Customers' programs ask the booking service every few seconds whether their data has landed; Furrowlight's asks every five seconds, 17,280 times a day. Yusra wants the service to tell them instead, within a minute of LOS, and Mirela wants every message signed, so that nobody else can tell a customer their data has landed.

## Debrief

Friday at 14:42:39, a minute after LOS, Furrowlight's receiver heard that Furrow-A's image had landed, signed and checked. Hollis answered in capitals: "WE TURNED OFF THE FIVE-SECOND LOOP." In the Fishbowl, Ingrid has printed the launch checklist: eleven days to go.

# signed-in

## Briefing

Friday 21 January, 15:00 UTC, at the backend desks. Right after his capitals about the five-second loop, Hollis wrote again: he wants to see Furrowlight's bookings himself, in a browser, without asking an engineer. Yusra wants a small signed-in page beside the booking service before rehearsal week starts on Monday. Mirela will read every line: a signed-in browser can be made to cancel a booking from any other website.

## Debrief

At 18:20 UTC Hollis signed in and saw Furrowlight's bookings and nobody else's, while Mirela's page with the hidden form was refused. Yusra sent the three beta customers the page's address, and Hollis's capitals came straight back: "NO MORE ASKING ENGINEERS."

# launch-day

## Briefing

Monday 24 January, in the Fishbowl. Ingrid runs launch like a pass: a checklist, rehearsals all week, a full rehearsal on Monday 31 January, then the booking service opens to every customer at 14:00 UTC on Tuesday 1 February. One line of her checklist has never been tried: restart the service during a pass. If a restart drops a booking or a message, a customer finds out first.

## Debrief

The 14:05 restart dropped nothing, and Furrowlight heard about Furrow-A's image at 14:13:03. By evening every customer had booked through the service, and Yusra archived her spreadsheet. Ingrid wants the next passes on every screen in the building. From Ystmark, Sunniva: "New dish controllers in March. No HTTP."

# the-first-week-open

## Briefing

Wednesday 2 February, the first morning after launch, at the backend desks. Real customers do what the beta never did: they attach command files to bookings, fetch ten-gigabyte images over lines that drop, tell the service where to call them back, and book past their plans. Yusra brings each surprise in as it comes; Mirela reads every change. What the service does this week, customers will rely on.

## Debrief

Friday's 17:00 deploy shipped the whole week. Hollis's field team fetched Tuesday's image through one link, Brinewatch's events reach Brinewatch, and its thirteenth booking got a clear no and a date. Yusra's TL;DR: "Week one: nobody missed the spreadsheet." On Monday, Ingrid has one more thing to ask.
