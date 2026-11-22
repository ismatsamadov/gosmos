Create the temporary file beside the destination, not in the system temporary directory. Arrange its removal as soon as it exists.

Only rename after both writing and closing have succeeded. Return the first operation error and leave the destination alone until the rename.

Deferred close and removal handle early returns. After a successful rename, the old temporary name is gone, so its deferred removal has nothing left to remove.
