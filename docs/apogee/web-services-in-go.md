# Web Services in Go: syllabus

Season 3 of Apogee, November 2027 to February 2028. Customers book passes by email; on
1 February Apogee opens a booking service to every customer, and the learner builds it. The
learner moves from reading a single request to a tested service that authenticates, limits,
calls customers back and restarts in the middle of a live pass without dropping anything.

**Authoring status:** all thirteen chapters are implemented as 115 lessons. Chapter 1 was read in
the app in light and dark before the rest was written; chapters 2–11 were written from this
plan two chapters at a time and each pair independently reviewed. The world, the characters and
the writing rules are in `story.md`, whose section "Season 3: open for business" fixes every
shared type, route, header, message and number. The lessons in `courses/apogee.sqlite` are the
source of truth: writing and review changed details such as function signatures, exercises and
test cases, and the larger changes are listed under "Changes in writing and review" at the end.

- **Title:** Web Services in Go
- **Slug:** `web-services-in-go` · **Kind:** Course
- **Short description:** Build Apogee's booking service with net/http, from a first request to
  a launch that drops nothing.
- **Description:** Customers book passes by email, and Yusra copies every booking into a
  spreadsheet by hand. On 1 February Apogee opens a booking service that customers' programs
  call day and night, and you build it, with Nkem reviewing. You learn HTTP from both ends:
  clients that never wait forever, handlers and routes, JSON checked field by field, middleware
  that gives every request an ID and a line in the log, bookings that survive a restart, keys
  and signed tokens, fair rate limits, webhooks that sign and retry, and a shutdown that drops
  nothing in the middle of a live pass. Standard library only, and every server runs on your
  own Mac.
- **Story:** an opening, each chapter's briefing and debrief, and the field guide, told as
  `docs/storytelling.md` says, following the bible's spine for the season.

## At a glance

| # | Chapter | Slug | Go topics | Lessons |
|---|---|---|---|---|
| 1 | What the Web Speaks | `what-the-web-speaks` | Read and build requests and responses: methods, paths and queries, status codes and headers, through `httptest` and one loopback round trip. | 9 |
| 2 | Asking Politely | `asking-politely` | Call a service with `http.Client`: timeouts, contexts, status checks, closed bodies, lenient decoding and errors that carry a status. | 9 |
| 3 | Your First Endpoint | `your-first-endpoint` | Answer requests with `http.Handler`, `HandlerFunc` and `ResponseWriter`, route paths with `ServeMux`, serve files and run an `http.Server`. | 8 |
| 4 | An Address for Every Pass | `an-address-for-every-pass` | Route by method and wildcard, read path values and query parameters, and let the mux answer 404 and 405. | 8 |
| 5 | The First Booking | `the-first-booking` | Decode request bodies strictly and within limits, validate them, and answer every error as JSON. | 9 |
| 6 | Around Every Request | `around-every-request` | Wrap handlers in middleware for headers, request IDs, access logs and panic recovery, chained in order. | 8 |
| 7 | Keeping Bookings | `keeping-bookings` | Keep bookings behind an interface, in memory with a mutex and on disk with safe writes, held to one contract test. | 9 |
| 8 | Who's Asking | `whos-asking` | Authenticate callers with hashed API keys and HMAC-signed tokens, compare in constant time, and authorize per customer. | 10 |
| 9 | Fair Use | `fair-use` | Limit each customer with a token bucket on an injected clock, answer 429 with Retry-After, and set server and request timeouts. | 8 |
| 10 | Calling Customers Back | `calling-customers-back` | Deliver signed webhooks with backoff, jitter and retries, verify them on arrival, and make a create safe to retry. | 10 |
| 11 | Signed In | `signed-in` | Sign people in with cookies and server-side sessions, refuse other sites' forms with `http.CrossOriginProtection`, let one trusted origin call the API with CORS, and send a browser's security headers. | 9 |
| 12 | Launch Day | `launch-day` | Shut a server down gracefully, carry contexts into every handler, keep pending work across a restart, and test the whole service. | 9 |
| 13 | The First Week Open | `the-first-week-open` | Read uploads from `multipart` forms within a limit and judge their bytes with `http.DetectContentType`, resume downloads with `http.ServeContent`, sign download links that expire, refuse webhook addresses inside a network in a `net.Dialer`'s `Control` hook, keep redirects on the site, and count a monthly plan on an injected clock. | 9 |

115 lessons: 4 code, 94 code_tests and 17 quizzes (85.2% hands-on). Calibration starts with
chapter 1; later chapters are written only after it has been imported and read in the app.
Difficulty rises from 4 to 10. Every server listens on 127.0.0.1 on a port the system picks;
clocks and random sources are injected; no test reaches the network, sleeps more than 50 ms or
touches real learner files.

## Chapter 1 · What the Web Speaks

`what-the-web-speaks` · 9 lessons

Requests and responses with `net/http` and `httptest`: the request line, methods, query
strings, status codes and headers, a handler called in memory, and a first round trip over the
loopback.

#### 1.1 What a Request Says
`what-a-request-says` · code · difficulty 4
- **Objective:** See that an HTTP request is lines of text: a method, a path and a version, then headers.
- **Scenario:** Nkem sketches the first thing Furrowlight's program will ever send Apogee: a request for the list of stations.
- **Exercise:** `main` builds a request for `http://127.0.0.1/` with `http.NewRequest` and prints it with the read-only `printRequest`; you change the URL to ask for `/stations` and add `Accept: application/json` with `req.Header.Set`.
- **Check:** `GET /stations HTTP/1.1` / `Host: 127.0.0.1` / `User-Agent: Go-http-client/1.1` / `Accept: application/json`. The helper prints the wire's CR LF line ends as newlines; Go adds Host and User-Agent itself.

#### 1.2 Methods Say What You Want
`methods-say-what-you-want` · code · difficulty 4
- **Objective:** Choose the method that names the action: GET reads, POST creates, DELETE removes.
- **Scenario:** Yusra forwards three of this morning's emails: one asks for passes, one books a pass, one cancels booking `bk-4494`.
- **Exercise:** `main` builds three requests, all with `http.MethodGet`, and prints each first line with the read-only `requestLine`; you use `http.MethodPost` for the booking and `http.MethodDelete` for the cancellation.
- **Check:** `GET /passes HTTP/1.1` / `POST /bookings HTTP/1.1` / `DELETE /bookings/bk-4494 HTTP/1.1`.

#### 1.3 The Query After the Path
`the-query-after-the-path` · code_tests · difficulty 4
- **Objective:** Build a query string with `url.Values` and `Encode`, which sorts the keys and escapes each value.
- **Scenario:** Hollis writes in capitals: "ONLY SALAR QUIETO. ONLY FURROW-A. FROM 09:00 OUR TIME."
- **Exercise:** `passesPath(sat, station, from string) string` glues the values onto `/passes?` by hand; you build the query with `url.Values`, leave out empty values, and return `/passes` alone when all three are empty.
- **Check:** `("furrow-a", "SQT", "")` → `/passes?sat=furrow-a&station=SQT`; `("", "SQT", "")` → `/passes?station=SQT`; `("", "", "")` → `/passes`; Hollis's `from` of `2027-12-06T09:00:00+01:00` → `/passes?from=2027-12-06T09%3A00%3A00%2B01%3A00&sat=furrow-a` (keys sorted), because a bare `+` in a query arrives as a space, the mistake the glued string makes.

#### 1.4 Quick Check: Requests
`quick-check-requests` · quiz · difficulty 4
- **Objective:** Recognize the parts of a request and what each method asks for.
- **Scenario:** Nkem checks your notes before Cirrine's engineers try the stand-in.
- **Exercise:** Three questions on a request line, the method that cancels a booking, and what follows `?` in a URL.
- **Check:** The request line is the method, the path with its query, and the version; DELETE removes; after `?` come `key=value` pairs joined by `&`.

#### 1.5 What an Answer Says
`what-an-answer-says` · code · difficulty 4
- **Objective:** Call a handler in memory with `httptest.NewRequest` and `httptest.NewRecorder`, and read the status line, headers and body it wrote.
- **Scenario:** Nkem's stand-in answers like the booking service will, about stations only, and you can ask it questions without any network.
- **Exercise:** `main` makes a recorder and a request for `/stations/PEL` and prints the result with the read-only `printResponse`, but never hands the request to the stand-in; you call `standIn().ServeHTTP(rec, req)`.
- **Check:** `HTTP/1.1 200 OK` / `Content-Type: application/json` / a blank line / `{"code":"PEL","name":"Port Ellery"}`. The starter prints only `HTTP/1.1 200 OK`: a recorder nobody wrote to still says 200.

#### 1.6 Every Answer Has a Status
`every-answer-has-a-status` · code_tests · difficulty 5
- **Objective:** Read a status code by its class: 2xx done, 4xx the request was wrong, 5xx the server failed, and 429 means wait.
- **Scenario:** Yusra's customer guide needs one line for each kind of answer a customer's program can get.
- **Exercise:** `advice(code int) string` returns `""`; you return `done` for 200–299, `wait and retry` for 429, `fix the request` for the other 400–499, `retry later` for 500–599 and `unexpected` otherwise.
- **Check:** 200, 201, 204, 299 → done; 400, 404, 409, 499 → fix the request; 429 → wait and retry; 500, 503, 599 → retry later; 199, 302, 600 → unexpected. The mistakes: `code > 200`, which misses 200, and treating 429 like the other 4xx.

#### 1.7 Headers Name the Body
`headers-name-the-body` · code_tests · difficulty 5
- **Objective:** Read a header with `Header.Get`, which matches names in any letter case, and parse `Content-Type` with `mime.ParseMediaType`.
- **Scenario:** Cirrine's test program threw away a good answer: its Content-Type said `application/json; charset=utf-8`, and their check compared the whole line.
- **Exercise:** `isJSON(h http.Header) bool` returns `false`; you return true when the media type is `application/json`, whatever its parameters and letter case.
- **Check:** `application/json`, `application/json; charset=utf-8` and `Application/JSON` → true; `text/plain`, `application/jsonp`, a missing header and an unparsable value → false; a header set as `content-type` is found. The mistake: comparing the whole value.

#### 1.8 Quick Check: Answers
`quick-check-answers` · quiz · difficulty 5
- **Objective:** Recognize status classes, header names and what a recorder is for.
- **Scenario:** Nkem's three questions before the stand-in goes on a real port.
- **Exercise:** Three questions on 201 against 200, `Header.Get("content-type")`, and a recorder against a server.
- **Check:** 201 means something was created; `Get` finds `Content-Type` in any letter case; a recorder calls the handler in memory, with no connection at all.

#### 1.9 A Real Round Trip
`a-real-round-trip` · code · difficulty 5
- **Objective:** Serve a handler on a loopback port with `httptest.NewServer`, send it a request with `http.Get`, and read and close the response body.
- **Scenario:** Friday: Nkem wants to watch a request cross a real connection before Cirrine's engineers use the stand-in themselves.
- **Exercise:** `main` starts the stand-in with `httptest.NewServer` and asks it for two stations through `fetch(url string) (status, body string, err error)`, which returns empty strings; you call `http.Get`, read the body with `io.ReadAll`, close it, and return `resp.Status` and the body.
- **Check:** `200 OK {"code":"PEL","name":"Port Ellery"}` / `404 Not Found no station "ZZZ"`. `main` trims the bodies and never prints the server's address.

## Chapter 2 · Asking Politely

`asking-politely` · 9 lessons

HTTP clients that never wait forever: `http.Client` timeouts, requests with a context, status
checks, closing every body, decoding answers, and errors that carry the status.

#### 2.1 A Client That Gives Up
`a-client-that-gives-up` · code_tests · difficulty 5
- **Objective:** Give an `http.Client` a `Timeout`, because `http.Get` and the default client wait forever.
- **Scenario:** On Tuesday the go/no-go poll waited four minutes on the weather service while Cape Wren's mast was silent, and Furrow-B's checklist finished after AOS.
- **Exercise:** `newWeatherClient(baseURL string, timeout time.Duration) *WeatherClient` gives it `&http.Client{}`; you set `Timeout`. `Reading(ctx, code)` is written and doesn't use `ctx` yet.
- **Check:** Against a stalled test server, a 50 ms client returns an error within a second, and `errors.Is(err, context.DeadlineExceeded)`; against a calm server it returns the reading; the client is a new one, not `http.DefaultClient`, the mistake that changes every caller's client.

#### 2.2 The Poll's Own Deadline
`the-polls-own-deadline` · code_tests · difficulty 5
- **Objective:** Bind a request to a context with `http.NewRequestWithContext` and send it with `Client.Do`.
- **Scenario:** Ingrid's checklist must be done ten minutes before AOS, so every call the poll makes has to end by then, even when the client's own timeout is longer.
- **Exercise:** `Reading(ctx, code)` builds its request with `http.NewRequest` and never looks at `ctx`; you build it with `http.NewRequestWithContext`.
- **Check:** With a 2 s client and a context that ends after 50 ms, a stalled call returns within a second with `context.DeadlineExceeded`; an already cancelled context → `context.Canceled` at once; a calm server still answers.

