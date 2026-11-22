package main

import (
	"iter"
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
				yield(p) // TODO: stop once the loop has stopped
			}
			from = to
		}
	}
}
