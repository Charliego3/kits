package logger

import (
	"fmt"
	"log"
)

const (
	panicl = "[PANIC] "
	errorl = "[ERROR] "
	warnl  = "[WARN ] "
	debugl = "[DEBUG] "
	infol  = "[INFO ] "
)

var (
	async = false
)

func Async() {
	async = true
}

func Panic(format string, v ...interface{}) {
	logger(RedColor, panicl, format, v...)
}

func Error(format string, v ...interface{}) {
	logger(RedColor, errorl, format, v...)
}

func Warn(format string, v ...interface{}) {
	logger(YellowColor, warnl, format, v...)
}

func Debug(format string, v ...interface{}) {
	logger(BlueColor, debugl, format, v...)
}

func Info(format string, v ...interface{}) {
	logger("", infol, format, v...)
}

func R(format string, v ...interface{}) {
	logger(RedColor, "", format, v...)
}

func Y(format string, v ...interface{}) {
	logger(YellowColor, "", format, v...)
}

func B(format string, v ...interface{}) {
	logger(BlueColor, "", format, v...)
}

func logger(color, level, format string, v ...interface{}) {
	var message string
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = format
	}
	message = level + message
	switch color {
	case YellowColor:
		message = yellow(message)
	case BlueColor:
		message = blue(message)
	case RedColor:
		message = red(message)
	}

	if level == panicl {
		log.Panicln(message)
	} else {
		log.Println(message)
	}
}