#### 2.3 Quick Check: Waiting
`quick-check-waiting` · quiz · difficulty 5
- **Objective:** Tell apart a client's timeout, a request's context and the default client.
- **Scenario:** Tomás, who wrote the weather service in 2024, asks why the poll ever trusted it: "It was never slow. It was patient!"
- **Exercise:** Three questions.
- **Check:** `http.Get` uses `http.DefaultClient`, which has no timeout; `Client.Timeout` also covers reading the body; whichever ends first, the context or the timeout, ends the call.

#### 2.4 A 404 Is an Answer
`a-404-is-an-answer` · code_tests · difficulty 5
- **Objective:** Check `resp.StatusCode` yourself, because `Client.Do` returns a nil error for any answer the server sends.
- **Scenario:** A typo asked the weather service for `CRW` instead of `CWR`, and the poll reported `invalid character 'o' in literal null (expecting 'u')`: the decoder had read the plain-text answer `no station "CRW"` as JSON.
- **Exercise:** `Reading` decodes whatever comes back; before decoding, you return an error naming `resp.Status` unless the status is 200.
- **Check:** 200 → the reading; 404 → an error containing `404 Not Found`; 503 and 500 → errors with their status. The mistake: checking `err` alone, which is nil for all three.

#### 2.5 Close Every Body
`close-every-body` · code_tests · difficulty 6
- **Objective:** Close every response body on every path, with `defer resp.Body.Close()` right after checking `Do`'s error.
- **Scenario:** After a week of polls the poll machine held hundreds of unfinished connections, and Tomás's old service stopped taking new ones.
- **Exercise:** `Reading` closes the body only after a successful decode; you defer the close.
- **Check:** The test's client wraps every body and counts closes: after a 200, a 404, a 503 and a 200 with broken JSON, no body is left open. The mistake: closing on the success path only.

#### 2.6 Quick Check: Status and Bodies
`quick-check-status-and-bodies` · quiz · difficulty 5
- **Objective:** Reason about status codes and bodies on the client side.
- **Scenario:** Ingrid's check before the poll's new client meets a real night.
- **Exercise:** Three questions.
- **Check:** `Do` returns a nil error with a 503; a body must be closed even when it is never read; a deferred `Close` runs on every return after it.

#### 2.7 Decode the Reading
`decode-the-reading` · code_tests · difficulty 6
- **Objective:** Decode JSON straight from a response body with `json.NewDecoder`, leniently for a service you don't own.
- **Scenario:** The weather service added `humidity_pct` to its answers last spring, and the poll must not break when it adds more.
- **Exercise:** `decodeWeather(r io.Reader) (Weather, error)` returns a zero reading, and `Reading` calls it; you decode into `Weather` and require `station` and `at`.
- **Check:** A full reading → its values; an extra field → still accepted; a missing `at` or `station` → an error; broken JSON → an error. The mistake: `DisallowUnknownFields`, right for the booking service's own input but wrong here.

#### 2.8 Errors That Carry the Status
`errors-that-carry-the-status` · code_tests · difficulty 6
- **Objective:** Return a custom error type that callers read with `errors.AsType`, keeping only a little of an error body with `io.LimitReader`.
- **Scenario:** The poll must tell a silent mast (503: wait for the next reading) from a typo in a station code (404: a bug to fix).
- **Exercise:** `StatusError` and its `Error` method are given; `Reading` returns `fmt.Errorf("weather for %s: %s", code, resp.Status)`; you wrap a `*StatusError` with `%w` instead, its `Body` at most 512 bytes of the answer, trimmed.
- **Check:** `errors.AsType[*StatusError]` finds 503 and `mast silent since 2027-11-08T12:51:00Z`; the message is `weather for CWR: 503 Service Unavailable: mast silent since 2027-11-08T12:51:00Z`; a 1 MB error body is cut to 512 bytes; a timeout is not a `StatusError`.

#### 2.9 The Wind Check
`the-wind-check` · code_tests · difficulty 6
- **Objective:** Combine a timed client, a deadline, status checks and decoding into one call per station.
- **Scenario:** Friday: the poll checks the wind at each station before the night's passes, and every answer must come in time or say why not.
- **Exercise:** `windCall(ctx context.Context, c *WeatherClient, code string, now time.Time) string` returns `""`; you call `Reading` and return the call as the bible gives it.
- **Check:** One test weather server: PEL at 22 km/h → `GO: wind 22 km/h`; YST at 61 → `GO: wind 61 km/h, warning`; CWR at 74 → `NO-GO: wind 74 km/h, stow`; exactly 50 → the warning and exactly 72 → stow; LAN read 47 minutes ago → `NO-GO: reading 47m0s old`, and exactly 15 minutes ago → GO; SQT's 503 → `NO-GO: mast silent`; a stalled station → `NO-GO: no weather` within a second, with a 50 ms client.

## Chapter 3 · Your First Endpoint

`your-first-endpoint` · 8 lessons

The server side: `http.Handler` and `HandlerFunc`, headers, status and body through
`ResponseWriter`, paths routed by `http.ServeMux`, files served from an `fs.FS`, and an
`http.Server` on a loopback listener.

#### 3.1 A Type That Answers
`a-type-that-answers` · code_tests · difficulty 5
- **Objective:** Satisfy `http.Handler` by giving a type a `ServeHTTP(w, r)` method.
- **Scenario:** Cirrine's engineers start their client this week, and the first thing it will ask the booking service is which stations exist.
- **Exercise:** `stationList`, a `[]Station`, has a `ServeHTTP` that writes nothing; you encode the list to `w` with `json.NewEncoder`.
- **Check:** Through a recorder: 200 and the five stations in table order with their JSON names (`code`, `lat_deg`, `bands`); a `stationList{}` → `[]`. The mistake: `fmt.Fprint(w, s)`, which writes Go's own syntax, not JSON.

#### 3.2 Any Function Will Do
`any-function-will-do` · code_tests · difficulty 5
- **Objective:** Turn a function or a method value into a handler with `http.HandlerFunc`.
- **Scenario:** Nkem would rather keep every handler as a method of one `Service` than invent a type for each.
- **Exercise:** `(*Service).listStations` is written; `(*Service).stationsHandler() http.Handler` returns `http.NotFoundHandler()`; you return `http.HandlerFunc(s.listStations)`.
- **Check:** A GET through it → the stations; two Services with different stations each answer with their own, because the method value keeps its receiver; the starter's handler → 404.

#### 3.3 Quick Check: Handlers
`quick-check-handlers` · quiz · difficulty 5
- **Objective:** Recognize what makes a handler.
- **Scenario:** Nkem's whiteboard: three lines of code and a question for each.
- **Exercise:** Three questions.
- **Check:** A handler is any value with `ServeHTTP(http.ResponseWriter, *http.Request)`; `http.HandlerFunc(f)` is a conversion, not a call; `s.listStations` without parentheses is a method value that remembers `s`.

#### 3.4 Headers Before the Body
`headers-before-the-body` · code_tests · difficulty 6
- **Objective:** Set headers, then the status, then write the body, because the first `Write` sends the headers.
- **Scenario:** Cirrine's client refuses the station list: the service set `Content-Type` after writing, so the list went out as `text/plain`.
- **Exercise:** `writeJSON(w http.ResponseWriter, status int, v any)` writes the body before setting its header; you marshal first, set `Content-Type: application/json`, call `WriteHeader(status)`, then write the JSON and a newline, and answer 500 with the plain text `internal error` when marshaling fails.
- **Check:** `writeJSON(rec, 201, ...)` → 201, `application/json`, a body that decodes and ends in a newline; a value that can't be marshaled (a channel) → 500, the plain-text error and no JSON header; `listStations` answers through it. The mistake: setting the header after `Write`, which the recorder shows as lost.

#### 3.5 Say It With a Status
`say-it-with-a-status` · code_tests · difficulty 6
- **Objective:** Answer an unusable request with a status before writing anything, using `http.Error`.
- **Scenario:** Yusra's guide promises that `/stations?band=X` lists the stations a customer can book on X-band, and Furrowlight has already typed `band=x`.
- **Exercise:** `listStations` ignores `band`; you read it with `r.URL.Query().Get`, keep the stations that offer it, and answer `http.Error(w, "unknown band \"x\"", http.StatusBadRequest)` for anything but `UHF`, `S` or `X`.
- **Check:** `?band=X` → YST, SQT, CWR; `?band=UHF` → PEL, LAN, since receive-only yagis don't count; no band and an empty `band=` → all five; `?band=x` → 400 `unknown band "x"`. The mistake: writing the list first, which leaves the status at 200.

#### 3.6 One Mux, Many Paths
`one-mux-many-paths` · code_tests · difficulty 6
- **Objective:** Route paths to handlers with `http.ServeMux`, and see the mux answer 404 for paths nobody registered.
- **Scenario:** Ingrid's monitor will ask the service every minute whether it is up, at `/health`.
- **Exercise:** `(*Service).routes() *http.ServeMux` registers `/stations` only; you add `/health`, answering `{"status":"ok"}` with `writeJSON`.
- **Check:** `/stations` and `/health` answer; `/stations/` and `/statoins` → 404 from the mux, because a path without a trailing slash matches only itself. The mistake: registering `/` for the health check, which then answers every path.

#### 3.7 The Guide Beside the Service
`the-guide-beside-the-service` · code_tests · difficulty 6
- **Objective:** Serve a folder of files under a path prefix with `http.FileServerFS` and `http.StripPrefix`.
- **Scenario:** Yusra wants her customer guide at `/guide/` on the service itself, so every customer reads the same version.
- **Exercise:** `routes` has no guide; you register `/guide/` with `http.StripPrefix("/guide/", http.FileServerFS(s.Guide))`.
- **Check:** With an `fstest.MapFS` guide: `/guide/` → `index.html` as `text/html`; `/guide/stations.html` → that page; `/guide/missing.html` → 404; `/guide` → 307 to `/guide/`. The mistake: no `StripPrefix`, so the file server looks for `guide/index.html` and every page is 404.

#### 3.8 Stations Go Live
`stations-go-live` · code_tests · difficulty 6
- **Objective:** Run a handler as an `http.Server` on a 127.0.0.1 listener whose port the system picks, and stop it cleanly.
- **Scenario:** Friday 19 November: Yusra sends Cirrine the address of the booking service's first test run.
- **Exercise:** `start(s *Service) (addr string, stop func() error, err error)` returns an error; you listen with `net.Listen("tcp", "127.0.0.1:0")`, serve `s.routes()` from an `http.Server` in a goroutine, and return the listener's address and the `stop` that the read-only `stopFunc(srv, served)` builds.
- **Check:** A real client gets `/stations` and `/health` from `addr`; `addr` starts with `127.0.0.1:` and its port isn't 0; after `stop`, a new request can't connect and `stop` returned nil. The mistake: listening on `:0`, which opens the port to every network.

## Chapter 4 · An Address for Every Pass

`an-address-for-every-pass` · 8 lessons

Routes with methods and wildcards: `GET /passes/{sat}/{station}/{orbit}`, `Request.PathValue`,
query parameters with defaults and limits, two methods on one path, and the 404s and 405s
`ServeMux` sends for you.

#### 4.1 Only the Right Method
`only-the-right-method` · code_tests · difficulty 6
- **Objective:** Put the method in a mux pattern, so a wrong method gets 405 and an `Allow` header.
- **Scenario:** Furrowlight's test program sent `POST /stations` by mistake and got the station list back, as if it had changed something.
- **Exercise:** `routes` registers `/stations`, `/health` and `/guide/` without methods; you register `GET /stations`, `GET /health` and `GET /guide/`.
- **Check:** GET → 200; HEAD → 200 with no body, through `httptest.NewServer`; `POST /stations` → 405 with `Allow: GET, HEAD`; `DELETE /health` → 405; an unknown path → 404.

#### 4.2 A Station by Its Code
`a-station-by-its-code` · code_tests · difficulty 6
- **Objective:** Capture a path segment with a `{code}` wildcard and read it with `r.PathValue`.
- **Scenario:** Yusra wants each station to have its own address, so the guide can link to Salar Quieto's page.
- **Exercise:** `getStation` answers 404 for every code; you register `GET /stations/{code}` and look up `r.PathValue("code")`.
- **Check:** `/stations/SQT` → Salar Quieto; `/stations/ZZZ` → 404 `no station "ZZZ"`; `/stations/sqt` → 404, since codes are capitals; `/stations/SQT/extra` and `/stations/` → 404 from the mux; `/stations` still lists all five.

