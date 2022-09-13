package main

import (
	"calendar"
	_ "crypto/sha512"
	"syscall/js"
)

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("Day", js.FuncOf(day))
	<-done
}

func day(this js.Value, args []js.Value) interface{} {
	return calendar.Day(1988, 1, 9, 3, args[0].Int(), args[1].Int(), args[2].Int())
}
