## A pass climbs to 40° at a station that can serve its satellite. Why can it still be dropped?

- [ ] Passes under 45° are dropped first when a day is busy
  > The rule keeps every candidate from 10° up and never ranks them by height; 40° is well above the line.
- [x] Every dish of its station is still busy with an earlier contact at its AOS
  > A pass takes the first dish free at its AOS; when `next` holds no such dish, `d` is -1 and the pass is dropped, however high it climbs.
- [ ] Its satellite already has a contact at another station then
  > Stations never share a dish, and a satellite can be in contact with two stations at once.

## Why can the rehearsal plan the twenty stations in any order?

- [ ] Because `assign` sorts the passes it is given
  > `assign` expects them already in AOS order and sorts nothing.
- [ ] Because every station has two dishes
  > Most stations have one; the number of dishes has nothing to do with the order.
- [x] Because no two stations share a dish, so one station's plan can't change another's
  > Each station's dishes are its own, so planning Ystmark first or last gives exactly the same contacts.

## Why must a booking match a contact's AOS to the second, and not just to the minute?

- [x] A booking carries no pass ID, so its satellite, station and AOS are all that name its pass
  > Those three are the only link to one pass: a booking a second off names a pass that doesn't exist, and that is worth reporting.
- [ ] bookingd writes every AOS rounded to the minute, so the seconds are always 00
  > bookingd keeps the AOS of the pass that was booked, seconds included, such as 16:03:46.
- [ ] Contacts at one station always start whole minutes apart
  > AOS times fall on any second, and two dishes at one station can start contacts seconds apart.