#### 4.3 Pass IDs Are Paths
`pass-ids-are-paths` · code_tests · difficulty 6
- **Objective:** Match several path segments with several wildcards in one pattern.
- **Scenario:** Hollis wants a link to each pass in his emails, and a pass ID such as `furrow-a/SQT/15864` already reads like a path.
- **Exercise:** `getPass` is registered at `GET /passes/{id}`, which can't match an ID with slashes; you register `GET /passes/{sat}/{station}/{orbit}` and find the catalog pass whose ID is the three values joined by `/`.
- **Check:** `/passes/furrow-a/SQT/15864` → its PassView, `id` beside the pass's own fields; `/passes/furrow-a/SQT/15865` → 404 `no pass furrow-a/SQT/15865`; `/passes/furrow-a/SQT` → 404 from the mux; `/passes/pipit-1/PEL/08154` → Pipit-1's pass.

#### 4.4 Quick Check: Patterns
`quick-check-patterns` · quiz · difficulty 6
- **Objective:** Predict what `ServeMux` answers for a method and a path.
- **Scenario:** Nkem's quiz before Hollis's links go out.
- **Exercise:** Three questions about the patterns shown in the readme.
- **Check:** `GET /stations` also answers HEAD; with only `GET /stations/{code}` registered, `POST /stations/SQT` gets 405 and `/stations/SQT/x` gets 404; `{code}` matches exactly one non-empty segment.

#### 4.5 Filters in the Query
`filters-in-the-query` · code_tests · difficulty 6
- **Objective:** Read optional query parameters with `r.URL.Query()`, refuse impossible values with 400, and answer an empty list as `[]`.
- **Scenario:** Furrowlight's program asks for `/passes?sat=furrow-a&station=SQT` and gets every pass in the catalog.
- **Exercise:** `listPasses` returns the whole catalog; you keep the passes matching `sat` and `station` when they are given, and answer 400 for a satellite or a station the service doesn't know.
- **Check:** Both filters; one; none → the whole catalog in order; `sat=thistle` → 400 `unknown satellite "thistle"`; `station=ZZZ` → 400 `unknown station "ZZZ"`; a filter matching nothing → `[]`. The mistake: a nil slice, which encodes as `null`.

#### 4.6 A Week at a Time
`a-week-at-a-time` · code_tests · difficulty 7
- **Objective:** Parse a time and a count from the query, with defaults from the service's injected clock and with limits.
- **Scenario:** Yusra's TL;DR: "Furrowlight wants to book a week at a time."
- **Exercise:** `listPasses` lists every pass from now on; you read `from` (RFC 3339, default `s.Now()`) and `days` (1 to 14, default 7) and keep the passes whose AOS is in `[from, from + days)`.
- **Check:** The defaults with a fixed clock; `days=1`; `days=14`; `days=0`, `days=15` and `days=seven` → 400 `days must be a whole number from 1 to 14`; `from=monday` → 400 `from must be an RFC 3339 time`; a pass at exactly `from` is in and one at exactly `from + days` is out.

#### 4.7 Two Methods, One Address
`two-methods-one-address` · code_tests · difficulty 7
- **Objective:** Register two methods on one path with different handlers, and answer `204 No Content` with no body.
- **Scenario:** Brinewatch emails every week to cancel `bk-4496` and book its pass again, and Yusra wants a cancellation to be one request.
- **Exercise:** `GET /bookings/{id}` is registered; you register `DELETE /bookings/{id}`, remove the booking from `s.Store` and answer 204.
- **Check:** Through `httptest.NewServer`: DELETE → 204 with an empty body; GET afterwards → 404 `booking bk-4496: not found`; a second DELETE → 404; `PUT /bookings/bk-4496` → 405 with `Allow: DELETE, GET, HEAD`.

#### 4.8 Every Route in One Place
`every-route-in-one-place` · code_tests · difficulty 7
- **Objective:** Gather a service's routes in one `routes` method and check them all with one table of requests.
- **Scenario:** Friday: Hollis's first email with a link to a pass goes out, and Nkem wants the whole address book checked first.
- **Exercise:** `routes` misses `GET /passes` and registers `/bookings/{id}` without methods; you complete it as the bible's route table gives it so far.
- **Check:** A table of 16 requests: each route's normal answer, a 404 for each wildcard, a 405 with its `Allow` header for each path, and a HEAD of a list.

## Chapter 5 · The First Booking

`the-first-booking` · 9 lessons

JSON in and out: decoding request bodies, errors as JSON from one function, strict decoding,
`http.MaxBytesReader`, the media type a body must have, validation rules, and `201 Created`
with a `Location`.

#### 5.1 Read the Booking
`read-the-booking` · code_tests · difficulty 6
- **Objective:** Decode a request's JSON body into a struct with `json.NewDecoder(r.Body)`.
- **Scenario:** Hollis wants to send his first booking as JSON instead of an email in capitals.
- **Exercise:** `createBooking` answers 201 with an empty booking; you decode the body into a `Booking`, keep it with `s.Store.Create(r.Context(), b)` and answer 201 with what the store returns.
- **Check:** A valid booking → 201 with `id` `bk-4501` and the fields sent; the next → `bk-4502`; broken JSON and an empty body → 400 and nothing stored. The mistake: ignoring the decode error and storing an empty booking.

#### 5.2 Errors in JSON Too
`errors-in-json-too` · code_tests · difficulty 6
- **Objective:** Write every error as JSON through one function that reads the status from the error value with `errors.AsType`.
- **Scenario:** Furrowlight's program reads JSON, not sentences, and half the service's errors are still plain text.
- **Exercise:** `writeError(w http.ResponseWriter, r *http.Request, err error)` answers 500 to everything; you answer an `*APIError`, wrapped or not, with its status, code and message as an `ErrorBody`, and anything else with 500 `internal` and `internal error`.
- **Check:** An `*APIError` 400 → its body; the same wrapped with `%w` → the same; a plain error → 500, and its own text never appears; `application/json` every time. The mistake: a type assertion, which misses the wrapped error.

#### 5.3 Refuse What You Don't Know
`refuse-what-you-dont-know` · code_tests · difficulty 6
- **Objective:** Decode a request strictly: no unknown fields, and nothing after the JSON value.
- **Scenario:** Furrowlight sent `"sat":"furrow-a"` instead of `"sat_id"`, and the service tried to book a pass for no satellite at all.
- **Exercise:** `decodeJSON(w, r, dst any) error` decodes leniently; you call `DisallowUnknownFields`, require a second `Decode` to return `io.EOF`, and return 400 `bad_request` errors with the bible's messages.
- **Check:** A valid body → nil; `"sat"` → `request body: json: unknown field "sat"`; two objects → `request body: unexpected data after the JSON value`; a trailing newline and spaces → fine; an empty body → `request body: empty`.

#### 5.4 Quick Check: Request Bodies
`quick-check-request-bodies` · quiz · difficulty 6
- **Objective:** Reason about strict decoding and error values.
- **Scenario:** Mirela reads the decoder before Furrowlight's first real booking.
- **Exercise:** Three questions.
- **Check:** `DisallowUnknownFields` turns a misspelled field into an error; one `Decode` reads one value, so trailing data needs a second call; an `*APIError` wrapped with `%w` is still found by `errors.AsType`.

#### 5.5 A Limit on Every Body
`a-limit-on-every-body` · code_tests · difficulty 7
- **Objective:** Cap a request body with `http.MaxBytesReader`, and answer 413 when a caller sends too much.
- **Scenario:** A Furrowlight test sent a 40 MB note by accident, and the service read every byte of it into memory.
- **Exercise:** `decodeJSON` reads as much as it is sent; you wrap `r.Body` in `http.MaxBytesReader(w, r.Body, 65536)` and answer a `*http.MaxBytesError` from either decode with 413 `too_large`.
- **Check:** A valid body of exactly 65,536 bytes, padded with spaces → accepted; 65,537 → 413 `request body over 65536 bytes`; 40 MB → 413 quickly. The mistake: answering 400, because the error's type was never checked.

#### 5.6 Say It's JSON
`say-its-json` · code_tests · difficulty 6
- **Objective:** Check a request's `Content-Type` before reading its body, and answer 415 for any other format.
- **Scenario:** Brinewatch's first test posted a booking as a form, `sat_id=brine-3&station=LAN`, and the decoder's complaint made no sense to them.
- **Exercise:** `decodeJSON` decodes whatever arrives; you parse the media type with `mime.ParseMediaType` first and answer 415 `unsupported_media_type` `send the body as application/json`.
- **Check:** `application/json`, `application/json; charset=utf-8` and `Application/JSON` → accepted; `application/x-www-form-urlencoded` and `text/plain` → 415; no header → 415; a refused body is never decoded.

#### 5.7 Shape Before Substance
`shape-before-substance` · code_tests · difficulty 7
- **Objective:** Validate a decoded value field by field and answer 422 with a message that names the first problem.
- **Scenario:** Hollis's program sends bookings without `los`, and the service keeps them.
- **Exercise:** `checkShape(b Booking) error` returns nil; you return an `*APIError` 422 `invalid_booking` for the first broken rule, in the bible's order.
- **Check:** One case per rule; two problems → the first message; LOS equal to AOS → `los must be after aos`; a 200-character note containing `ñ` → fine, and 201 characters → refused. The mistake: `len`, which counts bytes.

#### 5.8 Only Passes That Exist
`only-passes-that-exist` · code_tests · difficulty 7
- **Objective:** Check a request against the service's own data and an injected clock.
- **Scenario:** A Furrowlight test booked a pass of Thistle, which re-entered in July, and another that had started an hour before.
- **Exercise:** `checkPass(b Booking, catalog []PassView, now time.Time) error` returns nil; you return the bible's four 422s, in order: unknown satellite, unknown station, no predicted pass, closed.
- **Check:** `furrow-a/SQT/15864` booked on 3 December → nil; Thistle → `unknown satellite "thistle"`; LOS one second off → `no predicted pass of furrow-a over SQT at 2027-12-06T14:02:10Z`; `now` at exactly 13:52:10 on 6 December → `bookings for this pass closed at 2027-12-06T13:52:10Z`, and one second earlier → nil.

#### 5.9 Furrowlight Books a Pass
`furrowlight-books-a-pass` · code_tests · difficulty 7
- **Objective:** Put decoding, checks, the store and errors together in one handler that answers `201 Created` with a `Location`.
- **Scenario:** Friday 3 December: Hollis's program books Monday's 14:02 pass over Salar Quieto, the first booking that never touches Yusra's spreadsheet.
- **Exercise:** `createBooking` decodes and stores; you run `decodeJSON`, `checkShape`, `checkPass` and `Create`, answer `ErrConflict` with 409 `conflict`, and set `Location: /bookings/<id>` on the 201.
- **Check:** Hollis's booking → 201, `Location: /bookings/bk-4501` and the booking; the same pass again → 409 `furrow-a over SQT at 2027-12-06T14:02:10Z: already booked`; each earlier failure keeps its status (415, 413, 400, 422) and stores nothing; then the week's other Furrow-A and Furrow-B passes → `bk-4502` onward.

## Chapter 6 · Around Every Request

`around-every-request` · 8 lessons

Middleware: functions that wrap an `http.Handler` to add headers, carry a request ID in the
request's context, record the status through a wrapped `ResponseWriter`, write an access log,
recover panics, and chain in a fixed order.

#### 6.1 A Handler Around a Handler
`a-handler-around-a-handler` · code_tests · difficulty 7
- **Objective:** Write middleware: a function that takes an `http.Handler` and returns one that does its part, then calls it.
- **Scenario:** An office proxy at Cirrine kept an old answer to `GET /bookings/bk-4511` and showed a cancelled booking for a day.
- **Exercise:** `noStore(next http.Handler) http.Handler` returns `next`; you return a `HandlerFunc` that sets `Cache-Control: no-store` and then calls `next`.
- **Check:** The header on a 200, on the mux's 404 and on a 405; the wrapped handler runs exactly once per request. The mistake: setting the header after `next`, when the answer has already gone.

#### 6.2 Every Request Gets an ID
`every-request-gets-an-id` · code_tests · difficulty 7
- **Objective:** Carry a value through a request with `context.WithValue` and `r.WithContext`, under a key type of your own.
- **Scenario:** When Cirrine reports a problem, Yusra has no way to point at the one request it was.
- **Exercise:** `requestID` numbers requests and sets `X-Request-Id`, but doesn't pass the ID on; you put it in the request's context under an unexported key type, write `requestIDFrom`, and `writeError` copies it into `request_id`.
- **Check:** `req-000001` then `req-000002`, the header and the error body agreeing; the handler sees the same ID; a context that never went through `requestID` → `""`; a value stored under the plain string key `"request_id"` is not taken for it.

