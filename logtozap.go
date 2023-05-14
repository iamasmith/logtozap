package logtozap

import (
	"log"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logtozap struct {
	depth int
	level zapcore.Level
	zap   zapcore.Core
}

// Capture any specified loggers to the Zap Sugared logger.
// If non provided will capture the default logger.
// level specifies the log severity level for the Zap logger messages.
func ToSugared(z *zap.SugaredLogger, level zapcore.Level, lp ...*log.Logger) {
	route(z.Desugar().Core(), level, 0, lp)
}

func ToSugaredWithSkip(z *zap.SugaredLogger, level zapcore.Level, addSkip int, lp ...*log.Logger) {
	route(z.Desugar().Core(), level, addSkip, lp)
}

// Capture any specified loggers to the Zap Unsugared logger.
func ToLogger(z *zap.Logger, level zapcore.Level, lp ...*log.Logger) {
	route(z.Core(), level, 0, lp)
}

// Capture any specified loggers to the Zap Unsugared logger.
func ToLoggerWithSkip(z *zap.Logger, level zapcore.Level, addSkip int, lp ...*log.Logger) {
	route(z.Core(), level, addSkip, lp)
}

func route(z zapcore.Core, level zapcore.Level, addSkip int, lp []*log.Logger) {
	var l logtozap
	l.level = level
	l.depth = tuner{}.unwind() + addSkip
	l.zap = z
	if len(lp) == 0 {
		log.SetFlags(0)
		log.SetOutput(l)
		return
	}
	for _, e := range lp {
		e.SetFlags(0)
		e.SetOutput(l)
	}
}

func (l logtozap) Write(bytes []byte) (int, error) {
	entry := zapcore.Entry{
		Level:   l.level,
		Time:    time.Now(),
		Message: strings.TrimSpace(string(bytes)),
		Caller:  zapcore.NewEntryCaller(runtime.Caller(l.depth)),
	}
	l.zap.Write(entry, []zapcore.Field{})
	return len(bytes), nil
}

// We know this will probably always be 3 but we test it anyway in case
// of any changes to the log package in future versions
type tuner struct {
	depth  int
	caller string
}

func (t tuner) unwind() int {
	pc, _, _, _ := runtime.Caller(0)
	t.caller = runtime.FuncForPC(pc).Name()
	log.SetOutput(&t)
	log.Print("")
	return t.depth
}

func (t *tuner) Write([]byte) (int, error) {
	for i := 2; ; i++ {
		pc, _, _, _ := runtime.Caller(i)
		if runtime.FuncForPC(pc).Name() == t.caller {
			t.depth = i
			return 0, nil
		}
	}
}
