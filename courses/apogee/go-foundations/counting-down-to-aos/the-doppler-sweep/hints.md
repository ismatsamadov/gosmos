What should the strongest-so-far start as? Strengths here are negative dBm, so a starting value of 0 would beat every real reading.

Start the best offset at `fromHz`, with its own strength. Then loop with a three-part `for` that adds `stepHz` each round while the offset is at most `toHz`, and replace the best only when a strength is greater.

Shape: two variables for the best offset and its strength, set from `fromHz` before the loop; inside it, `if ___ > bestDBm {` updates both; after it, return the offset. `>` rather than `>=` keeps the first of a tie.
