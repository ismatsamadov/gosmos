`Satellite` now has a `String() string` method, so it satisfies `fmt.Stringer`, and `Println` prints what that method returns instead of the struct's fields.
