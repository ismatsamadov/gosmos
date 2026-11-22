Keep the command in flight together with the line you sent, its bytes exactly. On every WELCOME, compare its `last` with that command's number before anything else goes up: at or above it, the command is done and the next one can follow.

The copy goes out on the new call before any later command, with the satellite's lock still held, and its ACK or NAK settles it as the first copy's would have. Mark it as a copy: only then does a NAK `replayed sequence number` mean it ran.

If a PING ran twice, a lost command went again under a new number, or went again though the WELCOME's `last` had reached it. If the dropped call's line came back different, it was signed or sealed afresh instead of sent as it was. A backhaul pass that failed `refused` met the replay and didn't count it as run.
