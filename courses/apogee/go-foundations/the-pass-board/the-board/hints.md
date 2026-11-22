Do the header first and compare it with the target; then get one pass right before you change the other three.

Each `fmt.Println(` becomes `fmt.Printf(` with the format string as its first argument; the values stay where they are, after it.

A pass has six verbs for its six values: `%02d:%02d` takes the hour and then the minutes, and `%6.1f` takes the elevation.
