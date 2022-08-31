package main

import (
	"calendar"
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:4]
	y, _ := strconv.Atoi(args[0])
	m, _ := strconv.Atoi(args[1])
	d, _ := strconv.Atoi(args[2])
	x := calendar.Day(1988, 1, 9, 3, y, m, d)
	fmt.Println(x)
}
