`timeStep(c.now())` gives the console's step. Loop from that number minus 1 to that number plus 1, and ask `hotp` for each step's code.

Compare each with `subtle.ConstantTimeCompare`, as the draft does, and return true at the first match; after the loop, return false.

Keep the draft's first lines: an account without a secret has no code. Without them the kiosk's secret is nil, and the code anyone can make from an empty secret would pass.
