Parse each timestamp with `time.Parse` and the `time.RFC3339` layout, returning 0 and the error as soon as one fails.

Check the order before subtracting, with `Before`, and return 0 and an error of your own, made with `fmt.Errorf`, when the LOS comes first. Otherwise subtract with `Sub`, across midnight or not.

A negative duration means the subtraction runs the wrong way round: call `Sub` on the later time, with the earlier one as its argument.
