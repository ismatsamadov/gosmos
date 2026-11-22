`r.URL.Query().Get("band")` returns the band as text, or an empty string when there is none. Convert it with `Band(...)` so it compares with `UHF`, `SBand` and `XBand`.

Deal with the refusals before anything is written: an unknown band answers through `http.Error` and returns at once. Only then build the list, starting from `[]Station{}` and keeping each station whose `Bands` contain the band (`slices.Contains` checks that).

If `band=x` gets a 200, the list was written before the refusal: once `writeJSON` has run, the status is 200 for good, so the check has to come first.
