The compiler says `Station` has no field `DiameterM`. Where does that field live now, and what would let `yst` reach it directly?

An embedded field is a type written alone on its line, with no field name in front. Replace each `Equipment Dish` line with one.

Once `Dish` is embedded, the key `Equipment:` no longer exists. In the literals, name the embedded field `Dish`, or set `DiameterM` and `XBand` directly.
