Keep what `srv.Serve(ln)` returns in a variable instead of returning it straight away.

`errors.Is(err, http.ErrServerClosed)` is true for the error a deliberate stop gives `Serve`, even when it comes wrapped: return nil then.

Every other error must come back as it is, so the caller can tell a broken listener from a stop. If the case with the listener closed underneath returns nil, you are swallowing every error.
