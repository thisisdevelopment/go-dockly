package xlogger

import (
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testHook is a simple hook that records which levels it was fired for
type testHook struct {
	mu     sync.Mutex
	fired  []logrus.Level
	levels []logrus.Level
}

func newTestHook(levels []logrus.Level) *testHook {
	return &testHook{
		levels: levels,
		fired:  make([]logrus.Level, 0),
	}
}

func (h *testHook) Levels() []logrus.Level {
	return h.levels
}

func (h *testHook) Fire(entry *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fired = append(h.fired, entry.Level)
	return nil
}

func (h *testHook) getFired() []logrus.Level {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]logrus.Level, len(h.fired))
	copy(result, h.fired)
	return result
}

func (h *testHook) wasFiredFor(level logrus.Level) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, l := range h.fired {
		if l == level {
			return true
		}
	}
	return false
}

func TestNew(t *testing.T) {
	t.Run("creates logger with default config", func(t *testing.T) {
		logger, err := New(nil)
		require.NoError(t, err)
		assert.NotNil(t, logger)
		assert.Equal(t, "debug", logger.Level())
	})

	t.Run("creates logger with custom config", func(t *testing.T) {
		logger, err := New(&Config{Level: "info", Format: "json"})
		require.NoError(t, err)
		assert.NotNil(t, logger)
		assert.Equal(t, "info", logger.Level())
	})

	t.Run("returns error for invalid log level", func(t *testing.T) {
		_, err := New(&Config{Level: "invalid"})
		assert.Error(t, err)
	})
}

func TestAddHook(t *testing.T) {
	t.Run("hook is added and accessible", func(t *testing.T) {
		logger, err := New(&Config{Level: "debug"})
		require.NoError(t, err)

		hook := newTestHook(logrus.AllLevels)
		logger.AddHook(hook)

		assert.NotEmpty(t, logger.Hooks())
	})
}

func TestHooksFireForAllLogLevels(t *testing.T) {
	// This test verifies the fix: hooks should fire for ALL log levels,
	// including Error, Fatal, and Panic (which previously used errLog)

	tests := []struct {
		name     string
		logFunc  func(l *Logger)
		level    logrus.Level
		skipExit bool // skip for Fatal/Panic which would exit
	}{
		{
			name:    "Debug fires hook",
			logFunc: func(l *Logger) { l.Debug("test") },
			level:   logrus.DebugLevel,
		},
		{
			name:    "Debugf fires hook",
			logFunc: func(l *Logger) { l.Debugf("test %s", "msg") },
			level:   logrus.DebugLevel,
		},
		{
			name:    "Debugln fires hook",
			logFunc: func(l *Logger) { l.Debugln("test") },
			level:   logrus.DebugLevel,
		},
		{
			name:    "Info fires hook",
			logFunc: func(l *Logger) { l.Info("test") },
			level:   logrus.InfoLevel,
		},
		{
			name:    "Infof fires hook",
			logFunc: func(l *Logger) { l.Infof("test %s", "msg") },
			level:   logrus.InfoLevel,
		},
		{
			name:    "Warn fires hook",
			logFunc: func(l *Logger) { l.Warn("test") },
			level:   logrus.WarnLevel,
		},
		{
			name:    "Warnf fires hook",
			logFunc: func(l *Logger) { l.Warnf("test %s", "msg") },
			level:   logrus.WarnLevel,
		},
		{
			name:    "Warnln fires hook",
			logFunc: func(l *Logger) { l.Warnln("test") },
			level:   logrus.WarnLevel,
		},
		{
			name:    "Warning fires hook",
			logFunc: func(l *Logger) { l.Warning("test") },
			level:   logrus.WarnLevel,
		},
		{
			name:    "Warningf fires hook",
			logFunc: func(l *Logger) { l.Warningf("test %s", "msg") },
			level:   logrus.WarnLevel,
		},
		{
			name:    "Warningln fires hook",
			logFunc: func(l *Logger) { l.Warningln("test") },
			level:   logrus.WarnLevel,
		},
		{
			name:    "Error fires hook",
			logFunc: func(l *Logger) { l.Error("test") },
			level:   logrus.ErrorLevel,
		},
		{
			name:    "Errorf fires hook",
			logFunc: func(l *Logger) { l.Errorf("test %s", "msg") },
			level:   logrus.ErrorLevel,
		},
		{
			name:    "Errorln fires hook",
			logFunc: func(l *Logger) { l.Errorln("test") },
			level:   logrus.ErrorLevel,
		},
		{
			name:    "Print fires hook",
			logFunc: func(l *Logger) { l.Print("test") },
			level:   logrus.InfoLevel,
		},
		{
			name:    "Printf fires hook",
			logFunc: func(l *Logger) { l.Printf("test %s", "msg") },
			level:   logrus.InfoLevel,
		},
		{
			name:    "Println fires hook",
			logFunc: func(l *Logger) { l.Println("test") },
			level:   logrus.InfoLevel,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := New(&Config{Level: "debug"})
			require.NoError(t, err)

			hook := newTestHook(logrus.AllLevels)
			logger.AddHook(hook)

			tc.logFunc(logger)

			assert.True(t, hook.wasFiredFor(tc.level),
				"expected hook to fire for level %s", tc.level)
		})
	}
}

