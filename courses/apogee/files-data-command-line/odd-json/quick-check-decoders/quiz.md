## The file is one settings object and a line break. What does next hold?

- [x] io.EOF, the clean end of the input
  > Token skips the line break after the object and finds nothing else.
- [ ] A syntax error for the line break
  > A decoder skips spaces and line breaks between values; they are never an error.
- [ ] nil, with the line break as its token
  > A line break is white space, not a token, so Token reads past it.

## The file spells min_elev_deg as min_elev_dg. What does err hold?

- [ ] nil, and MinElevDeg stays 0
  > That is what happens without DisallowUnknownFields: the field is skipped in silence.
- [x] An error that names min_elev_dg
  > DisallowUnknownFields turns a field that matches no struct field into an error.
- [ ] nil, and MinElevDeg holds the file's value
  > The names differ, so no struct field matches and the value is never stored.

## The file holds two settings objects. Which result shows it?

- [ ] err, from the first Decode
  > The first object is complete and decodes without an error.
- [ ] The error DisallowUnknownFields adds
  > That check looks at the field names inside an object, not at what follows the object.
- [x] next, because it is not io.EOF
  > The input has not ended after the first object, so the file must be refused.
