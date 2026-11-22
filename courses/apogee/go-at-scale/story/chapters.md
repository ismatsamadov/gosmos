# stage-by-stage

## Briefing

Monday 5 March, at the backend desks. tlmd reads every station's telemetry in one loop; with twenty stations it will fall behind, and when two stations hear one satellite, the newest frame the Board and the pager read can go backward. Nkem wants it rebuilt stage by stage before Wednesday's drill of twenty feeds. On Thursday, Thistle-2's team watches the Board for the satellite's first beacon.

## Debrief

Thursday, 17:08:39: Thistle-2's first beacon crossed your pipeline to the Board while Joss and Professor Marr watched from the Fishbowl. On Friday the drill ran again on four parsers, every frame in order. Nkem: "Same shape for passplan 2, only bigger."

# twelve-at-a-time

## Briefing

Monday 12 March, at the backend desks. On Friday night passplan 2's first draft started eight thousand jobs at once, one for every satellite and every day, and ran its test machine out of memory. Its new machine in Port Ellery's rack has twelve cores and runs the Board too: a planner that starves it stops the screens. Nkem: "Twelve planner slots. Never more."

## Debrief

Friday 16 March: passplan 2 predicted take-over day's 7,485 passes in twelve slots, the dish loaders took contacts at their own pace, and nothing waited on the Board. Nkem: "Twelve held. Next question: what happens when two planners touch the same pass?"

# once-and-only-once

## Briefing

Monday 19 March, at Ingrid's desk in the Fishbowl. At 08:00 the Board froze for six minutes while a hundred screens each loaded the table of stations and satellites for themselves. Beside it, a test Board shows passplan 2's plans, which become the real ones on take-over day, so it must never show anything but the newest. Ingrid: "One load. Newest plan. Every time."

## Debrief

Friday 23 March: the Board loads its table once for a hundred screens, and the test Board takes a plan only from a newer run, in whatever order the runs finish. Ingrid added a line to take-over day's checklist: "Board shows newest plan only."

# passes-on-demand

## Briefing

Monday 26 March, at Yusra's desk. Her TL;DR: customers want to book a whole season ahead, and Hollis agrees in capitals: "WE PLAN A SEASON NOT A WEEK." bookingd's catalog offers fourteen days; for the third quarter it would hold 672,292 passes, and the old passplan hands a copy of them all to every service that asks, about eighty-four megabytes each. A season has to be made as it's read.

## Debrief

Friday 30 March: bookingd's catalog now predicts each day only when someone reads that far, and the archive server closes every file it opens. Hollis: "THIRD QUARTER BOOKED. ALL TWENTY-SIX FURROWS."

# the-iterator-toolkit

## Briefing

Monday 2 April, in the Fishbowl. At 00:00 the first wave went live: Tirrin Rock, Varsk Island, Calidra Atoll, Ilha Brumosa and Ostrel Island. Sunniva: "Tirrin up. Two dishes. Ice on both." From today passplan 2 plans in shadow beside the old planner, and the two schedules must agree contact by contact before it flies a pass. The Board wants one station's contacts, bookingd one customer's, Ingrid the next twelve.

## Debrief

Friday 6 April: passplan 2 now hears when a station's dishes are parked, as the old planner does, and the shadow run compared Thursday again: 273 contacts each, not one difference. Nkem: "They agree. Now, what does it sort before it plans?"

# next-pass-first

## Briefing

Monday 9 April, at the backend desks. The shadow run sorts every pass of a hundred days into one list per station before planning the first, and two passes that rise in the same second trade places from one run to the next. Nkem wants the eighty satellites' passes merged a pass at a time, next pass first, before Cirrine's twelve new satellites leave their rocket on Tuesday 17 April.

## Debrief

Friday 13 April: the shadow run plans as the passes come instead of sorting a hundred days first, and two passes in the same second no longer trade places. Tirrin Rock's Tuesday: 120 contacts. Nkem: "Next pass first. Now the twelve can come."

