## A copy of Pipit-2's sign key leaves the building. What can its holder command?

- [ ] Every satellite, since all their keys come from one master
  > A derived key leads back to neither the master nor any other satellite's keys, so it signs for Pipit-2 alone.
- [x] Pipit-2 alone
  > Each satellite's keys are derived with its own ID in the info, so a copy of Pipit-2's checks out on no other satellite.
- [ ] Nothing, since every command also needs the right number
  > The number is no secret: anyone can count one past the last command they saw go up.

## A PING is in flight when the line drops, and the next WELCOME's `last` is below its number. What makes it run exactly once?

- [ ] Sending it again under the next free number
  > If the first copy still reaches the satellite, both numbers run: the PING runs twice.
- [ ] Sending nothing more until the satellite's next contact
  > It never runs twice, but if the first copy never arrived it never runs at all, and the pass fails `no_ack`.
- [x] Sending the same line again, under the same number
  > The satellite runs one command per number: if the first copy got there after all, the second is refused as replayed.

## Why does `k-2028-10` retire a week after `k-2029-07` activates, and not at the switch?

- [x] So a satellite that missed its LOAD_KEY can still be commanded, and loaded, under the old master
  > Until 1 August it holds a master that works, so its next contact can still bring up `k-2029-07`; after that it can't be commanded at all.
- [ ] Because a satellite needs a week to put a new key to use
  > A satellite holds a master's keys the moment their LOAD_KEY runs, and they are used from the instant the master activates.
- [ ] Because every command that week is signed under both masters
  > A command is signed or sealed under one master: the newest that has activated and that its satellite holds.
