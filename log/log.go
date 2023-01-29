package log

import (
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/mhd7966/darvazeh/config"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func LogInit() {
	Log = logrus.New()

	Log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-Jan-02 15:04:05",
	})

	if config.Cfg.Log.OutputType == "stdout" {
		Log.SetOutput(os.Stdout)

	} else if config.Cfg.Log.OutputType == "file" {
		file, err := os.OpenFile(config.Cfg.Log.OutputAdd, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			Log.Fatal(err)
		}
		Log.SetOutput(file)
	} else {
		Log.SetOutput(os.Stdout)

	}

	logLevel, err := logrus.ParseLevel(config.Cfg.Log.LogLevel)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	Log.SetLevel(logLevel)

	client, err := raven.New(config.Cfg.Sentry.DSN)
	if err != nil {
		Log.Fatal(err)
	}

	hook, err := logrus_sentry.NewWithClientSentryHook(client, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})

	if err == nil {
		Log.Hooks.Add(hook)
	}

}
