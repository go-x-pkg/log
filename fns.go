package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cihub/seelog"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	_multipleErrMsg = "Multiple errors without a key."
)

var (
	nope      = func(string, ...interface{}) {}
	nopeLevel = func(string, ...interface{}) {}

	toWriter = func(w io.Writer, msg string, a ...interface{}) {
		var b strings.Builder
		b.WriteString(fmt.Sprintf(msg, a...))
		b.WriteString("\n")
		w.Write([]byte(b.String()))
	}

	toDiscard = func(msg string, a ...interface{}) { toWriter(io.Discard, msg, a...) }
	toStderr  = func(msg string, a ...interface{}) { toWriter(os.Stderr, msg, a...) }
	toStdout  = func(msg string, a ...interface{}) { toWriter(os.Stdout, msg, a...) }
)

type invalidPair struct {
	position   int
	key, value interface{}
}

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
	case *zap.SugaredLogger:
		return func(msg string, a ...interface{}) {
			isFirstZapField := false
			if len(a) > 0 {
				_, isFirstZapField = a[0].(zap.Field)
			}
			if isFirstZapField {
				switch l {
				case Quiet:
				case Trace, Debug:
					lgr.Debugw(msg, a...)
				case Info:
					lgr.Infow(msg, a...)
				case Warn:
					lgr.Warnw(msg, a...)
				case Error:
					lgr.Errorw(msg, a...)
				case Critical:
					lgr.Fatalw(msg, a...)
				}
			} else {
				switch l {
				case Quiet:
				case Trace, Debug:
					lgr.Debugf(msg, a...)
				case Info:
					lgr.Infof(msg, a...)
				case Warn:
					lgr.Warnf(msg, a...)
				case Error:
					lgr.Errorf(msg, a...)
				case Critical:
					lgr.Fatalf(msg, a...)
				}
			}
		}
	case *zap.Logger:
		return func(msg string, a ...interface{}) {
			var param []zap.Field
			if len(a) != 0 {
				param = make([]zap.Field, 0, len(a))
				var (
					invalid   invalidPairs
					seenError bool
				)
				for i := 0; i < len(a); {
					// This is a strongly-typed field. Consume it and move on.
					if f, ok := a[i].(zap.Field); ok {
						param = append(param, f)
						i++
						continue
					}

					// If it is an error, consume it and move on.
					if err, ok := a[i].(error); ok {
						if !seenError {
							seenError = true
							param = append(param, zap.Error(err))
						} else {
							lgr.With().Error(_multipleErrMsg, zap.Error(err))
						}
						i++
						continue
					}

					// Make sure this element isn't a dangling key.
					if i == len(a)-1 {
						param = append(param, zap.Any("ignored", a[i]))
						break
					}

					// Consume this value and the next, treating them as a key-value pair. If the
					// key isn't a string, add this pair to the slice of invalid pairs.
					key, val := a[i], a[i+1]
					if keyStr, ok := key.(string); !ok {
						// Subsequent errors are likely, so allocate once up front.
						if cap(invalid) == 0 {
							invalid = make(invalidPairs, 0, len(a)/2)
						}
						invalid = append(invalid, invalidPair{i, key, val})
					} else {
						param = append(param, zap.Any(keyStr, val))
					}
					i += 2
				}

				// If we encountered any invalid key-value pairs, log an error.
				if len(invalid) > 0 {
					param = append(param, zap.Array("invalid", invalid))
				}
			}
			switch l {
			case Quiet:
			case Trace, Debug:
				lgr.Debug(msg, param...)
			case Info:
				lgr.Info(msg, param...)
			case Warn:
				lgr.Warn(msg, param...)
			case Error:
				lgr.Error(msg, param...)
			case Critical:
				lgr.Fatal(msg, param...)
			}

		}
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

type invalidPairs []invalidPair

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	for i := range ps {
		err = multierr.Append(err, enc.AppendObject(ps[i]))
	}
	return err
}

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("position", int64(p.position))
	zap.Any("key", p.key).AddTo(enc)
	zap.Any("value", p.value).AddTo(enc)
	return nil
}
