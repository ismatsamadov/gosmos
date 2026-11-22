Which kind of loop gives you each station's position, so you can point at the element in the slice instead of at a copy of it?

Loop over the indexes of `stations`. When the station at an index has the code you want, return the address of that element of the slice, as the example does with `sats`. Return `nil` only after the loop has checked every station.
