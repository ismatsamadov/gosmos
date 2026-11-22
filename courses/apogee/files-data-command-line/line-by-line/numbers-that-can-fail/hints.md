`strconv.ParseFloat` takes the text and a bit size; 64 gives a `float64`, and a leading minus sign needs no extra work.

Check the error before you trust the number. When it is not nil, return 0 with that same error, so `41..70` can never become an angle.
