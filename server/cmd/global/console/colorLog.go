package console

import (
	"fmt"
	"log"
)

type MsgType string

const (
	reset = "\x1b[0m"

	// bright     = "\x1b[1m"
	// dim        = "\x1b[2m"
	underscore = "\x1b[4m"
	// blink      = "\x1b[5m"
	// reverse    = "\x1b[7m"
	// hidden     = "\x1b[8m"

	// Foreground Colors
	fgBlack   = "\x1b[30m"
	fgRed     = "\x1b[31m"
	fgGreen   = "\x1b[32m"
	fgYellow  = "\x1b[33m"
	fgBlue    = "\x1b[34m"
	fgMagenta = "\x1b[35m"
	fgCyan    = "\x1b[36m"
	fgWhite   = "\x1b[37m"

	// Background Colors
	// bgBlack   = "\x1b[40m"
	// bgRed     = "\x1b[41m"
	// bgGreen   = "\x1b[42m"
	// bgYellow  = "\x1b[43m"
	// bgBlue    = "\x1b[44m"
	// bgMagenta = "\x1b[45m"
	// bgCyan    = "\x1b[46m"
	// bgWhite   = "\x1b[47m"

	// Type
	Info    MsgType = fgCyan
	Success MsgType = fgGreen
	Warning MsgType = fgYellow
	Error   MsgType = fgRed
	Neutral MsgType = fgWhite
)

var MsgTypeToString = map[MsgType]string{Info: "Info", Success: "Success", Warning: "Warning", Error: "Error", Neutral: "Neutral"}

func Log(msg any, priority MsgType, underline ...bool) {
	msgType := Neutral
	if isOkColor(priority) {
		msgType = priority
	}

	MsgTypeText := []any{string(fgMagenta), fmt.Sprintf("[%s]", MsgTypeToString[msgType]), string(reset)}
	Message := []any{string(msgType), msg, string(reset)}
	if len(underline) == 1 && underline[0] {
		Message = append([]any{string(underscore)}, Message...)
	}

	LogMessage := append(MsgTypeText, Message...)
	log.Println(LogMessage...)
}

func isOkColor(priority MsgType) bool {
	switch priority {
	case Info, Success, Warning, Error, Neutral:
		return true
	default:
		return false
	}
}
