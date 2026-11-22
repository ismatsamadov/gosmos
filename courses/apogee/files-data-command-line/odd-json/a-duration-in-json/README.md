---
slug: a-duration-in-json
title: A Duration in JSON
type: code_tests
position: 7
difficulty: 5
topics: interfaces, methods, testing, time
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
gctl saved Paz's settings with the timeout as `30000000000`, a count of nanoseconds. She wants the file to say `30s`, the way she would type it.

A type with a `MarshalJSON` method writes its own JSON, and an `UnmarshalJSON` method on its pointer reads it back. `Delay` is a `time.Duration` underneath: `time.Duration(d).String()` gives text such as `2m0s`, and `time.ParseDuration` reads it. Let `json.Marshal` and `json.Unmarshal` of a `string` add and remove the quotes.

```go
// MarshalJSON writes a code in capitals.
func (c Code) MarshalJSON() ([]byte, error) {
	up := strings.ToUpper(string(c))
	return json.Marshal(up)
}
```

## Your task

Give `Delay` a `MarshalJSON` method that writes duration text, and `*Delay` an `UnmarshalJSON` method that accepts only a JSON string `time.ParseDuration` can read. Return every error.

> [!WARNING]
> Trimming the quotes off by hand misreads escaped text: in JSON, `"\u0033\u0030s"` is `"30s"`.
