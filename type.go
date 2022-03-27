package log

type FnT func(Level, string, ...interface{})

func (f FnT) Log(l Level, msg string, a ...interface{}) {
	f(l, msg, a...)
}

type FnTShrt func(Level, string)

type Logger interface {
	Log(Level, string, ...interface{})
}
