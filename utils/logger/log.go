package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func init() {
	InitLogger()
}

func InitLogger() {
	once.Do(func() {
		var err error

		// Encoders
		consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Explicit color encoder
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

		fileEncoderConfig := zap.NewDevelopmentEncoderConfig()
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // No colors
		fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig)

		// Log file writer
		logPath := getLogFilePath()
		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		fileWriter := zapcore.AddSync(file)

		// Console writer
		consoleWriter := zapcore.AddSync(os.Stdout)

		// Combine cores
		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, fileWriter, zapcore.DebugLevel),
			zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel),
		)

		logger = zap.New(core, zap.AddCaller(),zap.AddCallerSkip(1))
	})
}

func getLogFilePath() string {
	currentDate := time.Now().Format("2006-01-02")
	return filepath.Join("logs", currentDate+".log")
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// Infof logs a formatted info message.
func Infof(format string, args ...interface{}) {
    msg := fmt.Sprintf(format, args...)
    logger.Info(msg)
}

// Debugf logs a formatted debug message.
func Debugf(format string, args ...interface{}) {
    msg := fmt.Sprintf(format, args...)
    logger.Debug(msg)
}

// Warnf logs a formatted warning message.
func Warnf(format string, args ...interface{}) {
    msg := fmt.Sprintf(format, args...)
    logger.Warn(msg)
}

// Errorf logs a formatted error message.
func Errorf(format string, args ...interface{}) {
    msg := fmt.Sprintf(format, args...)
    logger.Error(msg)
}

// Fatalf logs a formatted fatal message and exits the application.
func Fatalf(format string, args ...interface{}) {
    msg := fmt.Sprintf(format, args...)
    logger.Fatal(msg)
}

func Sync() {
	logger.Sync()
}

// AddWithLog adds delta to the WaitGroup and logs the action along with caller info.
func AddWithLog(wg *sync.WaitGroup, delta int, description string) {
    wg.Add(delta)
    file, line := getCallerInfo(1)
    Infof("WaitGroup: Add(%d) called for %s at %s:%d", delta, description, file, line)
}

// DoneWithLog marks the WaitGroup as done and logs the action along with caller info.
func DoneWithLog(wg *sync.WaitGroup, description string) {
    wg.Done()
    file, line := getCallerInfo(1)
    Infof("WaitGroup: Done() called for %s at %s:%d", description, file, line)
}

// getCallerInfo retrieves the caller's file name and line number.
func getCallerInfo(skip int) (string, int) {
    // skip=0: getCallerInfo
    // skip=1: caller of getCallerInfo (AddWithLog or DoneWithLog)
    // skip=2: caller of AddWithLog or DoneWithLog
    pc, file, line, ok := runtime.Caller(skip + 1)
    if !ok {
        return "unknown", 0
    }
    fn := runtime.FuncForPC(pc)
    if fn == nil {
        return file, line
    }
    return fmt.Sprintf("%s", fn.Name()), line
}