package shack

type modeCode int8

const (
	DebugMode modeCode = iota
	TestMode
	ReleaseMode
)

var mode modeCode


func Mode(m modeCode) {
	mode = m
}


func IsDebugging() bool {
	return mode == DebugMode
}


func IsTesting() bool {
	return mode == TestMode
}