package main

import "slices"

// changes compares the contacts the two planners have for one
// second, matched by pass ID: a contact only the old planner has
// is "old only", one only the new planner has is "new only", and
// one the new planner flies on another dish is "dish". It reports
// them in the planners' order, the old planner's first.
func changes(olds, news []Contact) []Change {
	is := func(id string) func(Contact) bool {
		return func(c Contact) bool { return c.ID == id }
	}
	var out []Change
	for _, o := range olds {
		i := slices.IndexFunc(news, is(o.ID))
		if i < 0 {
			out = append(out, Change{"old only", o})
		} else if news[i].Dish != o.Dish {
			out = append(out, Change{"dish", news[i]})
		}
	}
	for _, n := range news {
		if !slices.ContainsFunc(olds, is(n.ID)) {
			out = append(out, Change{"new only", n})
		}
	}
	return out
}