func TestHooksFireWithError(t *testing.T) {
	// Test that WithError().Error() fires hooks (this always worked)
	logger, err := New(&Config{Level: "debug"})
	require.NoError(t, err)

	hook := newTestHook(logrus.AllLevels)
	logger.AddHook(hook)

	logger.WithError(assert.AnError).Error("test error")

	assert.True(t, hook.wasFiredFor(logrus.ErrorLevel),
		"expected hook to fire for Error level when using WithError")
}

func TestHooksFireWithField(t *testing.T) {
	logger, err := New(&Config{Level: "debug"})
	require.NoError(t, err)

	hook := newTestHook(logrus.AllLevels)
	logger.AddHook(hook)

	logger.WithField("key", "value").Error("test error")

	assert.True(t, hook.wasFiredFor(logrus.ErrorLevel),
		"expected hook to fire for Error level when using WithField")
}

func TestHooksFireWithFields(t *testing.T) {
	logger, err := New(&Config{Level: "debug"})
	require.NoError(t, err)

	hook := newTestHook(logrus.AllLevels)
	logger.AddHook(hook)

	logger.WithFields(logrus.Fields{"key": "value"}).Error("test error")

	assert.True(t, hook.wasFiredFor(logrus.ErrorLevel),
		"expected hook to fire for Error level when using WithFields")
}

func TestHookLevelFiltering(t *testing.T) {
	// Test that hooks only fire for their configured levels
	logger, err := New(&Config{Level: "debug"})
	require.NoError(t, err)

	// Hook only for error level
	hook := newTestHook([]logrus.Level{logrus.ErrorLevel})
	logger.AddHook(hook)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	fired := hook.getFired()
	assert.Len(t, fired, 1)
	assert.Equal(t, logrus.ErrorLevel, fired[0])
}

func TestToFields(t *testing.T) {
	logger, err := New(&Config{Level: "debug"})
	require.NoError(t, err)

	t.Run("creates fields from pairs", func(t *testing.T) {
		fields := logger.ToFields("key1", "value1", "key2", "value2")
		assert.Equal(t, "value1", fields["key1"])
		assert.Equal(t, "value2", fields["key2"])
	})
}

func TestLogFormats(t *testing.T) {
	formats := []string{"text", "json", "gelf"}

	for _, format := range formats {
		t.Run("format "+format, func(t *testing.T) {
			logger, err := New(&Config{Level: "debug", Format: format})
			require.NoError(t, err)
			assert.NotNil(t, logger)
		})
	}
}
