# Apogee: the story bible

Every Gosmos course in the Apogee series follows this bible. It fixes the world (company,
stations, satellites, people), the code the learner grows, and how lessons are written. When
a lesson needs something this file doesn't have, add it here first, then use it. How the
story is told, from the opening to the field guide, is `docs/storytelling.md`.

Everything in this world is invented: the company, its customers, the people, the places and
the satellites. Only the physics is real.

## The parts

The bible is kept in parts, one file each in `story/`, so that whoever writes or reviews a
chapter reads only what it needs: the world, the style guide and the part for its season or
project. Go Foundations, Season 1, has no part of its own: the world's part holds all it
needs, the codebase included. Something new goes into the part it belongs to.

- [`story/world.md`](story/world.md): the company, the ground stations, the fleet and the
  people; the codebase every season shares, from the `apogee/groundctl` module and its
  canonical types to the wire formats and the limits; the seasons, the recurring situations
  and the glossary.
- [`story/season-2.md`](story/season-2.md): "Season 2: tools for the van" (Files, Data and
  the Command Line, July to October 2027): the road chapter by chapter, and Season 2's files.
- [`story/gctl-project.md`](story/gctl-project.md): "The gctl project: the night shift's
  shell" (Build gctl, October 2027): the program's files, commands, pass cache, settings and
  history, the night of Friday 29 October and its passes, and every exact text.
- [`story/season-3.md`](story/season-3.md): "Season 3: open for business" (Web Services in
  Go, November 2027 to February 2028): the season chapter by chapter, the service's types,
  the API, keys and tokens, rate limits, the store, webhooks, the account page, the weather
  service, and how lessons run servers and test them.
- [`story/board-project.md`](story/board-project.md): "The Board project: every screen in
  the building" (Build the Board, February 2028): the program, what it serves, the night
  into Friday 11 February, and how the lessons test it.
- [`story/season-4.md`](story/season-4.md): "Season 4: down to the wire" (Networking from
  the Socket Up, March to June 2028): Ystmark in spring, the controller protocol, the
  simulator, dishd's code, the watch port, the station server, the live port, and how lessons
  test the wire.
- [`story/season-5.md`](story/season-5.md): "Season 5: keeping secrets" (Cryptography in Go,
  July to November 2028): every key, its size and its encoding, the test keys and fixed
  vectors, a section for each chapter from random values to Rotation Day and more than a
  password, and how lessons test it.
- [`story/season-6.md`](story/season-6.md): "Season 6: on call" (Logging, Metrics and
  Profiling, November 2028 to February 2029): Lantern Atoll, the leak, the 02:46 page and the
  nights after it, then on call, logs, metrics, alerts, error budgets, passplan, traces and
  spans, and how lessons test it.
- [`story/season-7.md`](story/season-7.md): "Season 7: scale" (Go at Scale, March to June
  2029): the network grown to twenty stations and eighty satellites, the pass model, the
  planning rule, passplan 2's settings, the code chapter by chapter, take-over day, and how
  lessons test it.
- [`story/capstone.md`](story/capstone.md): "The capstone: Pass Day, Live" (July 2029): the
  rehearsal's world, folder and settings, the pass archive and the plan, bookings, the pass
  line, telemetry, commands, files, deliveries, logs, numbers and pages, the investors' week
  as Ingrid scripts it, how the lessons test it, and the reference program for writers.
- [`story/style.md`](story/style.md): the style guide (voice, lessons, quizzes, mix and
  rhythm, code style, starter code and solutions, tests, expected output, and what the app
  does that writers must know) and the naming rules.
