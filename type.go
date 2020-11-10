package log

type FnT func(Level, string, ...interface{})
type FnTShrt func(Level, string)

type Logger interface {
	Log(Level, string, ...interface{})
}
