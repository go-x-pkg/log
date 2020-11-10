package log

import "strings"

type Level uint8

const (
	/// A special log level used to turn off logging.
	Quiet Level = iota

	// For pervasive information on states of all elementary constructs.
	// Use 'Trace' for in-depth debugging to find problem parts of a function,
	// to check values of temporary variables, etc.
	Trace

	// For detailed system behavior reports and diagnostic messages
	// to help to locate problems during development.
	Debug

	// For general information on the application's work.
	// Use 'Info' level in your code so that you could leave it
	// 'enabled' even in production. So it is a 'production log level'.
	Info

	// For indicating small errors, strange situations,
	// failures that are automatically handled in a safe manner.
	Warn

	// For severe failures that affects application's workflow,
	// not fatal, however (without forcing app shutdown).
	Error

	// For producing final messages before application’s death.
	Critical
)

var levelText = map[Level]string{
	Quiet:    "quiet",
	Trace:    "trace",
	Debug:    "debug",
	Info:     "info",
	Warn:     "warn",
	Error:    "error",
	Critical: "critical",
}

func (l Level) String() string { return levelText[l] }

func (l *Level) UnmarshalYAML(unmarshal func(interface{}) error) error {
	v := ""
	unmarshal(&v)
	*l = NewLevel(v)
	return nil
}
func (l Level) MarshalYAML() (interface{}, error) { return l.String(), nil }

func NewLevel(v string) Level {
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "-", -1)
	v = strings.Replace(v, "_", "-", -1)
	v = strings.Replace(v, "(", "", -1)
	v = strings.Replace(v, ")", "", -1)

	switch v {
	case "q", "quiet", "off":
		return Quiet
	case "t", "trace":
		return Trace
	case "d", "debug":
		return Debug
	case "i", "info":
		return Info
	case "w", "warn", "warning":
		return Warn
	case "e", "err", "error":
		return Error
	case "c", "crit", "critical":
		return Critical
	}

	return Quiet
}
