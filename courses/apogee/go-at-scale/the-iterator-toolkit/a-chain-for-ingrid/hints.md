`slices.Values(firstWaveDay().Contacts)` is Monday's schedule as a sequence. Hand it to `Filter` with `furrowAtNew`, which keeps the Furrows' contacts at Tirrin Rock, Varsk Island, Calidra Atoll, Ilha Brumosa and Ostrel Island.

Give what `Filter` returns to `Map` with `line`, and what `Map` returns to `Take` with 5: the chain reads from the inside out.

`strings.Join` needs a slice, so `slices.Collect` the chain's lines, join them with `"\n"` and print the result. If all 260 lines still print, `main` still loops over the whole schedule.
