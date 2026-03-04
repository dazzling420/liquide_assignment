package logger

import (
	"liquide_assignment/internal/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Service interface {
	GetLogger() *zap.SugaredLogger

	Errorf(format string, args ...interface{})
	// Error(args ...interface{})
	Error(args ...interface{})

	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})

	Infof(format string, args ...interface{})
	Info(args ...interface{})

	Warnf(format string, args ...interface{})
	Warn(args ...interface{})

	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

type standardLogger struct {
	logger *zap.SugaredLogger
}

func Init(loggerConfig config.Logger) *standardLogger {
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zap.DebugLevel)

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		LineEnding:   "\n",
	}

	// File writer with rotation
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   loggerConfig.FileName,
		MaxSize:    loggerConfig.MaxSizeInMB,
		MaxBackups: loggerConfig.MaxBackups,
		MaxAge:     loggerConfig.MaxAgeInDays,
		Compress:   loggerConfig.Compress,
	})

	var core zapcore.Core

	if loggerConfig.ConsoleLoggingEnabled {
		consoleWriter := zapcore.AddSync(os.Stdout)
		multiWriter := zapcore.NewMultiWriteSyncer(fileWriter, consoleWriter)
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), multiWriter, atom)
	} else {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, atom)
	}

	logger := zap.New(core, zap.AddCaller())
	sugar := logger.Sugar()
	// AddCallerSkip(1) is important because we wrap the logger
	sugar = sugar.WithOptions(zap.AddCallerSkip(1))

	return &standardLogger{sugar}
}

func (s *standardLogger) GetLogger() *zap.SugaredLogger {
	return s.logger
}

func (s *standardLogger) Errorf(format string, args ...interface{}) {
	s.logger.Errorf(format, args...)
}

func (s *standardLogger) Error(args ...interface{}) {
	reponseMessage := "unknown"
	if err, ok := args[len(args)-1].(error); ok {
		errString := err.Error()
		if res, ok := err.(config.Errors); ok {
			errString = res.ErrorMessage()
			if errString == "" {
				errString = err.Error()
			}
		}
		args = append(args[:len(args)-1], " ", errString)
		reponseMessage = err.Error()
	}
	s.logger.WithOptions(zap.Fields(zap.Field{
		Key:    "response_message",
		Type:   15,
		String: reponseMessage,
	})).Error(args...)
}

func (s *standardLogger) Fatalf(format string, args ...interface{}) {
	s.logger.Fatalf(format, args...)
}

func (s *standardLogger) Fatal(args ...interface{}) {
	s.logger.Fatal(args...)
}

func (s *standardLogger) Infof(format string, args ...interface{}) {
	s.logger.Infof(format, args...)
}

func (s *standardLogger) Info(args ...interface{}) {
	s.logger.Info(args...)
}

func (s *standardLogger) Warn(args ...interface{}) {
	s.logger.Warn(args...)
}

func (s *standardLogger) Warnf(format string, args ...interface{}) {
	s.logger.Warnf(format, args...)
}

func (s *standardLogger) Debugf(format string, args ...interface{}) {
	s.logger.Debugf(format, args...)
}

func (s *standardLogger) Debug(args ...interface{}) {
	s.logger.Debug(args...)
}

func (s *standardLogger) Printf(format string, args ...interface{}) {
	s.logger.Infof(format, args...)
}

func (s *standardLogger) Println(args ...interface{}) {
	s.logger.Info(args...)
	s.logger.Info("\n")
}
