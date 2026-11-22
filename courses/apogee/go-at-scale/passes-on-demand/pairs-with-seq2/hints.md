Range over `sch.Contacts`: they are in the schedule's order already, by AOS then pass ID. Each `Contact` carries its pass, with its `Station`, and its `Dish`.

Skip the contacts at other stations, and call `yield` with the dish first and the contact second for the others.

Return as soon as `yield` returns false. The test also ranges twice over one `At("TIR")` and wants the same 48 pairs both times, so keep no count or position outside the function you return.
