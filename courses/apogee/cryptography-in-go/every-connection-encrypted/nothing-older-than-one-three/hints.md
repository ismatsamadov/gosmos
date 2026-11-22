`harden` receives a pointer, so setting a field through it changes the caller's config: assign to `cfg.MinVersion` and leave every other field alone.

The constant for the newest version lives in `crypto/tls`: `tls.VersionTLS13`. One assignment is the whole body of `harden`.

If `TestKeepsTheRest` fails, `harden` replaced the whole config, as `*cfg = tls.Config{...}` would, instead of setting one field.
