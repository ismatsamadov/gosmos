Only the last line needs to change. The time `t` is already the right instant; it is still set to the offset the text gave.

Convert `t` to UTC before formatting it. Rewriting the text, such as swapping `-07:00` for `Z`, keeps the van's hour and names the wrong instant.
