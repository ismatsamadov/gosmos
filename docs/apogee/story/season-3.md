# Apogee story bible: Season 3

### Season 3: open for business

After the gctl project week, the learner starts the next assignment on Monday 1 November
2027. Their new badge says **Web Services Engineer**, still No. 039, based at Port Ellery
headquarters, since Monday 1 November 2027. Customers book passes by email: Yusra copies about
forty bookings a day into a spreadsheet by hand, and in August one typo booked Furrow-A over
the wrong station. On Tuesday 1 February 2028 Apogee opens a booking service to every customer
at once: their programs ask it for passes, book them, and hear back when their data has
landed, day and night. You build it with Go's standard library. Nkem reviews every change,
Yusra brings the customers' wishes, Mirela reviews who may do what, and Ingrid runs launch
day like a pass.

Story time runs from Monday 1 November 2027 to Friday 4 February 2028: thirteen chapters,
all but chapters 11 and 13 of about a week each (chapter 11 is the afternoon of Friday 21 January,
chapter 13 the first days open, Wednesday 2 to Friday 4 February), with a
quiet holiday fortnight (Monday 20 December to Sunday 2 January)
between chapters 7 and 8. No new people, places, satellites or customers are introduced.
Customers speak through Yusra ("TL;DR: ...") and Furrowlight through Hollis's capitals.
Thistle re-entered in July 2027, so the service's fleet has seven satellites and Kessick
Polytechnic books nothing this season.

#### The season, chapter by chapter

Everything happens at Port Ellery headquarters. Customers try the service before launch ("the
beta"): Cirrine from chapter 3, Furrowlight from chapter 4, Brinewatch from chapter 5 (its
client goes wrong in chapter 9). Every change is tried first on a booking of a Pipit, as Apogee
tests everything.

| Ch | Week of | Where | Who needs what | What is at stake |
|---|---|---|---|---|
| 1 | Mon 1 Nov 2027 | the backend team's desks beside the Fishbowl | Nkem: before building anything, learn what a customer's program will send and get back, on a stand-in service she wrote | Yusra types forty emailed bookings a day; a wrong word in a program books the wrong pass |
| 2 | Mon 8 Nov | the Fishbowl, at the go/no-go poll | Ingrid: the go/no-go poll's wind check must answer within two seconds or say why not | On Tuesday 2 November the poll waited four minutes on the weather service while Cape Wren's mast was silent, and Furrow-B's checklist finished after AOS |
| 3 | Mon 15 Nov | the backend desks | Yusra: the booking service's first answer, `GET /stations`, for Cirrine's engineers, who start their client this week; her customer guide served beside it | The first thing a customer ever asks Apogee's own service. On Friday 19 November Yusra sends Cirrine the address |
| 4 | Mon 22 Nov | the backend desks | Hollis, through Yusra: a link to every pass for his emails, and "TL;DR: Furrowlight wants to book a week at a time" | A link that points at the wrong pass books the wrong pass |
| 5 | Mon 29 Nov | the backend desks | Hollis: Furrowlight's program books its first pass, on Friday 3 December, for Furrow-A over Salar Quieto on Monday 6 December, AOS 14:02:10 | A booking the service accepts must be a pass that exists. It becomes `bk-4501`, the first booking that never touches Yusra's spreadsheet |
| 6 | Mon 6 Dec | the backend desks | Nkem: on Monday a panic in a handler answered Cirrine's program with a dropped connection and left no trace | Every request gets an answer, an ID and a line in the log, as every failure at Apogee gets a name, a reason and a line in the log |
| 7 | Mon 13 Dec | the backend desks | Yusra: on Monday a deploy restarted the beta and Furrowlight's twelve test bookings vanished | Bookings must survive a restart before the holiday fortnight. On Friday 17 December the last deploy of the year goes out |
| 8 | Mon 3 Jan 2028 | Mirela's desk, then the backend desks | Mirela's review: anyone who can reach the service can read and cancel any booking | "Brinewatch must never see Cirrine's passes." She signs off on Friday 7 January, "LGTM once it fails closed" |
| 9 | Mon 10 Jan | the backend desks | Brinewatch's beta client lost its network on Sunday night and has retried ever since with no pause: a request every millisecond | One customer's bug must not slow everyone else; Cirrine's bookings wait half a minute behind the flood |
| 10 | Mon 17 Jan | the backend desks | Yusra: customers ask `GET /bookings` every few seconds to learn whether their data has landed; Furrowlight asks every five seconds, 17,280 times a day | "Your data has landed," a second after the file lands, signed so nobody else can say it |
| 11 | Fri 21 Jan, 15:00 to 18:20 UTC | the backend desks | Hollis, right after his capitals about the five-second loop: he wants to see Furrowlight's bookings himself, in a browser; Yusra wants a signed-in page beside the API before rehearsal week, and Mirela reviews it | A signed-in browser can be made to cancel a booking from any other website. At 18:20 Hollis signs in and sees only Furrowlight's bookings; Mirela's page with a hidden cancel form is refused |
| 12 | Mon 24 Jan, then Tue 1 Feb | the Fishbowl | Ingrid runs launch like a pass: a rehearsal week, a full rehearsal on Monday 31 January, then launch at 14:00 UTC on Tuesday 1 February | At 14:01 Yusra finds that the customer guide's sample request says `Authorization: Token` instead of `Bearer`, and customers are copying it. The guide is built into the service, so it must restart at 14:05, while Furrow-A is over Salar Quieto (AOS 14:02:10, LOS 14:11:52) |
| 13 | Wed 2 to Fri 4 Feb | the backend desks | The first week open: Cirrine's 1.9 GB software image arrives as a command file (Wed 09:10) and Furrowlight's is a PDF (Wed 15:00); Hollis's download of Tuesday's image breaks four times (Thu 08:10) and his field team needs a link without a key (Thu 11:00); Tomás finds Brinewatch's events in the weather service's log (Thu 15:00); Mirela's own phishing email goes through `/go` (Fri 08:30); Brinewatch is at 11 of its 12 February bookings (Fri 11:00) | What the service does in its first week, customers will rely on. Friday's 17:00 deploy ships the week |

The opening (Monday 1 November) hands over the new badge, shows Yusra's spreadsheet and the
launch date, introduces Nkem (reviewer), Yusra (the customers' voice) and Mirela (who will
review the service in January), and sends the learner to Nkem's stand-in.

The finale (chapter 12, Tuesday 1 February): the restart at 14:05 drops no request; the
bookings and the pending webhooks survive it; Cirrine's event from its 13:30 pass, still being retried
when the old process stopped, is delivered by the new one; Furrow-A's image lands at 14:13:02
and Furrowlight's receiver hears about it at 14:13:03. By evening every customer has booked
through the service and Yusra archives her spreadsheet. The debrief points at the Board project
(the next passes on every screen in the building) and at Season 4: Sunniva's new dish
controllers at Ystmark arrive in March, and they speak only raw sockets.

The epilogue (chapter 13, Wednesday 2 to Friday 4 February) is the first week open. Each lesson is a
surprise from real customers, and Friday's 17:00 deploy ships everything the week added. Its debrief
ends on Monday, when Ingrid has one more thing to ask: the Board project.

#### New systems

- **bookingd**, the booking service, which this season builds: add it to the Services list
  ("bookingd (books passes for customers and tells them when their data lands)"). In readmes
  it is "the booking service" first, `bookingd` once it has been glossed.
- **The weather service**, an internal service Tomás wrote in 2024. It reads each station's
  weather mast (the anemometer and a thermometer) every minute and answers
  `GET /weather/{code}`. It is old and quirky: its errors are plain text, a silent mast answers
  503, and when a station's backhaul is slow it can stall for minutes. Chapter 2 writes its
  client; the Board project reuses that client.
- **Nkem's stand-in** (chapter 1 only): a few lines of Go that answer like the booking service
  will, about stations only, so the learner can study HTTP before the service exists.
- **The customer guide**: Yusra's pages for customers' engineers, served by the service at
  `/guide/` and built into its binary.
