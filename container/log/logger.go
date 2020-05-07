package log

import (
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

var BaseLogger *Logger

type Logger struct {
	*zap.Logger
	lumberjackLogger *lumberjack.Logger
}

func NewLogger(maxSize, maxAge int, logDir, name string, console, debug bool, level ...zapcore.Level) *Logger {
	var l zapcore.Level
	switch len(level) {
	case 0:
		if debug {
			l = zap.DebugLevel
		}
	case 1:
		l = level[0]
	default:
		panic("level参数最多只有1个")
	}
	fmt.Printf("生成日志器 最大大小,单位为MB=%d,最长时间,单位为天=%d,日志目录=%s,日志名=%s,是否为调试模式=%v,是否在终端输出=%v,日志级别=%s\n", maxSize, maxAge, logDir, name, debug, console, l.String())
	fileName := initLogFileName(name)
	if debug {
		fileName = name + ".log"
	}
	var writeSync []zapcore.WriteSyncer

	fileWrite := zapcore.AddSync(&lumberjack.Logger{
		Filename:  filepath.Join(logDir, fileName),
		MaxSize:   maxSize,
		MaxAge:    maxAge,
		LocalTime: true,
		Compress:  true,
	})
	writeSync = append(writeSync, fileWrite)
	if console {
		stdWrite, closeFunc, err := zap.Open("stderr")
		if err != nil {
			panic(err)
		}
		defer closeFunc()
		writeSync = append(writeSync, stdWrite)
	}
	syncer := zapcore.NewMultiWriteSyncer(writeSync...)
	if debug {
		config := zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		}
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(config),
			syncer,
			l,
		)
		logger := zap.New(core)
		logger = logger.WithOptions(zap.AddCaller())
		return &Logger{Logger: logger}
	} else {
		config := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			syncer,
			l,
		)
		// todo : 添加性能指标方法
		zapcore.RegisterHooks(core)
		logger := zap.New(core)
		logger = logger.WithOptions(zap.AddCaller())
		return &Logger{Logger: logger}
	}

}

func SetLogger(newLogger *Logger) {
	BaseLogger = newLogger
}
func NewTestLogger(component string, level zapcore.Level) *Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	config.Level = zap.NewAtomicLevelAt(level)
	logger, _ := config.Build(zap.AddCallerSkip(0))
	if component != "" {
		logger = logger.With(zap.String("组件", component))
	}
	return &Logger{Logger: logger}
}
func (l Logger) Close() {
	if l.lumberjackLogger != nil {
		_ = l.lumberjackLogger.Close()
	}
}
func initLogFileName(name string) string {
	return fmt.Sprintf("%s_%s.log", name, time.Now().Format("20060102_150405"))
}
func (l Logger) With(fields ...zap.Field) *Logger {
	return &Logger{lumberjackLogger: l.lumberjackLogger, Logger: l.Logger.With(fields...)}
}
