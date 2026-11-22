The starter prints `21ns` for 9 minutes and 12 seconds: it adds the two numbers and counts the sum as nanoseconds. What should each number be multiplied by?

Convert each parameter with `time.Duration(...)`, multiply it by its unit, `time.Minute` or `time.Second`, and add the two products.
