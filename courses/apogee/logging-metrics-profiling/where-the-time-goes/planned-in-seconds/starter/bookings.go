package main

// tuesdayBookings returns the bookings of Tuesday 16 January's passes among passes, in their order.
func tuesdayBookings(passes []PassView) []Booking {
	ids := map[string]string{
		"furrow-a/SQT/22047":  "bk-5788",
		"brine-3/SQT/15440":   "bk-5789",
		"brine-3/CWR/15443":   "bk-5791",
		"cirrine-1/CWR/17437": "bk-5794",
		"cirrine-1/YST/17436": "bk-5795",
		"furrow-b/YST/22050":  "bk-5796",
		"furrow-a/CWR/22054":  "bk-5797",
		"brine-3/LAN/15447":   "bk-5799",
		"cirrine-1/LAN/17438": "bk-5800",
		"brine-3/PEL/15447":   "bk-5801",
		"brine-3/CWR/15442":   "bk-5803",
		"furrow-a/SQT/22055":  "bk-5804",
		"furrow-a/CWR/22055":  "bk-5806",
		"cirrine-1/PEL/17443": "bk-5808",
	}
	var bookings []Booking
	for _, p := range passes {
		if id, ok := ids[p.ID]; ok && p.AOS.Before(epoch.AddDate(0, 0, 1)) {
			bookings = append(bookings, Booking{ID: id, SatID: p.SatID, Station: p.Station, AOS: p.AOS, LOS: p.LOS})
		}
	}
	return bookings
}
