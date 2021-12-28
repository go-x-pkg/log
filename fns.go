package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cihub/seelog"
)

var (
	nope      = func(string, ...interface{}) {}
	nopeLevel = func(string, ...interface{}) {}

	toWriter = func(w io.Writer, msg string, a ...interface{}) {
		b := strings.Builder{}
		b.WriteString(fmt.Sprintf(msg, a...))
		b.WriteString("\n")
		w.Write([]byte(b.String()))
	}

	toDiscard = func(msg string, a ...interface{}) { toWriter(ioutil.Discard, msg, a...) }
	toStderr  = func(msg string, a ...interface{}) { toWriter(os.Stderr, msg, a...) }
	toStdout  = func(msg string, a ...interface{}) { toWriter(os.Stdout, msg, a...) }
)

func LogFn(any logger, l Level) func(string, ...interface{}) {
	switch lgr := any.(type) {
	case seelog.LoggerInterface:
		return func(msg string, a ...interface{}) {
			switch l {
			case Quiet:
			case Trace:
				lgr.Tracef(msg, a...)
			case Debug:
				lgr.Debugf(msg, a...)
			case Info:
				lgr.Infof(msg, a...)
			case Warn:
				lgr.Warnf(msg, a...)
			case Error:
				lgr.Errorf(msg, a...)
			case Critical:
				lgr.Criticalf(msg, a...)
			}

			lgr.Flush()
		}
		// case *zap.Logger:
		// 	return func(msg string) {
		// 		switch l {
		// 		case Quiet:
		// 		case Trace, Debug:
		// 			lgr.Debug(msg)
		// 		case Info:
		// 			lgr.Info(msg)
		// 		case Warn:
		// 			lgr.Warn(msg)
		// 		case Error:
		// 			lgr.Warn(msg) // TODO: research for non-panicable Error for zap
		// 		case Critical:
		// 			lgr.Fatal(msg)
		// 		}

		// 		lgr.Sync()
		// 	}
	}

	switch l {
	case Quiet:
		return toDiscard
	case Trace, Debug, Info:
		return toStdout
	case Warn, Error, Critical:
		return toStderr
	}

	return nope
}

func Log(any logger, l Level, msg string)       { LogFn(any, l)(msg) }
func LogToDiscard(msg string, a ...interface{}) { LogFn(nil, Quiet)(msg, a...) }

func LogStd(l Level, msg string, a ...interface{}) { LogFn(nil, l)(msg, a...) }

func Logf(any logger, l Level, msg string, a ...interface{}) { LogFn(any, l)(msg, a...) }

func LogfStd(l Level, msg string, a ...interface{}) { LogStd(l, msg, a...) }
