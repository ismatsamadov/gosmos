The flown passes are the booked ones less the lost: give `availability` that and `booked`, and multiply its answer by 100 for a percent. Call `budget` once and keep its answer, since the row needs it twice.

`fmt.Sprintf` takes the layout from the task and then the six values in its order: code, booked, lost, availability, budget, and the budget less the lost. Return the string without a newline, since `main` prints each row with `Println`.

If a column is out of line, compare your layout with the header's in `main`: the same widths, with `d` for whole numbers and `.2f` for the percent. If Lantern Atoll's last column is 11 instead of -11, it was worked out the other way round.
