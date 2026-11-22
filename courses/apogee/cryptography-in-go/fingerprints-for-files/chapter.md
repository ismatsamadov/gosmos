---
slug: fingerprints-for-files
title: Fingerprints for Files
position: 1
---
Hash bytes with `sha256.Sum256`, streams through `hash.Hash` with `io.Copy` and `io.MultiWriter`, pick SHA-256 or SHA3-256 by name, and check delivered files against a manifest.
