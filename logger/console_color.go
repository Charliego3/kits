package logger

const (
	Reset  = "\u001B[0m"
	Red    = "\u001B[1;31m"
	Blue   = "\u001B[1;34m"
	Yellow = "\u001B[1;33m"
	Bold   = "\u001B[1m"
)

func red(msg string) string {
	return Red + msg + Reset
}

func blue(msg string) string {
	return Blue + msg + Reset
}

func yellow(msg string) string {
	return Yellow + msg + Reset
}
