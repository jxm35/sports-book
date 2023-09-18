package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func InitialiseDevLogger(fields ...zap.Field) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.FunctionKey = "func"
	config.EncoderConfig.StacktraceKey = "stack"
	config.EncoderConfig.CallerKey = "caller"
	log, _ := config.Build()
	logger = log.Sugar()
	// consoleEncoder := zapcore.NewJSONEncoder(config)
	// core := zapcore.NewTee(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel))
	// logger = zap.New(core).With(fields...).Sugar()
}

func InitialiseStructuredLogger(fields ...zap.Field) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.FunctionKey = "func"
	config.EncoderConfig.StacktraceKey = "stack"
	config.EncoderConfig.CallerKey = "caller"
	log, _ := config.Build()
	logger = log.Sugar()
}

func Info(message string, args ...any) {
	if len(args)%2 != 0 {
		logger.Error(
			"invalid args for log",
			"num_args", len(args),
			"message", message,
		)
	}
	logger.Info(message, args)
}

func Debug(message string, args ...any) {
	if len(args)%2 != 0 {
		logger.Error(
			"invalid args for log",
			"num_args", len(args),
			"message", message,
		)
	}
	logger.Debug(message, args)
}

func Error(message string, args ...any) {
	if len(args)%2 != 0 {
		logger.Error(
			"invalid args for log",
			"num_args", len(args),
			"message", message,
		)
	}
	logger.Error(message, args)
}
