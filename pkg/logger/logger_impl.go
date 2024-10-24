package logger

import (
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	once     sync.Once
	instance *LogrusLogger
)

type LogrusLogger struct {
	log *logrus.Logger
}

func initLogger() {
	log := logrus.New()
	log.SetReportCaller(false) // Disable built-in caller reporting
	log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:       false,
		DisableHTMLEscape: true,
	})

	// Set log level based on the environment variable
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		viper.AutomaticEnv()
	}

	logLevel := viper.GetString("LOG_LEVEL")
	switch strings.ToLower(logLevel) {
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.ErrorLevel) // Default to WarnLevel if no valid env variable is set
	}

	instance = &LogrusLogger{log: log}

}

func NewLogger() *LogrusLogger {
	once.Do(initLogger)
	return instance
}

func GetLogger() *LogrusLogger {
	if instance == nil {
		NewLogger()
	}
	return instance
}

// logWithCallerAdjustment adjusts the caller depth to the actual caller in the service code
func (l *LogrusLogger) logWithCallerAdjustment(level logrus.Level, skip int, msg string, args ...interface{}) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		l.log.Error(COULD_NOT_RETRIEVE_CALLER_INFO)
		return
	}

	function := runtime.FuncForPC(pc).Name()
	entry := l.log.WithFields(logrus.Fields{
		FUNCTION: function,
		FILE:     file,
		LINE:     line,
	})

	if len(args) > 0 && msg == "" {
		msg = "%v"
	}

	switch level {
	case logrus.InfoLevel:
		if msg != EMPTY_STRING {
			entry.Infof(msg, args...)
		} else {
			entry.Info(args...)
		}
	case logrus.ErrorLevel:
		if msg != EMPTY_STRING {
			entry.Errorf(msg, args...)
		} else {
			entry.Error(args...)
		}
	case logrus.FatalLevel:
		if msg != EMPTY_STRING {
			entry.Fatalf(msg, args...)
		} else {
			entry.Fatal(args...)
		}
	case logrus.DebugLevel:
		if msg != EMPTY_STRING {
			entry.Debugf(msg, args...)
		} else {
			entry.Debug(args...)
		}
	}
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.logWithCallerAdjustment(logrus.InfoLevel, 2, EMPTY_STRING, args...)
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.logWithCallerAdjustment(logrus.InfoLevel, 2, format, args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.logWithCallerAdjustment(logrus.ErrorLevel, 2, EMPTY_STRING, args...)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logWithCallerAdjustment(logrus.ErrorLevel, 2, format, args...)
}

func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.logWithCallerAdjustment(logrus.FatalLevel, 2, EMPTY_STRING, args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logWithCallerAdjustment(logrus.FatalLevel, 2, format, args...)
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logWithCallerAdjustment(logrus.DebugLevel, 2, EMPTY_STRING, args...)
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logWithCallerAdjustment(logrus.DebugLevel, 2, format, args...)
}
