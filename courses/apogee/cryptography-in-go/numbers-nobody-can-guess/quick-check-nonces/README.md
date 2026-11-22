---
slug: quick-check-nonces
title: Quick Check: Nonces
type: quiz
position: 6
difficulty: 6
topics: crypto, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Thursday. The counter is gone from Nkem's sketch, and every nonce now comes from `crypto/rand`. Before signing off, Mirela looks at the sketch the way anyone listening to the uplink will: every locked command goes up with its 12-byte nonce beside it, in the open. She has two questions, one about what that nonce must be and one about what the counter got wrong.

## Your task

Answer Mirela's two questions: what is true of a random 12-byte nonce under one key, and what a counter kept in memory hands out after the program that sends the commands up, uplinkd, restarts.
