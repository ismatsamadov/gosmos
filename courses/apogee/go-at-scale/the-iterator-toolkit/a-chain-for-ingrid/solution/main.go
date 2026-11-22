package main

import (
	"fmt"
	"slices"
	"strings"
)

// firstWave is the five stations that went live on Monday 2 April.
var firstWave = []string{"TIR", "VSK", "CAL", "BRM", "OST"}

// furrowAtNew reports whether c is a Furrowlight contact at a
// station of the first wave. Furrowlight's satellites are the
// Furrows.
func furrowAtNew(c Contact) bool {
	furrow := strings.HasPrefix(c.SatID, "furrow-")
	return furrow && slices.Contains(firstWave, c.Station)
}

// line is c as the Board shows it: AOS, pass ID and dish.
func line(c Contact) string {
	aos := c.AOS.Format("15:04:05")
	return fmt.Sprintf("%s %s dish %d", aos, c.ID, c.Dish)
}

func main() {
	monday := slices.Values(firstWaveDay().Contacts)
	lines := Take(Map(Filter(monday, furrowAtNew), line), 5)
	fmt.Println(strings.Join(slices.Collect(lines), "\n"))
}
