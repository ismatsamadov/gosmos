Check one field at a time and return an error at the first one that breaks its rule. A small helper that takes the field's name, its text and a test for one character saves writing the same loop three times.

Every rule but the sequence number's is about characters: compare each rune with ranges such as `r >= 'a' && r <= 'z'`, and remember that `!` to `~` already leaves out the space and the tab. An empty field is refused too.

Name each argument by its place, counting from 1, so the first is `argument 1` and the second `argument 2`. For the bar, the error must read exactly `argument 1 contains "|"`, which `%q` of the character as a string gives you.
