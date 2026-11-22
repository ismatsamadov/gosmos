package main

import "testing/fstest"

// fixtureFS returns the small in-memory disk that Mirela's tests read instead of the laptop's.
func fixtureFS() fstest.MapFS {
	return fstest.MapFS{
		"stations.csv":               {Data: []byte("code,name,lat_deg,lon_deg,mask_deg,bands\nCWR,Cape Wren,-53.9,-95.0,5,S X\n")},
		"templates/notices/stow.txt": {Data: []byte("{station} stowed: {reason}\n")},
	}
}
