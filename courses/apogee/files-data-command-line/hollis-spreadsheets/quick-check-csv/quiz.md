## How many cells does strings.Split make of bk-2117's line, cutting at each comma?

- [ ] 3
  > That is what a csv.Reader returns: it keeps the quoted comma inside the note.
- [x] 4
  > Split cuts at every comma, the one inside the quotes included.
- [ ] 2
  > There are three commas on the line, and Split cuts at each of them.

## How many records does a csv.Reader return for these lines?

- [ ] 4, one per line
  > The line break inside bk-2118's quoted note belongs to the note.
- [ ] 2, because the header is left out
  > The reader returns the header as an ordinary first record.
- [x] 3: the header and two bookings
  > bk-2118's note holds a line break inside quotes, so its two lines make one record.

## A header says STATION in two columns. What should the reader do?

- [x] Refuse the sheet with an error naming STATION
  > Either column could be the right one, so any choice would be a guess.
- [ ] Use the first STATION column
  > The second could hold the stations Hollis meant, and nothing in the sheet says which.
- [ ] Use the last STATION column
  > That is a guess too, and a wrong guess files bookings under the wrong station.
