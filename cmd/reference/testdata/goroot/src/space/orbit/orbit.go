// Package orbit plans the passes of a craft over a ground station.
//
// # Passes
//
// A pass starts when the craft rises above the horizon; see [Craft.Raise],
// [relay.Signal] and [strings.Builder]. The [Go spec] has the rules, and
// https://example.com/orbits explains more.
//
//	c := orbit.NewCraft("Wren")
//	c.Raise(2)
//
// [Go spec]: https://go.dev/ref/spec#Types
package orbit

import "space/relay"

// Altitudes a craft can hold, in km.
const (
	Low  = 400 // low orbit
	High = 1200
	low  = 1
)

// Drift is how far a craft wanders each pass.
var Drift = 3

// A Craft flies over the station. Its [Craft.Name] is shown & logged with ~care.
type Craft struct {
	Name string // what the station calls it
	fuel int
}

// NewCraft returns a craft called name.
func NewCraft(name string) *Craft { return &Craft{Name: name} }

// Raise lifts the craft by km; see [relay.Signal.Send].
func (c *Craft) Raise(km int) { c.fuel -= km }

// Plan plans a pass.
//
// Deprecated: Use [Craft.Raise] instead.
func Plan() {}

func hidden() relay.Signal { return nil }
