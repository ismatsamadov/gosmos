# Files, Data and the Command Line: syllabus

**Authoring status:** all eleven chapters are implemented as 99 lessons; chapter 10 was added later, between Files That Travel With the Program and The Top of the Road. Chapter 1 passed its
calibration checks before the rest was written; chapters 3–10 were written from this plan
two chapters at a time and each pair independently reviewed. The lessons in
`courses/apogee.sqlite` are the source of truth: writing and review changed details such as
function signatures, exercises and test cases, and the larger changes are listed under
"Changes in writing and review" at the end.

Season 2 of Apogee, July to October 2027. The van heads up the Quebrada Añil with no signal; Paz needs gctl to run a calibration day from files alone. The learner moves from one file to a tested command with explicit streams, configuration and embedded defaults.

- **Title:** Files, Data and the Command Line
- **Slug:** `files-data-command-line` · **Kind:** Course
- **Short description:** Build gctl, the team's offline command-line tool, on the test van's climb to Salar Quieto.
- **Description:** The test van climbs the Quebrada Añil to recalibrate Salar Quieto's dishes, with days at a time and no signal. With Paz in the van and Nkem reviewing the code, you build gctl from files alone: whole files, folders and safe saves; readers and writers; logs read line by line; bookings in JSON and spreadsheets in CSV, odd ones included; times that mean the same instant everywhere; one command with flags, subcommands and exit statuses; a station table that travels inside the program; and parcels in XML, zip and tar, unpacked safely. At the top of the road, a whole calibration day runs from gctl and files alone. Standard library only, with every test on isolated files and streams.

## At a glance

| # | Chapter | Slug | Go topics | Lessons |
|---|---|---|---|---|
| 1 | A Schedule on a Stick | `a-schedule-on-a-stick` | Read and write whole files, build paths, create folders, inspect permissions and replace a file through a temporary sibling. | 9 |
| 2 | Everything Is a Stream | `everything-is-a-stream` | Use Reader and Writer interfaces, copy streams and control buffering and limits. | 9 |
| 3 | Line by Line | `line-by-line` | Scan text safely, handle long records, parse fields and preserve useful line context. | 9 |
| 4 | JSON In, JSON Out | `json-in-json-out` | Encode and decode structured data, choose field names and distinguish missing, null and empty values. | 9 |
| 5 | Odd JSON | `odd-json` | Handle streamed values, strict input, variable payloads, exact numbers and custom JSON methods. | 9 |
| 6 | Hollis's Spreadsheets | `hollis-spreadsheets` | Read and write CSV with quoting, header validation and row-level diagnostics. | 9 |
| 7 | What Time Is It, Exactly? | `time-exactly` | Parse and format explicit layouts, offsets and durations without depending on the Mac’s zone files. | 9 |
| 8 | Flags and Subcommands | `flags-and-subcommands` | Build testable commands with explicit arguments, streams and exit statuses. | 9 |
| 9 | Files That Travel With the Program | `files-with-the-program` | Embed fixed resources and work against fs.FS, including test-only in-memory filesystems. | 9 |
| 10 | Parcels From Salar Quieto | `parcels-from-salar-quieto` | Decode and encode XML with struct tags, read zip and gzip'd tar archives, and unpack them with checked names, an os.Root and a cap on every file. | 9 |
| 11 | The Top of the Road | `the-top-of-the-road` | Compose the season’s helpers into an offline calibration tool tested through complete sessions. | 9 |

99 lessons. Calibration starts with chapter 1; later chapters are authored only after it has been imported and read in the app. Difficulty rises from 3 to 8. Every test uses an injected stream, fstest.MapFS or t.TempDir; no test touches real learner files. Times use UTC or time.FixedZone.

## Chapter 1 · A Schedule on a Stick

`a-schedule-on-a-stick` · 9 lessons

Read and write whole files, build paths, create folders, inspect permissions and replace a file through a temporary sibling.

#### 1.1 Read the Schedule
`read-the-schedule` · code_tests · difficulty 3
- **Objective:** Read an entire file with os.ReadFile.
- **Scenario:** Paz receives the van's schedule on a memory stick.
- **Exercise:** Implement readSchedule(path) returning its text and error.
- **Check:** Exact contents, empty files and missing files.

#### 1.2 Save the Report
`save-the-report` · code_tests · difficulty 3
- **Objective:** Write bytes with os.WriteFile and propagate errors.
- **Scenario:** Ingrid needs a report the van can carry offline.
- **Exercise:** Implement saveReport(path, text) with mode 0600.
- **Check:** New contents, replacement, empty text and missing parent.

#### 1.3 Quick Check: Files
`quick-check-files` · quiz · difficulty 3
- **Objective:** Distinguish returned bytes, errors and creation permissions.
- **Scenario:** Nkem checks the first two file helpers before the van leaves.
- **Exercise:** Answer three questions about ReadFile and WriteFile.
- **Check:** Missing paths return errors; string converts bytes; mode applies when creating.

#### 1.4 Paths That Fit
`paths-that-fit` · code · difficulty 3
- **Objective:** Build a local path with filepath.Join.
- **Scenario:** Paz needs one path per day's schedule.
- **Exercise:** Replace string concatenation with filepath.Join.
- **Check:** plans/2027-07-12/passes.txt and a cleaned parent path.

#### 1.5 Make Room for Reports
`make-room-for-reports` · code_tests · difficulty 3
- **Objective:** Create nested folders with os.MkdirAll.
- **Scenario:** The van has no reports folder on its first run.
- **Exercise:** Implement prepareReports(root) for a nested reports/daily folder.
- **Check:** New and existing directories succeed with no group or other permission bits; a blocking file returns an error.

