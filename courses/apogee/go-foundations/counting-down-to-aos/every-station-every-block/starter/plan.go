package main

// stationCode returns the code of station i, from 0 for Ystmark to 4 for Cape Wren.
func stationCode(i int) string {
	switch i {
	case 0:
		return "YST"
	case 1:
		return "PEL"
	case 2:
		return "SQT"
	case 3:
		return "LAN"
	case 4:
		return "CWR"
	default:
		return ""
	}
}

// visible reports whether a station sees Pipit-3 in a two-hour block; block 0 starts at 00:00 UTC.
func visible(station, block int) bool {
	switch station {
	case 0: // Ystmark sees almost every orbit
		return block != 2 && block != 6 && block != 9
	case 1: // Port Ellery
		return block <= 1 || (block >= 5 && block <= 7) || block == 11
	case 2: // Salar Quieto
		return block == 2 || block == 3 || block == 8 || block == 9
	case 3: // Lantern Atoll
		return block == 2 || block == 8
	case 4: // Cape Wren
		return (block >= 1 && block <= 3) || (block >= 7 && block <= 9)
	default:
		return false
	}
}
