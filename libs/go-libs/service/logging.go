package service

import (
	"io"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

func GetDefaultLogger(w io.Writer) logging.Logger {
	l := logrus.New()
	l.SetOutput(w)
	if IsDebug() {
		l.Level = logrus.DebugLevel
	}

	var formatter logrus.Formatter
	if viper.GetBool(JsonFormattingLoggerFlag) {
		jsonFormatter := &logrus.JSONFormatter{}
		jsonFormatter.TimestampFormat = "15-01-2018 15:04:05.000000"
		formatter = jsonFormatter
	} else {
		textFormatter := new(logrus.TextFormatter)
		textFormatter.TimestampFormat = "15-01-2018 15:04:05.000000"
		textFormatter.FullTimestamp = true
		formatter = textFormatter
	}

	l.SetFormatter(formatter)

	if viper.GetBool(otlptraces.OtelTracesFlag) {
		l.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		)))
	}
	return logging.NewLogrus(l)
}
