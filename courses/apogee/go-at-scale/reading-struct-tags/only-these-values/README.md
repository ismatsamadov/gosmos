---
slug: only-these-values
title: Only These Values
type: code_tests
position: 6
difficulty: 9
topics: reflection, runes, slices, strings
created: 2026-11-12T03:14:42Z
modified: 2026-11-12T03:14:42Z
---
Friday 1 June, afternoon. Ilha Calmaria's corrected form has the station's name now, and lists its radio bands as `S` and `K`, though Apogee books passes only on `UHF`, `S` and `X`.

`String()` gives a string field's text. A slice's `Value` has `Len()` and `Index(i)`, the `Value` of each element, so a rule written for one string can check each string of a `[]string` in turn. `oneof=UHF S X` lists the allowed words with spaces between them; `len=3` counts characters, and `á` is one character in two bytes.

```go
v := reflect.ValueOf([]string{"Tomás", "Nkem"})
for i := range v.Len() {
	s := v.Index(i).String()
	n := utf8.RuneCountInString(s)
	fmt.Println(len(s), n)
}
// 6 5
// 4 4
```

## Your task

Add `oneof` and `len` to `checkRule`. Each checks a string, or each string of a `[]string`, and reports the first that breaks it, in these forms:

```text
HostConfig.Bands: "K" is not one of UHF S X
HostConfig.Code: "VRLX" is not 3 characters
```

On a field of any other type, either rule is a mistake in the tag: return `misfit`'s problem.
