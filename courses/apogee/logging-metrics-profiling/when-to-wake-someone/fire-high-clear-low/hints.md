Only the first line of `whileFiring` changes: the comparison that decides when a firing rule is over.

A rule that has fired is over when its value falls below `r.ClearBelow`, the lower line, not below `r.FireAt`: compare `s.Value` with it, with `<`.

If exactly 85% resolves the rule, the comparison keeps the boundary: 0.85 isn't below 0.85, so it must stay firing.