- **The catalog**: passplan's predicted passes for the next 14 days, which the service offers
  for booking. A customer can book only a pass in the catalog.
- **Landings**: dlinkd tells the booking service when a booking's file has been delivered
  (a `Landing`), and the service tells the customer with a webhook.
- **The account page** (chapter 11): a small signed-in page that bookingd serves beside the API under
  `/account/`, where a customer's people sign in with a browser using their company's API key, see
  their own bookings and cancel one. Its sessions live only in memory, so a deploy or restart signs
  every browser out and nothing else is lost; the launch-day restart signs page users out too. In
  readmes: "the account page".
- **Furrowlight's planning page**: its engineers' own web page at `https://plan.furrowlight.example`
  (`.example` is reserved), whose script calls the API with a token that Furrowlight's server got. It
  is the one origin the API trusts for CORS (chapter 11).
- Customers reach the booking service only over HTTPS, which no lesson sets up; the cookie's `Secure`
  and the page's `Strict-Transport-Security` rely on it, and lessons use plain HTTP on 127.0.0.1.

#### The code, by chapter

Each lesson is a program of its own, so the read-only file a name sits in is the lesson's
choice (`data.go` and `types.go` are fine homes for fixtures and shared types); the names,
types, signatures, routes and messages below are fixed.

Each lesson is its own `package main` (module `lesson`, as the app writes it). A lesson gives
the pieces of earlier chapters that it needs as read-only files, under these names, and uses
only the fields it needs. In chapters 4 to 6 the store (`store.go`) is a hidden file, so the
learner first meets `MemStore`'s inside in chapter 7, where they write it.

