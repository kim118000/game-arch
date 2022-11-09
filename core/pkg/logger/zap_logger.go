package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Level = zapcore.Level
type Field = zapcore.Field

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

type Option func(logger *ZapLogger)

type ZapLogger struct {
	level      Level
	field      []Field
	showGid    bool
	outFormat  string
	fileName   string
	maxSize    int
	maxAge     int
	maxBackups int
	compress   bool
	flushLogs  func() error
	logger     Logger
}

func WithOutFormat(outformat string) Option {
	return func(logger *ZapLogger) {
		logger.outFormat = outformat
	}
}

func WithFileName(filename string) Option {
	return func(logger *ZapLogger) {
		logger.fileName = filename
	}
}

func WithMaxSize(maxsize int) Option {
	return func(logger *ZapLogger) {
		logger.maxSize = maxsize
	}
}

func WithMaxAge(maxage int) Option {
	return func(logger *ZapLogger) {
		logger.maxAge = maxage
	}
}

func WithCompress(compress bool) Option {
	return func(logger *ZapLogger) {
		logger.compress = compress
	}
}

func WithMaxBackups(maxbackups int) Option {
	return func(logger *ZapLogger) {
		logger.maxBackups = maxbackups
	}
}

func WithField(fields ...Field) Option {
	return func(logger *ZapLogger) {
		for _, f := range fields {
			logger.field = append(logger.field, f)
		}
	}
}

func WithShowGid(showGid bool) Option {
	return func(logger *ZapLogger) {
		logger.showGid = showGid
	}
}

func NewLogger(level int, options ...Option) *ZapLogger {
	zapLog := &ZapLogger{
		level: Level(level),
	}

	if len(options) > 0 {
		for _, op := range options{
			op(zapLog)
		}
	}

	syncer := make([]zapcore.WriteSyncer, 0, 2)
	std := zapcore.AddSync(os.Stdout)
	//zapcore.Lock(std)
	syncer = append(syncer, std)

	if len(zapLog.fileName) > 0 {
		// lumberjack.Logger is already safe for concurrent use, so we don't need to lock it.
		lumberJackLogger := &lumberjack.Logger{
			Filename:   zapLog.fileName,
			MaxSize:    zapLog.maxSize,   // 每个日志文件保存的大小 单位:M
			MaxAge:     zapLog.maxAge,     // 文件最多保存多少天
			MaxBackups: zapLog.maxBackups,    // 日志文件最多保存多少个备份
			Compress:   zapLog.compress, // 是否压缩
		}
		ws := zapcore.AddSync(lumberJackLogger)
		//zapcore.Lock(ws)
		syncer = append(syncer, ws)
	}

	encoder := zapLog.getEncoder()

	levelEnabler := zap.LevelEnablerFunc(func(level Level) bool {
		return level >= zapLog.level
	})

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(syncer...), levelEnabler)
	if len(zapLog.field) > 0 {
		core = core.With(zapLog.field)
	}
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	zapLog.logger = zapLogger.Sugar()
	zapLog.flushLogs = zapLogger.Sync

	return zapLog
}

func (z *ZapLogger) getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if z.outFormat == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}


// Debugf logs messages at DEBUG level.
func (z *ZapLogger) Debugf(format string, args ...interface{}) {
	z.logger.Debugf(format, args...)
}

// Infof logs messages at INFO level.
func (z *ZapLogger) Infof(format string, args ...interface{}) {
	z.logger.Infof(format, args...)
}

// Warnf logs messages at WARN level.
func (z *ZapLogger) Warnf(format string, args ...interface{}) {
	z.logger.Warnf(format, args...)
}

// Errorf logs messages at ERROR level.
func (z *ZapLogger) Errorf(format string, args ...interface{}) {
	z.logger.Errorf(format, args...)
}

// Fatalf logs messages at FATAL level.
func (z *ZapLogger) Fatalf(format string, args ...interface{}) {
	z.logger.Fatalf(format, args...)
}

// Cleanup does something windup for logger, like closing, flushing, etc.
func (z *ZapLogger) Cleanup() {
	if z.flushLogs != nil {
		_ = z.flushLogs()
	}
}