#### 1.6 Who Can Read It
`who-can-read-it` · code_tests · difficulty 3
- **Objective:** Inspect a file's permission bits with FileMode.Perm.
- **Scenario:** Mirela checks that private reports are readable only by their owner.
- **Exercise:** Implement ownerOnly(path), propagating Stat errors.
- **Check:** 0600 accepts; 0640 and 0644 refuse; missing file errors.

#### 1.7 Quick Check: Paths
`quick-check-paths` · quiz · difficulty 3
- **Objective:** Recognize path construction and directory operations.
- **Scenario:** Nkem checks the tool before its first run on the van.
- **Exercise:** Answer three questions about Join, MkdirAll and permission bits.
- **Check:** Join cleans; MkdirAll accepts existing directories; Perm isolates rwx bits.

#### 1.8 Replace It Whole
`replace-it-whole` · code_tests · difficulty 3
- **Objective:** Write a sibling temporary file before renaming it into place.
- **Scenario:** A failed save must not erase the last usable schedule.
- **Exercise:** Implement replaceSchedule(path, data), closing before Rename and cleaning temporary files.
- **Check:** Create, replace, a reader opened before the replacement still reads the old file, failed rename, unchanged destination and no leaked temporary files.

#### 1.9 The Stick Is Ready
`the-stick-is-ready` · code_tests · difficulty 4
- **Objective:** Compose path, directory and file helpers.
- **Scenario:** Paz needs the day's folder prepared and the schedule saved before departure.
- **Exercise:** Implement exportDay(root, day, text) using the supplied helpers.
- **Check:** Correct location and bytes, replaced report and propagated setup errors.

## Chapter 2 · Everything Is a Stream

`everything-is-a-stream` · 9 lessons

Use Reader and Writer interfaces, copy streams and control buffering and limits.

#### 2.1 Read Without a Path
`read-without-a-path` · code_tests · difficulty 3
- **Objective:** Accept an io.Reader rather than a filename.
- **Scenario:** The recorder must accept a file or an in-memory test feed.
- **Exercise:** Implement readText(r) with io.ReadAll.
- **Check:** Complete bytes, empty input and read errors.

#### 2.2 Write to Any Destination
`write-to-any-destination` · code_tests · difficulty 3
- **Objective:** Accept an io.Writer and return write errors.
- **Scenario:** Paz wants one report writer for the screen and a file.
- **Exercise:** Implement writeStatus(w, code, state).
- **Check:** Exact line, empty values and writer failure.

#### 2.3 Copy the Recorder
`copy-the-recorder` · code_tests · difficulty 3
- **Objective:** Transfer a stream with io.Copy.
- **Scenario:** The van must copy recordings without holding them all at once.
- **Exercise:** Implement copyRecording(dst, src), returning the byte count.
- **Check:** Full stream, empty stream and both error directions.

#### 2.4 Quick Check: Streams
`quick-check-streams` · quiz · difficulty 3
- **Objective:** Distinguish stream interfaces from concrete storage.
- **Scenario:** Nkem reviews the recorder's signatures.
- **Exercise:** Answer questions on Reader and Writer.
- **Check:** Readers produce bytes; writers consume them; Copy reports errors.

#### 2.5 Two Copies at Once
`two-copies-at-once` · code_tests · difficulty 3
- **Objective:** Write one stream to several writers with io.MultiWriter.
- **Scenario:** The operator needs the same report on screen and disk.
- **Exercise:** Implement mirrorReport(src, screen, archive).
- **Check:** Identical destinations and stopped writes after a destination error.

#### 2.6 Read and Keep a Copy
`read-and-keep-a-copy` · code_tests · difficulty 3
- **Objective:** Record consumed bytes with io.TeeReader.
- **Scenario:** Paz wants a trace of what the decoder actually read.
- **Exercise:** Implement captureText(src, log).
- **Check:** Consumed bytes logged, partial source failure and logging errors.

#### 2.7 Only the First Part
`only-the-first-part` · code_tests · difficulty 3
- **Objective:** Bound input consumption with io.LimitReader.
- **Scenario:** The preview must leave the rest of a recording unread.
- **Exercise:** Implement preview(r, n).
- **Check:** Zero, short and long inputs; remainder still readable.

#### 2.8 Flush Before Leaving
`flush-before-leaving` · code_tests · difficulty 3
- **Objective:** Flush a bufio.Writer and report flush errors.
- **Scenario:** A report looked complete but never reached the stick.
- **Exercise:** Implement bufferedReport(w, lines).
- **Check:** Final data delivered and delayed writer errors returned.

#### 2.9 The Recorder Pipeline
`the-recorder-pipeline` · code_tests · difficulty 4
- **Objective:** Combine limits, copying and multiple destinations.
- **Scenario:** The van needs a bounded diagnostic copy on disk and screen.
- **Exercise:** Implement recordPreview(src, limit, targets).
- **Check:** Byte limit, mirrored outputs, empty input and propagated errors.

## Chapter 3 · Line by Line

`line-by-line` · 9 lessons

Scan text safely, handle long records, parse fields and preserve useful line context.

#### 3.1 One Reading per Line
`one-reading-per-line` · code_tests · difficulty 4
- **Objective:** Scan lines with bufio.Scanner.
- **Scenario:** The controller log is too large to load whole.
- **Exercise:** Implement countLines(r).
- **Check:** Blank and final unterminated lines, empty input and scan errors.

