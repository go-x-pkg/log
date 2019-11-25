package log

// logger storage
type Loggers interface {
	ByName(string) interface{}
}

func Close(loggers Loggers) {
	if loggers == nil {
		return
	}

	if closer, ok := loggers.(interface{ Close() }); ok {
		closer.Close()
	}
}
