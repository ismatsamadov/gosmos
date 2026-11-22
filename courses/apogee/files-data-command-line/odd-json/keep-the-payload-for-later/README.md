---
slug: keep-the-payload-for-later
title: Keep the Payload for Later
type: code_tests
position: 4
difficulty: 5
topics: errors, structs, testing, types
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Furrowlight now says what reached it in delivery reports. Yusra forwarded the first batch at the last signal stop: one envelope per report, with a version, a kind, an ID and a payload whose shape depends on the kind.

A field of type `json.RawMessage` keeps a value's JSON exactly as it arrived, undecoded, so the envelope decodes now and the payload waits until the code knows its type. A JSON `null` arrives as the four bytes `null`.

```go
var msg struct {
	From string          `json:"from"`
	Body json.RawMessage `json:"body"`
}
in := `{"from":"YST","body":{"wind_kmh":64}}`
_ = json.Unmarshal([]byte(in), &msg)
fmt.Println(string(msg.Body)) // {"wind_kmh":64}
```

## Your task

1. Change the type of `Envelope.Payload` to `json.RawMessage`.
2. In `decodeEnvelope`, return an empty envelope and an error when the payload is missing or `null`.
