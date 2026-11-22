## Frames 68841, 68843 and 68842 arrive in that order. What does the live view show?

- [ ] 68841, 68842 and 68843: it puts them back in order
  > It never holds a frame back to wait for a straggler: 68843 is shown when it comes, and 68842 is older than that.
- [x] 68841 and 68843; 68842 is dropped
  > 68842 arrives after a newer frame was shown, so it's dropped; it was counted as missed when 68843 jumped over it.
- [ ] 68841, 68843, then 68842
  > Showing 68842 after 68843 would move the screen back in time.

## Why does the second try ask question 7 again, not question 8?

- [ ] So the recorder knows the first question was lost
  > The recorder can't tell and doesn't need to: it answers every question it gets.
- [ ] Because the recorder refuses a number it has seen before
  > The recorder answers any number, the same one twice included; the number is for the asker.
- [x] So an answer to the first try that arrives late still counts
  > Both tries ask the same thing, so an answer to either is the one the live view waits for.

## `STATUS 6 IDLE` arrives from the recorder while question 7 waits. What does `readAnswer` do?

- [ ] Returns `IDLE`: it came from the recorder
  > The right sender isn't enough: the number says it answers an older question, which may be a minute old.
- [x] Skips it and keeps reading until 7's answer or the deadline
  > It answers a question the live view already gave up on, and question 7's deadline stays where it was.
- [ ] Gives up on question 7 and asks again at once
  > Nothing has been lost yet: question 7's answer may be the very next datagram.