# fuzzing-the-parsers

## Briefing

Monday 16 April, at Mirela's desk. On Saturday, Thistle-2's student software sent a battery reading that wasn't a number, tlmd let it through, and the pager stayed quiet; at 05:10 a line with fields missing crashed the parser and stopped every station's telemetry for forty minutes. Joss apologizes. On Tuesday, Cirrine's twelve send their first telemetry through the same parser.

## Debrief

Tuesday 17 April, 09:41: Cirrine-3's first line came in through Varsk Island's new dish, a minute out of its rocket, and your parser passed it straight to the Board. Mirela: "LGTM once it fails closed. It does."

# a-day-in-a-millisecond

## Briefing

Monday 23 April, at the backend desks. passplan 2 keeps time four ways: a replan every ten minutes, a limit on every planning run, a wait for each contact's AOS and a pause for news to settle. Its tests wait for real: the suite takes eleven minutes and fails one run in five. The second wave goes live on Monday 7 May; Nkem wants tests that can't lie about time.

## Debrief

Friday 4 May: the timer suite that took eleven minutes now runs in under a second, and in it the 2,020 contacts of take-over day's plan fly in a few milliseconds, each at its AOS. Ingrid: "Copy. Every AOS on its second."

# benchmarks-that-tell-the-truth

## Briefing

Monday 7 May, at the Fishbowl's glass. The second wave went live at midnight, and Tomás has written beside his drawing: "passplan 2: forty times faster!" Nkem has written under it: "Faster at what?" Nobody can say how the forty was measured, and on take-over day the team will trust whatever the glass says. Nkem wants it measured again, honestly, by Friday 18 May.

## Debrief

Friday 18 May, your table on the glass: on one core, passplan 2 plans 18 May in about the old planner's time, with two to three times its allocations. Tomás wiped the forty off himself. Nkem: "Not faster. Heavier. That we can fix."

# fewer-allocations

## Briefing

Monday 21 May, at the backend desks. Friday's table showed passplan 2 making two to three times the old planner's allocations. With fifty-six satellites, every replan of a hundred days fills the planner's memory with garbage, and the test Board's updates stutter while it is cleared away. The third wave goes live on 4 June, and Nkem wants the planner's garbage halved before then.

## Debrief

Friday 25 May: planning take-over day at twenty stations, with its host reports and dish loader lines, now leaves 6 MB of garbage instead of 12, and most of what is left comes from predicting passes. Nkem: "Nearly half. The rest is prediction."

# reading-struct-tags

## Briefing

Monday 28 May, at the backend desks. The third wave's five hosts are sending in their forms: where each station stands, how many dishes it has and which radio bands it uses. A wrong form takes a station's passes away, and the wave goes live on Monday 4 June. Nkem wants one checker for every form and for passplan 2's own settings, with each field's rules written beside it.

## Debrief

Friday 8 June: Vorrel Island has flown five days on both its dishes, on a form your validator checked, and passplan 2 won't start on settings its own tags refuse. Nkem: "Twenty stations, every form right. Now Ingrid's day."

# a-hundred-pass-days

## Briefing

Monday 11 June, in the Fishbowl. Ingrid has set take-over day for Thursday 28 June: at 00:00 UTC passplan 2 takes over from the old planner on the heaviest day Apogee has ever flown, hours after Furrowlight's twenty-four new satellites are first heard. She will run it like a pass, with a full rehearsal in shadow on Monday 25 June. Nkem wants the planner's last pieces fitted before then.

## Debrief

Friday 29 June, 00:00, Ingrid's log: "passplan 2: GO for every day." Your planner has flown take-over day's 2,011 contacts through Cape Wren's wind, no dish touched by both planners. Hollis, through Yusra: "TWENTY-FOUR FURROWS. EVERY ONE HEARD." Nkem: "Answered. July: the whole ground segment."
