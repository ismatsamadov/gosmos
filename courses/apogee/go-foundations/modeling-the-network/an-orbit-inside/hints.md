Where does the inclination live: in `Satellite` itself, or in the `Orbit` inside it?

Reach it with two dots, `s.Orbit.` and then the field's name, and join two comparisons with `&&`.

Both ends count: 80 and 100 are polar, so compare with `>=` and `<=`, not `>` and `<`.
