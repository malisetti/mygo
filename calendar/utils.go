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

func DaysBtwDates(y0, m0, d0, y1, m1, d1 int) int {
	days := 0
	for i := y0; i < y1; i++ {
		days += 365
		if IsLeap(i) {
			days += 1
		}
	}

	for i := 0; i < len(months); i++ {
		if m0 > i {
			days -= DaysInMonth(y0, i)
		}
		if m1 > i {
			days += DaysInMonth(y1, i)
		}
	}

	days -= d0
	days += d1

	return days
}
