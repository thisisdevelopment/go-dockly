package xlogger

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/imdario/mergo"
	"github.com/logrusorgru/aurora"
	gelf "github.com/seatgeek/logrus-gelf-formatter"
	"github.com/stretchr/testify/suite"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Level       string
	Format      string
	HideFName   bool
	CallerDepth int
}

type Logger struct {
	log    *logrus.Logger
	errLog *logrus.Logger
	cfg    *Config
}

func New(inputCfg *Config) (*Logger, error) {
	cfg := &Config{
		Level:       "debug",
		Format:      "text",
		CallerDepth: 2,
	}

	err := merge(cfg, inputCfg)

	if err != nil {
		return nil, errors.Wrap(err, "failed to merge with default config.")
	}

	logLevel, err := logrus.ParseLevel(cfg.Level)

	if err != nil {
		return nil, errors.Wrap(err, "unsupported log level")
	}

	var formatter logrus.Formatter

	switch cfg.Format {
	case "gelf":
		formatter = &gelf.GelfTimestampFormatter{}
	case "text":
		formatter = &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     true,
		}
	case "json":
		fallthrough
	default:
		formatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.999Z07:00",
		}
	}

	l := &Logger{
		log: &logrus.Logger{
			Out:       os.Stdout,
			Hooks:     make(logrus.LevelHooks),
			Formatter: formatter,
			Level:     logLevel,
		},
		errLog: &logrus.Logger{
			Out:       os.Stderr,
			Hooks:     make(logrus.LevelHooks),
			Formatter: formatter,
			Level:     logLevel,
		},
		cfg: cfg,
	}

	return l, nil
}

func merge(dst, src interface{}) error {
	if src != nil && !reflect.ValueOf(src).IsNil() {
		return mergo.Merge(dst, src, mergo.WithOverride)
	}

	return nil
}

// DefaultTestLogger returns the log interface to the test suite
func DefaultTestLogger(s *suite.Suite) *Logger {
	cfg := new(Config)
	l, err := New(cfg)
	s.Require().NoError(err)
	return l
}

func (l *Logger) Log(log *logrus.Logger) *logrus.Entry {
	if pc, file, line, ok := runtime.Caller(l.cfg.CallerDepth); ok {
		fName := runtime.FuncForPC(pc).Name()

		currentDir, _ := os.Getwd()
		file = strings.Replace(file, currentDir+"/", "", -1)

		caller := fmt.Sprintf("%s:%v", file, line)

		if l.cfg.HideFName {
			return log.WithField("caller", caller)
		}
		return log.WithField("caller", caller).WithField("fName", fName)
	}
	return &logrus.Entry{}
}

// s.log.WithFields(s.log.toFields("domain", "test" "language", "nl")).Warn("test")
func (l *Logger) toFields(fields ...string) (f logrus.Fields) {

	f = make(logrus.Fields)

	if len(fields)%2 != 0 {
		l.Warnf("fields should always contain an even amount of elements")
	}

	for i, v := range fields {
		// in order of field:value, field:value
		if i%2 == 1 {
			continue
		}

		f[v] = fields[i+1]
	}

	return f
}

// WithField proxy method
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Log(l.log).WithField(key, value)
}

// WithFields proxy method
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Log(l.log).WithFields(fields)
}

// WithError proxy method
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Log(l.log).WithError(err)
}

var (
	// GitCommit holds short commit hash of source tree
	GitCommit string
	// GitBranch holds current branch name the code is built off
	GitBranch string
	// GitState shows whether there are uncommitted changes
	GitState string
	// BuildDate holds RFC3339 formatted UTC date (build time)
	BuildDate string
	// Version holds contents of ./VERSION file, if exists, or the value passed via the -version option
	Version string
)