#### 3.2 Keep the Last Error
`keep-the-last-error` · code_tests · difficulty 4
- **Objective:** Check Scanner.Err after iteration.
- **Scenario:** A damaged recording ended halfway through a report.
- **Exercise:** Implement collectLines(r), returning the scanning error.
- **Check:** Normal records and a failing reader after one line.

#### 3.3 Long Controller Lines
`long-controller-lines` · code_tests · difficulty 4
- **Objective:** Set a scanner buffer and maximum record size.
- **Scenario:** Paz's detailed calibration record exceeds the default limit.
- **Exercise:** Implement scanRecords(r) with a stated 256 KiB record bound.
- **Check:** A 100 KiB line succeeds; an over-limit record errors.

#### 3.4 Quick Check: Scanning
`quick-check-scanning` · quiz · difficulty 4
- **Objective:** Identify scanner limits and error handling.
- **Scenario:** Nkem checks the parser before the large log arrives.
- **Exercise:** Answer questions about Scan, Text, Buffer and Err.
- **Check:** Scan ending alone does not prove EOF; Buffer precedes scanning.

#### 3.5 Fields Across Spaces
`fields-across-spaces` · code_tests · difficulty 4
- **Objective:** Split whitespace-separated records with strings.Fields.
- **Scenario:** The controller alternates tabs and repeated spaces.
- **Exercise:** Implement parseColumns(line).
- **Check:** Tabs, repeated spaces and missing columns.

#### 3.6 Numbers That Can Fail
`numbers-that-can-fail` · code_tests · difficulty 4
- **Objective:** Parse fields with strconv and reject invalid numbers.
- **Scenario:** A typed letter must not become a zero angle.
- **Exercise:** Implement parseElevation(text).
- **Check:** Decimals, negative values and malformed text.

#### 3.7 Name the Bad Line
`name-the-bad-line` · code_tests · difficulty 4
- **Objective:** Wrap parsing failures with a line number.
- **Scenario:** Paz needs to find one broken record among 200000 lines.
- **Exercise:** Implement readElevations(r).
- **Check:** Correct values and errors naming the original failing line.

#### 3.8 Ignore Only Comments
`ignore-only-comments` · code_tests · difficulty 4
- **Objective:** Distinguish comments, blank lines and data.
- **Scenario:** The night shift adds notes without changing readings.
- **Exercise:** Implement dataLines(r), trimming only for classification.
- **Check:** Indented comments, blank lines and preserved data text.

#### 3.9 Summarize the Log
`summarize-the-log` · code_tests · difficulty 5
- **Objective:** Accumulate a streaming report without keeping all records.
- **Scenario:** Paz needs counts and the highest angle before calibration.
- **Exercise:** Implement summarize(r).
- **Check:** Empty input, repeated stations, maximum and bad-record errors.

## Chapter 4 · JSON In, JSON Out

`json-in-json-out` · 9 lessons

Encode and decode structured data, choose field names and distinguish missing, null and empty values.

#### 4.1 The Booking Export
`the-booking-export` · code_tests · difficulty 4
- **Objective:** Decode JSON into a struct.
- **Scenario:** Yusra's booking export must reach the van unchanged.
- **Exercise:** Implement decodeBooking(data).
- **Check:** Expected fields, bad syntax and wrong field types.

#### 4.2 Names on the Wire
`names-on-the-wire` · code_tests · difficulty 4
- **Objective:** Set JSON field names with struct tags.
- **Scenario:** The export says sat_id while Go says SatID.
- **Exercise:** Tag the supplied Booking fields.
- **Check:** Exact wire names and a round trip.

#### 4.3 Send a Booking Back
`send-a-booking-back` · code_tests · difficulty 4
- **Objective:** Encode a struct with json.Marshal.
- **Scenario:** The van must return its confirmed booking.
- **Exercise:** Implement encodeBooking(b).
- **Check:** Decoded contents, escaping and marshal errors from supplied invalid data.

#### 4.4 Quick Check: JSON Fields
`quick-check-json-fields` · quiz · difficulty 4
- **Objective:** Predict exported fields and tag behavior.
- **Scenario:** Nkem reviews the exchange format.
- **Exercise:** Answer three field and tag questions.
- **Check:** Unexported fields omitted; tags rename; syntax errors returned.

#### 4.5 A List of Passes
`a-list-of-passes` · code_tests · difficulty 4
- **Objective:** Decode and encode slices of records.
- **Scenario:** The day's booking file contains many passes.
- **Exercise:** Implement decodeBookings(data).
- **Check:** Several, empty, null and invalid entries.

#### 4.6 Optional Notes
`optional-notes` · code_tests · difficulty 4
- **Objective:** Use omitempty deliberately.
- **Scenario:** Hollis wants short exports without empty notes.
- **Exercise:** Tag optional fields and implement the stated encoder.
- **Check:** Empty note omitted; nonempty preserved; required false retained.

#### 4.7 Unknown Station Names
`unknown-station-names` · code_tests · difficulty 4
- **Objective:** Decode JSON objects into maps.
- **Scenario:** The customer adds station keys without asking for a new build.
- **Exercise:** Implement decodeCounts(data).
- **Check:** Arbitrary keys, zero values and invalid numeric values.

#### 4.8 Missing Is Not Zero
`missing-is-not-zero` · code_tests · difficulty 4
- **Objective:** Use a pointer field when presence matters.
- **Scenario:** An absent mask means keep the old one; zero means set zero.
- **Exercise:** Implement decodeMaskUpdate(data).
- **Check:** Absent, null, zero and positive values.

