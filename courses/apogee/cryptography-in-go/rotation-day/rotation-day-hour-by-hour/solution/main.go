package main

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// runDay flies every slot of d in AOS order and reports each pass.
func runDay(d Day) Report {
	slots := slices.Clone(d.Slots)
	slices.SortFunc(slots, func(a, b Slot) int { return a.AOS.Compare(b.AOS) })
	var r Report
	for _, s := range slots {
		head := fmt.Sprintf("%s %s %s", s.AOS.Format("15:04:05"), s.SatID, s.Station)
		tail, err := fly(d, s)
		if err != nil {
			r.Missed++
			r.Lines = append(r.Lines, head+" missed: "+err.Error())
			continue
		}
		r.Flown++
		r.Lines = append(r.Lines, head+" flown "+tail)
	}
	return r
}

// fly connects to the slot's station at AOS, then seals and delivers each of its PINGs. It returns the
// key IDs the PINGs were sealed with and the serial of the station's certificate.
func fly(d Day, s Slot) (string, error) {
	served := d.Old[s.Station]
	if switchAt, err := switchTime(d.Passes(), s.Station, d.SwitchFrom); err == nil && !s.AOS.Before(switchAt) {
		served = d.New[s.Station]
	}
	var list []byte
	if !s.AOS.Before(d.ListAt) {
		list = d.List
	}
	if err := d.Connect(s.Station, served, list, s.AOS); err != nil {
		return "", err
	}
	var ids []string
	for _, sealAt := range s.SealAt {
		v, err := signingKey(d.Ring.Versions(), sealAt)
		if err != nil {
			return "", err
		}
		if err := d.Ping(s.SatID, v.ID, sealAt.Add(5*time.Second)); err != nil {
			return "", err
		}
		ids = append(ids, v.ID)
	}
	return fmt.Sprintf("%s %v", strings.Join(ids, "+"), served.Leaf.SerialNumber), nil
}
