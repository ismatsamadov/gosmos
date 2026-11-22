A `switch` with nothing after the word `switch` runs the first `case` whose condition is true, so the order of the cases matters: 429 has to come before the case for the rest of the 4xx.

Each class is a range that includes both of its ends, so compare with `>=` and `<=`, joined by `&&`. A `default` case catches every code outside the three classes.

If 200, 299 or 599 fails, check the ends of your ranges: `>` and `<` leave the end itself out. If 429 says `fix the request`, its case comes too late.