#### 6.3 Quick Check: Middleware
`quick-check-middleware` · quiz · difficulty 7
- **Objective:** Reason about when middleware runs and what a request's context carries.
- **Scenario:** Nkem's questions before the wrapping starts.
- **Exercise:** Three questions.
- **Check:** Middleware runs its part when a request arrives, not when it wraps; `r.WithContext` returns a new request and leaves the old one alone; a key of an unexported type can't collide with another package's key.

#### 6.4 Remember the Status
`remember-the-status` · code_tests · difficulty 7
- **Objective:** Wrap `http.ResponseWriter` by embedding it, to record the status and the bytes a handler writes.
- **Scenario:** The access log must say what each request got, and a `ResponseWriter` tells nobody what was written to it.
- **Exercise:** `statusRecorder` embeds `http.ResponseWriter` and has `Status` and `Bytes`; you override `WriteHeader` and `Write` and add `Unwrap() http.ResponseWriter`.
- **Check:** An explicit 404; 200 when `Write` comes first; bytes summed over two writes; a second `WriteHeader` doesn't change the recorded status; `http.NewResponseController` on the recorder reaches the inner writer's `Flush` through `Unwrap`.

#### 6.5 A Line in the Log
`a-line-in-the-log` · code_tests · difficulty 7
- **Objective:** Write one access-log line per request from middleware, after the handler returns, with a `*log.Logger`.
- **Scenario:** Yusra pastes a request ID from Cirrine's email and wants the one line that says what happened to it.
- **Exercise:** `logRequests(l *log.Logger, now func() time.Time)` calls `next` and logs nothing; you wrap the writer in a `statusRecorder`, time the call with `now`, and log the bible's line.
- **Check:** `req-000007 POST /bookings 201 236B 4ms` with a fake clock that moves 4 ms; the mux's 404 logged as 404; a handler that only writes → 200; the query left out. The mistake: logging before `next` returns, with status 0.

#### 6.6 A Panic Gets an Answer
`a-panic-gets-an-answer` · code_tests · difficulty 8
- **Objective:** Recover a handler's panic in middleware and answer 500, while letting `http.ErrAbortHandler` through.
- **Scenario:** Monday's bug: a short orbit number in a pass ID made `getPass` panic, and Cirrine's program got a dropped connection instead of an answer.
- **Exercise:** `recoverPanics(l *log.Logger)` only calls `next`; you defer a function that recovers, logs `req-000003 panic: <value>`, and answers through `writeError` with an internal error, panicking again with `http.ErrAbortHandler` untouched.
- **Check:** A panicking handler → 500 JSON `internal` with the request ID, and the log line; a handler that doesn't panic is untouched; a panic with `http.ErrAbortHandler` goes on, and the test recovers it itself; through `httptest.NewServer`, the request after a panic still gets 200.

#### 6.7 Chain Them in Order
`chain-them-in-order` · code_tests · difficulty 7
- **Objective:** Compose middleware with a variadic `chain` in which the first one listed runs first.
- **Scenario:** Nkem wants the service's wrapping written as one list that reads top to bottom.
- **Exercise:** `chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler` returns `h`; you wrap from the last middleware to the first.
- **Check:** Three tracing middlewares record `a b c handler c b a`; no middleware → `h` itself; one middleware. The mistake: wrapping in list order, which reverses the trace.

#### 6.8 One Bug, No Outage
`one-bug-no-outage` · code_tests · difficulty 8
- **Objective:** Order a service's middleware so every request, even one that panics, gets an ID, an answer and a line in the log.
- **Scenario:** Friday: Cirrine sends the request that broke on Monday, and this time Yusra can quote its line in the log.
- **Exercise:** `(*Service).handler(l *log.Logger) http.Handler` chains `recoverPanics`, `logRequests`, `requestID` and `noStore` in that wrong order; you put them in the bible's order.
- **Check:** Through `httptest.NewServer`: the panicking route → 500 whose `request_id`, `X-Request-Id` header and log line `req-000002 GET … 500` agree; ordinary requests logged with their statuses; the next request answered. In the starter's order the panic leaves no line and the other lines have no ID.

## Chapter 7 · Keeping Bookings

`keeping-bookings` · 9 lessons

Storage behind an interface: a `Store` the handlers depend on, a map behind a mutex, one
contract test for every store, a file store that saves through a temporary file and reopens
after a restart, and store errors mapped to statuses.

#### 7.1 One Interface for Bookings
`one-interface-for-bookings` · code_tests · difficulty 7
- **Objective:** Declare an interface for what the handlers need, so any store can stand behind them.
- **Scenario:** On Monday a deploy restarted the beta and Furrowlight's twelve test bookings vanished; a store on disk is coming, and the handlers must not care which store they get.
- **Exercise:** `Service.Store` is a `*MemStore`; you declare `Store` with `Create`, `Get`, `List` and `Delete` as the bible gives them and make the field a `Store`. The starter fails to build with `cannot use &failingStore{} (value of type *failingStore) as *MemStore value in struct literal` (topics include `compiler-errors`).
- **Check:** The handlers work as before with the `MemStore`; with the test's failing store, `GET /bookings/bk-4501` and `POST /bookings` → 500 `internal error`, the store's own error text never shown.

#### 7.2 Bookings in Memory, Safely
`bookings-in-memory-safely` · code_tests · difficulty 7
- **Objective:** Implement an interface with a map behind a `sync.Mutex`, checking and inserting under one lock and handing out copies.
- **Scenario:** Until now a hidden placeholder kept the bookings; Nkem wants the service's own store, and on launch day two customers may try to book one pass in the same second.
- **Exercise:** `MemStore`'s four methods return zero values; you write them: IDs from `bk-4501`, `ErrConflict` for a second booking of the same satellite, station and AOS, `ErrNotFound` wrapped with the ID, and `List` sorted by AOS then ID in a new slice.
- **Check:** IDs in order; a missing ID → `booking bk-5499: not found` and `errors.Is(err, ErrNotFound)`; the conflict; changing `List`'s result leaves the store unchanged; 50 goroutines booking 50 passes → 50 different IDs, and 50 goroutines booking one pass → exactly one success; a cancelled context → `context.Canceled`.

#### 7.3 Quick Check: Stores
`quick-check-stores` · quiz · difficulty 7
- **Objective:** Reason about interface values, locks and copies.
- **Scenario:** Mirela, reading over Nkem's shoulder, asks three questions about the store.
- **Exercise:** Three questions.
- **Check:** A `*MemStore` satisfies `Store` because its methods have pointer receivers; the conflict check and the insert need one lock, or two customers can both win; `List` returns a new slice so no caller can change the store.

#### 7.4 One Test, Two Stores
`one-test-two-stores` · code_tests · difficulty 8
- **Objective:** Write one contract test as a function that checks any implementation of an interface.
- **Scenario:** The file store is next, and Nkem wants it held to every promise the memory store keeps.
- **Exercise:** `checkStore(t *testing.T, open func() Store)` in `contract.go` is empty; you write subtests: IDs from `bk-4501`, `Get` after `Create`, a conflict on the same pass, `Delete` then `Get` → `ErrNotFound`, and `List` in AOS order.
- **Check:** Your checks pass against `MemStore`, and in a re-executed test binary each of three planted stores fails at least one of them: one that books a pass twice, one whose `List` is unsorted, one whose `Delete` deletes nothing.

#### 7.5 Save Every Change
`save-every-change` · code_tests · difficulty 8
- **Objective:** Save a store's whole state after every change through a temporary file renamed into place, and undo the change when saving fails.
- **Scenario:** Yusra won't retire her spreadsheet until a booking survives the service stopping.
- **Exercise:** `FileStore` wraps a `*MemStore` and its `save()` does nothing; you marshal `List`'s bookings as one JSON array and write it with the read-only `replaceFile` after every `Create` and `Delete`.
- **Check:** After a `Create` the file holds that booking with its JSON names; after a `Delete`, `[]`; the file's mode is 0600; with the folder made read-only, `Create` returns an error and a `Get` of its ID → `ErrNotFound`. The mistake: a store that keeps a booking its file never got.

#### 7.6 Open Where You Left Off
`open-where-you-left-off` · code_tests · difficulty 8
- **Objective:** Load saved state at start, treating a missing file as empty and a broken one as an error.
- **Scenario:** After a restart the service must know every booking and never hand out `bk-4505` twice.
- **Exercise:** `OpenFileStore(path string) (*FileStore, error)` always starts empty; you read the file, decode it strictly, restore the bookings with the `MemStore`'s `load`, and continue the IDs after the highest.
- **Check:** A missing file → an empty store whose first ID is `bk-4501`; a file with `bk-4501` to `bk-4504` → next `bk-4505`; a file holding only Yusra's `bk-4491` → next `bk-4501`; duplicate IDs → an error; a broken file → an error and no store. The mistake: an empty store, whose next save writes over every booking.

#### 7.7 Quick Check: Files and Restarts
`quick-check-files-and-restarts` · quiz · difficulty 8
- **Objective:** Reason about safe saves and starts.
- **Scenario:** Nkem's last questions before the holiday deploy.
- **Exercise:** Three questions.
- **Check:** A rename swaps the whole file at once, so a reader sees all of the old file or all of the new; a broken file must stop the start; after a restart the next ID comes from the highest one saved.

#### 7.8 Store Errors, Right Statuses
`store-errors-right-statuses` · code_tests · difficulty 8
- **Objective:** Map a store's sentinel errors to statuses in one place with `errors.Is`.
- **Scenario:** Cirrine's program got a 500 for a booking that simply didn't exist, and their on-call engineer opened an incident.
- **Exercise:** `writeError` knows only `*APIError`; you answer `ErrNotFound` with 404 `not_found` and `ErrConflict` with 409 `conflict`, keeping the error's message, and take the special cases out of the handlers.
- **Check:** GET of a missing booking → 404 `booking bk-5499: not found`; a double booking → 409 with the store's message; a failing disk → 500 `internal error`; an `*APIError` still wins.

#### 7.9 Survive a Restart
`survive-a-restart` · code_tests · difficulty 8
- **Objective:** Wire a file-backed store into the service and show that bookings outlive the process that made them.
- **Scenario:** Friday 17 December: the last deploy before the holidays, and Furrowlight's January bookings must still be there on 3 January.
- **Exercise:** `newService(cfg Config) (*Service, error)` builds the Service on a `MemStore`; you open a `FileStore` at `bookings.json` in `cfg.Dir`, and add `GET /bookings`.
- **Check:** Through HTTP: three bookings; a second Service from the same folder lists the three in AOS order; its next booking is `bk-4504`; a cancellation survives another restart; a broken file stops `newService` with an error.

## Chapter 8 · Who's Asking

`whos-asking` · 10 lessons

Authentication and authorization: credentials in the `Authorization` header, constant-time
comparison with `crypto/subtle`, keys kept only as SHA-256 hashes, expiring HMAC-signed
tokens, middleware that fails closed, answers about the caller's own satellites only, and a
review written as tests.

#### 8.1 Who Is Calling?
`who-is-calling` · code_tests · difficulty 8
- **Objective:** Read a credential from the `Authorization` header: a scheme, one space, then the credential.
- **Scenario:** Mirela's review opens with one line: "Today anyone who can reach this service can read every booking. Who is calling?"
- **Exercise:** `credentials(r *http.Request, scheme string) (string, error)` returns `""`, nil; you return the credential after the expected scheme, matched in any letter case, and a 401 `*APIError` `missing credentials` otherwise.
- **Check:** `Key furrowlight.5b0e…` → the key; `key furrowlight.5b0e…` → the key too; no header, `Bearer x` when `Key` is expected, `Key` alone and `Key ` with nothing after → `missing credentials`. The mistake: `strings.HasPrefix(h, "Key ")`, which refuses `key`.

#### 8.2 Compare Without a Clue
`compare-without-a-clue` · code_tests · difficulty 8
- **Objective:** Compare secrets with `subtle.ConstantTimeCompare`, which takes as long however many bytes match.
- **Scenario:** Mirela times the key check: one that stops at the first wrong byte answers a little sooner the more of a guess is wrong.
- **Exercise:** `checkKey(keys map[string]string, key string) (string, bool)`, over a lesson-only table from customer ID to key, returns `"", false`; you take the customer ID before the first dot and compare the whole key with the stored one in constant time.
- **Check:** Each fixture key → its customer; a wrong secret, an unknown customer, an empty key, a key without a dot and a key one digit short → false. Because `==` gives the same answers, the test also reads `auth.go` with `go/parser` and requires a call to `subtle.ConstantTimeCompare`.

