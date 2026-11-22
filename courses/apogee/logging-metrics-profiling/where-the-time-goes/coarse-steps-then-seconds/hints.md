Only the loop's stepping changes: what it does with a changed answer stays. Move 30 seconds at a time, and keep the moment of the question before, because `crossing` needs both ends of the gap.

The old scan's last question comes one second before `to`. When a 30-second step would go past that second, ask at that second instead, so a pass that sets in the window's last half minute still ends inside it.

If passes come out up to half a minute off, `crossing` still gets moments one second apart: give it the moment of the question before and of this one. If a pass that sets just before `to` is missing, the last question was never asked; if a pass appears that the scan leaves out, a question went to `to` itself.
