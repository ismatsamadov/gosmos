## Before the WELCOME, what does the missing answer tell dishd about JOG 122?

- [ ] That the jog never reached the controller
  > The answer can be lost on its way back after the jog ran; silence can't tell the two apart.
- [x] Nothing about whether it ran: the jog or only its answer may have been lost
  > A command and its answer cross the line separately, so dishd learns the jog's fate only from the controller's count after the reconnect.
- [ ] That the controller refused the jog
  > A refusal is a NAK, an answer that arrives; here nothing came back at all.

## Why does a command sent again keep its own sequence number?

- [ ] Because the controller accepts only the number right after its `last_seq`
  > It runs any number above `last_seq`: 130 would run as well as 123.
- [ ] Because the controller keeps each command by its number and looks it up
  > The controller judges a number by `last_seq` alone, the number of the last command it ran; there is nothing to look up.
- [x] So that the controller refuses it as a duplicate if it has already run
  > The controller runs only numbers above its `last_seq`, so a command that already ran, sent again under its own number, moves nothing.

## Which command, run twice, leaves the dish somewhere other than where one run would?

- [x] JOG A +0.35°
  > A jog moves the dish from where it is, so a second run adds another 0.35°, as on Tuesday.
- [ ] POINT 183.25° 41.70°
  > POINT turns the dish to a place, and the second run finds it there already.
- [ ] STOW
  > STOW turns the dish straight up and holds it there; a stowed dish stays where it is.
