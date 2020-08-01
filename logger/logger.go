package logger

import (
	"fmt"
	"log"
)

const (
	pl = "[PANIC] "
	el = "[ERROR] "
	wl = "[WARN ] "
	dl = "[DEBUG] "
	il = "[INFO ] "
)

var (
	async = false
)

func Async() {
	async = true
}

func Panic(format string, v ...interface{}) {
	output(Red, pl, format, v...)
}

func Error(format string, v ...interface{}) {
	output(Red, el, format, v...)
}

func Warn(format string, v ...interface{}) {
	output(Yellow, wl, format, v...)
}

func Debug(format string, v ...interface{}) {
	output(Blue, dl, format, v...)
}

func Info(format string, v ...interface{}) {
	output(Bold, il, format, v...)
}

func R(format string, v ...interface{}) {
	output(Red, "", format, v...)
}

func Y(format string, v ...interface{}) {
	output(Yellow, "", format, v...)
}

func B(format string, v ...interface{}) {
	output(Blue, "", format, v...)
}

func output(color, level, format string, v ...interface{}) {
	var message string
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = format
	}
	message = color + level + message + Reset

	if level == pl {
		log.Panicln(message)
	} else {
		log.Println(message)
	}
}
