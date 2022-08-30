package calendar

func Day(y0, m0, d0, dd, y1, m1, d1 int) string {
	n := DaysBtwDates(y0, m0, d0, y1, m1, d1)
	cd := append(days[dd:], days[:dd]...)
	return cd[n%7]
}
