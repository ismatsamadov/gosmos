In `forLobby`, make a new slice for the passes and copy each pass into it with `ID` and `Note` set to `""`. The view's own slice belongs to the Board: blanking it in place would blank the Fishbowl's screen too.

In `servePage`, build the `pageData` first, with `Lobby: isLobby(r)`, and swap its `View` for `forLobby()` when it is set. In `board.html`, wrap the on-call line in `{{if not .Lobby}}` and `{{end}}`, all on that one line.

In `serveJSON`, for the lobby, marshal `v.forLobby()` and build the tag with `-lobby` before the closing quote. In `serveEvents`, note the lobby once, and pass each view through `forLobby()` before `writeEvent`; the copy keeps the version, so ids stay the same.
