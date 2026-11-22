Ask `n.Store` for the booking with `Get`: when it is missing, the error already names it, so return that error as it is. `n.owner` finds the customer that operates the booking's satellite.

Fill an `Event` field by field: the `ID` formats `n.events.Add(1)` as the example does, the `Type` is `"data.landed"`, `BookingID`, `SatID` and `Station` come from the booking, and `File`, `Bytes` and `LandedAt` from the landing.

Hand `n.Sender.Deliver` a `*Delivery` for the customer's ID and the event, and return what it returns. If the second landing gets `evt-000001` again, the counter never moved: `Add(1)` both moves it and returns the new number.
