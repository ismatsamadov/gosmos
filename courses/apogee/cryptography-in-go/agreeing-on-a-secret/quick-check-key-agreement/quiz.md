## Backhaul company kept both shares. What can it compute?

- [ ] The secret, by combining the two shares
  > Combining takes a private key, and neither private key ever crossed the backhaul.
- [x] Nothing of the secret: each private key stayed home
  > Each end combines its own private key with the other's share, and with today's computers nobody holding only the two shares can repeat that.
- [ ] Half of the secret, one share's worth
  > The secret isn't built from pieces of the shares: without a private key, the shares give nothing of it.

## A share arrives marked YST. What does that prove?

- [ ] That Ystmark sent it, since only Ystmark's key makes that share
  > Anybody can make a key pair of their own and send its share marked YST.
- [ ] That nobody changed it on the way
  > A share swapped on the way looks exactly as valid as the one that was sent.
- [x] Nothing about who sent it
  > A share is 32 bytes anybody can make; proving who sent it takes more, such as a signature only Ystmark can make.

## Which part of a station's key may cross the backhaul?

- [x] Only its share, the public half
  > The share is public by design; the private key never leaves its station.
- [ ] Both halves, sealed under the key they agree on
  > The private key never travels at all, sealed or not: nobody but its own station ever needs it.
- [ ] The private key too, once the share has been checked
  > Whoever read the private key could compute every secret the station agrees on.
