Which call takes a name and a function, and which field of each case holds the name?

`t.Run` takes the subtest's name and a `func(t *testing.T)`; the body of that function is where the check now goes.

Move the whole `if` into the function you pass to `t.Run`, and call its parameter `t`: that `t` shadows the outer one, so a `t.Errorf` on it reports the failure under the case's own name.
