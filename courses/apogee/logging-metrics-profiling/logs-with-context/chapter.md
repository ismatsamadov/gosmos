---
slug: logs-with-context
title: Logs With Context
position: 2
---
Carry pass and request IDs in `context.Context`, write a `slog.Handler` that adds them to every record and survives `With`, pass them between programs in a header, and read the lines back in order.
