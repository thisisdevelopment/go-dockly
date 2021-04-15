package xlogger

import (
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Deaurorize struct {
	Level logrus.Level
}

func NewDeAurorizeHook(lvl logrus.Level) logrus.Hook {
	return &Deaurorize{
		Level: lvl,
	}
}

func (d *Deaurorize) Fire(entry *logrus.Entry) error {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	// remove trailing new line
	entry.Message = strings.TrimRight(entry.Message, "\n")
	// remove any ansi
	entry.Message = regexp.MustCompile(ansi).ReplaceAllString(entry.Message, "")
	return nil
}

func (d *Deaurorize) Levels() (levels []logrus.Level) {
	for _, level := range logrus.AllLevels {
		if level <= d.Level {
			levels = append(levels, level)
		}
	}
	return levels
}
