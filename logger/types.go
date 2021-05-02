package logger

type Interface interface {
	Verbose(f interface{}, s ...interface{})
	Debug(f interface{}, s ...interface{})
	Info(f interface{}, s ...interface{})
	Warn(f interface{}, s ...interface{})
	Error(f interface{}, s ...interface{})
	Fatal(f interface{}, s ...interface{})
	FatalGo(f interface{}, s ...interface{})
}

type Empty struct{}

func (*Empty) Verbose(f interface{}, s ...interface{}) {}
func (*Empty) Debug(f interface{}, s ...interface{})   {}
func (*Empty) Info(f interface{}, s ...interface{})    {}
func (*Empty) Warn(f interface{}, s ...interface{})    {}
func (*Empty) Error(f interface{}, s ...interface{})   {}
func (*Empty) Fatal(f interface{}, s ...interface{})   {}
func (*Empty) FatalGo(f interface{}, s ...interface{}) {}

type Levels int

const (
	// LevelDebug самый подробный лог
	LevelDebug Levels = iota
	// LevelVerbose подробный лог, но без дебагг-инфо
	LevelVerbose
	// LevelInfo только ошибки, предупреждения и информация
	LevelInfo
	// LevelWarn только ошибки и предупреждения
	LevelWarn
	// LevelError только ошибки
	LevelError
	// LevelFatal ошибка приводит к завершению приложения
	LevelFatal
)

const (
	_ int = iota + 90 // fgHiBlack
	fgHiRed
	fgHiGreen
	_ // fgHiYellow
	fgHiBlue
	fgHiMagenta
	fgHiCyan
	_ // fgHiWhite
)

type Type struct {
	Prefix string
	Color  int
}

var Types = map[Levels]Type{
	LevelFatal:   {"ftl", fgHiRed},
	LevelError:   {"err", fgHiRed},
	LevelWarn:    {"wrn", fgHiMagenta},
	LevelInfo:    {"inf", fgHiGreen},
	LevelVerbose: {"vrb", fgHiCyan},
	LevelDebug:   {"dbg", fgHiBlue},
}
