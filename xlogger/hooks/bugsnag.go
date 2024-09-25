package hooks

import (
	"errors"
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/sirupsen/logrus"
)

type BugsnagHook struct {
	levels []logrus.Level
}

func NewBugsnagHook(levels []logrus.Level, cfg ...bugsnag.Configuration) *BugsnagHook {
	if len(cfg) > 0 {
		bugsnag.Configure(cfg[0])
	}

	if levels == nil {
		levels = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel}
	}

	return &BugsnagHook{
		levels: levels,
	}
}

func logrusLevelToBugsnagSeverity(level logrus.Level) any {
	switch level {
	case logrus.PanicLevel:
		fallthrough
	case logrus.FatalLevel:
		fallthrough
	case logrus.ErrorLevel:
		return bugsnag.SeverityError
	case logrus.WarnLevel:
		return bugsnag.SeverityWarning
	default:
		return bugsnag.SeverityInfo
	}
}

func (h *BugsnagHook) Levels() []logrus.Level {
	return h.levels
}

func (h *BugsnagHook) Fire(entry *logrus.Entry) error {
	var data []any

	if entry.Context != nil {
		data = append(data, entry.Context)
	}

	data = append(data, logrusLevelToBugsnagSeverity(entry.Level), bugsnag.MetaData{"data": entry.Data})

	return bugsnag.Notify(errors.New(entry.Message), data...)
}
