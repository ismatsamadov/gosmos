The build error names the first thing the other teams can't reach. Which names in `telemetry/telemetry.go` do `main.go` and the tests use?

Capitalize the type, its three fields and the function, then fix every use of the old names inside `telemetry` too, such as the composite literal at the end of `parse`.

Doc comments start with the names they describe, so they change as well. `battery` keeps its lowercase name: only `Parse` calls it, and the last test checks that nothing else is exported.
