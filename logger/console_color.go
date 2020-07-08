package logger

const (
	ResetColor  = "\u001B[0m"
	RedColor    = "\u001B[1;31m"
	BlueColor   = "\u001B[1;34m"
	YellowColor = "\u001B[1;33m"
)

func red(msg string) string {
	return RedColor + msg + ResetColor
}

func blue(msg string) string {
	return BlueColor + msg + ResetColor
}

func yellow(msg string) string {
	return YellowColor + msg + ResetColor
}
