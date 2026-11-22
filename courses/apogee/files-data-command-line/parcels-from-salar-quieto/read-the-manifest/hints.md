Each tag goes after the field's type, in backquotes, as `xml:"…"` with no space after the colon. Look at the manifest in the test to see whether each value is an element of its own or an attribute inside a start tag.

`version` and `size` sit inside start tags, so their tags end in `,attr`. `Files` gathers every `file` element inside `files`, which is the `a>b` form.

`Entry` stands for one `<file>` element: its size is an attribute of that element, and its name is the element's own text, which the `,chardata` form reads, with nothing before the comma.