#### 4.9 The Round Trip
`the-round-trip` · code_tests · difficulty 5
- **Objective:** Compose decoding, validation and encoding.
- **Scenario:** Yusra needs an accepted booking returned with a stable identifier.
- **Exercise:** Implement confirm(data).
- **Check:** Preserved fields, rejected missing IDs and valid output JSON.

## Chapter 5 · Odd JSON

`odd-json` · 9 lessons

Handle streamed values, strict input, variable payloads, exact numbers and custom JSON methods.

#### 5.1 Bookings in a Stream
`bookings-in-a-stream` · code_tests · difficulty 5
- **Objective:** Decode successive JSON values with json.Decoder.
- **Scenario:** The recorder appends one booking at a time.
- **Exercise:** Implement decodeStream(r).
- **Check:** Several values, whitespace, clean EOF and truncated final value.

#### 5.2 Refuse Misspelled Fields
`refuse-misspelled-fields` · code_tests · difficulty 5
- **Objective:** Enable DisallowUnknownFields for configuration.
- **Scenario:** A misspelled safety limit must not be silently ignored.
- **Exercise:** Implement strictConfig(r).
- **Check:** Known fields accepted and unknown fields refused.

#### 5.3 One Object Means One
`one-object-means-one` · code_tests · difficulty 5
- **Objective:** Require EOF after the expected JSON object.
- **Scenario:** Two concatenated settings objects must not pick the first silently.
- **Exercise:** Implement decodeOne(r).
- **Check:** Trailing whitespace accepted; a second value and garbage refused.

#### 5.4 Quick Check: Decoders
`quick-check-decoders` · quiz · difficulty 5
- **Objective:** Recognize strict decoding and stream boundaries.
- **Scenario:** Mirela reviews the input boundary.
- **Exercise:** Answer questions about EOF and unknown fields.
- **Check:** EOF after a complete object is expected; second values are rejected explicitly.

#### 5.5 Keep the Payload for Later
`keep-the-payload-for-later` · code_tests · difficulty 5
- **Objective:** Delay decoding with json.RawMessage.
- **Scenario:** Furrowlight sends several report shapes in one envelope.
- **Exercise:** Implement decodeEnvelope(data).
- **Check:** Raw payload retained and malformed envelopes refused.

#### 5.6 Choose by Report Kind
`choose-by-report-kind` · code_tests · difficulty 5
- **Objective:** Decode RawMessage according to a discriminator.
- **Scenario:** Image reports and station reports need different checks.
- **Exercise:** Implement decodeReport(data).
- **Check:** Both known kinds, unknown kind and invalid payload.

#### 5.7 Numbers Without Rounding
`numbers-without-rounding` · code_tests · difficulty 5
- **Objective:** Keep large integers exact with Decoder.UseNumber.
- **Scenario:** A frame identifier is larger than float64 can represent exactly.
- **Exercise:** Implement readSequence(r).
- **Check:** Large exact values, fractions rejected and overflow reported.

#### 5.8 A Duration in JSON
`a-duration-in-json` · code_tests · difficulty 5
- **Objective:** Implement custom JSON methods for duration text.
- **Scenario:** Paz wants settings that say 30s rather than raw nanoseconds.
- **Exercise:** Implement MarshalJSON and UnmarshalJSON on Delay.
- **Check:** Round trips, invalid duration, wrong JSON type and escaped text.

#### 5.9 The Changing Report
`the-changing-report` · code_tests · difficulty 6
- **Objective:** Combine strict envelopes, exact IDs and typed payloads.
- **Scenario:** Furrowlight changes its export while the van remains offline.
- **Exercise:** Implement normalizeReport(r).
- **Check:** Known versions, exact IDs and errors for unsupported versions.

## Chapter 6 · Hollis's Spreadsheets

`hollis-spreadsheets` · 9 lessons

Read and write CSV with quoting, header validation and row-level diagnostics.

#### 6.1 Read the First Sheet
`read-the-first-sheet` · code_tests · difficulty 5
- **Objective:** Read CSV records with csv.Reader.
- **Scenario:** Hollis sends bookings as a spreadsheet export.
- **Exercise:** Implement readSheet(r).
- **Check:** Rows, empty sheet and malformed quotes.

#### 6.2 Commas Inside Names
`commas-inside-names` · code_tests · difficulty 5
- **Objective:** Let the CSV parser handle quoted fields.
- **Scenario:** A note contains commas and a line break.
- **Exercise:** Implement readNotes(r).
- **Check:** Quoted commas, escaped quotes and multiline fields.

#### 6.3 Headers Have Meaning
`headers-have-meaning` · code_tests · difficulty 5
- **Objective:** Map columns by header rather than position.
- **Scenario:** Hollis moves the station column.
- **Exercise:** Implement stationColumn(r).
- **Check:** Reordered headers, missing header and duplicate header.

#### 6.4 Quick Check: CSV
`quick-check-csv` · quiz · difficulty 5
- **Objective:** Distinguish CSV records from string splitting.
- **Scenario:** Nkem reviews the spreadsheet parser.
- **Exercise:** Answer questions about quoting, records and writer errors.
- **Check:** Fields may contain commas and newlines; Split is insufficient.

#### 6.5 The Wrong Number of Cells
`the-wrong-number-of-cells` · code_tests · difficulty 5
- **Objective:** Validate record width and report row errors.
- **Scenario:** One booking is missing its orbit number.
- **Exercise:** Implement checkedRows(r).
- **Check:** Expected width, short rows and long rows.

