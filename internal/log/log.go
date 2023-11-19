package log

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var rawJSON = []byte(`{
	  "level": "info",
	  "encoding": "json",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stdout"],
	  "encoderConfig": {
	    "messageKey": "msg",
	    "levelKey": "level",
		"timeKey": "timestamp"
	  }
	}`)

type config struct {
	zap zap.Config
	set bool
	sync.Mutex
}

var cfg config

type key struct{}

func setup() {
	env := os.Getenv("ENV")
	version := os.Getenv("VERSION")
	app := os.Getenv("APP")

	cfg.Lock()
	defer cfg.Unlock()

	if err := json.Unmarshal(rawJSON, &cfg.zap); err != nil {
		panic(err)
	}

	cfg.zap.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	cfg.zap.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg.zap.InitialFields = map[string]interface{}{}

	cfg.zap.InitialFields["app"] = app
	cfg.zap.InitialFields["env"] = env
	cfg.zap.InitialFields["version"] = version
}

func With(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(key{}).(*zap.SugaredLogger); ok {
		return l
	}

	if !cfg.set {
		setup()

		cfg.set = true
	}

	l, err := cfg.zap.Build()
	if err != nil {
		panic(err)
	}

	sugar := l.Sugar()

	return sugar
}

func Add(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, key{}, l)
}

func AddFields(ctx context.Context, fields ...interface{}) context.Context {
	if len(fields) == 0 {
		return ctx
	}

	return Add(ctx, With(ctx).With(fields...))
}

func Error(ctx context.Context, msg string, fields ...interface{}) {
	With(ctx).Errorw(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...interface{}) {
	With(ctx).Infow(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...interface{}) {
	With(ctx).Debugw(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...interface{}) {
	With(ctx).Warnw(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...interface{}) {
	With(ctx).Fatalw(msg, fields...)
}

func Err(err error) zap.Field {
	return zap.Error(err)
}

func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

func Time(key string, val time.Time) zap.Field {
	return zap.Time(key, val)
}
