---
slug: go-at-scale
title: Go at Scale
kind: Course
position: 8
short: Rebuild Apogee's pass planner for twenty stations and eighty satellites, with Go that stays correct and fast under load.
---
Apogee is growing from five dishes to twenty and from eight satellites to eighty, and the program that decides which dish listens to which satellite must now plan a hundred days instead of one. You lead its rewrite, with Nkem reviewing every change. You build pipelines that stop cleanly, limit how much runs at once, load things once, share state without locks and know why that is safe, walk half a million passes without holding them, write generic queues and merges, fuzz the parsers a customer's satellite can break, test a whole day of timers in a millisecond, measure honestly, cut the garbage, and read struct tags with reflection. The season ends on take-over day, when your planner takes over from the old one on a live day, under the heaviest schedule Apogee has ever flown. Standard library only, and every check runs on your Mac.
