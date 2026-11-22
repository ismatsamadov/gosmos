## Why does every satellite get `k-2028-10` the week before uplinkd seals with it?

- [x] So no command sealed after the switch reaches a satellite that lacks the key
  > Verifiers get a version first, so the first command sealed at 00:00:00 finds every satellite ready.
- [ ] So the satellites can start sealing with it a week early
  > Satellites check commands and seal nothing, so `Activate` doesn't apply to them.
- [ ] Because a key must be a week old before anyone may use it
  > No such rule exists: the week is for loading the key, each satellite on its first pass.

## After 17 October `r-2028-08` signs no more receipts. Why does it stay in the customer guide?

- [ ] So bookingd can go back to it if `r-2028-10` is ever revoked
  > `r-2028-08` retires at 00:00 on the 17th, and a retired version never signs again.
- [x] Receipts it signed before then must still verify
  > A receipt is judged by its `signed_at`, so customers need the old public key for as long as they keep those receipts.
- [ ] Customers' tools refuse a guide that lists only one key
  > A guide may list one key; old ones stay for the receipts they signed.

## How is a revoked key refused, compared with a retired one?

- [ ] Both are refused from their retire instant
  > Only retirement waits for its instant; revocation waits for nothing.
- [ ] A revoked key stays accepted until the overlap ends
  > The overlap protects a key replaced on schedule, never one that is revoked.
- [x] A revoked key is refused at once, a retired one from its retire instant
  > Revocation wins over every window the moment it is known, and it is never undone.
