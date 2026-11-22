---
slug: sealed-messages
title: Sealed Messages
position: 5
---
Encrypt and authenticate with `aes.NewCipher` and `cipher.NewGCM`, give every message a fresh nonce, bind a readable header as associated data, and write, parse and open the sealed command line in a fixed order.
