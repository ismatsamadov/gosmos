---
slug: quick-check-secure-connections
title: Quick Check: Secure Connections
type: quiz
position: 3
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T07:46:46Z
---
Wednesday. Lantern Atoll's link is the first to switch to TLS, tomorrow morning, and Nkem reads its client config before it goes, one question at a time:

```go
cfg := &tls.Config{
	RootCAs:    apogeeRoots,
	ServerName: "lan.apogee.internal",
	MinVersion: tls.VersionTLS13,
}
conn, err := tls.Dial("tcp", "10.4.0.7:8443", cfg)
```

## Your task

Answer Nkem's three questions: which name this client checks Lantern Atoll's certificate against, which roots it would trust without its `RootCAs` line, and what `InsecureSkipVerify` turns off.
