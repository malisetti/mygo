package calendar

func IsLeap(x int) bool {
	return (x%4 == 0) && (x%100 != 0 || x%400 == 0)
}

func DaysInYear(x int) int {
	if IsLeap(x) {
		return 366
	}
	return 365
}

func DaysInMonth(year, month int) int {
	if month == 1 {
		if !IsLeap(year) {
			return monthNdays[months[1]][0]
		}
		return monthNdays[months[1]][1]
	}

	return monthNdays[months[month]][0]
}

func CompareDates(y0, m0, d0, y1, m1, d1 int) int {
	y0v := y0*100 + m0*10 + d0
	y1v := y1*100 + m1*10 + d1

	if y0v == y1v {
		return 0
	}

	if y0v < y1v {
		return -1
	}

	return 1
}

func DaysBtwDates(y0, m0, d0, y1, m1, d1 int) (int, int) {
	var x, xm, xd, y, ym, yd int
	cv := CompareDates(y0, m0, d0, y1, m1, d1)
	switch cv {
	case -1:
		x = y0
		xm = m0
		xd = d0

		y = y1
		ym = m1
		yd = d1
	case 1:
		x = y1
		xm = m1
		xd = d1

		y = y0
		ym = m0
		yd = d0
	default:
		return 0, cv
	}

	days := 0
	for i := x; i < y; i++ {
		days += 365
		if IsLeap(i) {
			days += 1
		}
	}

	for i := 0; i < len(months); i++ {
		if xm > i {
			days -= DaysInMonth(x, i)
		}
		if ym > i {
			days += DaysInMonth(y, i)
		}
	}

	days -= xd
	days += yd

	return days, cv
}
