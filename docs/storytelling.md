# Storytelling

People finish what they care about. Gosmos teaches Go inside a story so that a learner cares
from the first minute: you are the newest engineer at a small space company, every lesson is a
task from your job, and the people and satellites in it come to matter to you. This page says
what that story must do for the learner, how a course tells it, and how the app shows it. Every
course follows it; a series' story bible (`docs/apogee/story.md` for Apogee) supplies the world.

Everything is original. Other courses and platforms can show how much ground a path covers or
how a story carries a learner through it; their lessons, worlds, names and words are never read
for writing, mirrored or adapted.

## What the learner lives through

Each of these is a promise. A lesson, chapter or course that breaks one gets rewritten.

1. **"I work here."** Within a minute of opening the app you have a badge with your name, a
   role, a desk at headquarters and a buddy who reviews your code. The lessons never name you;
   they say "you", and the app's own screens use the name on your badge.
2. **"My code matters today."** Every task has a reason right now: a satellite rises in ten
   minutes, a customer is waiting, a storm is coming. When you pass, the world moves: the
   ops room's screen shows your work, a satellite answers, a customer gets their data.
3. **"I know these people."** A small cast with distinct voices returns all season. They ask,
   explain, worry and celebrate; they never lecture.
4. **"I understand every word."** Plain words come first and space words second. Any word of
   the world is one click from its meaning, everywhere it appears.
5. **"I want to know what happens next."** A chapter opens with a briefing and closes with a
   debrief. After every pass the app shows the next lesson's first line. A season builds to
   moments worth reaching: a launch, a storm, a satellite's last pass, your first week on call.
6. **"I'm getting good, and I can see it."** Your code grows into one real program. Your badge
   collects a patch for every chapter, your rank rises with your XP, and the tasks grow from
   printing a line to running a whole day of passes at once.
7. **"This is the only place I need."** The path goes from a first `Println` to Go a team
   would trust in production: the language, its standard library, testing, tools, services,
   protocols, security, observability and performance, all on this Mac, offline.

## The shape of a story

| Level | In the library | Story time | It is |
|---|---|---|---|
| Series | the library, such as Apogee | years | the world: the company, its places, its people, its machines |
| Season | a course | about four months | one stretch of your career, with a beginning, a crisis and a finale |
| Episode | a chapter | about a week | one problem the team has, solved with one area of Go |
| Beat | a lesson | minutes to a day | one task, one idea, one small change in the world |

A season has an **opening** (once), then for each chapter a **briefing**, its lessons' **beats**
and a **debrief**. The last chapter's debrief is the season's finale and sets up the next
season.

### The opening

The opening is the hook: the first thing a learner sees, before the season's first lesson. It
answers four questions in under a minute, and ends by opening the first lesson.

1. **Where am I?** The world in one breath, in plain words, with one concrete image.
2. **Who am I here?** The learner's badge: their name (typed in, optional, kept on the Mac),
   their role, their start date, their number.
3. **Who and what will I meet?** The people and places, as cards, a sentence each.
4. **What do I do first, and why today?** The first assignment, with a reason that is happening
   now.

Rules:
- Four to six pages. A page says one thing in at most 60 words, plus its cards.
- Every page names its button with what happens next: "Meet the team", not "Next".
- Only the words a learner needs on day one, each explained on the page that uses it. The
  opening is where the five words of the first week are taught: for Apogee, pass, ground
  station, AOS and LOS, telemetry, uplink.
- It can be skipped, and replayed from Mission Control.
- It says what the season holds without spoiling it: a launch is coming; how it goes is not
  said.

### Briefings

A briefing opens a chapter: "this week at Apogee". It shows on Mission Control while the
chapter is next, and above its first lesson.

- Two to four sentences, at most 70 words.
- When and where, what is wrong or needed, who needs it, and what is at stake: "Furrowlight's
  photos arrive with holes in them, and their farmers plan a week of spraying on them."
- No Go words. The chapter's description lists the Go it teaches; the briefing says why the team
  needs it.

### Beats

A beat is a lesson's first one or two sentences (the story bible's style guide gives the whole
readme's shape). It says why this task matters right now, and it is what the app shows as the
teaser after the previous lesson's pass, so it has to make the learner want to start.

- A character or a moment, a number, a reason: "Pipit-1 rises over Port Ellery in ten minutes,
  and the checklist printer has shuffled Ingrid's steps."
- Never a recap of the last lesson, never a lecture.

### Debriefs

A debrief closes a chapter, on the panel that greets the pass that finishes it, beside the
chapter's patch.

- One to three sentences, at most 50 words.
- What changed in the world because of the learner's work, often in a character's words: "Your
  words came back from space at 04:19 UTC. Ingrid pinned the printout above the Board."
