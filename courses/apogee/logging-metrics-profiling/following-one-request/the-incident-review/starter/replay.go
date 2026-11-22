package main

import "time"

// The night's three passes at Lantern Atoll, as the bench replayed them.
const (
	cirrine = "cirrine-1/LAN/17759"
	furrowB = "furrow-b/LAN/22379"
	furrowA = "furrow-a/LAN/22382"
)

// at returns hh:mm:ss UTC on Wednesday 7 February 2029, the night the bench replayed.
func at(hh, mm, ss int) time.Time {
	return time.Date(2029, 2, 7, hh, mm, ss, 0, time.UTC)
}

// span returns one finished span of the replay.
func span(traceID, spanID, parentID, name, passID string, start, end time.Time) Span {
	return Span{TraceID: traceID, SpanID: spanID, ParentID: parentID, Name: name, PassID: passID, Start: start, End: end}
}

// benchReplay holds the spans of Tuesday's bench replay, in the order they ended: the night of
// 6–7 February on the old uplinkd, three of its four workers held from the start, the restart at
// 02:10 included. Each worker span runs from pass prep until the pass had a worker or gave up.
var benchReplay = []Span{
	span("c4f5a58656eef510", "c8d04605", "9dcec3ad", "worker", cirrine, at(1, 16, 8), at(1, 16, 8)),
	span("c4f5a58656eef510", "9dcec3ad", "", "pass", cirrine, at(1, 16, 8), at(1, 27, 13)),
	span("c4f5a58656eef510", "cbedc0dc", "9dcec3ad", "fly", cirrine, at(1, 16, 8), at(1, 33, 18)),
	span("3bf98798cae97950", "7ffa3780", "0a8c6d7f", "worker", furrowB, at(1, 36, 52), at(1, 45, 23)),
	span("3bf98798cae97950", "0a8c6d7f", "", "pass", furrowB, at(1, 36, 52), at(1, 45, 23)),
	span("730ad2626b1c2f8e", "2f0901a1", "21ff2330", "worker", furrowA, at(2, 22, 30), at(2, 22, 30)),
	span("730ad2626b1c2f8e", "a9735a3c", "21ff2330", "fly", furrowA, at(2, 22, 30), at(2, 25, 42)),
	span("730ad2626b1c2f8e", "21ff2330", "", "pass", furrowA, at(2, 22, 30), at(2, 25, 42)),
}
