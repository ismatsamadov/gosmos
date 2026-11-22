Which voltages sit right at the edges of the bands, and which band does the real `batteryBand` give each of them?

A small table of voltages and wanted bands, with a loop over it, keeps `checkBand` short: one `t.Errorf` inside the loop reports every wrong answer.

Check each threshold from both sides: at 8.2 V the band is FULL, and 0.01 V below, at 8.19 V, it is OK. Do the same at 7.2 and 6.6 V.
