## What makes these six digits different from the console's old six-digit codes?

- [ ] HMAC-SHA-1 gives each code 160 bits of strength.
  > Six digits hold 19.9 bits however they are made; the 160 bits are in the secret behind them, which only the phone and the console hold.
- [x] Each code counts only beside the right password, within its 90 seconds, and once.
  > The old codes were the whole secret and never changed. A new code is a second proof that expires with its window and can't be used again.
- [ ] The console stores them hashed, so a stolen copy reveals none.
  > The console stores no codes at all: it makes them from each secret when it checks one.

## Why does the console also accept the step before and the step after?

- [x] Clocks drift apart, and reading and typing a code takes seconds.
  > A phone 20 seconds slow, like Ingrid's, shows the code of the step before for most of each step.
- [ ] Each extra step makes a code harder to guess.
  > The opposite: each step accepted is one more code a guess could match, which is why it stops at one either side.
- [ ] The phone shows three codes at once.
  > An authenticator app shows one code, the one for the current step by its own clock.

## Why does the console keep each phone's secret itself, when it keeps only digests of passwords and tokens?

- [ ] A 20-byte secret is too short to hash.
  > SHA-256 takes input of any length; the secret's size has nothing to do with it.
- [ ] A digest of the secret would let anyone make the codes.
  > A digest makes no codes at all, which is exactly why it can't stand in for the secret.
- [x] To check a code, it must make the same code the phone makes, from the secret itself.
  > So the secret is kept whole, and guarded like a key.
