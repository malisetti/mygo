package main

import (
	"calendar"
	"fmt"
)

func main() {
	c := calendar.New()
	if !c.IsWeekend() {
		panic("today is weekend")
	}
	const today = "2022-28-08"
	if x := c.Today(); x != today {
		fmt.Printf("wanted %s but got %s", today, x)
	}

	fmt.Println(c.Day())
}
