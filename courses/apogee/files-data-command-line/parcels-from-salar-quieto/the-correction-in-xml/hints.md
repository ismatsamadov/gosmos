The `XMLName` field's type is `xml.Name`, and its tag holds the root element's name, in the same backquoted form as the other fields' tags.

`MarshalIndent` takes the value, a prefix for every line and the indent for each level: here the prefix is empty and the indent is two spaces.

`xml.Header` already ends in a newline, so print it with `fmt.Print`; then print the marshaled bytes with `fmt.Println`, which adds the newline after `</correction>`.