#### 8.3 Keep Only a Hash
`keep-only-a-hash` · code_tests · difficulty 8
- **Objective:** Keep a key's SHA-256 hash instead of the key, and check a presented key by hashing it.
- **Scenario:** Mirela: "If a backup of the customer table leaks, nobody should be able to book with what is in it."
- **Exercise:** `checkKey(customers map[string]Customer, key string) (Customer, bool)` compares with a table of plain keys that is gone; you hash the key with `sha256.Sum256`, hex-encode it and compare it with `KeyHash` in constant time.
- **Check:** Furrowlight's test key → Furrowlight, which the mistake of comparing the key itself with `KeyHash` refuses; a wrong secret, an unknown customer, a key in capitals and the hash itself sent as a key → refused.

#### 8.4 Quick Check: Keys
`quick-check-keys` · quiz · difficulty 8
- **Objective:** Reason about secrets, hashes and the two kinds of refusal.
- **Scenario:** Mirela's questions before any customer gets a key.
- **Exercise:** Three questions.
- **Check:** `==` on strings can stop at the first different byte, so its timing leaks how much was right; a table of hashes is useless to whoever steals it; 401 means "we don't know who you are", 403 "we know, and you may not".

#### 8.5 A Token That Expires
`a-token-that-expires` · code_tests · difficulty 8
- **Objective:** Sign a short-lived token with HMAC-SHA256 over everything it claims.
- **Scenario:** Mirela wants keys to cross the network once an hour, not with every request.
- **Exercise:** `signToken(customer string, expires time.Time, key []byte) string` returns `""`; you return the bible's `<customer>.<expiry>.<signature>`.
- **Check:** Furrowlight expiring at 2028-01-05T11:00:00Z → `furrowlight.1830682800.` and the HMAC the test computes under the test token key; another expiry → another signature, the mistake of signing the customer alone; another key → another signature.

#### 8.6 Check the Token
`check-the-token` · code_tests · difficulty 9
- **Objective:** Verify a signed token: its shape, then its signature with `hmac.Equal`, then its expiry against an injected clock.
- **Scenario:** A Brinewatch engineer asks Mirela what stops them from editing their token to say `cirrine`.
- **Exercise:** `verifyToken(tok string, key []byte, now time.Time) (string, error)` returns `""`, nil; you check the three parts, recompute and compare the signature, then check the expiry.
- **Check:** A valid token → its customer; the customer or the expiry edited → `invalid token`; another key → `invalid token`; two or four parts, and a signature that isn't hex → `invalid token`; one second before expiry → valid; exactly at expiry → `token expired`.

#### 8.7 Trade Keys for Tokens
`trade-keys-for-tokens` · code_tests · difficulty 8
- **Objective:** Write the handler that checks an API key and answers with a signed token and its expiry.
- **Scenario:** Furrowlight's program will call `POST /tokens` once an hour, and every other call will carry the token.
- **Exercise:** `issueToken` answers 501; you read the `Key` credential, check it with `checkKey`, and answer 200 with a `TokenResponse` valid for one hour from `s.Now()`, or 401 with `WWW-Authenticate: Key`.
- **Check:** Furrowlight's key at 10:00:00.6 → a token `verifyToken` accepts and `expires_at` 11:00:00, truncated to the second; a wrong key → 401 `invalid api key`; no header → 401 `missing credentials`; a token sent as the key → `invalid api key`.

#### 8.8 Fail Closed
`fail-closed` · code_tests · difficulty 9
- **Objective:** Write authentication middleware that refuses by default and passes the caller on in the request's context.
- **Scenario:** Mirela's sign-off rule: "LGTM once it fails closed."
- **Exercise:** `authenticate(customers, key, now)` calls `next` for everyone; you read the Bearer token, verify it, find its customer, answer 401 with `WWW-Authenticate: Bearer` on any failure, and otherwise call `next` with `withCustomer`. `routes` puts the passes and bookings behind `s.private`.
- **Check:** No header → 401 `missing credentials` and `next` never called; a bad token → `invalid token`; an expired one → `token expired`; a valid one → `next` sees Furrowlight through `customerFrom`; `/stations` and `/guide/` stay public. The mistake: writing the 401 and calling `next` anyway, which the test's handler counts.

#### 8.9 Only Your Own
`only-your-own` · code_tests · difficulty 9
- **Objective:** Authorize every request against its caller: show only their own satellites' passes and bookings, and refuse the rest.
- **Scenario:** Mirela's headline: "Brinewatch must never see Cirrine's passes."
- **Exercise:** `listPasses`, `getPass`, `createBooking`, `listBookings`, `getBooking` and `cancelBooking` ignore the caller; you write `owns(c Customer, satID string) bool` and use it with `customerFrom` so that the lists hold only the caller's, another customer's pass or booking answers 404 exactly like a missing one, and booking another customer's satellite answers 403.
- **Check:** Brinewatch lists only Brine-3's passes and bookings; Brinewatch's GET of Cirrine's `bk-4510` → 404 with the very message of a missing booking; its DELETE → 404, and the booking is still there; Brinewatch booking `cirrine-1` → 403 `cirrine-1 is not one of your satellites`; Cirrine still sees its own.

#### 8.10 Mirela's Review
`mirelas-review` · code_tests · difficulty 9
- **Objective:** Write a review of a whole service as a test: requests through `httptest` that must pass on the real service and fail on broken ones.
- **Scenario:** Friday 7 January: Mirela asks you to write her review as code, and plants broken services for it to catch.
- **Exercise:** `review(t *testing.T, base string, advance func(time.Duration))` in `review.go` is empty; over HTTP you check that every private route refuses a request without credentials, that Furrowlight's key gets a token that lists only Furrowlight's bookings, that a key sent as a Bearer token is refused, that a token is `token expired` after `advance(time.Hour)`, and that Brinewatch can't read Cirrine's booking.
- **Check:** Your review passes against the real service and fails, in a re-executed test binary, against each of four planted services: one that fails open, one that lists every customer's bookings, one whose tokens never expire, one that accepts keys as tokens.

## Chapter 9 · Fair Use

`fair-use` · 8 lessons

Rate limits and slow callers: a token bucket per customer on an injected clock, `429` with
`Retry-After`, clients that wait as told, `http.Server` timeouts, and a deadline on every
request.

#### 9.1 A Bucket That Refills
`a-bucket-that-refills` · code_tests · difficulty 8
- **Objective:** Implement a token bucket that allows ten requests at once and then one a second, from an injected time.
- **Scenario:** Brinewatch's beta client lost its network on Sunday night and has retried ever since with no pause: a request every millisecond.
- **Exercise:** `(*Bucket).Take(now time.Time) bool` returns true; you add the seconds since `Last` times the rate, cap at 10, and take a token if a whole one is there.
- **Check:** Ten at one instant → true and the eleventh → false; a second later → one more; half a second later → still false, and another half → true; an hour later → ten, not 3,600; a clock stepping backward adds nothing; a zero `Bucket` starts full.

#### 9.2 A Bucket for Each Customer
`a-bucket-for-each-customer` · code_tests · difficulty 8
- **Objective:** Keep per-key state in a map behind a mutex, safe for many requests at once.
- **Scenario:** Brinewatch's flood must not spend Cirrine's share.
- **Exercise:** `(*Limiter).Allow(customer string) (ok bool, remaining int, retryAfter time.Duration)` returns `true, 10, 0`; you keep a `*Bucket` per customer, made on first use, and return the whole tokens left and the time until one is back, rounded up to a second.
- **Check:** Brinewatch's eleventh refused while Cirrine's first is allowed; remaining 9, 8 … 0; `retryAfter` 1 s when empty; 20 goroutines × 10 calls for one customer at a frozen clock → exactly 10 allowed.

#### 9.3 Quick Check: Rate Limits
`quick-check-rate-limits` · quiz · difficulty 8
- **Objective:** Reason about buckets, rates and testing with a fake clock.
- **Scenario:** Yusra asks what to promise customers in the guide.
- **Exercise:** Three questions.
- **Check:** The bucket's size is the burst a customer may send at once; elapsed time times the rate refills it, up to the size; a frozen fake clock makes the test exact without sleeping.

#### 9.4 Tell Them When
`tell-them-when` · code_tests · difficulty 8
- **Objective:** Answer an over-limit request with 429 and a `Retry-After` header, and report what is left on every answer.
- **Scenario:** Brinewatch's engineers say they would gladly wait, if only the service said how long.
- **Exercise:** `limit(l *Limiter)` calls `next` for everyone; you call `Allow` with the customer from the context, set `X-RateLimit-Limit` and `X-RateLimit-Remaining`, and refuse with 429 `rate_limited` and `Retry-After` in whole seconds.
- **Check:** `10` and `9` on the first answer; the eleventh → 429 with `Retry-After: 1` and remaining `0`; with 0.6 of a token left, `Retry-After: 1`, not `0.4`; after the clock moves a second, allowed again; a request with no customer in its context → 500, never unlimited.

#### 9.5 Wait as Long as Told
`wait-as-long-as-told` · code_tests · difficulty 8
- **Objective:** Read `Retry-After` on the client side and wait that long, within sensible bounds.
- **Scenario:** Yusra sends Brinewatch a corrected sample client, and its retry loop must do what the header says.
- **Exercise:** `retryDelay(res *http.Response) time.Duration` returns 0; you parse whole seconds, use 1 s when the header is missing or unreadable, and cap it at 60 s.
- **Check:** `3` → 3 s; missing, `soon`, `0` and `-5` → 1 s; `600` → 60 s; an HTTP date → 1 s, because the booking service always sends seconds.

#### 9.6 A Slow Client Holds Nothing
`a-slow-client-holds-nothing` · code_tests · difficulty 8
- **Objective:** Protect a server from slow clients with `http.Server`'s timeouts and its header limit.
- **Scenario:** Brinewatch's client also sent its headers a byte at a time, and each of its connections sat open for minutes.
- **Exercise:** `newServer(h http.Handler, t Timeouts) *http.Server` sets only `Handler`; you set the four timeouts from `t` and `MaxHeaderBytes` to 16 KiB.
- **Check:** A raw connection that sends half a request and stops is closed within a second at a 50 ms header timeout; a kept-alive connection left idle is closed at a 50 ms idle timeout; a 32 KiB header → 431; an ordinary request is answered; each of the server's fields equals `t`'s.

#### 9.7 A Deadline for Every Request
`a-deadline-for-every-request` · code_tests · difficulty 9
- **Objective:** Give every request's context a deadline in middleware, and answer 503 when the work runs out of time.
- **Scenario:** During the flood, Cirrine's bookings waited half a minute for the store before anyone answered.
- **Exercise:** `withDeadline(d time.Duration)` calls `next` unchanged; you give the request a context from `context.WithTimeout`, cancel it when `next` returns, and make `writeError` answer `context.DeadlineExceeded` and `context.Canceled` with 503 `unavailable`.
- **Check:** A store that blocks until its context ends → 503 `request took too long` within a second at 50 ms; a fast handler is unaffected; the handler's context has a deadline about `d` away; `withDeadline(0)` adds none. The mistake: calling `next` with the old request instead of `r.WithContext(ctx)`, which the test's store catches.

#### 9.8 The Flood Test
`the-flood-test` · code_tests · difficulty 9
- **Objective:** Place the limiter and the deadline in the service so one customer's flood leaves the others untouched.
- **Scenario:** Friday 14 January: Nkem replays Brinewatch's flood against the service while Cirrine books passes.
- **Exercise:** `(*Service).private` puts `limit` outside `authenticate`, and `handler` has no `withDeadline`; you put `limit` inside `authenticate` and `withDeadline(s.Timeouts.Handler)` innermost around the mux.
- **Check:** Through `ServeHTTP`: 1,000 Brinewatch requests at a frozen clock → 10 answered and 990 × 429, each with `Retry-After`; Cirrine's five requests in between → all 200; requests without a token → 401, never 429, and they spend nobody's tokens; five seconds of fake time later Brinewatch gets five more.

## Chapter 10 · Calling Customers Back

`calling-customers-back` · 10 lessons

Webhooks: an event sent as a signed POST, signatures and timestamps checked on arrival,
exponential backoff with jitter from a seeded source, which failures to retry, delivery with
retries on injected waits, and idempotency keys that make a retried booking safe.

#### 10.1 Tell Them It Landed
`tell-them-it-landed` · code_tests · difficulty 8
- **Objective:** Send an event to a customer as a JSON `POST`, and count any 2xx answer as delivered.
- **Scenario:** Furrowlight's program asks `GET /bookings` every five seconds, 17,280 times a day, to learn whether an image has landed.
- **Exercise:** `post(ctx context.Context, c *http.Client, url string, secret []byte, ts int64, ev Event) error` returns nil without sending (its secret and timestamp are for the next lesson); you marshal the event, POST it with `Content-Type: application/json`, close the body, and return an error for any status outside 200–299.
- **Check:** The test's receiver gets the exact JSON and header; 200, 202 and 204 → nil; 404 and 500 → an error naming the status; a closed server → an error. The mistake: `!= 200`, which refuses a 204.

