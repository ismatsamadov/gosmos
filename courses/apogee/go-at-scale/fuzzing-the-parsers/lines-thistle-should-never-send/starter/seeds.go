package main

// lineSeeds are the lines FuzzParseLine runs on every go test, before it invents any.
var lineSeeds = []string{
	"cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=NOMINAL",
	"cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5",
	"cirrine-1 9012 2029-03-07T20:45:05Z batt_v=nan temp_c=11.5 mode=NOMINAL",
	"cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=nominal",
}
