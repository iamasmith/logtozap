package logtozap

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestToSugared(t *testing.T) {
	assert := assert.New(t)
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core).Sugar().With("myfield", "hassomething")
	ToSugared(logger, zapcore.WarnLevel)
	log.Print("many happy returns from the logger")
	assert.Equal(1, logs.Len(), "must be equal")
	ent := logs.All()[0]
	assert.Contains(ent.Message, "many happy returns", "must contain")
	assert.ElementsMatch(
		ent.Context,
		[]zap.Field{
			{Key: "myfield", Type: zapcore.StringType, String: "hassomething"},
		},
	)
}

func TestToSugaredWithSkip(t *testing.T) {
	assert := assert.New(t)
	core, logs := observer.New(zap.InfoLevel)
	lp := log.New(log.Writer(), "", 0)
	ctxLog := ctxWrapper{l: lp}
	logger := zap.New(core).Sugar().With("myfield", "hassomething")
	ToSugaredWithSkip(logger, zapcore.WarnLevel, 1, lp)
	ctxLog.Print(context.Background(), "many happy returns from the logger")
	assert.Equal(1, logs.Len(), "must be equal")
	ent := logs.All()[0]
	assert.Contains(ent.Message, "many happy returns", "must contain")
	assert.ElementsMatch(
		ent.Context,
		[]zap.Field{
			{Key: "myfield", Type: zapcore.StringType, String: "hassomething"},
		},
	)
}

func TestToLogger(t *testing.T) {
	assert := assert.New(t)
	core, logs := observer.New(zap.InfoLevel)
	fields := []zapcore.Field{
		{Type: zapcore.StringType, Key: "myfield", String: "hassomething"},
	}
	logger := zap.New(core).With(fields...)
	ToLogger(logger, zapcore.WarnLevel)
	log.Print("many happy returns from the logger")
	assert.Equal(1, logs.Len(), "must be equal")
	ent := logs.All()[0]
	assert.Contains(ent.Message, "many happy returns", "must contain")
	assert.ElementsMatch(
		ent.Context,
		[]zap.Field{
			{Key: "myfield", Type: zapcore.StringType, String: "hassomething"},
		},
	)
}

func TestToLoggerWithSkip(t *testing.T) {
	assert := assert.New(t)
	core, logs := observer.New(zap.InfoLevel)
	lp := log.New(log.Writer(), "", 0)
	ctxLog := ctxWrapper{l: lp}
	fields := []zapcore.Field{
		{Type: zapcore.StringType, Key: "myfield", String: "hassomething"},
	}
	logger := zap.New(core).With(fields...)
	ToLoggerWithSkip(logger, zapcore.WarnLevel, 1, lp)
	ctxLog.Print(context.Background(), "many happy returns from the logger")
	assert.Equal(1, logs.Len(), "must be equal")
	ent := logs.All()[0]
	assert.Contains(ent.Message, "many happy returns", "must contain")
	assert.ElementsMatch(
		ent.Context,
		[]zap.Field{
			{Key: "myfield", Type: zapcore.StringType, String: "hassomething"},
		},
	)
}

type ctxWrapper struct {
	l *log.Logger
}

func (r ctxWrapper) Print(ctx context.Context, s string) {
	r.l.Print(s)
}

func TestToNew(t *testing.T) {
	assert := assert.New(t)
	log_ := log.New(log.Writer(), "CUSTOM:", 2)
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core).Sugar().With("myfield", "hassomething")
	ToSugared(logger, zapcore.WarnLevel, log_)
	log_.Print("Testing")
	assert.Equal(1, logs.Len(), "must be equal")
	ent := logs.All()[0]
	assert.Equal(ent.Message, "CUSTOM:Testing", "must equal")
	assert.ElementsMatch(
		ent.Context,
		[]zap.Field{
			{Key: "myfield", Type: zapcore.StringType, String: "hassomething"},
		},
	)
}
