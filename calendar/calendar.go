package calendar

import (
	"fmt"
	"time"
)

type Calendar struct {
	time  time.Time
	ctime *Time
}

type Period int

const (
	PCalendar Period = iota
	PYear
	PMonth
	PWeek
	PDay
	PHour
	PMinute
	PSecond
)

type Time struct {
	name int
	kind Period
	next *Time
}

type Year struct {
	name   int
	months []*Month
}

type Month struct {
	name  int
	weeks []*Week
}

type Week struct {
	name int
	days []*Day
}

type Day struct {
	name  int
	hours []*Hour
}

type Hour struct {
	name    int
	minutes []*Minitue
}

type Minitue struct {
	name    int
	seconds []*Second
}

type Second struct {
	name     int
	duration time.Duration
}

var months = [12]string{
	time.January.String(),
	time.February.String(),
	time.March.String(),
	time.April.String(),
	time.May.String(),
	time.June.String(),
	time.July.String(),
	time.August.String(),
	time.September.String(),
	time.October.String(),
	time.November.String(),
	time.December.String(),
}

func daysInMonth(year int, month time.Month) int {
	if month == time.February {
		if isLeap(year) {
			return monthNdays[months[time.February-1]][0]
		}
		return monthNdays[months[time.February-1]][1]
	}

	return monthNdays[months[time.February-1]][0]
}

var days = [7]string{
	time.Sunday.String(),
	time.Monday.String(),
	time.Tuesday.String(),
	time.Wednesday.String(),
	time.Thursday.String(),
	time.Friday.String(),
	time.Saturday.String(),
}

func isLeap(x int) bool {
	if x%4 != 0 {
		return false
	}
	if x%100 != 0 {
		return false
	}
	if x%400 != 0 {
		return false
	}
	return true
}

func daysInYear(x int) int {
	if isLeap(x) {
		return 366
	}
	return 365
}

var monthNdays = map[string][]int{
	months[0]:  {31},
	months[1]:  {28, 29},
	months[2]:  {31},
	months[3]:  {30},
	months[4]:  {31},
	months[5]:  {30},
	months[6]:  {31},
	months[7]:  {31},
	months[8]:  {30},
	months[9]:  {31},
	months[10]: {30},
	months[11]: {31},
}

func nextof(p Period) (Period, int, int) {
	switch p {
	case PYear:
		return PMonth, 0, 12
	case PMonth:
		return PWeek, 4, 5
	case PWeek:
		return PDay, 0, 7
	case PDay:
		return PHour, 0, 24
	case PHour:
		return PMinute, 0, 60
	case PMinute:
		return PSecond, 0, 60
	default:
		return PSecond, 0, 1
	}
}

func makeTime(t time.Time) *Time {
	tx := &Time{
		kind: PCalendar,
		name: int(PCalendar),
		next: makeSubTime(t, PCalendar),
	}
	return tx
}

func makeSubTime(t time.Time, h Period) *Time {
	switch h {
	case PCalendar:
		return &Time{
			name: t.Year(),
			kind: PCalendar,
			next: makeSubTime(t, PYear),
		}
	case PYear:
		return &Time{
			name: t.Year(),
			kind: PYear,
			next: makeSubTime(t, PMonth),
		}
	case PMonth:
		return &Time{
			name: int(t.Month()),
			kind: PMonth,
			next: makeSubTime(t, PWeek),
		}
	case PWeek:
		return &Time{
			name: 0,
			kind: PWeek,
			next: makeSubTime(t, PDay),
		}
	case PDay:
		return &Time{
			name: t.Day(),
			kind: PDay,
			next: makeSubTime(t, PHour),
		}
	case PHour:
		return &Time{
			name: t.Hour(),
			kind: PHour,
			next: makeSubTime(t, PMinute),
		}
	case PMinute:
		return &Time{
			name: t.Minute(),
			kind: PMinute,
			next: makeSubTime(t, PSecond),
		}
	default:
		return &Time{
			name: t.Second(),
			kind: PSecond,
			next: nil,
		}
	}
}

func New() *Calendar {
	n := time.Now()
	return &Calendar{
		time:  n,
		ctime: makeTime(n),
	}
}

func (c Calendar) IsWeekend() bool {
	day := c.ctime.next.next.next.next.name
	if time.Weekday(day) == time.Sunday || time.Weekday(day) == time.Saturday {
		return true
	}

	return false
}

func (c Calendar) Today() string {
	dt := c.ctime.next.next.next.next.next
	mt := c.ctime.next.next.next
	return fmt.Sprintf("%d-%02d-%02d", c.ctime.next.name, dt.name, mt.name)
}

func (c Calendar) IsLeap() bool {
	return isLeap(c.ctime.next.name)
}

func (c Calendar) Day() string {
	return days[c.ctime.next.next.next.next.name]
}
