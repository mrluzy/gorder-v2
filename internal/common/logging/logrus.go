package logging

import (
	"context"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mrluzy/gorder-v2/common/tracing"
	"github.com/rifflock/lfshook"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
	setOutput(logrus.StandardLogger())
	logrus.AddHook(&traceHook{})
}

func setOutput(logger *logrus.Logger) {
	var (
		folder    = "./log/"
		filePath  = "app.log"
		errorPath = "errors.log"
	)

	if err := os.MkdirAll(folder, 0750); err != nil && !os.IsExist(err) {
		panic(err)
	}
	file, err := os.OpenFile(folder+filePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	_, err = os.OpenFile(folder+errorPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	logrus.SetOutput(file)

	rotateInfo, err := rotatelogs.New(
		folder+filePath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName("app.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour*1),
	)
	if err != nil {
		panic(err)
	}

	rotateError, err := rotatelogs.New(
		folder+errorPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName("errors.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour*1),
	)
	if err != nil {
		panic(err)
	}

	rotationMap := lfshook.WriterMap{
		logrus.InfoLevel:  rotateInfo,
		logrus.DebugLevel: rotateInfo,
		logrus.WarnLevel:  rotateError,
		logrus.PanicLevel: rotateError,
		logrus.ErrorLevel: rotateError,
		logrus.FatalLevel: rotateError,
	}
	logrus.AddHook(lfshook.NewHook(rotationMap, &logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	}))
}

type traceHook struct {
}

func (t traceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (t traceHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		entry.Data["trace"] = tracing.TraceID(entry.Context)
		entry = entry.WithTime(time.Now())
	}
	return nil
}

func infoWithCost(ctx context.Context, start time.Time, fields logrus.Fields, format string, args ...any) {
	fields[Cost] = time.Since(start).Milliseconds()
	Infof(ctx, fields, format, args...)
}

func logf(ctx context.Context, level logrus.Level, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Logf(level, format, args...)
}

func Infof(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Infof(format, args...)
}

func Warnf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Warnf(format, args...)
}

func Panicf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Panicf(format, args...)
}

func Errorf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Errorf(format, args...)
}

func SetFormatter(logger *logrus.Logger) {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyMsg:   "message",
		},
	})
	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
		//logger.SetFormatter(&prefixed.TextFormatter{
		//	ForceColors:     true,
		//	ForceFormatting: true,
		//	TimestampFormat: time.RFC3339,
		//})
	}
}
