`chain` runs the middleware in the order you list it, the first one outermost. In `private`, the guard that finds the customer has to be listed before the one that counts the customer's requests.

In `handler`, add one more argument to `chain`, after `noStore`: the last one listed sits right around the mux, inside the request ID, the log and the panic recovery.

If every private request answers 500, `limit` still runs before `authenticate` and finds no customer. If the stalled store's request takes half a second and answers 500, no deadline reached the handlers: `withDeadline` is missing from `handler`.
