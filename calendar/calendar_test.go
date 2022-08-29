package calendar

import "testing"

func TestCalendar(t *testing.T) {
	c := New()
	if !c.IsWeekend() {
		t.Error("today is weekend")
	}
	const today = "2022-28-08"
	if x := c.Today(); x != today {
		t.Errorf("wanted %s but got %s", today, x)
	}
	if c.IsLeap() {
		t.Error("this year is not leap")
	}
}