#### 6.6 Numbers from Cells
`numbers-from-cells` · code_tests · difficulty 5
- **Objective:** Parse typed CSV values with useful context.
- **Scenario:** Hollis's elevation column contains one word.
- **Exercise:** Implement readPassRows(r).
- **Check:** Valid values and errors naming row and column.

#### 6.7 Write a Sheet Back
`write-a-sheet-back` · code_tests · difficulty 5
- **Objective:** Encode CSV with csv.Writer and check Error.
- **Scenario:** Yusra needs the corrected file for the customer.
- **Exercise:** Implement writeSheet(w, rows).
- **Check:** Quoting, final flush and writer failure.

#### 6.8 A Different Separator
`a-different-separator` · code_tests · difficulty 5
- **Objective:** Configure the documented CSV dialect.
- **Scenario:** The van's controller exports semicolons.
- **Exercise:** Implement readControllerCSV(r).
- **Check:** Semicolons, quoted separators and rejected bad quoting.

#### 6.9 Hollis Gets the Report
`hollis-gets-the-report` · code_tests · difficulty 6
- **Objective:** Read, validate and write a complete spreadsheet transformation.
- **Scenario:** The missing-pass report must survive a round trip through a spreadsheet.
- **Exercise:** Implement summarizeSheet(r, w).
- **Check:** Reordered columns, correct totals, quoting and errors.

## Chapter 7 · What Time Is It, Exactly?

`time-exactly` · 9 lessons

Parse and format explicit layouts, offsets and durations without depending on the Mac’s zone files.

#### 7.1 Read the Timestamp
`read-the-timestamp` · code_tests · difficulty 6
- **Objective:** Parse RFC3339 timestamps with time.Parse.
- **Scenario:** The booking file states the exact start of a pass.
- **Exercise:** Implement parseStart(text).
- **Check:** UTC, explicit offsets and invalid text.

#### 7.2 Write the Station Clock
`write-the-station-clock` · code_tests · difficulty 6
- **Objective:** Format time using Go's reference layout.
- **Scenario:** Paz wants the same readable heading on every laptop.
- **Exercise:** Implement heading(t).
- **Check:** Exact date and minute formatting.

#### 7.3 It Is the Same Instant
`it-is-the-same-instant` · code_tests · difficulty 6
- **Objective:** Normalize an instant to UTC.
- **Scenario:** The van's local clock differs from the network clock.
- **Exercise:** Implement networkTime(text).
- **Check:** Equivalent offsets yield the same UTC instant.

#### 7.4 Quick Check: Layouts
`quick-check-layouts` · quiz · difficulty 6
- **Objective:** Recognize the reference time and offsets.
- **Scenario:** Nkem checks the timestamp conversion.
- **Exercise:** Answer questions about layouts, UTC and instants.
- **Check:** The reference date is 2006-01-02; offsets do not change an instant.

#### 7.5 The Van's Fixed Offset
`the-van-s-fixed-offset` · code_tests · difficulty 6
- **Objective:** Use time.FixedZone without system zone files.
- **Scenario:** The offline laptop needs a fixed display offset.
- **Exercise:** Implement vanTime(t, secondsEast).
- **Check:** Positive and negative offsets and preserved instant.

#### 7.6 How Long to Wait
`how-long-to-wait` · code_tests · difficulty 6
- **Objective:** Parse human-readable duration text.
- **Scenario:** Paz changes a timeout from the settings file.
- **Exercise:** Implement positiveDelay(text).
- **Check:** 30s and 2m accepted; zero, negative and malformed refused.

#### 7.7 The Pass Across Midnight
`the-pass-across-midnight` · code_tests · difficulty 6
- **Objective:** Compute duration from complete timestamps.
- **Scenario:** The last pass begins before midnight and ends after it.
- **Exercise:** Implement passDuration(start, end).
- **Check:** Across midnight, offsets and end-before-start errors.

#### 7.8 Quick Check: Time Values
`quick-check-time-values` · quiz · difficulty 6
- **Objective:** Distinguish timestamps, offsets and elapsed durations.
- **Scenario:** Ingrid checks the calibration-day clock.
- **Exercise:** Answer three time-value questions.
- **Check:** Sub measures elapsed time; UTC is display; durations carry no date.

#### 7.9 The Day in One Zone
`the-day-in-one-zone` · code_tests · difficulty 7
- **Objective:** Group passes by an explicitly chosen date.
- **Scenario:** The van and the network must agree which day owns a pass.
- **Exercise:** Implement dayBookings(passes, day, zone).
- **Check:** Boundary instants, sorting and no system-zone dependency.

## Chapter 8 · Flags and Subcommands

`flags-and-subcommands` · 9 lessons

Build testable commands with explicit arguments, streams and exit statuses.

#### 8.1 A Function Behind Main
`a-function-behind-main` · code_tests · difficulty 6
- **Objective:** Put command behavior in run with injected streams.
- **Scenario:** Paz wants command tests without starting another program.
- **Exercise:** Implement run(args, stdin, stdout, stderr).
- **Check:** Expected output and success code.

#### 8.2 Its Own Flag Set
`its-own-flag-set` · code_tests · difficulty 6
- **Objective:** Parse flags with flag.NewFlagSet and ContinueOnError.
- **Scenario:** The tool needs a station filter.
- **Exercise:** Implement parseOptions(args, errors).
- **Check:** Defaults, explicit flags and invalid values.

#### 8.3 Help Is a Success
`help-is-a-success` · code_tests · difficulty 6
- **Objective:** Handle flag.ErrHelp and direct usage output.
- **Scenario:** The night shift asks how to run the tool.
- **Exercise:** Implement help handling.
- **Check:** Help exits zero and invalid options exit two.

