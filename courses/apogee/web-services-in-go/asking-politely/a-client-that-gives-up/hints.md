An `http.Client` is a struct, and its `Timeout` field is a `time.Duration`. A client made with `&http.Client{}` leaves that field at zero, which means no limit.

Set the field where the client is made, inside the composite literal, to the `timeout` the function receives. Keep the client new: don't hand out `http.DefaultClient`.

If the test of Cape Wren still waits a whole second, the client you return has no `Timeout`; if `TestNewWeatherClient` names `http.DefaultClient`, make a client of your own instead of changing the shared one.
