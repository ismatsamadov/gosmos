---
slug: passwords-at-rest
title: Passwords at Rest
position: 3
---
Salt passwords, slow every guess with `pbkdf2.Key`, store self-describing records, compare with `subtle.ConstantTimeCompare`, upgrade old records at login, and refuse unknown users and wrong passwords alike.