#### 10.2 Sign What You Send
`sign-what-you-send` · code_tests · difficulty 8
- **Objective:** Sign a request with HMAC-SHA256 over a timestamp and the exact body bytes, sent in headers.
- **Scenario:** Mirela: "Anyone can POST to Furrowlight's address. Only we can sign."
- **Exercise:** `signature(secret []byte, ts int64, body []byte) string` returns `""`; you return the hex HMAC of `<ts>.<body>`, and in `post` you set `Apogee-Timestamp` to `ts` and `Apogee-Signature` to the signature of the bytes you send.
- **Check:** The receiver's own computation matches; the signature covers the very bytes sent; a new timestamp gives a new signature. The mistake: signing the body alone, which lets an old event be sent again under a new time.

#### 10.3 Check It on Arrival
`check-it-on-arrival` · code_tests · difficulty 9
- **Objective:** Verify a signed request on the receiving side, a timestamp within five minutes and then the signature with `hmac.Equal`, before trusting anything in its body.
- **Scenario:** Yusra is sending customers a sample receiver, and Hollis will copy it line for line.
- **Exercise:** `verifyEvent(r *http.Request, secret []byte, now time.Time) (Event, error)` decodes without checking; you read at most 64 KiB, check the timestamp, recompute and compare the signature, and only then decode.
- **Check:** A valid request → the event; the body or the timestamp changed → an error; exactly 300 s old → accepted; 301 s old or 301 s ahead → refused; a missing header or a signature that isn't hex → an error.

#### 10.4 Quick Check: Signed Requests
`quick-check-signed-requests` · quiz · difficulty 8
- **Objective:** Reason about what a signature and a timestamp protect.
- **Scenario:** Mirela's questions after reading the sample receiver.
- **Exercise:** Three questions.
- **Check:** The timestamp is signed so a captured request can't be sent again later; the receiver checks the signature before trusting anything in the body; the five-minute window allows for clocks a little apart.

#### 10.5 Wait Longer Each Time
`wait-longer-each-time` · code_tests · difficulty 8
- **Objective:** Compute an exponential backoff with a cap and jitter from a seeded `math/rand/v2` source.
- **Scenario:** When Cirrine's receiver came back from maintenance, a sender retrying every second would have knocked it over again.
- **Exercise:** `backoff(attempt int, r *rand.Rand) time.Duration` returns 0; you double from 2 s per earlier failure, cap at 5 minutes, and multiply by a factor drawn from [0.5, 1.0).
- **Check:** With `rand.New(rand.NewPCG(1, 2))` the first seven waits fall in 1–2 s … 64–128 s and equal the test's own computation from the same seed; attempts 9, 20 and 70 fall in 150–300 s, the mistake being a shift that overflows; no wait is ever 0.

#### 10.6 Which Failures to Retry
`which-failures-to-retry` · code_tests · difficulty 8
- **Objective:** Decide from a status or an error whether trying again can help.
- **Scenario:** Brinewatch's receiver answers `410 Gone` from an address they retired, and Yusra doesn't want it knocked on eight times.
- **Exercise:** `retryable(status int, err error) bool` returns true; you retry transport errors, 408, 429 and 5xx, and nothing else.
- **Check:** A refused connection and a timeout → true; a cancelled context → false; 408, 429, 500, 502, 503 → true; 400, 401, 404, 410, 422 → false; 301 and 302 → false; 200 → false.

#### 10.7 Deliver With Retries
`deliver-with-retries` · code_tests · difficulty 9
- **Objective:** Retry a delivery on its schedule with the waits injected, so tests run in no time.
- **Scenario:** Cirrine's receiver answers 503 twice during its maintenance, and the event must arrive on the third try.
- **Exercise:** `(*Sender).Deliver(ctx context.Context, d *Delivery) error` tries once; you loop: sign with a fresh timestamp, post, stop on success or on a failure not worth retrying, otherwise wait `backoff(attempt)` through `s.Wait`, at most 8 attempts.
- **Check:** 503, 503, 204 → nil after 3 attempts with waits in 1–2 s and 2–4 s; always 503 → `ErrUndelivered` after 8 attempts and 7 waits; 404 → `ErrUndelivered` after 1; every attempt carries the same event ID and its own timestamp from the fake clock; a context cancelled during a wait → `context.Canceled`, with `d.Attempts` kept.

#### 10.8 Quick Check: Retries
`quick-check-retries` · quiz · difficulty 8
- **Objective:** Reason about backoff, jitter and what to retry.
- **Scenario:** Nkem's questions before the first real webhook goes out.
- **Exercise:** Three questions.
- **Check:** Jitter keeps many senders from retrying in step; a 4xx other than 408 and 429 won't change by asking again; a receiver must expect the same event twice, because a lost answer means a retry.

#### 10.9 Same Key, Same Answer
`same-key-same-answer` · code_tests · difficulty 9
- **Objective:** Make a create request safe to retry: keep the answer under the caller's `Idempotency-Key` and give it again.
- **Scenario:** Furrowlight's program timed out waiting for a booking, retried, and got a 409 for the pass it had just booked itself.
- **Exercise:** `idempotent(keys *keyStore)` passes every request through; you fingerprint the body with SHA-256 and follow the bible's rules: replay with `Idempotent-Replayed: true`, 409 for another body or a request still running, 400 for a malformed key, 24 hours on the injected clock.
- **Check:** First request → 201 `bk-4501`; the same key and body → 201 `bk-4501` again, with the header and no new booking; the same key with another body → 409; key `!!` → 400; 24 hours and a second later the key is new; Brinewatch's `k1` and Cirrine's `k1` don't collide; a 500 is not kept; no key → the handler as before.

#### 10.10 Your Data Has Landed
`your-data-has-landed` · code_tests · difficulty 9
- **Objective:** Turn a landing into a signed, retried webhook for the right customer, end to end.
- **Scenario:** Friday 21 January: Furrow-A's image lands at Salar Quieto, and Furrowlight's receiver should know within a minute.
- **Exercise:** `(*Notifier).Landed(ctx context.Context, l Landing) error` returns nil; you find the booking in the store and its customer by satellite, build the `data.landed` event with the next `evt-` ID, and deliver it with the `Sender`.
- **Check:** Furrowlight's test receiver gets one event with the booking's fields and a valid signature; a receiver failing once gets it on attempt 2 with the same ID; an unknown booking → an error and nothing sent; IDs count `evt-000001`, `evt-000002`.

## Chapter 11 · Signed In

`signed-in` · 9 lessons

Browsers: a session cookie with `HttpOnly`, `Secure`, `SameSite`, `Path` and `MaxAge`, a server-side session behind it with an ID from `crypto/rand.Text` that ends 30 minutes after its last use, reading it back with `r.Cookie`, signing out on both sides, `http.CrossOriginProtection` against other sites' forms, CORS for one trusted origin, and the headers that tell a browser how to treat a page.

#### 11.1 The Cookie's Fine Print
`the-cookies-fine-print` · code_tests · difficulty 8
- **Objective:** Set a cookie with `http.SetCookie` and the attributes that limit it: `Path`, `MaxAge` in seconds, `HttpOnly`, `Secure` and `SameSite`.
- **Scenario:** Friday 21 January, 15:00 UTC: before a single cookie leaves the service, Mirela wants to know who can read it, where the browser sends it, and for how long.
- **Exercise:** `sessionCookie(id string) *http.Cookie` in `cookie.go` returns a bare `session` cookie; you add `Path: "/account"`, a `MaxAge` of `sessionLife` in seconds, `HttpOnly`, `Secure` and `SameSite: http.SameSiteLaxMode`.
- **Check:** Each field; `http.SetCookie` writes `session=<id>; Path=/account; Max-Age=1800; HttpOnly; Secure; SameSite=Lax`. The mistake: `MaxAge: 30`, minutes where seconds are counted.

#### 11.2 A Session Behind the Cookie
`a-session-behind-the-cookie` · code_tests · difficulty 9
- **Objective:** Keep a server-side session behind the cookie: an ID from `crypto/rand.Text`, an end on an injected clock, renewed by every use.
- **Scenario:** In the draft a session's ID is the customer's ID, and Nkem asks what the page shows when Hollis edits his cookie to say `cirrine`.
- **Exercise:** `(*Sessions).Start(customer string) string` returns the customer's ID and `Lookup(id string) (string, bool)` trusts any ID; you keep a `session` in `s.byID` under a `rand.Text()` ID, ending `sessionLife` after `s.now()`, and in `Lookup` renew a live session, delete an ended one and refuse unknown IDs, all under `s.mu`.
- **Check:** IDs of at least 26 characters from `A`–`Z` and `2`–`7`, two different; `Lookup("cirrine")`, `Lookup("furrowlight")` and `Lookup("")` → false; 29m59s after `Start` → live, exactly 30 minutes → ended and deleted; used at 15:50 and 16:15 → live, at 16:45 → ended; 50 browsers at once → 50 different IDs, each found. The mistakes: `now.After(end)`, which keeps a session alive at its last instant, and a renewal never stored back in the map.

#### 11.3 Signed In or Sent Back
`signed-in-or-sent-back` · code_tests · difficulty 8
- **Objective:** Read a cookie with `r.Cookie` and treat `http.ErrNoCookie` as a browser that isn't signed in.
- **Scenario:** Yusra signs in with Apogee's key to try the page on her test booking `bk-4784`, Pipit-1 over Port Ellery on Monday, and the page sends her back to the sign-in form.
- **Exercise:** `(*Account).signedIn(r *http.Request) (string, Customer, bool)` says nobody is signed in; you read the `session` cookie, look it up in `a.Sessions` and `a.Customers`, and return the cookie's value, the customer and true. The read-only `page` answers 303 to `/account/sign-in` on false and renews the cookie otherwise.
- **Check:** Through `page`: Yusra's cookie → 200 with `bk-4784` only and the cookie renewed under her ID; the same after a `units=km` cookie; no cookie, another cookie only, an unknown ID, `session=apogee`, an empty value and a session of a customer the page doesn't know → 303; Hollis 30 minutes after his last visit → 303. The mistakes: reading the whole `Cookie` header, a nil cookie (the test reports the panic), and returning the customer's ID where the session's belongs.

#### 11.4 Signing Out for Real
`signing-out-for-real` · code_tests · difficulty 8
- **Objective:** Sign out on both sides: end the session on the service, and clear the browser's cookie with `Max-Age=0` under the same name and path.
- **Scenario:** Mirela signs out, puts her old cookie back by hand and reloads: Apogee's bookings are back.
- **Exercise:** `(*Account).signOut` only redirects to `/account/sign-in`; before that you end the cookie's session with `a.Sessions.End` and send `sessionCookie("")` with `MaxAge: -1`.
- **Check:** 303 to `/account/sign-in`; `Lookup` of the old ID → false, and `page` with the old cookie pasted back → 303; a session cookie with `Path=/account`, an empty value and `Max-Age=0`; Yusra's and Hollis's sessions still live; no cookie → 303 and no panic. The mistakes: clearing only the cookie, a clearing cookie with `Path=/`, and the cookie set after `http.Redirect`.

#### 11.5 Forms From Other Sites
`forms-from-other-sites` · code_tests · difficulty 9
- **Objective:** Refuse cross-site request forgery with `http.CrossOriginProtection`, and see why a cookie, which the browser sends by itself, makes it possible.
- **Scenario:** Mirela's page on her laptop has one button, "See Furrow-A's new image", over a form that posts to `/account/bookings/bk-4776/cancel`: Furrowlight's booking of Furrow-B over Salar Quieto on Monday.
- **Exercise:** `noCrossOrigin(h http.Handler) http.Handler` returns `h`; you return it behind a new `http.CrossOriginProtection`.
- **Check:** Through `noCrossOrigin(a.routes())` with Hollis's cookie: a cross-site form, a same-site form and an older browser's form with a foreign `Origin` → 403 and `bk-4776` kept; the page's own form, an older browser's same-origin form and a program with neither header → 303 and the booking canceled; `GET /account/` from a link elsewhere → 200; the sign-in form posted from elsewhere → 403 and no cookie. The mistake: checking `Sec-Fetch-Site: cross-site` by hand, which misses the same site and older browsers.