#### 8.4 Quick Check: Commands
`quick-check-commands` · quiz · difficulty 6
- **Objective:** Distinguish stdout, stderr and process status.
- **Scenario:** Nkem reviews the shell-facing contract.
- **Exercise:** Answer three command-contract questions.
- **Check:** Data uses stdout; diagnostics stderr; run returns status.

#### 8.5 Choose a Subcommand
`choose-a-subcommand` · code_tests · difficulty 6
- **Objective:** Dispatch subcommands before parsing their flags.
- **Scenario:** gctl grows passes and stations.
- **Exercise:** Implement dispatch(args, streams).
- **Check:** Known and unknown commands and no-command usage.

#### 8.6 Read from the Pipe
`read-from-the-pipe` · code_tests · difficulty 6
- **Objective:** Use stdin when a filename is absent.
- **Scenario:** The recorder feeds gctl directly.
- **Exercise:** Implement input selection.
- **Check:** Injected stdin read, explicit file used and open errors returned.

#### 8.7 Errors Belong on Stderr
`errors-belong-on-stderr` · code_tests · difficulty 6
- **Objective:** Keep machine-readable output free of diagnostics.
- **Scenario:** A script must distinguish bad input from an empty report.
- **Exercise:** Implement command error reporting.
- **Check:** No success data on errors; nonzero code and readable diagnostics.

#### 8.8 Quick Check: Exit Codes
`quick-check-exit-codes` · quiz · difficulty 6
- **Objective:** Choose stable statuses for usage and runtime failure.
- **Scenario:** Ingrid checks the tool's failure behavior.
- **Exercise:** Answer questions on successful help, bad flags and failed input.
- **Check:** Zero success, two usage, one runtime failure.

#### 8.9 Three Useful Commands
`three-useful-commands` · code_tests · difficulty 7
- **Objective:** Compose parsers and dispatch for gctl passes, stations and stow.
- **Scenario:** Paz needs the van's daily tasks under one command.
- **Exercise:** Implement the supplied command handlers.
- **Check:** Command-specific flags, validation, output and error codes.

## Chapter 9 · Files That Travel With the Program

`files-with-the-program` · 9 lessons

Embed fixed resources and work against fs.FS, including test-only in-memory filesystems.

#### 9.1 The Table Travels Too
`the-table-travels-too` · code_tests · difficulty 7
- **Objective:** Embed a file with a go:embed directive.
- **Scenario:** The van cannot depend on a station-table file beside the binary.
- **Exercise:** Embed stations.csv and expose its contents.
- **Check:** Bundled bytes and no working-directory lookup.

#### 9.2 A Filesystem in the Binary
`a-filesystem-in-the-binary` · code_tests · difficulty 7
- **Objective:** Embed a directory as embed.FS.
- **Scenario:** The van needs several fixed report templates.
- **Exercise:** Implement template loading from embedded data.
- **Check:** Known template, missing template and exact contents.

#### 9.3 Read Any Filesystem
`read-any-filesystem` · code_tests · difficulty 7
- **Objective:** Accept fs.FS and read with fs.ReadFile.
- **Scenario:** Tests should supply their own station table.
- **Exercise:** Implement readStationsFS(fsys).
- **Check:** Embedded FS and in-memory FS behave alike.

#### 9.4 Quick Check: Embedded Files
`quick-check-embedded-files` · quiz · difficulty 7
- **Objective:** Recognize build-time resources and fs paths.
- **Scenario:** Nkem checks what ships in the tool.
- **Exercise:** Answer questions on embedding and paths.
- **Check:** Files are included at build time; FS paths use slash and relative names.

#### 9.5 Walk the Templates
`walk-the-templates` · code_tests · difficulty 7
- **Objective:** Visit resources with fs.WalkDir.
- **Scenario:** Paz wants a list of templates that ship with gctl.
- **Exercise:** Implement listTemplates(fsys).
- **Check:** Nested files, directory filtering and walk errors.

#### 9.6 A Disk of Test Data
`a-disk-of-test-data` · code_tests · difficulty 7
- **Objective:** Build fixtures with fstest.MapFS.
- **Scenario:** Mirela wants tests isolated from the real disk.
- **Exercise:** Construct a MapFS for the supplied reader.
- **Check:** Contents, nested paths and missing-file errors.

#### 9.7 Views of One Folder
`views-of-one-folder` · code_tests · difficulty 7
- **Objective:** Restrict a resource view with fs.Sub.
- **Scenario:** A renderer should see templates without their top-level prefix.
- **Exercise:** Implement templateFS(fsys).
- **Check:** Relative reads and errors for missing resources.

#### 9.8 Quick Check: Filesystem Views
`quick-check-filesystem-views` · quiz · difficulty 7
- **Objective:** Distinguish fs.FS from writable operating-system paths.
- **Scenario:** Nkem reviews the resource boundary.
- **Exercise:** Answer questions on ReadFile, Sub and MapFS.
- **Check:** FS permits reading; Sub changes the root; MapFS is test data.

#### 9.9 The Offline Resource Pack
`the-offline-resource-pack` · code_tests · difficulty 8
- **Objective:** Combine embedded defaults, overrides and FS-based tests.
- **Scenario:** The van needs fixed defaults plus a supplied calibration file.
- **Exercise:** Implement loadResources(fsys, override).
- **Check:** Default data, valid override and surfaced override errors.

## Chapter 10 · Parcels From Salar Quieto

`parcels-from-salar-quieto` · 9 lessons

