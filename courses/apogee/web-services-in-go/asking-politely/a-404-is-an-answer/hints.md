`Do` gave you `resp` without an error, but that only means an answer arrived. Look at `resp.StatusCode` before you hand `resp.Body` to the decoder.

Compare the code with `http.StatusOK`. When they differ, return an empty `Weather` and an error made with `fmt.Errorf`, the station's code and `resp.Status`, which already reads `404 Not Found`.

If the 404 case still reports `invalid character`, your check runs after the decoder: it has to come between the error check of `Do` and the decoding.
