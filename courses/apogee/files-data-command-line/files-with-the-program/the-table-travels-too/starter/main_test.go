package main

import "testing"

// table is stations.csv, byte for byte.
const table = `code,name,lat_deg,lon_deg,mask_deg,bands
YST,Ystmark,77.8,2.0,5,S X
PEL,Port Ellery,46.9,-38.1,5,UHF S
SQT,Salar Quieto,-23.4,-104.6,3,S X
LAN,Lantern Atoll,-1.8,82.5,8,UHF S
CWR,Cape Wren,-53.9,-95.0,5,S X
`

func TestStationTable(t *testing.T) {
	// An empty folder, like a bare stick: a table read from the disk at run time is not here.
	t.Chdir(t.TempDir())
	if got := stationTable(); got != table {
		t.Errorf("stationTable() = %q, want %q", got, table)
	}
}