Decode and encode XML with struct tags, read zip and gzip'd tar archives, and unpack them with checked names, an os.Root and a cap on every file.

**Story:** Friday 15 to Sunday 17 October 2027, the last weekend at 3,700 m before Monday's climb. A stick comes down the road from Salar Quieto with two parcels from La Grande's dish controller: its settings as a zip, its log since July as a gzip'd tar. Calibration day will replace La Grande's correction, so Paz wants the old one safe on the stick; Nkem reviews the XML, and Mirela sends test parcels that try to write outside the folder or swell past the cap.

#### 10.1 Read the Manifest
`read-the-manifest` · code_tests · difficulty 7
- **Objective:** Map XML onto a struct with element, `,attr`, `a>b` and `,chardata` tags.
- **Scenario:** La Grande's controller sends a parcel whose manifest.xml lists what it holds.
- **Exercise:** Tag the supplied Manifest and Entry fields.
- **Check:** Version from an attribute, dish and time from elements, each file's name and size from nested elements; the unknown `<firmware>` skipped.

#### 10.2 What the Decoder Skips
`what-the-decoder-skips` · code_tests · difficulty 7
- **Objective:** Check required elements after xml.Unmarshal, which skips unknown elements and leaves missing ones at zero.
- **Scenario:** Nkem asks what gctl reads when `<dish>` is missing.
- **Exercise:** Complete the checks in readManifest.
- **Check:** Unknown element accepted; missing or empty dish, missing time, no files and correction.xml by mistake refused by name; broken XML refused.

#### 10.3 The Correction in XML
`the-correction-in-xml` · code · difficulty 7
- **Objective:** Write XML with an XMLName field, xml.MarshalIndent and xml.Header.
- **Scenario:** On Thursday La Grande's controller takes its new correction only in the shape of the parcel's correction.xml.
- **Exercise:** Name the root, indent by two spaces and print the header first.
- **Check:** Output equal to the parcel's correction.xml, byte for byte.

#### 10.4 Quick Check: XML
`quick-check-xml` · quiz · difficulty 7
- **Objective:** Predict how tags decode, what decoding does with unknown and missing elements, and what MarshalIndent writes.
- **Scenario:** Nkem reviews the XML code at Friday's late signal.
- **Exercise:** Answer three questions about a Stow type.
- **Check:** `,attr` reads an attribute matched exactly; unknown elements skipped and missing ones zero, with no error; the header comes from xml.Header.

#### 10.5 Open the Parcel
`open-the-parcel` · code_tests · difficulty 7
- **Objective:** Read one file out of a zip archive with zip.OpenReader, File and Open.
- **Scenario:** gctl must open lg-settings-2027-10-15.zip itself instead of a copy Paz took out by hand.
- **Exercise:** Implement parcelManifest(path).
- **Check:** The manifest from a parcel in t.TempDir, first or last in it; no manifest, a manifest without a dish, a file that is not a zip and a missing parcel refused.

#### 10.6 Names That Stay Inside
`names-that-stay-inside` · code_tests · difficulty 7
- **Objective:** Refuse archive names that would leave the target folder with filepath.IsLocal.
- **Scenario:** Mirela sends test parcels whose names point outside.
- **Exercise:** Implement checkNames(files).
- **Check:** Absolute, `..` and wandering-then-climbing names refused by name, the first of two; names that wander but stay inside and file names with two dots accepted.

#### 10.7 The Log Parcel
`the-log-parcel` · code_tests · difficulty 7
- **Objective:** Read a gzip'd tar with compress/gzip and archive/tar, capping each file with io.LimitReader.
- **Scenario:** The controller log since July comes as a .tar.gz; Mirela warns that zeros pack a thousandfold.
- **Exercise:** Implement readLogs(r, limit).
- **Check:** Every file's text; an empty parcel; exactly the limit kept, one byte more refused by name; 32 MiB of zeros refused after reading a small part of the parcel; plain text and a cut-off parcel refused.

#### 10.8 Unpack Under a Root
`unpack-under-a-root` · code_tests · difficulty 7
- **Objective:** Write files only inside one folder through os.Root, under the size cap.
- **Scenario:** Mirela asks what happens when the stick's folder already holds a link leading out.
- **Exercise:** Complete unpackFile(root, f, limit).
- **Check:** Nested files with exact bytes; a link out of the folder and a climbing name refused with nothing written outside; exactly the limit kept, a swollen file stopped one byte past it.

#### 10.9 Unpack the Whole Parcel
`unpack-the-whole-parcel` · code_tests · difficulty 8
- **Objective:** Combine zip reading, name checks, os.Root, size caps and the manifest into one safe unpacking.
- **Scenario:** Sunday night, before Monday's climb, Paz wants La Grande's settings on the stick so the old correction can be put back.
- **Exercise:** Implement unpackParcel(path, dir) with the supplied helpers.
- **Check:** Manifest and files from a good parcel; an unsafe name writes nothing; a link refused with nothing outside; every file capped; no manifest, a bad manifest, a non-zip and a missing parcel refused; the current folder untouched.

## Chapter 11 · The Top of the Road

`the-top-of-the-road` · 9 lessons

Compose the season’s helpers into an offline calibration tool tested through complete sessions.

#### 11.1 One Configuration File
`one-configuration-file` · code_tests · difficulty 7
- **Objective:** Load strict configuration with defaults.
- **Scenario:** Paz prepares the van before the final climb.
- **Exercise:** Implement loadConfig(r).
- **Check:** Defaults, overrides, unknown fields and extra JSON values.

