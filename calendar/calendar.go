package calendar

func Day(y0, m0, d0, dd, y1, m1, d1 int) string {
	n, c := DaysBtwDates(y0, m0, d0, y1, m1, d1)
	switch c {
	case -1:
		cd := append(days[dd:], days[:dd]...)
		return cd[n%7]
	case 1:
		cd := append(days[dd:], days[:dd]...)
		for i, j := 0, len(cd)-1; i < j; i, j = i+1, j-1 {
			cd[i], cd[j] = cd[j], cd[i]
		}
		xn := n % 7
		if xn == 0 {
			return cd[len(cd)-1]
		}
		return cd[xn-1]
	default:
		return days[dd]
	}
}
