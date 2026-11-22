`%w` prints the wrapped error's text the same way `%v` does, so the message is unchanged, but it also keeps `ErrDuplicate` inside the new error, where `errors.Is` finds it.
