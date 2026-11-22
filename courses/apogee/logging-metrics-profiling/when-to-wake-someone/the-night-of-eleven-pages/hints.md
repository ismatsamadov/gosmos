Keep the three names and change the settings, one rule per row of Ingrid's table. A field left out of a `Rule` literal is zero, which is what `For` 0 and `Repeat` 0 mean.

A duration such as `90 * time.Second` is a number times a unit of `time`. `failed-pass` is a `Warning`, and the other two rules are `Critical`.

If the replay still pages Tomás at 19:42, `failed-pass` is still Critical; if it pages at 00:46 and not 20:21, `pass-failures` still fires at 3; if its RESOLVED comes before 08:45, its `ClearBelow` is still above 1.