- It may end on a small hook into the next briefing; it never starts a new problem.

### The finale

The last chapter of a season builds on everything before it and ends on the season's biggest
moment. Its debrief closes the season's story and points at the next one: a promotion, a new
station, a new kind of trouble.

## Plain words

The learner is a smart person who has never worked in space or in software. Write for them.

- **Plain words first, the term second**, the first time: "the moment the satellite comes into
  view, which ops calls AOS". After that the term is fine.
- **Every word of the world is in the course's field guide**, with its meaning in a sentence or
  two of plain language. The app links a readme's first use of a term to its meaning, and the
  field guide is always one click away. A term that lessons also use as an everyday word, such as
  pass ("the tests pass"), is in the field guide but not linked.
- **Titles use plain words.** A course, chapter or lesson title may use a term the learner has
  already met; a title seen before its chapter starts, as Mission Control shows the coming
  chapters, avoids acronyms.
- **Go's words are taught, not avoided.** Slice, goroutine and interface are what the learner
  came for: each lesson's idea names its Go precisely. The field guide is for the world's words;
  Go's words have the Go documentation, built into the app.
- **Numbers over adjectives**: "seven kilometers a second", not "very fast".
- **The app's words are plain too**: lesson, chapter, course, hint, solution, review, badge,
  patch, XP (experience points). No word appears on screen that a newcomer would have to look up
  to use the app.

## The learner, the cast and the stakes

- **The learner** is competent and growing. They are never the joke, never blamed, never
  hurried. Mistakes are normal at Apogee: a failure gets a name, a reason and a line in the log.
- **The cast** is small: a beat uses one person, rarely two. Each has one voice, kept exactly as
  the story bible gives it, and each reacts to the learner's work. A buddy who reviews your code
  carries the season; others come and go with their chapters.
- **The stakes are human-sized**: late data, a missed pass, a customer's deadline, a student
  team's last chance to hear their satellite. Nobody is hurt, nothing explodes, and no disaster
  is played for drama.
- **Urgency comes from the world, never from the app.** A pass in ten minutes is a story; the
  app has no timers, streaks to lose, guilt or countdowns.
- **Humor belongs to the characters** (an engineer's altitude jokes, a station chief's
  telegraphic weather), never at the learner's expense, never as filler.

## Continuity

- Small facts return: the 29 bytes of your first message, Pipit-1's beacon sitting 1.2 kHz high,
  the dish nicknames at Salar Quieto. Each return costs nothing and makes the world feel real.
- The code the learner writes grows into one program (for Apogee, `groundctl`), with the story
  bible's types and names, never renamed.
- Story time moves forward with the chapters, and never backward.
- Nothing contradicts the story bible. A writer who needs a new fact adds it to the bible
  first, then uses it.

## How the app tells it

All of it comes from the course's library (`docs/library-format.md`): a course without a story
looks as it always did.

| Moment | What the learner sees |
|---|---|
| First launch, or Begin course | The opening, full window, ending in the first lesson |
| Mission Control | Their badge and patches; the briefing of the chapter they're in, with the next lesson |
| A chapter's first lesson | The chapter's briefing above the readme |
| Any readme | The first use of each linked term underlined; a click shows its meaning |
| A pass | XP, and the next lesson's beat as a teaser |
| The pass that finishes a chapter | The debrief, and the chapter's patch added to the badge |
| Field Guide, in the sidebar | Every term of the course's world, in plain words |

Rewards follow progress and nothing else: XP, rank, milestones and patches are worked out from
the lessons done, never bought, never lost by being away.

## What a story never does

- Copies another course, platform or world, or uses a real company, product, person or flown
  mission.
- Explains its world in a wall of text. Lore arrives one beat at a time.
- Hypes ("awesome", "super easy", "simply"), scolds or guilts.
- Puts an unexplained acronym in front of a newcomer.
- Spoils a coming chapter in a title.
- Lets the story crowd out the Go: a beat is one or two sentences, and the lesson is the task.

## A writer's checklist

Before a chapter goes to review:

- [ ] The briefing says when, where, what's wrong, who needs it and what's at stake, in at most
      70 words and no Go.
- [ ] Every beat gives a reason for today and makes the next lesson worth starting.
- [ ] The debrief shows what the learner's work changed, in at most 50 words.
- [ ] Every word of the world used is in the field guide, glossed in plain words at first use,
      and linked unless it's an everyday word.
- [ ] Each character sounds like their entry in the story bible.
- [ ] No timers, guilt, hype, real names or spoilers; nothing at the learner's expense.
- [ ] The chapter moves story time forward and contradicts nothing in the story bible.
