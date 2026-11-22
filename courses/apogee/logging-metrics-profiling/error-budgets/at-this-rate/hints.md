Answer the two plain cases first: nothing remaining (0 or less) gives 0 and true whatever was lost, and then nothing lost this week gives 0 and false.

For the rest, turn `remaining` into a `time.Duration`, multiply it by 7 days, and divide by `lostLast7d`, also as a `time.Duration`, last of all.

If 5 left at 3 a week comes out as 168h0m0s, the division ran before the multiplication and kept only whole weeks. If an overspent budget gives false, check `remaining` before `lostLast7d`.
