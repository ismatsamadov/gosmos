Walk every version. One may sign at `now` when `now` isn't before its `Activate`, it isn't `Revoked`, and its `Retire` is zero or after `now`.

Among those, keep the one whose `Activate` is latest, comparing with `After`, and remember whether you found any at all: a zero `KeyVersion` is no answer.

If none may sign, return `fmt.Errorf` with `now.Format(time.RFC3339)`. If the reversed ring gives other answers, the code still returns the first version that fits instead of the newest.
