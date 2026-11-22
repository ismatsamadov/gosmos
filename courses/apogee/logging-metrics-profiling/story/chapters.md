# logs-a-machine-can-read

## Briefing

Your second night on call, at 02:46 on Wednesday, the pager woke you: one pass in five at Lantern Atoll was failing, and Cirrine-1's command never went up. The station's log said `error`, three times; a restart at 03:14 fixed it, and nobody knows why. Nkem: "`error` is not a log line." By Friday at 14:00, the station's lines must say what failed, where, and on which pass.

## Debrief

Friday, 14:00: Lantern Atoll's uplinkd writes lines that name the pass, its elevation, its wait and its whole error, in JSON the night shift can search. Ingrid adds a line to her checklist, "until we know why": "Monday 08:00: restart LAN uplinkd."

# levels-and-groups

## Briefing

Monday 13 November, in the Fishbowl. Since Friday's deploy, Lantern Atoll's uplinkd writes down every detail, to catch the next failure: 2.1 million lines a night, and the ops console gives up before the night shift finds a single pass. Ingrid wants a night she can read at a glance, everything when someone asks, and no restart to change it.

## Debrief

Friday 17 November, 14:00: the night shift's view shows where each pass starts and how it ends, and can turn a station up or down without a restart. Ingrid: "Night view, two lines a pass. LAN's file keeps everything at DEBUG, until we know why."

# logs-with-context

## Briefing

Monday 20 November, at the backend desks. You and Nkem spend the whole day lining up what four programs wrote about Cirrine-1's failed pass of 8 November: the booking service, the planner passplan, and uplinkd at headquarters and at Lantern Atoll. Next time a pass fails at three in the morning, the night shift must find its whole story in one search, in a minute.

## Debrief

Friday 24 November: one search finds the pass's seven lines in all four programs, in seconds. Cirrine-1's command reached Lantern Atoll at 01:48:30 and never left: the pass died inside the station's uplinkd, before the command reached the radio. Mirela asks for a week of bookingd's new lines.

# never-in-a-log

## Briefing

Monday 27 November, Mirela's desk. Reading last Wednesday's bookingd lines as anyone with the ops console could, Mirela finds Furrowlight's API key in 14 of them: whoever read it could have booked passes as Furrowlight. Yusra and Hollis replace the key by 16:00 (Hollis: "NEW KEY. WHY."). On Friday at 16:00, Mirela plants twelve lines full of secrets, and no secret may reach a disk.

## Debrief

Friday, 16:00: all twelve of Mirela's planted lines reach the disk with every secret blanked and every address cut to one letter and its domain. Mirela signs off: "LGTM once it fails closed." At Lantern Atoll, uplinkd's single log file, still at DEBUG, keeps growing.

# logs-that-fit-the-disk

## Briefing

Monday 4 December, 08:10, ten minutes into your week on call: the pager goes off, and again every ten minutes. Lantern Atoll's disk is 98% full, 41 GB of it one log file: since 13 November uplinkd has started at DEBUG, "until we know why." A full disk leaves the recorder nowhere to keep the next pass, as Ystmark's did in May. By Friday, every station's logs need a budget.

## Debrief

Friday, 14:00: the new log files reach every station. At Lantern Atoll the first line rotates the 41 GB file out and the budget deletes it: its logs now hold at 1 GiB, its disk about two thirds full. On Monday, Yusra brings a question no log line can answer.

# counting-what-matters

## Briefing

Monday 11 December, at Yusra's desk. Brinewatch wants the exact number of its failed passes in November, and every week from now on. Ingrid's log says 5, uplinkd's lines 3, passplan's 6 and bookingd's records 4. Logs are not counts: by Tuesday at 10:00 every station must count its own passes, and Brinewatch's first exact week goes out on Tuesday 19 December.

## Debrief

Tuesday 19 December, 10:00: Brinewatch's first exact week goes out, 62 passes booked, 61 flown and 1 failed, in Wednesday's downpour at Lantern Atoll. Yusra: "TL;DR: one number, and nobody argued." Tomás takes the pager for Christmas week.

# when-to-wake-someone

## Briefing

Tuesday 2 January, the Fishbowl. On the night of 27–28 December Lantern Atoll's passes failed from evening to dawn, and the pager woke Tomás eleven times for one problem; at 07:28 he muted the station, with no end. Ingrid wants new rules live on Monday 8 January at 08:00: one page per problem, to the person on call, and silences that end.

## Debrief

Friday 5 January: replayed through the new rules, the night pages Tomás once, at 20:21, and tells him at 08:45 that it's over, while six single failures wait for the morning list. Tomás: "One bell for one buoy."

# error-budgets

## Briefing

Monday 8 January, the Fishbowl. Lantern Atoll lost 22 of its 1,199 passes last quarter, and nobody can say whether that broke Apogee's promise, because "every pass counts" was never written as a number. Ingrid, Tomás and Yusra agree the numbers this morning, and on Friday Yusra sends every customer the quarter's report. The generator's new part arrives by boat on Thursday.

## Debrief

Friday afternoon, Yusra's report reaches every customer, with Lantern Atoll's row eleven passes over its budget. On Thursday the generator got its new part, fitted from 09:00 to 11:00 with the station off the air. Tomás: "A new injector and steady power: now watch that burn rate fall!"

# where-the-time-goes

## Briefing

Tuesday 16 January, 05:10: wind at 74 km/h stows Cape Wren's dish, and passplan takes 3 min 52 s to plan the rest of the day. Until 05:14 the Board still shows Brine-3's 05:29 pass at Cape Wren as booked, and the night shift starts its checklist for a pass that can't fly. Ingrid: "Four minutes is most of a checklist." By Friday, a day must plan in seconds.

## Debrief

Friday 19 January, 14:00: passplan plans a day in 6.1 s, the same plan as before. On Monday Ingrid struck a line from her checklist: "Generator fixed. No more Monday restarts."

# where-the-memory-goes

## Briefing

Monday 5 February. Lantern Atoll is failing passes again, though its generator was fixed in January, and Ingrid has left a note on the Board: "LAN failing: call No. 039 before anyone restarts anything." Every restart since October has wiped out whatever was going wrong. This time, before any restart, uplinkd's memory is saved and read.

## Debrief

Friday 9 February, 10:00: at every station, a failed pass now says why. Sunniva: "Ystmark: one stuck. Wind 40. Coffee holding." Paz: "Salar Quieto, zero stuck. Even the goroutines are short of breath up here."

# following-one-request

## Briefing

Monday 12 February. The incident review on Friday 23 February at 15:00 must explain everything that went wrong at Lantern Atoll since 17 October, and Nkem wants it in one picture. On Tuesday the lab bench replays the night of 6–7 February on the old program, recording every step; from Monday 19 February every station's uplinkd records its own.

## Debrief

Friday 23 February, 15:00: one picture explains three months. Since the fix, Lantern Atoll's uplinkd has run fourteen days without a restart, its goroutines at 14 every hour, and Ingrid strikes the runbook's last line, "restart uplinkd". From March: twenty stations, eighty satellites, and passplan's rewrite.
