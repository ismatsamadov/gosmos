## How does `html/template` print `{{.Call}}` in the class attribute and in the script?

- [ ] The same way in both: `<`, `>`, `&` and quotes become HTML entities
  > Entities are right for HTML text and attributes, but inside a script the value has to be written as JavaScript.
- [x] By where it lands: as attribute text in the class, as a quoted JavaScript string in the script
  > The template works out the context of each value and escapes it for that place, so the script gets `"track"`, quotes included.
- [ ] Only in the class: a value inside a script element is printed as it is
  > Values inside a script are escaped too, as JavaScript, which is what keeps a note from running as code.

## What changes when a note is converted to `template.HTML`?

- [x] Escaping is off for that one value, so its tags reach the page as tags
  > `template.HTML` says the value is trusted HTML, so `<b>` makes the text bold instead of showing up as text.
- [ ] The note is escaped twice, so `&` shows up as `&amp;amp;`
  > Nothing is escaped twice: `template.HTML` turns escaping off for the value, not on.
- [ ] Escaping stops for every value after it in the template
  > Only that one value is trusted; every other value is still escaped where it lands.

## Why does `servePage` execute the template into a `bytes.Buffer` first?

- [ ] A buffer makes the template run faster than writing to `w`
  > The page is made once either way; what matters is when its first byte leaves.
- [ ] `Execute` can only write to a `bytes.Buffer`
  > `Execute` takes any `io.Writer`, the `ResponseWriter` included.
- [x] Nothing is sent until the page is complete, so an error can still answer 500
  > Writing straight to `w` would send the first half with status 200 before the error is known.
