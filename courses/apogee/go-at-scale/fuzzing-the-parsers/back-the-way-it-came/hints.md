Keep the value checks. After them, call `formatLine(f)` on the frame `parse` returned, and compare what it writes with `line`.

When the two differ, fail with `t.Errorf`, showing the line and what `formatLine` wrote for it, so the difference shows at a glance.

If a mutant survives, the round trip doesn't run on every accepted line, or compares less than the whole line: compare the two strings with `!=`, after the error check and nowhere else.