#### 11.6 One Trusted Origin
`one-trusted-origin` · code_tests · difficulty 9
- **Objective:** Let one trusted origin call the API from a browser with CORS: answer its preflight `OPTIONS`, set `Access-Control-Allow-Origin`, `-Methods` and `-Headers`, and add `Vary: Origin` to every answer.
- **Scenario:** Yusra's TL;DR: "Furrowlight's planning page wants Monday's passes straight from the API." Its script at `https://plan.furrowlight.example` calls `GET /passes` with a token, and the browser hides every answer.
- **Exercise:** `cors(origin string) func(http.Handler) http.Handler` passes everything on; you add `Vary: Origin` and, for exactly `origin`, set `Access-Control-Allow-Origin`; its `OPTIONS` gets 204, `Access-Control-Allow-Methods: GET, POST, DELETE` and `Access-Control-Allow-Headers: Authorization, Content-Type, Idempotency-Key` without reaching `next`.
- **Check:** Against a stand-in API that refuses requests without a token: the planning page's preflight → 204 with the three headers and `Vary`, never reaching the API; its `GET /passes` → the API's 200 with `Access-Control-Allow-Origin`; `https://elsewhere.example`, `http://plan.furrowlight.example`, `https://plan.furrowlight.example.elsewhere.example` and `https://www.furrowlight.example` → no CORS header and the API's answers otherwise unchanged; a program without `Origin` → untouched but for `Vary`. The mistakes: `*`, echoing any `Origin`, a prefix match, and a preflight passed on to the API.

#### 11.7 Quick Check: Other Sites
`quick-check-other-sites` · quiz · difficulty 8
- **Objective:** Tell apart what guards a cookie's page and what CORS does.
- **Scenario:** Nkem's three questions before the page goes to the beta.
- **Exercise:** Three questions.
- **Check:** The browser sends a cookie by itself but a token only where a page's code puts it, so only the account page needs a guard against other sites' forms; CORS decides which pages' scripts may read answers, not which requests arrive; `Vary: Origin` keeps a cache from handing one origin's answer to another.

#### 11.8 Four Headers for the Page
`four-headers-for-the-page` · code_tests · difficulty 8
- **Objective:** Set a browser's security headers on every answer in middleware: `Content-Security-Policy` with `frame-ancestors`, `Strict-Transport-Security`, `X-Content-Type-Options` and `Referrer-Policy`.
- **Scenario:** Mirela's last list for the page: another site's frame, a guessed type, an address that leaks in a link, and a visit that starts over plain HTTP.
- **Exercise:** `securityHeaders(next http.Handler) http.Handler` returns `next`; before calling `next` you set `Content-Security-Policy: default-src 'self'; frame-ancestors 'none'`, `Strict-Transport-Security: max-age=31536000`, `X-Content-Type-Options: nosniff` and `Referrer-Policy: same-origin`.
- **Check:** Through `securityHeaders(noCrossOrigin(a.routes()))`: the signed-in page (200), the page signed out (303), the sign-in form (200), a wrong key (401), a form from another site (403) and an unknown address (404) all carry the four, with the policy's directives in any order. The mistake: setting them after `next`, when the answer has gone.

