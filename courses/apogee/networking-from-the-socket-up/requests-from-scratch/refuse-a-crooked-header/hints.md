`strings.Cut` already reports in its third result whether the line had a colon. Refuse when it didn't, when `isToken` turns the name down, or when `isFieldValue` turns the value down; a leading space or tab makes the name fail `isToken` too.

`readLine` returns `errLineTooLong` for a line that didn't fit its 1024 bytes: here that becomes the 431 about a header line. For the field count, keep a counter in the loop and refuse the line that would be the 33rd field.

If 33 lines of the same name get through, the count is of the map's names, and `len(h)` stays 1: count each line you read instead.
