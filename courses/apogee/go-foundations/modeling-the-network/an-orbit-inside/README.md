---
slug: an-orbit-inside
title: An Orbit Inside a Satellite
type: code_tests
position: 2
difficulty: 4
topics: structs
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
The planner needs to know which satellites pass near Earth's poles, following polar orbits. Ystmark sees them on almost every trip around Earth, or revolution. Tomás sketches a satellite in his notebook: an ID, a name, and an orbit with its own altitude and inclination.

A struct field can hold another struct. Build it with a literal inside a literal, and reach in with one dot per level:

```go
type Wind struct{ GustKmh float64 }
type Site struct {
	Code string
	Wind Wind
}

cwr := Site{Code: "CWR", Wind: Wind{GustKmh: 81}}
fmt.Println(cwr.Wind.GustKmh) // 81
```

## Your task

Make `isPolar` report whether a satellite's orbit is inclined from 80 to 100 degrees, both included.

> [!NOTE]
> The tilt of an orbit to the equator is its inclination. It caps how far north or south a satellite flies.
