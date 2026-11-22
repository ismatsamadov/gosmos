Register the defer before calling the decoder. Call `recover` inside that deferred function itself, not a helper it calls.

Change the named `err` result when a panic value is present. A normal return leaves it nil; a panic interrupts the assignment from `decode`.
