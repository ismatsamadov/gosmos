A variable declared inside the returned function would start over on every call. Declare the count in `makeCounter` instead.

Each call to `makeCounter` creates its own count. The function it returns must update that captured variable before returning it.
