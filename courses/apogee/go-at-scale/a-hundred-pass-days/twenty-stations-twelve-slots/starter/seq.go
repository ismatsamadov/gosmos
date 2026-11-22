package main

import (
	"iter"
	"slices"
	"time"
)

// Predictor predicts a satellite's passes over a station with AOS in [from, to), by AOS.
type Predictor func(e Elements, s Site, from, to time.Time) []PassView

// Passes yields e's passes over s in w, by AOS, asking predict for one day at a time.
func Passes(predict Predictor, e Elements, s Site, w Window) iter.Seq[PassView] {
	return func(yield func(PassView) bool) {
		for from := w.From; from.Before(w.To); {
			to := from.AddDate(0, 0, 1)
			if w.To.Before(to) {
				to = w.To
			}
			for _, p := range predict(e, s, from, to) {
				if !yield(p) {
					return
				}
			}
			from = to
		}
	}
}

// Catalog yields every pass in w at a station that can serve its
// satellite, in byAOS order, predicting and sorting one day at a
// time.
func Catalog(n *Network, w Window, predict Predictor) iter.Seq[PassView] {
	return func(yield func(PassView) bool) {
		for from := w.From; from.Before(w.To); {
			to := from.AddDate(0, 0, 1)
			if w.To.Before(to) {
				to = w.To
			}
			var day []PassView
			for i, e := range n.Elements {
				for _, s := range n.Sites {
					if canServe(s, n.Satellites[i]) {
						ps := predict(e, s, from, to)
						day = append(day, ps...)
					}
				}
			}
			slices.SortFunc(day, byAOS)
			for _, p := range day {
				if !yield(p) {
					return
				}
			}
			from = to
		}
	}
}

// Bookable yields the passes Catalog would yield for the operator's
// satellites alone, predicting no other satellite.
func Bookable(n *Network, w Window, predict Predictor, operator string) iter.Seq[PassView] {
	fleet := &Network{Sites: n.Sites}
	for i, sat := range n.Satellites {
		if sat.Operator == operator {
			fleet.Satellites = append(fleet.Satellites, sat)
			fleet.Elements = append(fleet.Elements, n.Elements[i])
		}
	}
	return Catalog(fleet, w, predict)
}