func (l *Logger) BuildInfo(banner, cfgPath, version, commit, branch, state, date string) {

	log.Printf("loaded config from `%s`\n", aurora.Cyan(cfgPath))

	dir, _ := os.Getwd()

	log.Printf("running from `%s`\n", aurora.Cyan(dir))

	fmt.Println(aurora.Yellow(banner))

	fmt.Printf(`LogLevel: %s
Version: %s
Commit: %s
Branch: %s
Status: %s
BuildDate: %s
	
`, aurora.Cyan(l.log.Level), aurora.Yellow(version), aurora.Yellow(commit), branch, state, aurora.Yellow(date))
}

func (l *Logger) Level() string {
	return l.log.Level.String()
}

func (l *Logger) AddHook(hook logrus.Hook) {
	l.log.Hooks.Add(hook)
}

func (l *Logger) Dump(v ...interface{}) {
	spew.Dump(v...)
	l.Log(l.log).Info(aurora.Yellow("^^^dump^^^"))
}

func (l *Logger) Printf(format string, v ...interface{}) {
	if l.log.Level >= logrus.InfoLevel {
		l.Log(l.log).Printf(format, v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.log.Level >= logrus.InfoLevel {
		l.Log(l.log).Infof(format, v...)
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.log.Level >= logrus.FatalLevel {
		l.Log(l.errLog).Fatalf(format, v...)
	}
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	if l.log.Level >= logrus.PanicLevel {
		l.Log(l.errLog).Panicf(format, v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.log.Level >= logrus.DebugLevel {
		l.Log(l.log).Debugf(format, v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.log.Level >= logrus.WarnLevel {
		l.Log(l.log).Warnf(format, v...)
	}
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	if l.log.Level >= logrus.WarnLevel {
		l.Log(l.log).Warningf(format, v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.log.Level >= logrus.ErrorLevel {
		l.Log(l.errLog).Errorf(format, v...)
	}
}

func (l *Logger) Print(v ...interface{}) {
	if l.log.Level >= logrus.InfoLevel {
		l.Log(l.log).Print(v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.log.Level >= logrus.InfoLevel {
		l.Log(l.log).Info(v...)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.log.Level >= logrus.FatalLevel {
		l.Log(l.errLog).Fatal(v...)
	}
}

func (l *Logger) Panic(v ...interface{}) {
	if l.log.Level >= logrus.PanicLevel {
		l.Log(l.errLog).Panic(v...)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.log.Level >= logrus.DebugLevel {
		l.Log(l.log).Debug(v...)
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.log.Level >= logrus.WarnLevel {
		l.Log(l.log).Warn(v...)
	}
}

func (l *Logger) Warning(v ...interface{}) {
	if l.log.Level >= logrus.WarnLevel {
		l.Log(l.log).Warning(v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.log.Level >= logrus.ErrorLevel {
		l.Log(l.errLog).Error(v...)
	}
}

func (l *Logger) Println(v ...interface{}) {
	if l.log.Level >= logrus.InfoLevel {
		l.Log(l.log).Println(v...)
	}
}

func (l *Logger) Fatalln(v ...interface{}) {
	if l.log.Level >= logrus.FatalLevel {
		l.Log(l.errLog).Fatalln(v...)
	}
}

func (l *Logger) Panicln(v ...interface{}) {
	if l.log.Level >= logrus.PanicLevel {
		l.Log(l.errLog).Panicln(v...)
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	if l.log.Level >= logrus.DebugLevel {
		l.Log(l.log).Debugln(v...)
	}
}

func (l *Logger) Warnln(v ...interface{}) {
	if l.log.Level >= logrus.WarnLevel {
		l.Log(l.log).Warnln(v...)
	}
}

func (l *Logger) Warningln(v ...interface{}) {
	if l.log.Level >= logrus.WarnLevel {
		l.Log(l.log).Warningln(v...)
	}
}

func (l *Logger) Errorln(v ...interface{}) {
	if l.log.Level >= logrus.ErrorLevel {
		l.Log(l.errLog).Errorln(v...)
	}
}
