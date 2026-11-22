---
slug: forty-at-once
title: Forty at Once
type: code_tests
position: 2
difficulty: 9
topics: channels, goroutines
created: 2026-11-07T06:23:23Z
modified: 2026-11-07T06:23:23Z
---
Every morning at 08:00 the cleaners switch on all forty screens in the same minute as a refresh. Nkem rehearses it this afternoon, and three screens show the 07:59 Board until the next change.

`serveEvents` reads the current view, writes it, and only then subscribes. A view published between the read and `Subscribe` goes to channels that exist, and this stream's doesn't exist yet. Subscribe first, and nothing published afterwards is missed. Now a view can arrive twice, as the current view and again on the channel, so a stream remembers the version it wrote last and skips any that is no newer.

## Your task

In `serveEvents`:

1. Subscribe, with the deferred `Unsubscribe`, before reading the current view.
2. Keep the version of the last event written, and skip a view from the channel whose version is no newer.

Each of forty screens then reads every version once, in order.