#### 11.2 Choose the Input Format
`choose-the-input-format` · code_tests · difficulty 7
- **Objective:** Dispatch CSV and JSON decoders explicitly.
- **Scenario:** The van receives both booking formats.
- **Exercise:** Implement readBookings(r, format).
- **Check:** Both formats, unsupported format and propagated errors.

#### 11.3 Validate Before Booking
`validate-before-booking` · code_tests · difficulty 7
- **Objective:** Reject invalid records before writing any output.
- **Scenario:** An impossible station code must not enter the day's plan.
- **Exercise:** Implement validateBookings(bookings, stations).
- **Check:** Missing IDs, times and stations; valid batch.

#### 11.4 Sort the Calibration Day
`sort-the-calibration-day` · code_tests · difficulty 7
- **Objective:** Order normalized passes without changing input.
- **Scenario:** Paz needs a stable plan at the top of the road.
- **Exercise:** Implement orderedPasses(passes).
- **Check:** UTC ordering, deterministic ties and input preserved.

#### 11.5 Save Only a Complete Report
`save-only-a-complete-report` · code_tests · difficulty 7
- **Objective:** Use atomic replacement for final exports.
- **Scenario:** A failed final save must preserve yesterday's plan.
- **Exercise:** Implement saveDay(path, report).
- **Check:** Exact report, replacement and cleanup on failure.

#### 11.6 Quick Check: The Whole Tool
`quick-check-the-whole-tool` · quiz · difficulty 7
- **Objective:** Choose safe boundaries for formats, streams and disk writes.
- **Scenario:** Nkem reviews the final wiring before departure.
- **Exercise:** Answer questions about validation and file replacement.
- **Check:** Validate before writing; stderr for failure; temporary sibling before rename.

#### 11.7 A Golden Session
`a-golden-session` · code_tests · difficulty 7
- **Objective:** Check complete command output against a fixture.
- **Scenario:** Ingrid wants the van's exact report checked before the trip.
- **Exercise:** Implement the run function under a supplied golden test.
- **Check:** Stable stdout, diagnostics and exit status across fixtures.

#### 11.8 Failures from End to End
`failures-from-end-to-end` · code_tests · difficulty 7
- **Objective:** Preserve errors and avoid partial success at command boundaries.
- **Scenario:** A broken stick must produce an explanation, not an empty plan.
- **Exercise:** Complete error paths in run.
- **Check:** Bad flags, bad data, failed output and missing resources.

#### 11.9 La Grande Points True
`la-grande-points-true` · code_tests · difficulty 8
- **Objective:** Finish the offline gctl calibration-day flow.
- **Scenario:** Paz reaches Salar Quieto and needs the whole day run from files.
- **Exercise:** Complete runDay using the tested season helpers.
- **Check:** A complete offline session, stable report, correct statuses and no real-user files touched.

## Changes in writing and review

Where the lessons as built differ from the entries above:

- **3.1–3.3** The starter of 3.1 has the scanner and returns its error, so 3.1 is the loop and
  3.2 is `Err`; 3.3's `scanRecords` returns a record count, and a sweep of exactly 256 KiB with
  its newline still scans. A sweep is kept in its own file, not in the 200,000-line log.
- **3.5, 3.7, 3.9** Errors read `3 columns instead of 4` and `line N: …`; 3.9 counts readings
  per dish (`LG`, `CH`), not per station.
- **4.1, 4.3, 4.6, 4.8** 4.1's Booking has no tags yet (ID, Station, AOS, LOS); 4.3's marshal
  error is a time in year 20270; 4.6's kept `false` is a lesson-only `Confirmation`; 4.8 is
  `decodeMaskUpdate(data, oldDeg)`, returning the mask after the update.
- **5.5, 5.7, 5.9** The learner turns the payload into a `json.RawMessage` and refuses a
  missing or `null` one; 5.7 is `readReportID` (a report ID above 2^53, not a frame
  identifier); 5.9 accepts versions 1 and 2 and refuses version 3.
- **6.3–6.6, 6.8, 6.9** The learner writes `columnIndex` behind a given `stationColumn`; the
  6.4 quiz leaves writer errors to 6.7, which teaches them; 6.5's broken row lacks its LOS cell
  and every bad line is reported with `errors.Join`; 6.6 uses a lesson-only pass sheet
  `ID,STATION,ELEV,NOTE`; 6.8 also sets `Comment = '#'`; 6.9 writes `STATION,MISSED` then
  `TOTAL`.
- **7.3, 7.7–7.9** `networkTime` returns a string; `passDuration` takes the two timestamps as
  text; 7.8's second question is about `In` with the van's fixed zone; `dayBookings` returns a
  new sorted slice and no error.
- **8.1–8.9** 8.1's `run` prints the station table; 8.3's help handling is a shared
  `parseFlags(fs, args, stdout, stderr) (status int, ok bool)`; 8.5 dispatches inside `run`;
  8.6 is `loadBookings(name, stdin)`; 8.7 is `report(lines, err, stdout, stderr) int`, which
  never prints partial results; 8.9's learner writes `runStow`, the other commands given.
- **9.6, 9.7, 9.9** The learner's `fstest.MapFS` is `fixtureFS()` in `fixtures.go`;
  `templateFS` checks the folder with `fs.Stat` before `fs.Sub`; `loadResources` falls back to
  the built-in `calibration.json` only when the stick has none.
- **11.5, 11.7–11.9** `saveDay(path, report func(io.Writer) error)`; chapter 10's gctl has the
  `day` command only, checked by golden sessions; 10.8's error paths are in `planDay`; 10.9 is
  `runDay(dir, format, stdout)`.
