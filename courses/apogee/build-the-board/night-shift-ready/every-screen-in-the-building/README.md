---
slug: every-screen-in-the-building
title: Every Screen in the Building
type: code_tests
position: 3
difficulty: 9
topics: http, security, templates
created: 2026-11-07T06:23:23Z
modified: 2026-11-07T06:23:23Z
---
The night into Friday: forty screens, Cape Wren's wind rising, a deploy at 05:00, and Hollis at reception at 09:00, where no visitor may read another customer's notes.

The reception screen asks for `?screen=lobby` on `/`, `/board.json` and `/events`: the same view, with every booking's ID and note blank and no on-call name. Its ETag is `"v3-lobby"`, so a lobby copy is never taken for the Fishbowl's. `testdata/night.txt` then plays the whole night against golden files of what each answer must be.

## Your task

1. In `board.go`, write `isLobby(r)`, true for `?screen=lobby`, and `forLobby`, a method of `View` returning a copy with every pass's `ID` and `Note` blank and `OnCall` empty.
2. With `?screen=lobby`, `servePage`, `serveJSON` and `serveEvents` show `forLobby()`. `pageData` gets `Lobby bool`, and `board.html` leaves out the on-call line when it is set.
3. The lobby's ETag is the version and `-lobby`, in quotes.

A lobby row ends with two empty cells:

```text
<td>pipit-1</td><td></td><td></td></tr>
```
