package dockle

import "strings"

const (
	fatal int = iota
	warn
	info
	skip
	pass
)

var ErrorLevel = map[string]int{
	"fatal": fatal,
	"warn":  warn,
	"info":  info,
	"skip":  skip,
	"pass":  pass,
}

func ConvertErrorLevelToNumber(errorLevel string) int {
	number, _ := ErrorLevel[strings.ToLower(errorLevel)]
	return number
}
