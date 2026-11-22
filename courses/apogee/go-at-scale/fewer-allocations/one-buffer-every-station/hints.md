Move the line that makes `buf` above the loop over the stations: from then on, the loop only reuses it.

At the top of each turn, set `buf`'s length back to 0 with a slice expression, so the new station's lines start at the front of the array.

If a station's write holds the stations before it too, `buf` wasn't emptied; if the count is still near twenty, something inside the loop still makes a slice.
