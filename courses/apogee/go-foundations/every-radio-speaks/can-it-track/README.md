---
slug: can-it-track
title: Can It Track?
type: code_tests
position: 6
difficulty: 6
topics: interfaces, types
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T14:56:15Z
---
The S-band and X-band radios correct Doppler shift on their own. The UHF radio can't, so before a UHF pass the routine shifts the frequency for it.

Radios that track satisfy a second interface, `DopplerTracker`. A type assertion can ask whether an interface value also satisfies another interface. Its comma-ok form never panics: `ok` says whether it does.

```go
type Stower interface{ Stow() }

func parkIfYouCan(e any) {
	if s, ok := e.(Stower); ok {
		s.Stow()
	}
}
```

## Your task

Make `needsPrecompensation` return `false` when `r` is also a `DopplerTracker`, and `true` otherwise.

> [!WARNING]
> Without `, ok`, an assertion that fails panics, and the first UHF pass would crash the routine.
