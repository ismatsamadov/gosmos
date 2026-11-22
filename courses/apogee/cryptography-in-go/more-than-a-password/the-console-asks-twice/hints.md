Order matters: until `g.console.Check` passes, try no code at all, or a wrong password could spend a recovery code. A locked account is refused before any code is tried, for the same reason.

Ask `g.codes.Check(user, code)` first, and `g.recovery.Use(user, code)` only when it says no, since `Use` spends what it accepts. `g.misses[user]++` works on a missing entry, and `delete(g.misses, user)` clears one.

If a wrong password locks Mirela out, the miss is counted before `g.console.Check`; if the old password still opens a session after a reset, `ResetPassword` went on past `Redeem`'s error or skipped `setPassword`.
