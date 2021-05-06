package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	Level          Levels
	UseColors      bool
	CallStackAdder int
	customFilename string
	out            []io.Writer
	sync.Mutex
}

func NewLogger(out io.Writer, level Levels) *Logger {
	return &Logger{out: []io.Writer{out}, Level: level, UseColors: true}
}

var Log = NewLogger(os.Stdout, LevelDebug)

func itoa(buf *[]byte, i int, width int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || width > 1 {
		width--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Logger) writeToOut(level Levels, message string) {

	now := time.Now()
	logType := Types[level]
	_, file, line, ok := runtime.Caller(3 + l.CallStackAdder)
	buf := make([]byte, 0, 1024)
	if l.UseColors {
		buf = append(buf, "\x1b[1;"...)
		itoa(&buf, logType.Color, 2)
		buf = append(buf, 'm')
	}
	buf = append(buf, "    "...)
	buf = append(buf, logType.Prefix...)
	buf = append(buf, ' ')

	year, month, day := now.Date()
	itoa(&buf, year, 4)
	buf = append(buf, '-')
	itoa(&buf, int(month), 2)
	buf = append(buf, '-')
	itoa(&buf, day, 2)
	buf = append(buf, ' ')

	hour, min, sec := now.Clock()
	itoa(&buf, hour, 2)
	buf = append(buf, ':')
	itoa(&buf, min, 2)
	buf = append(buf, ':')
	itoa(&buf, sec, 2)
	buf = append(buf, '.')
	itoa(&buf, now.Nanosecond()/1e3, 6)

	if l.customFilename != "" {
		buf = append(buf, " ["...)
		buf = append(buf, l.customFilename...)
		l.customFilename = ""
		buf = append(buf, "]"...)
	} else if ok {
		buf = append(buf, " ["...)
		buf = append(buf, filepath.Base(file)...)
		buf = append(buf, ':')
		itoa(&buf, line, -1)
		buf = append(buf, "]"...)
	}

	if l.UseColors {
		buf = append(buf, "\x1b[0m"...)
	}
	buf = append(buf, ' ')
	buf = append(buf, message...)
	buf = append(buf, '\n')

	l.Lock()
	for idx, writer := range l.out {
		if writer != nil {
			if _, err := writer.Write(buf); err != nil {
				if err != os.ErrClosed {
					fmt.Printf("log write error: %v\n", err)
				}
				l.out[idx] = nil
			}
		}
	}
	l.Unlock()
}

func (l *Logger) log(level Levels, f interface{}, s ...interface{}) (err error) {
	if level >= l.Level && f != nil {
		if f == nil {
			return
		}
		if false {
			_ = fmt.Sprintf(f.(string), s)
		}
		if first, ok := f.(string); ok && len(s) > 0 && strings.Contains(first, "%") {
			l.writeToOut(level, fmt.Sprintf(first, s...))
		} else {
			l.writeToOut(level, fmt.Sprint(f, fmt.Sprint(s...)))
		}
	}
	return
}

func (l *Logger) Print(f interface{}, s ...interface{}) {
	l.log(l.Level, f, s...)
}

func (l *Logger) Verbose(f interface{}, s ...interface{}) {
	l.log(LevelVerbose, f, s...)
}

func (l *Logger) Debug(f interface{}, s ...interface{}) {
	l.log(LevelDebug, f, s...)
}

func (l *Logger) Info(f interface{}, s ...interface{}) {
	l.log(LevelInfo, f, s...)
}

func (l *Logger) Warn(f interface{}, s ...interface{}) {
	l.log(LevelWarn, f, s...)
}

func (l *Logger) Error(f interface{}, s ...interface{}) {
	l.log(LevelError, f, s...)
}

func (l *Logger) Fatal(f interface{}, s ...interface{}) {
	l.log(LevelFatal, f, s...)
	os.Exit(-1)
}

func (l *Logger) FatalGo(f interface{}, s ...interface{}) {
	custom := ""
	_, file, line, ok := runtime.Caller(l.CallStackAdder + 1)
	if ok {
		custom = fmt.Sprintf("%v:%v", filepath.Base(file), line)
	}
	go func() {
		l.customFilename = custom
		l.Fatal(f, s...)
	}()
}

func (l *Logger) AddWriter(writer io.Writer) {
	l.Lock()
	l.out = append(l.out, writer)
	l.Unlock()
}

func (l *Logger) AddFileWriter(fileName string) (err error) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = Errorf(err)
	} else {
		l.AddWriter(f)
	}
	return
}

func GetCurrentFileAndLine(opts ...int) string {
	depth := 1
	if len(opts) > 0 {
		depth = opts[0]
	}
	_, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("%v:%v", filepath.Base(file), line)
}

func Errorf(f interface{}, s ...interface{}) (err error) {
	if f == nil {
		return
	}
	if false {
		_ = fmt.Sprintf(f.(string), s)
	}
	if first, ok := f.(string); ok && len(s) > 0 && strings.Contains(first, "%") {
		// str = fmt.Sprintf(first, s[1:]...)
		srcl := fmt.Sprintf("[%v] ", GetCurrentFileAndLine(2))
		err = fmt.Errorf(srcl+first, s...)
	} else if err, ok = f.(error); ok && len(s) == 0 {
		if err != nil {
			if strings.HasPrefix(err.Error(), "[") {
				err = fmt.Errorf("[%v]%w", GetCurrentFileAndLine(2), err)
			} else {
				err = fmt.Errorf("[%v]: %w", GetCurrentFileAndLine(2), err)
			}
		}
	} else {
		// str = fmt.Sprint(s...)
		err = fmt.Errorf("[%v]: %v %v", GetCurrentFileAndLine(2), f, fmt.Sprint(s...))
	}
	return
}