#### 11.9 Hollis Signs In
`hollis-signs-in` · code_tests · difficulty 9
- **Objective:** Put each guard where it belongs: the page's guards around the page and CORS around the API, in an order that gives refusals the headers too.
- **Scenario:** Friday 18:00 UTC: Hollis is waiting to sign in, and Yusra wants the page in the beta tonight.
- **Exercise:** `site(api http.Handler, a *Account, origin string) http.Handler` splits `/account/` from the API with a mux but guards nothing; you put the account routes behind `securityHeaders`, `noStore` and `noCrossOrigin`, in that order, and the API behind `cors(origin)`.
- **Check:** Against a stand-in API: Hollis signs in with Furrowlight's key and sees `bk-4776` and `bk-4783` but not `bk-4781` or `bk-4784`, and every page answer carries the four headers and `no-store`; Mirela's cross-site cancel → 403 with the headers, and the booking is kept; Yusra cancels her Pipit booking from the page; the planning page's preflight, `GET /passes` and cross-origin `POST /bookings` with a token → 204, 200 and 201; the account page never sends `Access-Control-Allow-Origin`; `PUT /passes` → the API's own 405; after signing out, the old cookie → 303. The mistakes: `noCrossOrigin` around the API (the planning page's booking gets 403), `cors` around the page, and the headers inside `noCrossOrigin` (its 403 lacks them).

## Chapter 12 · Launch Day

`launch-day` · 9 lessons

Stopping well: `http.Server.Shutdown`, `http.ErrServerClosed`, running until a context ends,
draining before shutdown, base contexts that reach every handler, an outbox that survives a
restart, and one test of the whole service.

#### 12.1 Stop Taking New Calls
`stop-taking-new-calls` · code_tests · difficulty 9
- **Objective:** Stop a server with `Shutdown`, which closes its listener and waits for requests in flight.
- **Scenario:** Ingrid's launch checklist has a line nobody has tried yet: "restart bookingd during a pass".
- **Exercise:** `stop(srv *http.Server, within time.Duration) error` calls `srv.Close()`; you call `Shutdown` with a context that ends after `within`.
- **Check:** A request in flight when `stop` starts finishes with 200, and `stop` returns nil after it; a new connection after `stop` is refused; a request that outlives `within` → `context.DeadlineExceeded`. The mistake: `Close`, which cuts the request in flight.

#### 12.2 Closed Is Not Failed
`closed-is-not-failed` · code_tests · difficulty 9
- **Objective:** Treat `http.ErrServerClosed` from `Serve` as a clean stop and anything else as a failure.
- **Scenario:** The rehearsal's log says `bookingd failed: http: Server closed` every time someone stopped it on purpose, and the pager woke Nkem twice.
- **Exercise:** `serve(srv *http.Server, ln net.Listener) error` returns whatever `Serve` returns; you return nil for `http.ErrServerClosed`, found with `errors.Is`.
- **Check:** `Serve` then `Shutdown` → nil; `Serve` on a listener already closed → that error. The mistake: comparing error strings.

#### 12.3 Run Until Told to Stop
`run-until-told-to-stop` · code_tests · difficulty 9
- **Objective:** Run a server until a context ends, then shut it down within a limit, from one function that both `main` and the tests call.
- **Scenario:** On launch day the deploy tool stops the service with a signal, and `main` turns it into a cancelled context with `signal.NotifyContext`.
- **Exercise:** `run(ctx, ln, app, t)` serves and returns at once; you serve in a goroutine, wait for `ctx.Done()` or a `Serve` failure, then shut down within `t.Shutdown` and return the result.
- **Check:** Requests answered while `ctx` lives; cancel → `run` returns nil within a second; a request in flight at cancel finishes; a failure of `Serve`, its listener closed underneath it, is returned without waiting for `ctx`.

#### 12.4 Quick Check: Shutdown
`quick-check-shutdown` · quiz · difficulty 9
- **Objective:** Reason about the order of a clean stop.
- **Scenario:** Ingrid's go/no-go questions for the launch rehearsal.
- **Exercise:** Three questions.
- **Check:** `Shutdown` refuses new connections at once and waits for the ones still answering; `Serve` returns `ErrServerClosed` as soon as `Shutdown` begins, so `main` must wait for `Shutdown` itself; `signal.NotifyContext` turns a signal into a cancelled context.

#### 12.5 Draining, Not Down
`draining-not-down` · code_tests · difficulty 9
- **Objective:** Announce that a server is going away before it stops: a health check that answers 503 while every other request still succeeds.
- **Scenario:** Yusra's guide tells customers' programs to check `/health` before sending a batch of bookings; during a restart they should wait, not start a batch that gets cut in half.
- **Exercise:** `run` shuts down as soon as `ctx` ends; you set the service's `draining` flag first, keep serving for `t.Drain`, then shut down, and `health` answers 503 `{"status":"draining"}` while the flag is set.
- **Check:** Before cancel, `/health` → 200; within a 300 ms drain after cancel, `/health` → 503 `draining` while `/stations` → 200; `run` returns after the drain and the shutdown. The mistake: setting the flag after `Shutdown`, when no request can see it.

#### 12.6 Contexts All the Way Down
`contexts-all-the-way-down` · code_tests · difficulty 10
- **Objective:** Give every handler's context a base that `run` cancels when shutdown runs out of time, so slow work stops instead of holding the old process.
- **Scenario:** In rehearsal, one `GET /passes` stuck behind a slow disk held the old process for the whole twenty seconds.
- **Exercise:** In `run`, handlers' contexts come from the server's default; you set `srv.BaseContext` to a context that `run` cancels when `Shutdown`'s deadline passes, then `Close` the server and return `shutdown: context deadline exceeded`.
- **Check:** A store blocked on its context ends soon after the shutdown limit, and `run` returns an error for which `errors.Is(err, context.DeadlineExceeded)`; the store saw its context cancelled; with no slow request, `run` returns nil. The mistake: cancelling the base context before `Shutdown`, which cuts requests that would have finished.

#### 12.7 The Outbox Survives
`the-outbox-survives` · code_tests · difficulty 9
- **Objective:** Keep pending work in a file saved on every change, and pick it up again after a restart.
- **Scenario:** Cirrine's receiver was down when the 13:40 image landed, and its retries must not die with the process at 14:05.
- **Exercise:** `Outbox` keeps deliveries in memory only; you save `outbox.json` through `replaceFile` on every change, and `OpenOutbox(path)` loads it and continues event IDs from `next_event`.
- **Check:** Two deliveries added, then reopened → both pending with their attempts and `next_at`; a delivered event leaves the file; after a reopen with nothing pending, the next ID still continues (`evt-000118` after 117); a missing file → empty, with next 1; a broken file → an error.

#### 12.8 One Test for Everything
`one-test-for-everything` · code_tests · difficulty 10
- **Objective:** Write one test that drives the whole service through `httptest`, the way a customer's program would.
- **Scenario:** Monday 31 January, the full rehearsal: Nkem wants one test that would have caught every bug of the season.
- **Exercise:** `checkService(t *testing.T, base string)` in `check.go` is empty; over HTTP you check `/health`, get tokens for Furrowlight and Brinewatch, list a week of passes, book one with an `Idempotency-Key` and send it again, check that Brinewatch can't see it, and send requests until a 429.
- **Check:** Your test passes against the assembled service (`newApp`, read-only) and fails, in a re-executed test binary, against each of four planted services: one that ignores idempotency keys, so the repeat gets 409 instead of the first answer, one that shows every customer's bookings, one that never limits, one that answers errors as plain text.

#### 12.9 Open for Business
`open-for-business` · code_tests · difficulty 10
- **Objective:** Restart a whole service in the middle of traffic without dropping a request, a booking or an event.
- **Scenario:** Tuesday 1 February, 14:05 UTC: Furrow-A is over Salar Quieto, customers are booking, and one wrong word in the guide forces a restart.
- **Exercise:** `run` stops the server but leaves the notifier running and never resumes the outbox; you stop in the bible's order (drain, shut down, stop the notifier, which saves the outbox) and start the notifier from the outbox.
- **Check:** The first process: bookings made, a request in flight at the stop that finishes, Cirrine's `evt-000117` saved with its attempts. The second process, from the same folder: the bookings intact and the next ID continuing, Cirrine's event delivered with the same ID, and the 14:13:02 landing delivered to Furrowlight on its first attempt as `evt-000118`. The whole run takes a few seconds at most.

## Chapter 13 · The First Week Open

`the-first-week-open` · 9 lessons

The first customers' week: file uploads read from `multipart` forms within a limit and sniffed
with `http.DetectContentType`, ranged downloads that resume through `http.ServeContent`, download
links signed with HMAC that expire, webhooks refused inside the network by a `net.Dialer`'s
`Control` hook, redirects that stay on the site, and a monthly plan on an injected clock.

**Story:** Wednesday 2 to Friday 4 February 2028, the first week open, at the backend desks, between Launch Day and the Board project. Real customers do what the beta never did. Cirrine attaches Cirrine-1's 1.9 GB software image as a command file and Furrowlight sends a PDF; Hollis's download of Tuesday's 10.9 GB image breaks four times, and his field team needs a link without a key; Tomás, on call for launch week, finds Brinewatch's events in his weather service's log; Mirela writes a phishing email of her own through Yusra's `/go` links; Brinewatch reaches 11 of the 12 bookings its plan allows in February. Friday's 17:00 deploy ships the week.

#### 13.1 Forms That Carry Files
`forms-that-carry-files` · code_tests · difficulty 9
- **Objective:** Read one file out of a `multipart/form-data` body with `r.MultipartReader()` and `NextPart`, behind `http.MaxBytesReader`.
- **Scenario:** Wednesday 2 February, 09:10: Cirrine's engineer attached Cirrine-1's whole software image, 1.9 GB, to Thursday's booking instead of its command file, and the draft read every byte into memory.
- **Exercise:** `readCommands(w http.ResponseWriter, r *http.Request) ([]byte, error)` returns the whole body; you cap the body at `maxBodyBytes` and return the bytes of the part named `file`, errors as they come, and `http.ErrMissingFile` when no part has that name.
- **Check:** Cirrine's form, a `note` part first → exactly the command file; a 65,000-byte file → all of it; a text/plain body → `http.ErrNotMultipart`; a form with only a note → `http.ErrMissingFile`; the 1.9 GB form, made as it is read → a `*http.MaxBytesError` with at most 65,537 bytes read. `r.FormFile` passes too. The mistakes: the limit set after `MultipartReader`, which reads the body it found, and an `io.LimitReader` on the part, which stops without saying the body was too large.

#### 13.2 The Bytes, Not the Label
`the-bytes-not-the-label` · code_tests · difficulty 9
- **Objective:** Judge an upload by its bytes with `http.DetectContentType`, not by the `Content-Type` its sender wrote.
- **Scenario:** Wednesday, 15:00: Furrowlight's command file for Friday's 14:36 pass over Salar Quieto (`bk-4921`) is a PDF of the command list, labeled `text/plain` like everything its script sends.
- **Exercise:** `checkCommands(data []byte) error` accepts every file; you return a 415 `unsupported_media_type`, `command file must be plain text, not <type>`, unless the sniffed type is exactly `text/plain; charset=utf-8`. The read-only `upload` handler reads, checks and keeps the file.
- **Check:** A command list, and one with a UTF-8 byte-order mark → nil; a PDF, a web page, XML, a zip, a PNG, UTF-16 text and raw bytes → 415 naming each sniffed type; through `upload`, a PDF labeled `text/plain` → 415 and nothing kept, the list labeled `application/octet-stream` → 204 and kept. The mistakes: `text/` as a prefix (HTML and XML get through), a comparison with `text/plain` alone (every list refused), `text/plain` as a prefix (UTF-16 gets through).

#### 13.3 Pick Up Where It Stopped
`pick-up-where-it-stopped` · code_tests · difficulty 9
- **Objective:** Serve a file with `http.ServeContent`, which answers `Range` with 206 and `Content-Range`, and `If-Range` by the file's modification time.
- **Scenario:** Thursday 3 February, 08:10. Hollis: "FURROW-A 0201 1402: FOUR TRIES, NEVER WHOLE." Every retry of Tuesday's 10,912,500,000-byte image got the whole file again.
- **Exercise:** `sendFile(w http.ResponseWriter, r *http.Request, file *os.File)` copies the whole file with `io.Copy`; you answer with `http.ServeContent` and the name and modification time from `file.Stat()`. The read-only `download` handler finds the caller's landed file.
- **Check:** Against a 64 KiB stand-in landed at 14:13:02: the whole file with `Accept-Ranges: bytes` and `Last-Modified`; `Range: bytes=40000-` → 206, `Content-Range: bytes 40000-65535/65536`, and the two parts make the image; `bytes=65536-` → 416; `If-Range` with the image's own date → 206, with an older one → 200 and the whole file; a booking whose pass hasn't flown → 404 `no file for bk-4921 yet`. The mistake: a zero or wall-clock modification time, which never matches `If-Range`.

#### 13.4 A Link That Runs Out
`a-link-that-runs-out` · code_tests · difficulty 9
- **Objective:** Check a download link signed with HMAC-SHA256 over its path and expiry: the signature in constant time with `hmac.Equal`, then the expiry.
- **Scenario:** Thursday, 11:00: Hollis wants Furrowlight's field team, agronomists with laptops and no API key, to fetch Tuesday's image; a link from `POST /bookings/{id}/link` works for 24 hours.
- **Exercise:** `checkLink(path, expires, sig string, key []byte, now time.Time) error` accepts any link; you recompute `linkSig`, compare in constant time, and return `errBadLink` for a changed or malformed link and `errLinkExpired` from the expiry's second on.
- **Check:** Hollis's link on Thursday and one second before expiry → nil; at its expiry and on Saturday → `link expired`; its expiry moved on (also when read after the real expiry), bk-4912's signature on bk-4913's path, another key, no or non-hex signature, a non-numeric or missing expiry → `invalid link`; `go/parser` finds `hmac.Equal` (or `subtle.ConstantTimeCompare`) in `checkLink`; through the handlers, the minted `/files/bk-4912?expires=1833274800&sig=…` downloads whole and from byte 40000, and is refused when edited and on Friday at 11:00.

#### 13.5 Quick Check: Files and Links
`quick-check-files-and-links` · quiz · difficulty 9
- **Objective:** Reason about resumes, body limits and what a link's signature covers.
- **Scenario:** Thursday afternoon: Nkem's three questions before the field team gets its first link.
- **Exercise:** Three questions.
- **Check:** A 200 to a resume is the whole file, so the client starts again from it; `http.MaxBytesReader` ends the read at 64 KiB with an error the service answers with 413; an expiry outside the signature could be moved on by anyone holding the link.

#### 13.6 Webhooks Stay Outside
`webhooks-stay-outside` · code_tests · difficulty 10
- **Objective:** Refuse to connect to loopback, private, link-local and unspecified addresses from a `net.Dialer`'s `Control` hook, which sees the address actually dialed.
- **Scenario:** Thursday, 15:00: Tomás, on call this week, finds Brinewatch's events in his weather service's log; Brinewatch gave its receiver's inside address, `http://10.20.0.31/apogee-events`, which inside Apogee's network is the weather service.
- **Exercise:** `checkDial(network, address string, _ syscall.RawConn) error` allows everything; you parse the address with `netip.ParseAddrPort` and return an error wrapping `errInside` for an inside address or one that doesn't parse. The read-only `webhookClient` installs it.
- **Check:** 203.0.113.7, 198.51.100.20, 172.32.0.1 and 2001:db8::7 → nil; 10.20.0.31, 172.16.5.4, 172.31.255.254, 192.168.1.20, 127.0.0.1, 127.8.9.10, ::1, 169.254.169.254, fe80::1, fd12:3456::7, 0.0.0.0, :: and `weather-service:80` → `errInside`; `webhookClient` posting to a test server at 127.0.0.1 and to `localhost` (resolved from the Mac's hosts file) → `errInside`, and the server gets no request. The mistakes: string prefixes (172.32.0.1 is outside), IPv6 forgotten, 0.0.0.0, which reaches this very machine, and a parse error let through.

#### 13.7 Only Our Own Pages
`only-our-own-pages` · code_tests · difficulty 9
- **Objective:** Keep a redirect's `next` to a path on the same site: parsed with `url.Parse`, starting with one `/`, with no `\`.
- **Scenario:** Friday 4 February, 08:30: Mirela forwards an email she wrote herself whose link, `/go?next=https://apogee-bookings.example/renew`, starts with the booking service's own address and ends on a page that asks for the API key.
- **Exercise:** `localPath(next string) bool` accepts everything; you accept only a `next` that parses, starts with `/` but not `//`, and holds no backslash. The read-only `goNext` redirects to it with 303, or to `/guide/`, and logs where it sent the browser.
- **Check:** `/guide/files.html`, a page with a query and a fragment, and `/account/` → themselves; an empty `next`, Mirela's link, `//…`, `///…`, `/\…`, `\\…`, a tab or a line break between the slashes, `https:/…`, `javascript:` and `guide/files.html` → `/guide/`; no log line names Mirela's site. The mistakes: a leading `/` alone, and an empty host after parsing, which lets `/\…` and `https:/…` through.

#### 13.8 A Plan for the Month
`a-plan-for-the-month` · code_tests · difficulty 9
- **Objective:** Count usage per customer per calendar month on an injected clock, refuse past the plan with a clear 403, and start again at the month boundary in UTC.
- **Scenario:** Friday, 11:00. Yusra: "TL;DR: Brinewatch has made 11 of the 12 bookings its plan allows in February, and asks what number 13 gets."
- **Exercise:** `(*usage).take(customer string, plan int, now time.Time) error` counts nothing; you start a new count when `now`, in UTC, is in another month than `u.Month`, refuse a used-up plan with 403 `quota_exceeded`, `plan of 12 bookings a month used up: the count starts again at 2028-03-01T00:00:00Z`, and otherwise count. The read-only `Quota` keeps the count in `usage.json` and skips customers without a plan.
- **Check:** Brinewatch's 12th counted and its 13th refused; at 2028-02-29T23:59:59Z refused, at 2028-03-01T00:00:00Z counted from 1; Cirrine's 41st refused with its own plan, and counted beside Brinewatch's full plan; a clock at 21:30 UTC−3 on 29 February → March; February 2029 → a new count; December's refusal names 2029-01-01; through `Quota`, the 13th still refused after reopening the file, Apogee's 100 bookings never counted, and 50 bookings at once with 30 of 40 made → exactly 10. The mistakes: `>` for `>=`, a month without its year, the clock's own zone, a reset date from `now.AddDate`.

#### 13.9 One Week Open
`one-week-open` · code_tests · difficulty 10
- **Objective:** Put each of the week's features behind the guard it needs, beside the API and inside its middleware.
- **Scenario:** Friday, 16:00: the week's work goes out in tonight's deploy, and Nkem wants it served as one handler.
- **Exercise:** `firstWeek(app *App, f *Files, q *Quota, l *log.Logger) http.Handler` registers the new routes unguarded inside the service's middleware; you put the upload, download and link routes behind `s.private`, wrap `s.Store` in `metered` with `q`, and give the notifier's `Sender` the client from `webhookClient()`.
- **Check:** Through `httptest.NewServer` at Friday 16:00 with the assembled service: Furrowlight's command file for `bk-4933` → 204 and kept, without a token → 401, a PDF → 415; `bk-4912`'s file from byte 60000 → 206; a link until Saturday 16:00, fetched with no token → 200 and with a range → 206, its expiry moved on → 403; `/go` to Mirela's site → `/guide/`; Brinewatch's 12th February booking → 201 `bk-4937`, its 13th → 403 `quota_exceeded`, and `usage.json` at 12; a delivery to Brinewatch's old inside address → `ErrUndelivered` wrapping `errInside`, with no request reaching the stand-in weather service; `/stations` still answered, and the upload in the access log.

## Changes in writing and review

- **Booking IDs** start 1,500 higher than first planned, so they keep growing after Season 2's:
  Yusra's spreadsheet ends at `bk-4498`, the service's first booking is `bk-4501`, and launch
  day is near `bk-4912`.
- **Passes** all come from one orbit model. Three catalog passes were longer than their orbits
  allow and now end earlier (`pipit-1/PEL/08154` at 04:21:27, `cirrine-1/YST/11364` at
  07:40:04, `brine-3/LAN/09213` at 09:22:02). Passes and events added while writing chapters 6,
  7, 10 and 11 moved to real passes, with file sizes each pass's link can carry.
- **Given code** keeps the text the learner last saw, lesson after lesson. Every change to a
  given declaration is named by a readme or the story: the fields and routes each chapter's
  feature adds, and 6.8's planted orbit-number bug, which chapter 6's debrief says Nkem fixed.
- **Chapter 5** `invalidBooking(msg string) *APIError` lives in `errors.go` for the whole
  season. From 5.2 every handler error goes through `writeError`; 5.1 still answers with
  `http.Error`.
- **Chapter 6** `logRequests` logs 200 for a handler that writes nothing, and every error body
  carries `request_id` from 6.2 on.
- **Chapter 8** checks that the caller operates the satellite before asking the catalog, so a
  real pass of another customer's satellite and a missing one both answer 403 and give nothing
  away. From chapter 8 on, booking `thistle` answers 403.
- **Chapter 9** A zero timeout leaves its limit unset; a zero `ReadHeader` or `Idle` falls back
  to `Read`, as Go's server does.
- **Chapter 10** `post` returns an `*AnswerError` for an answer outside 200–299 and `statusOf`
  reads its status; `Sender` holds the customers. 10.7 also tests a stop that comes while an
  attempt is in flight, and 10.9 asks for `claim` rather than the whole `idempotent`
  middleware, which runs about forty lines.
- **Chapter 12** Cirrine's launch-day event comes from Cirrine-1's 13:30 pass over Ystmark,
  whose file lands at 14:04:57 over Ystmark's slow backhaul. 11.6's stuck request is a
  `POST /bookings` saving behind the slow disk, the one request that touches it. The notifier
  keeps or drops a delivery by its outcome, `Notifier.Start` and `Stop` and the outbox's
  methods are part of the API, and 11.8's request helper is called `request`.
- **Spelling** is American throughout, `canceled` included.
- **Chapter 11** was added after the season shipped. It takes Friday 21 January's afternoon, the only story time between chapter 10's Friday and Launch Day's Monday, and adds new declarations only: no given text of another chapter changes, and Launch Day's lessons leave the account page out of `newApp`.
- **Chapter 13** was added after the season shipped, as its epilogue: Wednesday 2 to Friday 4 February 2028, after Launch Day's finale and before the Board project. It adds new declarations and files only; no given text of another chapter changes, and 13.9 wraps Launch Day's `newApp` instead of changing it. The plan counts bookings made in a calendar month and keeps them in `usage.json`, so a restart forgets nothing.
