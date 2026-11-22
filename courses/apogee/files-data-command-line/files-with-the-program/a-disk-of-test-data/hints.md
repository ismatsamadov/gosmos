Each entry of the map has a path as its key and a `MapFile` holding the bytes in `Data` as its value, as in the example. The template's key keeps its folders, `templates/notices/stow.txt`, with no leading slash.

The table's text is two lines, each ending in `\n`: the header, then Cape Wren's line, exactly as the task spells them. Leave every other station out, because the test checks that SQT is unknown.

If `fstest.TestFS` reports a file it expected and did not find, compare that key with the task character by character. If the notice comes out without its newline, add `\n` to the template's text.
