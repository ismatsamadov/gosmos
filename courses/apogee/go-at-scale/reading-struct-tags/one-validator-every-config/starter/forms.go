package main

// hostForms are the fifteen hosts' forms, in table order, each as
// its host first sent it.
var hostForms = []HostConfig{
	form("TIR", "Tirrin Rock", 69.5, -12, 5, 2, "S", "X"),
	form("VSK", "Varsk Island", 51, -160, 7, 1, "S", "X"),
	form("CAL", "Calidra Atoll", 0.5, -150, 8, 1, "S"),
	form("BRM", "Ilha Brumosa", -30, -15, 5, 1, "S", "X"),
	form("OST", "Ostrel Island", -42.5, 70, 5, 1, "S", "X"),
	form("ILV", "Ilvane Spit", 72, -165, 3, 2, "S", "X"),
	form("CSL", "Cape Solvane", 35, -45, 5, 1, "S", "X"),
	form("SRH", "Saltreach Island", 17, 135, 6, 1, "S"),
	form("CRP", "Corith Point", 10, 65, 6, 1, "S"),
	form("NIM", "Nimrie Point", -45, -150, 5, 1, "S", "X"),
	form("HSK", "Hasku Point", 32, 170, 4, 1, "S", "X"),
	form("ICM", "", 8, -30, 7, 1, "S", "K"),
	form("QND", "Quindle Island", -25, 95, 5, 1, "S", "X"),
	form("PVT", "Punta Ventisca", -60.5, 40, 5, 1, "S", "X"),
	form("VRL", "Vorrel Island", -64, -140, 4, 0, "S", "X"),
}

// form is a host's form: the station's code and name, its latitude,
// longitude and mask in degrees, its dishes and its bands.
func form(code, name string, lat, lon, mask float64, dishes int,
	bands ...string) HostConfig {
	return HostConfig{Code: code, Name: name, LatDeg: lat,
		LonDeg: lon, MaskDeg: mask, Dishes: dishes, Bands: bands}
}
