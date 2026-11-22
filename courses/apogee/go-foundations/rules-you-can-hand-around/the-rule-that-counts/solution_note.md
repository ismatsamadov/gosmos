The closure captures the count variable, not a frozen copy of its initial value. Each call to `makeCounter` creates a separate variable, while repeated calls to one returned function share its count.
