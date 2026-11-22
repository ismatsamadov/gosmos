---
slug: nothing-older-than-one-three
title: Nothing Older Than 1.3
type: code_tests
position: 2
difficulty: 9
topics: crypto, http, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
An old tool at Salar Quieto still speaks TLS 1.2 and nothing newer, and every station's server would quietly let it in. Mirela wants it refused: "One version. The newest."

A Go server accepts TLS 1.2 unless its config says otherwise. Over 1.3, every connection agrees fresh keys, by default with `X25519MLKEM768`: the X25519 and ML-KEM of the backhaul week working together, which 1.2 can't use. `ConnectionState` says what a connection agreed:

```go
st := resp.TLS
fmt.Println(tls.VersionName(st.Version))
fmt.Println(st.CurveID)
// TLS 1.3
// X25519MLKEM768
```

## Your task

Make `harden` set `cfg.MinVersion` to `tls.VersionTLS13`, and change nothing else in `cfg`. The test hardens Ystmark's config, then connects once as an ordinary client and once as the old tool.
