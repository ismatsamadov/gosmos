# Apogee story bible: the gctl project

### The gctl project: the night shift's shell

The van comes home to Port Ellery on Saturday 23 October 2027, and by Monday every engineer
wants gctl on their own Mac. Ingrid asks for one more thing: an interactive mode, so the
night shift can type commands during a pass without starting the tool again each time. The
guided project **Build gctl** runs from Monday 25 to Friday 29 October 2027 at Port Ellery:
one chapter a day or so, ending on the night of Friday 29 October, when the night shift runs
a whole pass from your prompt. Nkem reviews; Ingrid is the customer; Cape Wren's wind supplies
the night's trouble. The badge is unchanged (Ground Tools Engineer, No. 039).

`gctl shell` is one `package main` program that grows lesson by lesson, each lesson's
starter being the previous lesson's solution:

- `shell.go`: `Shell`, its loop (`Run(in io.Reader, out io.Writer) error`) and dispatch;
  `commands.go`: the command table, a `map[string]command` from a name to its help text and
  its function; `cache.go`: the pass cache; `settings.go` and `history.go` in chapter 4;
  `data.go` (read-only): the night's passes and the station table; `main.go` wires them up.
- The prompt is `gctl> ` (no newline), written before each line is read. A line is split
  into words with `strings.Fields`: the command, then its arguments. A blank line does
  nothing. An unknown command prints `unknown command "<name>"; type help`. A command's error
  prints `error: <message>` and the session goes on; a panicking command prints
  `error: internal error in <name>` and the session goes on. `exit` or the end of the input
  prints `bye` and ends the session; `Run` returns only a failure of the input itself.
- Commands: `help` (every command and its help text, sorted by name), `exit`,
  `passes [CODE]` (the night's passes, all or one station's, one per line as
  `04:12-04:21 PEL pipit-1 max 61.4`, ordered by AOS, ending in ` (stowed)` at a stowed
  station), `stations` (`PEL Port Ellery`, and
  `stowed: <reason>` after a stowed one), `stow CODE [reason...]`, `next` (the next pass
  after now at a station that isn't stowed: `next: pipit-1 over PEL at 04:12, in 18m0s`, or
  `no more passes tonight`), `history` (every command typed this session and before,
  numbered from 1).
- The pass cache loads the night's passes through a loader function on first use, keeps
  them for a time to live (10 minutes unless the settings say otherwise), loads again when
  they are stale, and is safe for a background refresher that runs beside the commands. Its
  clock is injected (`now func() time.Time`) so tests move time without waiting.
- Settings: `gctl.json`, read strictly (Season 2's rules): `prompt`, `cache_ttl` as duration
  text, and `station`, the default for `passes`. A missing file means the defaults; a broken
  one is an error. History is appended to a file named by the program, one line per command,
  and read back at the next start.
- The night of Friday 29 October 2027 runs past midnight, so its passes are on Saturday
  30 October, UTC. Their lengths fit each orbit, mask and maximum elevation:

  | Satellite | Station | AOS–LOS | Max elevation |
  |---|---|---|---|
  | Pipit-1 | PEL | 04:12–04:21 | 61.4° (as on your first Board) |
  | Pipit-2 | YST | 04:31–04:40 | 47.2° |
  | Furrow-A | CWR | 04:58–05:06 | 33.5° |
  | Brine-3 | LAN | 05:20–05:26 | 22.8° |
  | Cirrine-1 | YST | 05:47–05:57 | 71.9° |

  At 04:40 Cape Wren's wind reaches 74 km/h and its dish stows itself (it stows from
  72 km/h); the night shift records it with `stow CWR wind 74 km/h`, and `next` must skip
  Furrow-A's pass. The relief shift takes over at 06:00, and Ingrid writes one line in the
  log: "gctl shell: GO for every shift." On Monday Yusra's message opens Season 3: "TL;DR:
  customers want to book their own passes."
- Earlier in the week, on Wednesday, a storm cuts Ystmark's fiber and Sunniva runs gctl over
  the slow satellite backhaul ("Storm. Fiber down. Nine seconds a command."), which is why the
  shell keeps its passes in a cache.
- Exact texts: `stow` confirms `CWR stowed: wind 74 km/h`; `history` prints `%3d  %s` lines;
  errors read `unknown station "XYZ"`, `passes takes at most one station code`,
  `stations takes no arguments`, `stow needs a station code`, `CWR is already stowed`,
  `next takes no arguments`, `history takes no arguments`, `internal error in <name>` and
  `history: <err>`; a broken settings file reads `gctl.json: more data after the settings`,
  `gctl.json: cache_ttl "10" is not a positive duration` or
  `gctl.json: cache_ttl 600 is not duration text`. The settings file is `gctl.json` and the
  history file `gctl_history`, both beside the program.
