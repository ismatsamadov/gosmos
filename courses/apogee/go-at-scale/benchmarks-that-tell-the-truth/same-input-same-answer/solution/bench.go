package main

import "fmt"

// sameContacts returns nil when old and new hold the same contacts,
// by pass ID and dish, in the same order, and otherwise an error
// that names the first difference.
func sameContacts(old, new []Contact) error {
	if len(old) != len(new) {
		return fmt.Errorf("old has %d contacts, new has %d", len(old), len(new))
	}
	for i := range old {
		o, n := old[i], new[i]
		if o.ID != n.ID || o.Dish != n.Dish {
			return fmt.Errorf("contact %d: old %s dish %d, new %s dish %d", i, o.ID, o.Dish, n.ID, n.Dish)
		}
	}
	return nil
}
