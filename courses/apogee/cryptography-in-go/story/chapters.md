# numbers-nobody-can-guess

## Briefing

Week one, at Mirela's desk in Port Ellery: "Every secret this season starts as random bytes. Show me where ours come from." The draft nonces for Cirrine's commands and the ops console's session tokens both come from the clock, so anyone who knows roughly when one was made can make it again. On Friday at 16:00, Mirela makes the season's master key with your code.

## Debrief

Friday, 16:00: Mirela makes `k-2028-07` with your helper, 32 bytes nobody can guess, and locks its printout in the safe. "LGTM: nothing in it reads a clock," she writes under your last review. On Monday, Yusra brings Cirrine's second condition: proof for every file.

# fingerprints-for-files

## Briefing

Week two, at the backend desks, where delivered files land. Yusra: "TL;DR: Cirrine wants proof that every file is exactly what left the satellite." The check number each delivery carries, a CRC, catches radio noise, not a person: whoever changes bytes on purpose can make it match again. On Friday, Cirrine's first delivery with a fingerprint for every file goes out, and its archive will check each one.

## Debrief

Friday 21 July: Cirrine's first delivery with fingerprints lands, and every line of Yusra's report reads `ok`. Her note to Cirrine is one line: "TL;DR: every file now arrives with a fingerprint your archive can check." Next week, Mirela reads how the ground signs its commands.

# signed-by-the-ground

## Briefing

Week three, at Mirela's desk. Before Cirrine's commands are locked, Mirela reviews how Apogee signs them and finds three weak spots: a homemade signature in Pipit-3's software, a bar that lets one command pass for another, and a replay guard that counts a command before checking it. Each could let a forged command through or lock a satellite out. On Friday, twelve planted lines test the fix.

## Debrief

Friday 28 July: of Mirela's twelve planted lines, the checker accepts her two good ones and refuses the other ten, and not one refusal moves the guard. Under the review she writes: "LGTM once it fails closed. It does."

# passwords-at-rest

## Briefing

Week four, in the Fishbowl. On Sunday the spare laptop vanished from the night-shift desk with a copy of the ops console's user file: 26 accounts, each password stored unsalted. Mirela: "Assume whoever has it has a fast computer and a list of common passwords." By Friday every password is changed, and from 16:00, when the login reopens, every record the console saves makes each guess cost 600,000 rounds.

## Debrief

Friday, 16:00: the console's login reopens. Mirela tries the thousand common passwords on a copy of the new user file, and not one record gives way. Ingrid closes the laptop's entry in the log: "Lost: one laptop; given away: nothing."

# keys-from-keys

## Briefing

Week five, at the lab bench in Port Ellery. One key signs commands for all seven satellites, and Mirela finds it in four places: uplinkd, a spare Pipit-2 computer on the bench and two old test scripts. Whoever reads it out of any of them commands the whole fleet. By Thursday's passes each satellite gets keys of its own; on Friday at 00:00 UTC the old key retires.

## Debrief

Thursday's passes carry every satellite's own keys up, and at 00:00 UTC on Friday `k-2027-03` retires. Mirela wipes it from the spare Pipit-2 on the bench: "One satellite's key now opens one satellite." On Monday, Tomás wants to show you what the test van's dish can hear.

# sealed-messages

## Briefing

Weeks six and seven, at Port Ellery. Tomás parks the test van below the roof's dish, points the van's own at it and reads Pipit-2's morning commands off the air, word for word. From Friday 25 August, Cirrine's contract says nobody but Cirrine and Apogee may read Cirrine-1's. Mirela wants them sealed before its 07:30 pass over Ystmark that Friday.

## Debrief

Friday 25 August, 07:30 UTC: Cirrine-1's first sealed command, `PING`, goes up from Ystmark, and Cirrine-1 answers. A dish listening there would read the header and not one word of the command. Yusra: "TL;DR: Cirrine's first condition is met." On Monday, Hollis writes in capitals.

# signatures-anyone-can-check

## Briefing

Week of 28 August, at Yusra's desk. Monday, 09:02, Hollis: "FURROW-A 0827 1402 NEVER ARRIVED." Apogee's delivery program says Sunday's file landed at 14:13:05; Furrowlight's receiver says nothing. The signed message Apogee sent can't settle it, because Furrowlight holds the same secret. Yusra: "TL;DR: receipts nobody but Apogee can sign and anyone can check." Receipts start on Friday.

