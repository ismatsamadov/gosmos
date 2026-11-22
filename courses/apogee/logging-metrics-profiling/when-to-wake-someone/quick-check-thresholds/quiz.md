## From 00:45 to 02:45, `pass-failures` paged at every LOS. Why?

- [ ] Each LOS brought a new problem to page about
  > It was one problem all night; only the rule's check came round again.
- [x] It checks again at each LOS and keeps no memory of having paged
  > A rule with no state pages every time its check is true, so one problem pages as often as it is checked.
- [ ] The station sent each pass's result twice
  > Each result arrived once; the repeats came from the rule.

## What does holding `disk-full` for 15 minutes trade?

- [x] A real problem pages 15 minutes later, and a disk that only touches 90% never pages
  > A recording that lands and is collected within the quarter hour wakes nobody, while a disk that stays full still pages, 15 minutes on.
- [ ] Nothing: real problems page as early as before, and only false pages go
  > A real problem waits out the 15 minutes too, before anyone hears of it.
- [ ] Fewer pages, but only while the disk is under 90%
  > Under 90% the rule never fires at all; the hold only matters at or above it.

## At 00:45:26, one failed pass paged Tomás twice. Why?

- [ ] The pager sent each page twice, to be sure
  > The pager sent each page once; there were two pages.
- [ ] Two passes failed at the same moment
  > Only `cirrine-1/LAN/17147` ended then.
- [x] Two rules, `pass-failures` and `lan-failed-pass`, watch the same failures
  > Each rule pages on its own, so two rules about one problem page twice; the new `failed-pass` only puts a line on the morning list.
