Range over the bookings with their index, and check each one in turn. Look a station up with the two-value form of a map index: its second value tells a missing code from any other.

Order the checks as the task lists them and return at the first that fails. `IsZero` finds a time that was never set, and `After` is false for a LOS that is missing, equal to its AOS or earlier.

Only after the loop has passed every booking is the batch valid. Put the booking's ID in the message with `fmt.Errorf`, or its index plus 1 when the ID is empty.