## Debrief

Friday 1 September: Furrowlight checks Sunday's receipt against the file in its own storage, and it matches. The file had arrived; their receiver had missed the message. Hollis: "RECEIPT CHECKS OUT. THANKS." Sunniva, from Ystmark: "Storm Sunday. Fiber at risk."

# agreeing-on-a-secret

## Briefing

Week of 4 September, over Ystmark's backhaul. On Sunday a storm cut the station's fiber, and everything now crosses a slow satellite link that another company runs. Sunniva: "Fiber cut. Backhaul up. Everything on it readable." This week's command file, and next month Rotation Day's keys, must cross a link whose owner can read and change anything.

## Debrief

Thursday night, the week's command file crosses the backhaul sealed under a key Port Ellery and Ystmark agreed in the open; the middle Mirela played got nothing. Sunniva: "File in. Three lines. Nobody else read them." On Monday, Mirela asks who else can reach a station.

# who-are-you-server

## Briefing

Week of 11 September, Mirela's audit: anything that reaches a station's network can talk to its uplinkd and its dish controllers, and nothing asks who is calling. On Tuesday Paz finds an old test script at Salar Quieto still asking Chiquita's controller for its position every minute: harmless, and nobody knew. By Friday 22 September every station and server must carry a certificate from Apogee's own authority.

## Debrief

Friday 22 September: eight certificates issued, every key file readable by its owner alone, and the root untouched in the safe. Mirela: "LGTM. Every station can now say who it is. Next week, it proves it on every connection."

# every-connection-encrypted

## Briefing

Monday 25 September, in the Fishbowl. Tomás shows that uplinkd's traffic to Lantern Atoll crosses another company's network in plain text, readable by anyone along the way. Mirela's answer: every connection encrypted, both ends proven with Apogee's certificates. Rotation Day's new keys will travel these links in three weeks, so by Friday 6 October every station must refuse anyone else.

## Debrief

Friday 6 October: every station speaks TLS 1.3 and proves itself, uplinkd too, and a sealed Cirrine-1 command crossed to Ystmark byte for byte, unread by the station that queued it. Sunniva: "Snow 35 cm. Link encrypted. Both ends proven. Coffee holding." Next week, Mirela starts rehearsing Rotation Day.

# brinewatchs-older-keys

## Briefing

Friday 6 October, at Yusra's desk. Brinewatch, the ship-tracking customer, can install Rotation Day's new secrets only in its monthly maintenance window, tomorrow from 06:00 UTC, and only the way its older systems allow. Mirela suspects the pad its secrets have traveled on since November. Miss the window, and on Rotation Day Brinewatch is either locked out of the booking service or keeps one old secret.

## Debrief

Saturday, 06:00 UTC: Yusra runs your handover in Brinewatch's window, and both new secrets open on Brinewatch's side under their own labels, readable by Brinewatch alone. Mirela wipes Apogee's copy of the pad and signs off: "LGTM. It failed closed on every wrong key I sent it."

# rotation-day

## Briefing

Week of 9 October, the Fishbowl. On Tuesday 17 October every key in the network is replaced, at every station, without missing a single pass. Ingrid will run the day like a pass, from the Cirrine-1 pass that crosses midnight to Furrow-A over Salar Quieto at 14:02. Before then, every satellite needs its new key, and Mirela wants every switch rehearsed.

## Debrief

Ten passes flown, none missed, and certificate 1005 refused from 10:12. Mirela signs the log: "Every key replaced. Every pass flown. It failed closed every time we asked it to." In November you carry the pager again, and Lantern Atoll's generator has started to sound different.

# more-than-a-password

## Briefing

Wednesday 25 October, the day after the old key retired. Cirrine's yearly review passes Rotation Day with one question, which Yusra forwards: could someone holding an engineer's console password stow a dish at Ystmark during Cirrine-1's pass? Today, yes. Mirela: "Friday 3 November, 16:00. From then, a password alone opens nothing." Before then: the reset links, a code on every phone, and a way in for a lost one.

## Debrief

Friday 3 November, 16:00: the console asks every engineer for a password and a code. Mirela tries her password without her phone, and is refused: "LGTM. It failed closed." Yusra to Cirrine: "TL;DR: a stolen password alone now opens nothing." On Monday at 08:00, Mirela hands you the pager.