| Chapter | Names it adds (and files) |
|---|---|
| 1 | `standIn() http.Handler` (Nkem's, read-only `standin.go`); `printRequest`, `requestLine`, `printResponse` (read-only `wire.go`); `passesPath`, `advice`, `isJSON`, `fetch` |
| 2 | `Weather`, `WeatherClient`, `newWeatherClient`, `(*WeatherClient).Reading`, `decodeWeather`, `StatusError`, `windCall` (`weather.go`) |
| 3 | `Service`, `(*Service).routes`, `listStations`, `health`, `writeJSON`, the guide route; `start` (`service.go`, `stations.go`) |
| 4 | `getStation`, `listPasses`, `getPass`, `getBooking`, `cancelBooking`; `PassView`; `MemStore` (hidden `store.go`) |
| 5 | `createBooking`, `APIError`, `ErrorBody`, `writeError`, `decodeJSON`, `bodyError`, `maxBodyBytes`, `checkShape`, `checkPass`, `invalidBooking` (in `errors.go`), `knownSat`, `knownStation`, `window`, `findStation` (`errors.go`, `decode.go`, `rules.go`, `data.go`) |
| 6 | `noStore`, `requestID`, `requestIDFrom`, `statusRecorder`, `logRequests`, `recoverPanics`, `chain`, `(*Service).handler` (`middleware.go`) |
| 7 | `Store`, `ErrNotFound`, `ErrConflict`, `MemStore`, `NewMemStore`, `FileStore`, `OpenFileStore`, `firstNumber`, `(*MemStore).load`, `replaceFile`, `listBookings`, `Config`, `newService` (`store.go`, `filestore.go`) |
| 8 | `Customer`, `credentials`, `checkKey`, `signToken`, `verifyToken`, `issueToken`, `TokenResponse`, `authenticate`, `withCustomer`, `customerFrom`, `owns`, `(*Service).private`, `unauthorized`, `tokenMAC`, `findBooking` (`auth.go`, `bookings.go`) |
| 9 | `Bucket`, `Limiter`, `newLimiter`, `limit`, `retryDelay`, `Timeouts`, `newServer`, `withDeadline`, `bucketSize`, `refillPerSec`, `errRateLimited`, `(*Bucket).Take` (`limit.go`, `server.go`) |
| 10 | `Event`, `Landing`, `post`, `AnswerError`, `statusOf`, `signature`, `verifyEvent`, `backoff`, `retryable`, `Sender`, `Delivery`, `ErrUndelivered`, `Notifier`, `idempotent`, `keyStore`, `claim`, `validKey`, `replay`, `answerRecorder` (`webhook.go`, `idempotency.go`) |
| 11 | `cookieName`, `sessionLife`, `sessionCookie` (`cookie.go`); `session`, `Sessions`, `newSessions`, `(*Sessions).Start`, `Lookup`, `End` (`sessions.go`); `Account`, `(*Account).routes`, `signedIn`, `page`, `signInForm`, `signIn`, `signOut`, `cancel` (`account.go`); `pageHTML`, `signInHTML`, `render`, `renderPage`, `renderSignIn` (`page.go`); `noCrossOrigin`, `cors`, `securityHeaders` (`browser.go`); `site` (11.9's `main.go`) |
| 12 | `stop`, `serve`, `run`, `Outbox`, `OpenOutbox`, `(*Outbox).NewEventID`, `Put`, `Remove`, `Pending`, `outboxFile`, `(*Notifier).Start`, `Stop`, `(*Notifier).deliver`, `sleep`, `App`, `newApp`; 12.8's `client.go` (`request`, `answer`, `errorCode`, with chapter 8's `send` and `token`) (`run.go`, `outbox.go`, `webhook.go`, `app.go`) |
| 13 | `readCommands`, `uploadError` (`commands.go`); `checkCommands` (`sniff.go`); `Files`, `(*Files).booking`, `upload`, `landed`, `download`, `link`, `shared` (`files.go`); `sendFile` (`send.go`); `linkLife`, `Link`, `errBadLink`, `errLinkExpired`, `linkSig`, `linkFor`, `checkLink` (`link.go`); `errInside`, `checkDial`, `webhookClient` (`dial.go`); `localPath`, `goNext` (`next.go`); `usage`, `(*usage).take` (`usage.go`); `Quota`, `openQuota`, `(*Quota).Take`, `save` (`quota.go`); `metered`, `(*Quota).Give` (`metered.go`); `firstWeek` (13.9's `main.go`) |

Signatures every chapter after the first one that defines them uses exactly:

```go
// chapters 1 and 2
func passesPath(sat, station, from string) string
func advice(code int) string
func isJSON(h http.Header) bool
func fetch(url string) (status, body string, err error)
func newWeatherClient(baseURL string, timeout time.Duration) *WeatherClient
func (c *WeatherClient) Reading(ctx context.Context, code string) (Weather, error)
func decodeWeather(r io.Reader) (Weather, error)
func windCall(ctx context.Context, c *WeatherClient, code string, now time.Time) string

// the booking service
func start(s *Service) (addr string, stop func() error, err error) // chapter 3 only
func writeJSON(w http.ResponseWriter, status int, v any)
func writeError(w http.ResponseWriter, r *http.Request, err error)
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error
func checkShape(b Booking) error
func checkPass(b Booking, catalog []PassView, now time.Time) error
func (s *Service) routes() *http.ServeMux
func (s *Service) handler(l *log.Logger) http.Handler // the mux inside every middleware
func (s *Service) private(h http.HandlerFunc) http.Handler // authenticate, then limit
func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler
func noStore(next http.Handler) http.Handler
func requestID(next http.Handler) http.Handler
func requestIDFrom(ctx context.Context) string
func logRequests(l *log.Logger, now func() time.Time) func(http.Handler) http.Handler
func recoverPanics(l *log.Logger) func(http.Handler) http.Handler
func credentials(r *http.Request, scheme string) (string, error)
func checkKey(customers map[string]Customer, key string) (Customer, bool)
func signToken(customer string, expires time.Time, key []byte) string
func verifyToken(tok string, key []byte, now time.Time) (string, error)
func authenticate(customers map[string]Customer, key []byte, now func() time.Time) func(http.Handler) http.Handler
func withCustomer(ctx context.Context, c Customer) context.Context
func customerFrom(ctx context.Context) (Customer, bool)
func owns(c Customer, satID string) bool
func limit(l *Limiter) func(http.Handler) http.Handler
func withDeadline(d time.Duration) func(http.Handler) http.Handler
func newServer(h http.Handler, t Timeouts) *http.Server
func retryDelay(res *http.Response) time.Duration
func post(ctx context.Context, c *http.Client, url string, secret []byte, ts int64, ev Event) error // a non-2xx answer is an *AnswerError
func statusOf(err error) int // the status an *AnswerError in err carries, or 0 when no answer came
func signature(secret []byte, ts int64, body []byte) string
func verifyEvent(r *http.Request, secret []byte, now time.Time) (Event, error)
func backoff(attempt int, r *rand.Rand) time.Duration
func retryable(status int, err error) bool
func (s *Sender) Deliver(ctx context.Context, d *Delivery) error
func (n *Notifier) Landed(ctx context.Context, l Landing) error
func (n *Notifier) Start()      // chapter 12: resumes the outbox's deliveries, each when due
func (n *Notifier) Stop() error // chapter 12: ends deliveries in progress, keeping their attempts
func newKeyStore(now func() time.Time) *keyStore
func idempotent(keys *keyStore) func(http.Handler) http.Handler
func newService(cfg Config) (*Service, error)
func stop(srv *http.Server, within time.Duration) error
func serve(srv *http.Server, ln net.Listener) error
func OpenOutbox(path string) (*Outbox, error)
func (o *Outbox) NewEventID() (string, error) // "evt-000118", once the count is saved
func (o *Outbox) Put(d Delivery) error         // adds or replaces d's delivery, and saves
func (o *Outbox) Remove(id string) error       // drops event id's delivery, and saves
func (o *Outbox) Pending() []Delivery          // a copy of the deliveries still to make
func newApp(cfg Config) (*App, error)
func run(ctx context.Context, ln net.Listener, app *App, t Timeouts) error

// chapter 11: the account page
func sessionCookie(id string) *http.Cookie
func newSessions(now func() time.Time) *Sessions
func (s *Sessions) Start(customer string) string     // a rand.Text() ID, ending sessionLife after now
func (s *Sessions) Lookup(id string) (string, bool) // its customer; renews it; an ended one is deleted
func (s *Sessions) End(id string)
func (a *Account) routes() *http.ServeMux
func (a *Account) signedIn(r *http.Request) (string, Customer, bool) // session ID, customer, ok
func noCrossOrigin(h http.Handler) http.Handler
func cors(origin string) func(http.Handler) http.Handler
func securityHeaders(next http.Handler) http.Handler
func site(api http.Handler, a *Account, origin string) http.Handler

// chapter 13: the first week open
func readCommands(w http.ResponseWriter, r *http.Request) ([]byte, error) // the form's "file" part, the body capped at maxBodyBytes
func checkCommands(data []byte) error                                     // 415 unless sniffed as exactly text/plain; charset=utf-8
func sendFile(w http.ResponseWriter, r *http.Request, file *os.File)    // http.ServeContent with the file's name and ModTime
func linkSig(path string, expires int64, key []byte) string             // hex HMAC-SHA256 of "<path>?expires=<unix>"
func linkFor(path string, expires time.Time, key []byte) string
func checkLink(path, expires, sig string, key []byte, now time.Time) error // errBadLink, then errLinkExpired
func checkDial(network, address string, c syscall.RawConn) error        // a net.Dialer's Control hook
func webhookClient() *http.Client
func localPath(next string) bool
func goNext(l *log.Logger) http.Handler
func (u *usage) take(customer string, plan int, now time.Time) error
func openQuota(path string, plans map[string]int, now func() time.Time) (*Quota, error)
func (q *Quota) Take(customer string) error
func (q *Quota) Give(customer string) error
func firstWeek(app *App, f *Files, q *Quota, l *log.Logger) http.Handler

// checks the learner writes, in editable non-test files
func checkStore(t *testing.T, open func() Store)                             // chapter 7, contract.go
func review(t *testing.T, base string, advance func(time.Duration))          // chapter 8, review.go
func checkService(t *testing.T, base string)                                 // chapter 12, check.go
```

`App` (chapter 12) is the assembled service: `Service *Service`, `Notifier *Notifier` and
`Handler http.Handler`, the mux inside every middleware. `Notifier` holds the `Store`, the
customers, the `Sender` and (chapter 12) the `Outbox`, and counts event IDs.

#### The service's types

The canonical types keep their fields and types. On the wire, Station and Pass carry these
tags; Booking is Season 2's, unchanged.

```go
type Station struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	LatDeg     float64 `json:"lat_deg"`
	LonDeg     float64 `json:"lon_deg"`
	MaskDeg    float64 `json:"mask_deg"`
	Bands      []Band  `json:"bands"` // the bands a customer can book, as in Season 2's stations.csv
	Stowed     bool    `json:"stowed"`
	StowReason string  `json:"stow_reason,omitempty"`
}

type Pass struct {
	SatID      string    `json:"sat_id"`
	Station    string    `json:"station"`
	AOS        time.Time `json:"aos"`
	LOS        time.Time `json:"los"`
	MaxElevDeg float64   `json:"max_elev_deg"`
}

// PassView is a predicted pass as the service shows it. Pass is embedded, so its
// fields sit beside "id" in the JSON.
type PassView struct {
	ID string `json:"id"` // the pass ID: "furrow-a/SQT/15864"
	Pass
}

// Booking is one reserved pass (Season 2, unchanged).
type Booking struct {
	ID      string    `json:"id"`      // "bk-4501", given by the service
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"`
	LOS     time.Time `json:"los"`
	Note    string    `json:"note,omitempty"`
}

// ErrorBody is the JSON body of every error the service writes.
type ErrorBody struct {
	Code      string `json:"error"`                // "not_found": see the table of errors
	Message   string `json:"message"`              // for people: lowercase, no final period
	RequestID string `json:"request_id,omitempty"` // from chapter 6
}

// APIError is an error that says exactly what the caller gets.
type APIError struct {
	Status  int    // 422
	Code    string // "invalid_booking"
	Message string // "los must be after aos"
}

func (e *APIError) Error() string { return e.Message }

// Customer is one company that books passes through the service (chapter 8 on).
type Customer struct {
	ID         string   // "furrowlight": lowercase, one word
	Name       string   // "Furrowlight"
	SatIDs     []string // the satellites it operates
	KeyHash    string   // lowercase hex SHA-256 of its whole API key
	WebhookURL string   // where its events go (chapter 10)
	WebhookKey []byte   // the secret its events are signed with (chapter 10)
}

// TokenResponse is the body of POST /tokens.
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Service is the booking service: its data, its clock and its handlers.
type Service struct {
	Stations  []Station
	Catalog   []PassView          // ordered by AOS, then ID
	Guide     fs.FS               // the customer guide's pages
	Store     Store               // a *MemStore until chapter 7 makes it this interface
	Customers map[string]Customer // by ID (chapter 8)
	TokenKey  []byte              // signs tokens (chapter 8)
	Limiter   *Limiter            // chapter 9
	Timeouts  Timeouts            // chapter 9
	Keys      *keyStore           // idempotency keys (chapter 10)
	Now       func() time.Time
	draining  atomic.Bool // chapter 12
}

// Config is everything the service is built from. Tests fill it from the fixtures.
type Config struct {
	Dir       string           // data folder: bookings.json, outbox.json (t.TempDir() in tests)
	Stations  []Station
	Catalog   []PassView
	Guide     fs.FS            // embedded in production, an fstest.MapFS in tests
	Customers []Customer       // chapter 8
	TokenKey  []byte           // chapter 8
	Now       func() time.Time
	Rand      *rand.Rand       // chapter 10: jitter, seeded in tests
	Log       io.Writer        // chapter 6: the access log
	Timeouts  Timeouts         // chapters 9 and 12
}
```

#### Fixtures every chapter shares

- **Stations**, in the bible's order, with Season 2's bookable bands: YST `S X`, PEL `UHF S`,
  SQT `S X`, LAN `UHF S`, CWR `S X`. A receive-only UHF yagi is never bookable.
  `{"code":"SQT","name":"Salar Quieto","lat_deg":-23.4,"lon_deg":-104.6,"mask_deg":3,"bands":["S","X"],"stowed":false}`
- **The fleet the service knows**: `pipit-1`, `pipit-2`, `pipit-3`, `furrow-a`, `furrow-b`,
  `brine-3`, `cirrine-1`. `thistle` is unknown to it (re-entered in July).
- **Customers and their test keys.** An API key is the customer ID, a dot, then 32 lowercase
  hex digits (16 random bytes in production). Fixtures use exactly these:

  | ID | Name | Satellites | Test API key | Test webhook secret |
  |---|---|---|---|---|
  | `apogee` | Apogee | pipit-1, pipit-2, pipit-3 | `apogee.9e4d2a7c1b0f3e5d8a6c4b2f0e1d3a5c` | `apogee webhook secret` |
  | `furrowlight` | Furrowlight | furrow-a, furrow-b | `furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b` | `furrowlight webhook secret` |
  | `brinewatch` | Brinewatch | brine-3 | `brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0` | `brinewatch webhook secret` |
  | `cirrine` | Cirrine Atmospherics | cirrine-1 | `cirrine.c3a1f7e95d2b4086a7e1c9f3b5d8a024` | `cirrine webhook secret` |

  The test token key is the 32 ASCII bytes `bookingd test token key 32 bytes`. Production keys
  and secrets are 32 random bytes and never appear in a lesson.
- **The catalog.** Pass IDs as the bible gives them (satellite, station, orbit number padded
  to five digits). Orbit numbers grow by about 15 a day (Furrow-A 15.21, Pipit-1 15.16,
  Brine-3 15.3, Cirrine-1 14.92 a day), continuing from Season 1's Board; a pass between
  two the catalog lists is predicted from a model of its satellite's orbit, never guessed. Furrow-A crosses Salar Quieto at
  14:02 UTC on the days this season needs. Passes every chapter may use:

  | Pass ID | AOS | LOS | Max elev |
  |---|---|---|---|
  | `pipit-1/PEL/08154` | 2027-12-06T04:12:05Z | 2027-12-06T04:21:27Z | 61.4 |
  | `cirrine-1/YST/11364` | 2027-12-06T07:30:00Z | 2027-12-06T07:40:04Z | 44.2 |
  | `brine-3/LAN/09213` | 2027-12-06T09:14:20Z | 2027-12-06T09:22:02Z | 38.5 |
  | `furrow-a/SQT/15864` | 2027-12-06T14:02:10Z | 2027-12-06T14:11:52Z | 47.9 |
  | `furrow-b/SQT/15862` | 2027-12-06T14:49:30Z | 2027-12-06T14:58:58Z | 44.0 |
  | `furrow-a/SQT/16731` | 2028-02-01T14:02:10Z | 2028-02-01T14:11:52Z | 47.9 |

  Lessons add passes of their own by the same rules (passes last 5 to 15 minutes; Brine-3
  never rises over Ystmark).
- **Booking IDs** are `bk-` and four digits. Yusra's spreadsheet ends at `bk-4498`; in
  chapters 4 to 6 the service holds her last imported bookings, `bk-4491` to `bk-4498`. The
  service's own count starts at `bk-4501` (Furrowlight's first, chapter 5) and is near
  `bk-4912` on launch day (Furrow-A's 14:02 pass), and reaches `bk-4937` in the first week open (Brinewatch's 12th February booking, Friday 4 February).

#### The API

Every answer the service writes is JSON (`Content-Type: application/json`) except the guide's
pages and the router's own 404 and 405. Times are RFC 3339 in UTC. A list is `[]` when empty,
never `null`. The router's answers are plain text: `404 page not found` for a path no route
matches, and `Method Not Allowed` with an `Allow` header (`Allow: DELETE, GET, HEAD`) for a
route's path with another method. Never register a catch-all `/`: it would answer every wrong
method with its own 404 instead of 405.

| Route | Answer | Guard | Chapter |
|---|---|---|---|
| `GET /health` | 200 `{"status":"ok"}`; 503 `{"status":"draining"}` while shutting down | public | 3; 12 |
| `GET /stations` | 200 every Station in table order; `?band=UHF\|S\|X` keeps the stations that offer it | public | 3 |
| `GET /stations/{code}` | 200 the Station; 404 `no station "ZZZ"`. Codes match exactly: `sqt` is unknown | public | 4 |
| `GET /guide/` | the customer guide's pages (`http.FileServerFS` under `http.StripPrefix("/guide/", …)`); `/guide` redirects to `/guide/` with 307 | public | 3 |
| `POST /tokens` | 200 TokenResponse | API key | 8 |
| `GET /passes` | 200 PassView list, ordered by AOS then ID, with the query below | token, limit | 4; 8 |
| `GET /passes/{sat}/{station}/{orbit}` | 200 the PassView whose ID is `{sat}/{station}/{orbit}`; 404 `no pass furrow-a/SQT/15865` | token, limit | 4; 8 |
| `POST /bookings` | 201 the stored Booking, with `Location: /bookings/bk-4501` | token, limit, idempotency key | 5; 8; 10 |
| `GET /bookings` | 200 the caller's bookings, ordered by AOS then ID | token, limit | 7; 8 |
| `GET /bookings/{id}` | 200 the Booking; 404 `booking bk-5499: not found` | token, limit | 4; 8 |
| `DELETE /bookings/{id}` | 204 with no body; 404 as for GET | token, limit | 4; 8 |

`GET /passes` query, all optional, other parameters ignored:

- `sat`: a satellite of the fleet, else 400 `unknown satellite "thistle"`.
- `station`: a station code, else 400 `unknown station "ZZZ"`.
- `from`: an RFC 3339 time, default the service's `Now()`, else 400
  `from must be an RFC 3339 time`.
- `days`: a whole number from 1 to 14, default 7, else 400
  `days must be a whole number from 1 to 14`.
- It lists passes whose AOS is at or after `from` and before `from` plus `days` × 24 hours.

In chapters 3 and 4, before errors are JSON, handlers answer their errors with `http.Error`
and the same message text. From lesson 5.2 every handler error goes through `writeError` as an
`ErrorBody` (5.1 still answers with `http.Error`):

`{"error":"not_found","message":"booking bk-5499: not found","request_id":"req-000042"}`

| Status | `error` | Messages, exactly |
|---|---|---|
| 400 | `bad_request` | `request body: empty`; `request body: ` followed by the decoder's error, such as `request body: json: unknown field "sat"`; `request body: unexpected data after the JSON value`; the query messages above; `unknown band "x"`; `idempotency key must be 1 to 64 letters, digits, - or _`; `form has no "file" part` (chapter 13) |
| 401 | `unauthorized` | `missing credentials`; `invalid api key`; `invalid token`; `token expired` |
| 403 | `forbidden`; `quota_exceeded` | `cirrine-1 is not one of your satellites`; `invalid link`; `link expired`; `plan of 12 bookings a month used up: the count starts again at 2028-03-01T00:00:00Z` (the last three chapter 13) |
| 404 | `not_found` | `no station "ZZZ"`; `no pass furrow-a/SQT/15865`; `booking bk-5499: not found` |
| 409 | `conflict` | `furrow-a over SQT at 2027-12-06T14:02:10Z: already booked`; `idempotency key "wk49-mon" was used for a different request`; `a request with idempotency key "wk49-mon" is in progress` |
| 413 | `too_large` | `request body over 65536 bytes` |
| 415 | `unsupported_media_type` | `send the body as application/json`; `send the command file as multipart/form-data`; `command file must be plain text, not application/pdf` (chapter 13) |
| 422 | `invalid_booking` | the booking rules below |
| 429 | `rate_limited` | `slow down: 10 requests at once, then 1 a second` |
| 500 | `internal` | `internal error`, whatever went wrong; the error's own text never reaches the caller |
| 503 | `unavailable` | `request took too long`; `request canceled` |

`writeError` decides the status from the error value: an `*APIError`, wrapped or not
(`errors.AsType`), gives its own status, code and message (chapter 5); `ErrNotFound` gives 404
and `ErrConflict` 409 with the error's message (chapter 7); `context.DeadlineExceeded` and
`context.Canceled` give 503 (chapter 9); anything else gives 500 `internal`. From chapter 6 it
fills `request_id` from the request's context.

`writeJSON` marshals first; if that fails it answers 500 with the plain text
`internal error`. Otherwise it sets `Content-Type: application/json`, writes the status, then
the JSON and a newline.

**Creating a booking.** The body is a Booking without `id`. The service checks, in this order,
and the first failure answers:

1. `authenticate` (chapter 8): 401.
2. `limit` (chapter 9): 429.
3. `idempotent` (chapter 10): 400 for a malformed key, 413 for a body over the limit, a
   replay, or 409.
4. `decodeJSON`: the media type of `Content-Type` must be `application/json` (any parameters,
   any letter case), else 415; at most 65,536 bytes (`http.MaxBytesReader`; a
   `*http.MaxBytesError` from either decode call is 413); one JSON object with no unknown field
   (`DisallowUnknownFields`), then nothing but white space (a second `Decode` must return
   `io.EOF`), else 400.
5. `checkShape`, 422 `invalid_booking`, in this order: `id is chosen by the service`;
   `sat_id is required`; `station is required`; `aos and los are required`;
   `los must be after aos`; `note is longer than 200 characters` (characters, not bytes).
6. Ownership (chapter 8): the caller must operate the satellite, else 403. It comes before
   the catalog check, so a customer learns nothing about another customer's passes, not even
   whether one exists (a 422 for a missing pass beside a 403 for a real one would tell).
   Nobody operates `thistle`, so from chapter 8 on booking it answers 403 `thistle is not
   one of your satellites`; checkPass's `unknown satellite` is reachable only in chapters 5
   to 7 (`GET /passes?sat=thistle` still answers 400).
7. `checkPass`, 422 `invalid_booking`, in this order: `unknown satellite "thistle"`;
   `unknown station "ZZZ"`; `no predicted pass of furrow-a over SQT at 2027-12-06T14:02:10Z`
   (the satellite, station, AOS and LOS must equal a catalog pass's, times compared with
   `Equal`); `bookings for this pass closed at 2027-12-06T13:52:10Z`. Bookings close ten
   minutes before AOS, when Ingrid's checklist must be done ("AOS 04:12. Checklist done by
   04:02."): at exactly AOS minus 10 minutes the pass is closed.
8. `Store.Create`: a booking of the same pass (satellite, station and AOS) already kept is
   409 `conflict`.
9. 201 with `Location: /bookings/<id>` and the stored Booking.

A cancellation (`DELETE`) has no time rule: it answers 204 whenever the booking exists and is
the caller's.

#### Keys and tokens (chapter 8)

- **API keys** are given to customers by Yusra and sent only to `POST /tokens`, and from chapter 11 once per sign-in in the account page's form (field
  `key`), in
  `Authorization: Key <api key>`. The service keeps only `KeyHash`, the lowercase hex SHA-256
  of the whole key, and checks a key by hashing it and comparing the two hex strings with
  `subtle.ConstantTimeCompare`. A key that doesn't parse (no dot, a secret that isn't 32
  lowercase hex digits), an unknown customer and a wrong secret all answer the same 401
  `invalid api key`.
- **Tokens** are what every other customer route takes, in `Authorization: Bearer <token>`.
  A token is `<customer ID>.<expiry in Unix seconds>.<signature>`; the signature is the 64
  lowercase hex digits of HMAC-SHA256, under the service's 32-byte token key, of
  `<customer ID>.<expiry>`. A token lives one hour: `expires_at` is `Now()` plus one hour,
  truncated to the second. Furrowlight's token expiring at 2028-01-05T11:00:00Z is
  `furrowlight.1830682800.` followed by its 64 hex digits.
- `verifyToken` checks, in order: three parts separated by dots; a signature of 64 hex digits
  that `hmac.Equal` accepts; an expiry in the future. A token is valid while `now` is before
  its expiry; at the expiry second it is `token expired`. Any other failure, including a
  customer the service doesn't know, is `invalid token`.
- `credentials` reads `Authorization`: the scheme (matched with `strings.EqualFold`), one
  space, then a non-empty credential. A missing header, another scheme or an empty credential
  answers 401 `missing credentials`.
- Every 401 carries `WWW-Authenticate: Key` on `POST /tokens` and `WWW-Authenticate: Bearer`
  everywhere else. `authenticate` fails closed: it never calls the next handler after a 401.
- `authenticate` puts the Customer in the request's context (`withCustomer`, read back with
  `customerFrom`) under an unexported key type.
- **Authorization.** A customer sees and books only its own satellites (`owns`): `GET /passes`
  and `GET /bookings` list only its satellites' passes and bookings (`?sat=` of another
  customer's satellite gives `[]`); another customer's pass or booking answers 404 exactly as a
  missing one, and its cancellation changes nothing; booking another customer's satellite
  answers 403. Apogee's own key sees only the Pipits.
- `GET /health`, `GET /stations`, `GET /stations/{code}` and the guide stay public. `POST
  /tokens` is not rate limited.

#### Rate limits and timeouts (chapter 9)

- A **token bucket** per customer, applied after `authenticate` on every token route: it holds
  at most 10 tokens, refills continuously at 1 token per second, and starts full (a zero
  `Bucket` refills to full on first use). A request takes one token. Time comes from the
  injected `now`; when the clock steps backward nothing is added. Readmes say "token bucket"
  and explain once that its tokens are not chapter 8's.
- Every answer that passed the limiter carries `X-RateLimit-Limit: 10` and
  `X-RateLimit-Remaining: <whole tokens left after this request>`. A refused request answers
  429 `rate_limited` with `Retry-After: <seconds until one token is back, rounded up, at least
  1>` and `X-RateLimit-Remaining: 0`. Requests refused by `authenticate` are never counted. A
  request that reaches `limit` without a customer in its context is a wiring mistake and
  answers 500.
- Clients wait as `Retry-After` says: whole seconds, 1 second when the header is missing or
  unreadable, at most 60 seconds (`retryDelay`).

```go
// Bucket is one customer's token bucket.
type Bucket struct {
	Tokens float64   // tokens left, at most 10
	Last   time.Time // when Tokens was last brought up to date
}

// Limiter keeps a bucket per customer; it is safe for many requests at once.
type Limiter struct {
	now     func() time.Time
	mu      sync.Mutex
	buckets map[string]*Bucket
}

func newLimiter(now func() time.Time) *Limiter
func (l *Limiter) Allow(customer string) (ok bool, remaining int, retryAfter time.Duration)

// Timeouts are the service's limits, production values in comments. A zero leaves that limit unset.
type Timeouts struct {
	ReadHeader time.Duration // 5s: a client has this long to send its headers
	Read       time.Duration // 15s: ... to send the whole request
	Write      time.Duration // 30s: ... to take the whole answer
	Idle       time.Duration // 2m: a kept-alive connection may sit unused this long
	Handler    time.Duration // 10s: the deadline of every request's context (withDeadline)
	Drain      time.Duration // 5s: /health says draining this long before shutdown begins
	Shutdown   time.Duration // 20s: how long shutdown waits for requests in flight
}
```

A test sets the limits it checks, in milliseconds (50 ms; up to a second for Drain and
Shutdown), and leaves the others zero, so a request it holds open on purpose is never cut by
a limit it isn't testing.

`newServer` sets `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout` and `IdleTimeout` from
the Timeouts and `MaxHeaderBytes` to 16 KiB (16384; Go allows 4 KiB more before it answers
431, so tests send a 32 KiB header). A zero `ReadTimeout` or `WriteTimeout` means no limit; a zero `ReadHeaderTimeout` or
`IdleTimeout` falls back to `ReadTimeout`, as Go's server does. `withDeadline(0)`
adds no deadline, a zero `Drain` skips draining, and a zero `Shutdown` waits for every
request in flight.

#### Middleware, request IDs and the log (chapter 6)

- Middleware is `func(http.Handler) http.Handler`. `chain(h, a, b, c)` returns
  `a(b(c(h)))`: the first one listed runs first.
- Around the mux, outermost first: `requestID`, `logRequests`, `recoverPanics`, `noStore`,
  then `withDeadline` (chapter 9). Each private route is `chain(h, authenticate, limit)`, and
  `POST /bookings` adds `idempotent` inside them (chapter 10).
- `noStore` sets `Cache-Control: no-store` on every answer the service writes; Go's file server
  removes it from its own error answers (a missing guide page), so no test expects it there.
- `requestID` numbers requests from `req-000001` (six digits, counting up per process, from an
  `atomic.Uint64`), sets the `X-Request-Id` answer header before calling the next handler, and
  puts the ID in the request's context under an unexported key type. `requestIDFrom` returns
  `""` for a context without one.
- `statusRecorder` embeds `http.ResponseWriter`, records the first status written (200 when
  `Write` comes first, and 200 when the handler writes nothing, as the server sends) and the
  bytes written, and has `Unwrap() http.ResponseWriter` so
  `http.NewResponseController` still finds `Flush`.
- `logRequests` writes one line per request with a `*log.Logger` made with no prefix and no
  flags (safe for many goroutines), after the handler returns:
  `req-000007 POST /bookings 201 236B 4ms`. That is the request ID, the method, the URL path
  (never the query), the status, the bytes, and the time from the injected clock rounded to the
  millisecond in `time.Duration`'s own format (`0s`, `4ms`, `1.204s`).
- `recoverPanics` recovers a handler's panic, logs `req-000003 panic: <value>`, and answers 500
  `internal` through `writeError`. A panic with `http.ErrAbortHandler` is panicked again, so
  the server cuts the connection as that value asks.

#### The store and its file (chapter 7)

```go
var (
	ErrNotFound = errors.New("not found")      // wrapped: "booking bk-5499: not found"
	ErrConflict = errors.New("already booked") // wrapped: "furrow-a over SQT at 2027-12-06T14:02:10Z: already booked"
)

// Store keeps bookings. Its methods are safe to call from many goroutines at once.
type Store interface {
	// Create gives b the next ID and keeps it. It returns ErrConflict when a booking
	// of the same pass (satellite, station and AOS) is already kept.
	Create(ctx context.Context, b Booking) (Booking, error)
	// Get returns the booking with this ID, or ErrNotFound.
	Get(ctx context.Context, id string) (Booking, error)
	// List returns every booking, ordered by AOS and then ID, in a new slice.
	List(ctx context.Context) ([]Booking, error)
	// Delete removes the booking with this ID, or returns ErrNotFound.
	Delete(ctx context.Context, id string) error
}
```

- Every method returns `ctx.Err()` first when the context is already done.
- `MemStore` (`NewMemStore()`) is a map behind a `sync.Mutex`, with the next ID number; the
  conflict check and the insert happen under one lock. IDs count up from `bk-4501`. An
  unexported `load([]Booking) error` restores saved bookings with their IDs (duplicate IDs are
  an error) and sets the next number to the larger of the current one and one more than the
  highest held: never backward, never below 4501.
- `FileStore` (`OpenFileStore(path)`) keeps a `*MemStore` and the path of `bookings.json` in
  the service's data folder. After every successful Create or Delete it saves `List`'s
  bookings as one JSON array through `replaceFile(path, data)`: Season 2's write to a private
  temporary sibling, close, then rename, mode 0600. When saving fails it undoes the change and
  returns the error, so the store and the file never disagree. Opening: a missing file is an
  empty store; a file that doesn't decode, strictly, is an error, never an empty store.

#### Webhooks and idempotency keys (chapter 10)

```go
// Landing is dlinkd's word that a booking's file has been delivered.
type Landing struct {
	BookingID string    // "bk-4912"
	File      string    // "furrow-a-20280201-1402.tif": satellite, AOS day and time
	Bytes     int64
	At        time.Time // when the file was complete, UTC
}

// Event is what the service tells a customer when something happens to a booking.
type Event struct {
	ID        string    `json:"id"`         // "evt-000118": the same on every attempt
	Type      string    `json:"type"`       // "data.landed", the only type so far
	BookingID string    `json:"booking_id"`
	SatID     string    `json:"sat_id"`
	Station   string    `json:"station"`
	File      string    `json:"file"`
	Bytes     int64     `json:"bytes"`
	LandedAt  time.Time `json:"landed_at"`
}

// Delivery is one event on its way to one customer.
type Delivery struct {
	Customer string    `json:"customer"` // "furrowlight"
	Event    Event     `json:"event"`
	Attempts int       `json:"attempts"` // made so far
	NextAt   time.Time `json:"next_at"`  // when the next attempt is due
}

var ErrUndelivered = errors.New("undelivered")

// AnswerError is a receiver's answer outside 200 to 299: the event did not arrive.
type AnswerError struct {
	StatusCode int    // 503
	Status     string // "503 Service Unavailable"
}

// Sender delivers events to customers' webhooks.
type Sender struct {
	Client    *http.Client                                     // Timeout 5s; never follows a redirect
	Customers map[string]Customer                              // by ID: each one's WebhookURL and WebhookKey
	Now       func() time.Time                                 // each attempt's timestamp
	Rand      *rand.Rand                                       // jitter; used only under mu
	Wait      func(ctx context.Context, d time.Duration) error // sleeps in production; tests record d
	mu        sync.Mutex
}
```

- **The request.** `POST <the customer's WebhookURL>` with the event as `json.Marshal` writes
  it (no newline) and these headers: `Content-Type: application/json`;
  `Apogee-Timestamp: <Unix seconds of this attempt>`, such as `1833027183`;
  `Apogee-Signature: <64 lowercase hex digits>`, the HMAC-SHA256 under the customer's webhook
  secret of the timestamp, a dot, then the exact body bytes.
  `{"id":"evt-000118","type":"data.landed","booking_id":"bk-4912","sat_id":"furrow-a","station":"SQT","file":"furrow-a-20280201-1402.tif","bytes":10912500000,"landed_at":"2028-02-01T14:13:02Z"}`
- **Delivered** means any 2xx answer. **Retried**: transport errors (refused connections,
  timeouts), 408, 429 and 500 to 599. **Not retried**: any other answer, 3xx included (the
  client never follows redirects); the event is undelivered at once. A canceled context (the
  service stopping) is not a failure: the delivery keeps its attempt count for later.
- **The schedule.** At most 8 attempts. The first goes out as soon as the landing arrives (on launch day
  a second after the file lands, 71 s after LOS). After failed attempt n (n = 1 to 7) the sender waits `backoff(n)`:
  2 s doubled for each earlier failure (`2s × 2^(n−1)`), capped at 5 minutes, times a jitter
  factor drawn uniformly from [0.5, 1.0) from the sender's seeded `math/rand/v2` source. So
  the waits fall in 1–2 s, 2–4 s, 4–8 s, 8–16 s, 16–32 s, 32–64 s and 64–128 s. Each attempt is
  signed again with its own timestamp. After the 8th failure the event is undelivered
  (`ErrUndelivered`) and goes on Yusra's morning list.
- **Receivers** (Yusra's sample for customers, `verifyEvent`): read at most 64 KiB of body;
  refuse a timestamp more than 300 seconds from their own clock either way (exactly 300 is
  accepted); recompute the signature and compare with `hmac.Equal` before trusting anything in
  the body; drop an event whose `id` they have already handled; answer 2xx within 5 seconds.
- **Event IDs** are `evt-` and six digits, counting up per service from `evt-000001`. From
  chapter 12 the count is kept in the outbox, so a restart never reuses an ID.
- **Idempotency keys.** `POST /bookings` may carry `Idempotency-Key`: 1 to 64 letters, digits,
  `-` or `_` (else 400). Keys belong to the customer (Brinewatch's `k1` is not Cirrine's
  `k1`) and are honored for 24 hours on the injected clock; after that the same key starts afresh. The service keeps the
  request body's SHA-256 with the answer's status, `Content-Type` and `Location` headers and
  body. The same key with the same body gets that answer again with
  `Idempotent-Replayed: true` and books nothing; the same key with another body answers 409;
  a second request while the first is still running answers 409 (in progress). Answers of 500
  and above are not kept, so a retry of those runs for real. `idempotent` reads the body
  through `http.MaxBytesReader` (413 when too long) and hands the handler a fresh reader of
  the same bytes.

#### The account page (chapter 11)

- **Routes:** `GET /account/{$}` → the page, or 303 to `/account/sign-in`; `GET /account/sign-in` → the form;
  `POST /account/sign-in` → form field `key`, body at most 4,096 bytes, checked with `checkKey`: 303 to
  `/account/` with the session cookie, or 401 with the form and `invalid api key`; `POST /account/sign-out`
  → ends the session, clears the cookie, 303 to the form; `POST /account/bookings/{id}/cancel` → only the
  caller's own booking (another customer's is a plain-text 404 `booking bk-5499: not found`, exactly
  like a missing one), then 303 to `/account/`. Pages are HTML (`text/html; charset=utf-8`), rendered into
  a buffer from a template; errors are plain text.
- **Cookie:** `session=<rand.Text() ID>; Path=/account; Max-Age=1800; HttpOnly; Secure; SameSite=Lax`,
  cleared with the same name and path and `Max-Age=0`.
- **Session:** valid while now is before its end; ended at that very instant and deleted; every
  successful lookup renews it for 30 minutes, and the page re-sends the cookie.
- **Guards:** `http.CrossOriginProtection` (Go's own 403) around every account route; `cors` around the
  API only: `Vary: Origin` everywhere and, for exactly the planning page, `Access-Control-Allow-Origin`,
  its preflight answered 204 with `Allow-Methods: GET, POST, DELETE` and `Allow-Headers: Authorization,
  Content-Type, Idempotency-Key`; on every account answer `Content-Security-Policy: default-src 'self';
  frame-ancestors 'none'`, `Strict-Transport-Security: max-age=31536000`, `X-Content-Type-Options:
  nosniff`, `Referrer-Policy: same-origin`, plus `Cache-Control: no-store`.
- **Wiring:** `site` is a mux with `/account/` → `chain(a.routes(), securityHeaders, noStore,
  noCrossOrigin)` and `/` → `cors(origin)(api)`. The API's own mux still answers its 404 and 405 (the
  rule against a catch-all `/` is about the API's mux). Launch Day's lessons leave the page out of
  `newApp`.
- **Fixtures:** new bookings `bk-4776` (Furrowlight, `furrow-b/SQT/16607`, Mon 24 Jan 14:19:53–14:29:53,
  booked before Friday), `bk-4784` (Yusra's test booking, `pipit-1/PEL/08904`, Mon 24 Jan
  15:08:30–15:17:50, made Friday afternoon) and `bk-4785` (the planning page's first booking, from 11.9's
  stand-in API); `bk-4779`, `bk-4781` and `bk-4783` come from chapter 10; `bk-4786` to `bk-4800` stay free
  before Launch Day's `bk-4801`. Furrowlight's test tokens expiring at 18:00 and 19:00 UTC are
  `furrowlight.1832090400.8e8c…` and `furrowlight.1832094000.86f0…`, under the test token key. Every
  pass was checked with `cmd/passes`.

#### Starting and stopping (chapter 12)

- `run` serves `app.Handler` on the given listener with `newServer`, until `ctx` is done. Then,
  in this order: the service starts draining (`GET /health` answers 503 `draining`, every
  other route still works) for `Timeouts.Drain`; `Shutdown` runs with a deadline of
  `Timeouts.Shutdown`; the notifier stops and saves the outbox; `run` returns nil. `Serve`'s
  `http.ErrServerClosed` is never returned as a failure.
- Handlers' contexts come from a base context (`BaseContext`) that `run` cancels only if
  `Shutdown`'s deadline passes; then it calls `Close` and returns the error wrapped as
  `shutdown: context deadline exceeded`. Handlers pass `r.Context()` to the store, never
  `context.Background()`; deliveries run under the notifier's own context, which `Stop` cancels.
- `main` turns SIGINT and SIGTERM into a canceled context with `signal.NotifyContext`, and is
  the only code that ever listens on a real address. No lesson runs it.
- **The outbox**: `outbox.json` in the data folder, saved whole through `replaceFile` on every
  change: `{"next_event":118,"pending":[<Delivery>...]}`. At start, pending deliveries whose
  `next_at` has passed are tried at once and the rest at their time; event IDs continue from
  `next_event` even when nothing is pending. A missing file is an empty outbox with
  `next_event` 1; a broken one stops the start.
- On launch day the outbox holds Cirrine's `evt-000117` (the file of Cirrine-1's 13:30 pass over Ystmark,
  `cirrine-1/YST/12218`, which landed at 14:04:57 over Ystmark's slow backhaul; its receiver
  was down for its own deploy, and two attempts were made) when the service restarts; Furrowlight's
  launch-day event is `evt-000118`.

#### The first week open (chapter 13)

- **Days:** Wed 2 Feb 09:10, Cirrine's 1.9 GB software image as the command file of `bk-4918`
  (`cirrine-1/YST/12245`, Thu 3 Feb 09:00:05–09:10:24); Wed 15:00, a PDF as Furrowlight's command file
  for `bk-4921` (`furrow-a/SQT/16777`, Fri 4 Feb 14:36:38–14:46:32); Thu 08:10, Hollis's four broken
  downloads of `furrow-a-20280201-1402.tif` (`bk-4912`, 10,912,500,000 bytes, landed 14:13:02 on 1 Feb);
  Thu 11:00, his field team's link (until Fri 11:00); Thu 15:00, Tomás (on call for launch week) finds
  Brinewatch's events in the weather service's log, and Brinewatch sends Yusra its public address that
  afternoon; Fri 08:30, Mirela's email; Fri 11:00, Brinewatch at 11 of 12; Fri 16:00, 13.9; 17:00, the
  deploy.
- **Routes**, all inside the service's middleware: `POST /bookings/{id}/commands` (token, limit): a
  `multipart/form-data` body of at most 65,536 bytes whose part `file` is kept as `commands/<id>.txt` in
  the data folder when `http.DetectContentType` calls it exactly `text/plain; charset=utf-8`; 204.
  `GET /bookings/{id}/file` (token, limit) and `GET /files/{id}?expires=&sig=` (the link alone): the
  landed file through `http.ServeContent` as `application/octet-stream`, 206 for a `Range`, 404 `no file
  for bk-4921 yet`. `POST /bookings/{id}/link` (token, limit): 200 `{"url":…,"expires_at":…}`, valid
  `linkLife` (24 hours) from `Now()` truncated to the second. `GET /go?next=` (public): 303 to `next` when
  it parses, starts with one `/` and holds no `\`, else to `/guide/`; it logs `req-000009 go
  /guide/files.html`.
- **dlinkd's delivery folder** holds each landed file in a folder named for its booking:
  `bk-4912/furrow-a-20280201-1402.tif`. From chapter 13 customers fetch landed files from the booking
  service.
- **Download links:** `<path>?expires=<Unix seconds>&sig=<64 hex>`, the HMAC-SHA256 under the link key (32
  random bytes of its own; in tests the ASCII `bookingd test link key, 32 bytes`) of `<path>?expires=<Unix
  seconds>`. The signature is checked first with `hmac.Equal` (`invalid link`), then the expiry (`link
  expired` from its second on), both 403 `forbidden`. With the test key, `/files/bk-4912?expires=1833274800`
  signs as `2757ac02d305eefdc043203dc622beab809b459066d7826fa50c58868ad3ac1d`.
- **Webhooks** go out through `webhookClient`: 5 s, no redirect followed, and a dialer whose `Control` hook
  refuses loopback, private, link-local and unspecified addresses and any address that doesn't parse
  (`errInside`, "webhook address is inside a network"). Launch Day's `newApp` still builds the old client;
  13.9's `firstWeek` swaps it in. A refused dial is a transport error, so `retryable` still retries it:
  eight attempts, none of which leaves the machine. Brinewatch's `evt-000121`, for `bk-4913`'s file on
  Wednesday, was the first event posted to 10.20.0.31.
- **Plans**, bookings made per calendar month in UTC: Furrowlight 60, Cirrine 40, Brinewatch 12; Apogee's
  own have none. `Quota` keeps `usage.json` in the data folder. Before Friday's deploy Yusra seeds it with
  `{"month":"2028-02-01T00:00:00Z","used":{"brinewatch":11,"cirrine":4,"furrowlight":9}}`, the 24
  bookings `bk-4913` to `bk-4936`. Past the plan the answer is 403 `quota_exceeded` `plan of 12 bookings a
  month used up: the count starts again at 2028-03-01T00:00:00Z`, never 429, since clients retry a 429
  within a minute. `metered` counts each booking before the store keeps it, and gives back one the store
  refuses.
- **Fixtures:** `bk-4913` (Brinewatch, `brine-3/LAN/10103`, Wed 12:27:50–12:34:37), `bk-4918`, `bk-4921`,
  `bk-4933` (Furrowlight, `furrow-a/SQT/16792`, Sat 5 Feb 14:16:50–14:26:55), `bk-4936` (Brinewatch,
  `brine-3/CWR/10136`, Fri 15:43:26–15:52:39), `bk-4937` (Brinewatch's 12th, `brine-3/PEL/10140`, Fri
  21:06:32–21:14:28); the refused 13th is `brine-3/PEL/10141` (Fri 22:43:27–22:52:45). 13.9's catalog adds
  `furrow-b/YST/16776` and `cirrine-1/YST/12265`. Every pass comes from `cmd/passes`.

#### The weather service (chapter 2, and the Board)

- `GET /weather/{code}` answers 200 with one reading:
  `{"station":"CWR","at":"2027-11-08T13:55:00Z","wind_kmh":64,"gust_kmh":81,"temp_c":4.5}`.
  It sometimes adds fields (`humidity_pct` since last spring), so its client ignores unknown
  fields: the opposite of the booking service's own strict input, on purpose.
- Its errors are plain text: 404 `no station "CRW"`; 503
  `mast silent since 2027-11-08T12:51:00Z` with `Retry-After: 60`.
- It listens at 10.20.0.31, port 80, on Apogee's office network (chapter 13).
- A reading is fresh for 15 minutes after its `at`: exactly 15 minutes old is still fresh.
- Its client never waits more than 2 seconds (`http.Client{Timeout: 2 * time.Second}`), and
  each call also takes the poll's context.
- The go/no-go wind call uses the sustained wind (`wind_kmh`), as Season 1 did: under 50 km/h
  `GO: wind 22 km/h`; from 50 `GO: wind 61 km/h, warning`; from 72
  `NO-GO: wind 74 km/h, stow`; a stale reading `NO-GO: reading 47m0s old`; a 503
  `NO-GO: mast silent`; any other failure `NO-GO: no weather`. Wind prints with no decimals.
  The rules apply in this order, the first that fits answering: a 503, any other failure, a
  stale reading, then the wind from 72 km/h down.
- `decodeWeather` decodes leniently and refuses a reading without its station or its time:
  `reading has no "station"`, `reading has no "at"`.

```go
// Weather is one station's latest reading from the weather service.
type Weather struct {
	Station string    `json:"station"`
	At      time.Time `json:"at"`
	WindKmh float64   `json:"wind_kmh"`
	GustKmh float64   `json:"gust_kmh"`
	TempC   float64   `json:"temp_c"`
}

// WeatherClient asks the weather service for station readings.
type WeatherClient struct {
	BaseURL string       // the service's address, in tests an httptest server's URL
	HTTP    *http.Client // with a Timeout
}

func (c *WeatherClient) Reading(ctx context.Context, code string) (Weather, error)

// StatusError is an answer from the weather service other than 200 OK.
type StatusError struct {
	StatusCode int    // 503
	Status     string // "503 Service Unavailable"
	Body       string // at most the first 512 bytes it sent, trimmed
}

func (e *StatusError) Error() string { return e.Status + ": " + e.Body }
```

`Reading` wraps its errors with the station: `weather for CWR: 503 Service Unavailable: mast
silent since 2027-11-08T12:51:00Z`.

#### Nkem's stand-in (chapter 1)

`standIn()` returns an `http.Handler`, a mux with two routes: `GET /stations` answers 200 with
`[{"code":"YST","name":"Ystmark"},{"code":"PEL","name":"Port Ellery"},...]` (code and name
only, table order) and `GET /stations/{code}` answers `{"code":"PEL","name":"Port Ellery"}`,
or 404 `no station "ZZZ"` in plain text. Both JSON answers end in a newline. Requests built
only to be printed use URLs on `http://127.0.0.1`, so a printed `Host` line reads
`127.0.0.1`; no lesson prints a port.

#### How lessons run servers and test them

- Every server is `httptest.NewServer` or `httptest.NewUnstartedServer` (which listen on
  127.0.0.1 on a port the system picks), or an `http.Server` serving
  `net.Listen("tcp", "127.0.0.1:0")`. Never a fixed port, never `localhost` or `:0` (all
  interfaces), never an address outside 127.0.0.1. Nothing reaches the network, and nothing
  asks macOS for permission.
- One handler call goes through `httptest.NewRecorder` and `httptest.NewRequest`. A check of
  a 204 or a HEAD answer goes through `httptest.NewServer`: a recorder keeps a body that a real
  server drops.
- Clients in tests are `srv.Client()` or have a timeout of at most 2 seconds. Every wait in a
  test has a guard of about one second, so a starter fails within a second or two and never
  hangs the run.
- `httptest.Server.Close` waits for requests in flight. A handler blocked on purpose waits on
  a channel the test closes or on `r.Context().Done()`, and streams are canceled first: cleanups
  run last-in first-out, so register `srv.Close` before the cleanups that release handlers.
- Clocks are `now func() time.Time` from a small fake clock with `Now` and `Add` (under a mutex
  when goroutines read it); randomness is `rand.New(rand.NewPCG(1, 2))`; sleeps are 50 ms at
  most; a timeout under test is 50 ms, and a request that must succeed uses the test server's
  `srv.Client()` or a client with a 2-second timeout, so a busy Mac never fails it.
- Tests compare JSON after decoding it, except where a lesson's point is the exact bytes.
  Messages read `call = got, want want`: `GET /stations/ZZZ = 200, want 404`.
- Data files live under `t.TempDir()`. The guide in tests is an `fstest.MapFS`.
- Where learners write checks themselves (chapters 7, 8 and 12), their code goes in an
  editable non-test file and the read-only test proves it against planted mutants in a
  re-executed test binary (`os.Args[0] -test.run=...` with `APOGEE_MUTANT` naming the
  mutant), as in Season 1's chapter 15.
- A flood of hundreds of requests goes through `ServeHTTP` with recorders, not the loopback.
- 13.6's test dials the name `localhost`, resolved from the Mac's own hosts file, to show that a name
  resolving inside is refused; no other test resolves a name. 13.1's 1.9 GB upload is made as it is read
  and stops after 1 MiB, so a starter that reads everything fails at once.

**Topics added this season** (front matter): `http`, `json`, `security`. Lessons keep using the
older ones (`context`, `errors`, `interfaces`, `mutexes`, `testing`, `time`...).

**New field guide words** (world words only; Go and HTTP words are taught in readmes):
booking, booking service, bookingd, the beta, catalog, customer guide, weather service, weather
mast, landing, account page, planning page, command file, plan (not linked: an everyday word).

**Name checks** to run before the season ships: `bookingd` (and the header names
`Apogee-Timestamp` and `Apogee-Signature`, which use the invented company's name). Nothing else
this season is a new proper name. The planning page's host, `plan.furrowlight.example`, uses the
reserved `.example`. Chapter 13 adds none: Mirela's look-alike site is `apogee-bookings.example`,
under the reserved `.example`.
