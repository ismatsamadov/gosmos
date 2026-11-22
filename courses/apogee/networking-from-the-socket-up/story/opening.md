# Spring at Ystmark

kind: text
eyebrow: YSTMARK · MONDAY 6 MARCH 2028
action: Take your new badge

The sun came back to Ystmark three weeks ago. It is −18 °C, and on the ridge above the station's one building, two white domes called radomes keep the wind off its dishes.

The computers that turn the dishes are worn out: Dish 2's lost its position twice in February. New ones came up on your plane, in a crate.

# A New Badge

kind: badge
eyebrow: APOGEE
action: Meet the team

Same number, a new role: until June, you write the software that talks to Ystmark's dishes, from a bench below the ridge.

## Station Software Engineer

label: Role

## No. 039

label: Badge

## Ystmark

label: Based at

## Monday 6 March 2028

label: Since

# Who You Work With

kind: text
eyebrow: THE TEAM
action: Open the crate

The station's one building is the hut: a bench computer, Sunniva's terminal and a coffee machine, below the ridge. Three people share the work on the dishes with you.

## Sunniva Ruud

label: SR
subtitle: Ystmark's station chief

Runs the station and gives the weather first, in as few words as the line allows: "Minus 18, wind light. Crate open. Coffee holding."

## Nkem Halloran

label: NH
subtitle: Staff engineer, your reviewer, at Port Ellery

Reads every change you make from headquarters, one question at a time.

## Tomás Ferreira-Lund

label: TF
subtitle: CTO and radio engineer

Flies up later in March, once Dish 1's new controller is in, to point the dish true.

# What Came in the Crate

kind: text
eyebrow: THE HUT · THE BENCH
action: See your first assignment

Each dish is turned by a dish controller, a computer in its radome. The new controllers speak in raw bytes over a network connection, and Apogee steers them from headquarters at Port Ellery, across the island's line to the mainland.

## The dish controllers

label: YST
subtitle: YST-1 for Dish 1, YST-2 for Dish 2

Dish 1's goes in first, from Monday 20 March.

## The manual

label: DOC
subtitle: The maker's book for the controllers

It is wrong on page one.

## The simulator

label: SIM
subtitle: The maker's program, on a stick

Behaves exactly like a controller, so your code can talk to one before the real ones are installed.

## The island's line

label: LINE
subtitle: Ystmark's link to the mainland

A fiber cable, 40 ms there and back. In storms it fails, and only the backhaul is left: a slow link through a satellite, 0.7 s there and back.

# Your First Assignment

kind: text
eyebrow: THE BENCH · DAY ONE
action: Connect to the simulator

Dish 1's old controller comes out on Monday 20 March. Before then, the program that will talk to the new ones, the dish service called dishd, must talk to the simulator on the bench. It runs on your Mac alone, and every connection stays inside it. Start by connecting and reading what it says first.
