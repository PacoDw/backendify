package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is the main object that can be used to log.
// Internally uses zap.Logger library.
type Logger struct {
	*zap.Logger

	config *Config
}

// NewLogger creates a new logger, internally uses Zap logger.
func NewLogger(config *Config) *Logger {
	// check the environment mode, by default is development mode
	var (
		levelMode = zap.DebugLevel
		pec       = zap.NewDevelopmentEncoderConfig()
	)

	if config.Environment == Production {
		pec = zap.NewProductionEncoderConfig()

		levelMode = zap.InfoLevel
	}

	// The encoder can be customized for each output
	pec.EncodeTime = zapcore.ISO8601TimeEncoder

	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.OutputFile,
		MaxSize:    config.MaxSize, // megabytes
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge, // days
	})

	// Tee different zap cores to writte into Console and a file
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(pec), zapcore.AddSync(w), zap.DebugLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(pec), zapcore.AddSync(os.Stdout), levelMode),
	)

	// Creates the logger including the stack trace for different error levels
	logger := zap.New(core,
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddStacktrace(zapcore.PanicLevel),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Warn(fmt.Sprintf("logger.Sync(): %s", err))
		}
	}()

	// Add all std logs into our logger
	zap.RedirectStdLog(logger)

	logger.Info(fmt.Sprintf("Initializing on %s Mode", config.Environment))

	return &Logger{
		Logger: logger,
		config: config,
	}
}
