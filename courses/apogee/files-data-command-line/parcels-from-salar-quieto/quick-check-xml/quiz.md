## Which XML fills Station with CWR?

- [ ] `<stow><station>CWR</station></stow>`
  > Station is tagged station,attr, so it reads an attribute; a station element matches no field and is skipped.
- [x] `<stow station="CWR"></stow>`
  > The ,attr flag reads the station attribute of the element the struct stands for, here stow.
- [ ] `<stow Station="CWR"></stow>`
  > Names are matched exactly, and the tag names the attribute station, in lowercase.

## A file holds `<stow station="CWR"><gust_kmh>81</gust_kmh></stow>`. What does xml.Unmarshal return?

- [ ] An error for the unknown gust_kmh element
  > Unmarshal skips any element that no field asks for, and it has no switch that refuses them.
- [ ] An error for the missing wind_kmh element
  > Unmarshal never requires an element; checking what must be there is the caller's job.
- [x] No error, with WindKmh left at 0
  > The unknown element is skipped, and the missing one leaves its field at the zero value.

## What does `xml.MarshalIndent(s, "", "  ")` write first?

- [x] `<stow station="CWR">`
  > The root element comes first, named by the XMLName tag, with Station as its attribute.
- [ ] `<?xml version="1.0" encoding="UTF-8"?>`
  > Marshaling never writes that line; it is xml.Header, which you print yourself.
- [ ] `<Stow station="CWR">`
  > The type's name is used only when there is no XMLName field to name the root.
