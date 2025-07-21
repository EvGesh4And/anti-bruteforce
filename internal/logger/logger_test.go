package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger_DebugLevel(t *testing.T) {
	var buf bytes.Buffer
	lg := New("debug", &buf, false)
	lg.Debug("event validation check")
	lg.Info("event added")
	lg.Warn("connection lost, trying to restore it")
	require.Contains(t, buf.String(), "event validation check")
	require.Contains(t, buf.String(), "event added")
	require.Contains(t, buf.String(), "connection lost, trying to restore it")
}

func TestLogger_InfoLevel(t *testing.T) {
	var buf bytes.Buffer
	lg := New("info", &buf, false)
	lg.Info("event added")
	lg.Warn("connection lost, trying to restore it")

	out := buf.String()
	require.Contains(t, out, "event added")
	require.Contains(t, out, "connection lost, trying to restore it")
	require.NotContains(t, out, "event validation check") // Debug не должен попасть
}

func TestLogger_WarnLevel(t *testing.T) {
	var buf bytes.Buffer
	lg := New("warn", &buf, false)
	lg.Info("event added")
	lg.Warn("connection lost, trying to restore it")

	out := buf.String()
	require.NotContains(t, out, "event added") // Info не должен попасть
	require.Contains(t, out, "connection lost, trying to restore it")
}

func TestLogger_ErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	lg := New("error", &buf, false)
	lg.Info("event added")
	lg.Error("database connection completely lost")
	lg.Warn("connection lost, trying to restore it")
	lg.Debug("event validation check")

	out := buf.String()
	require.Contains(t, out, "database connection completely lost")
	require.NotContains(t, out, "event added")
	require.NotContains(t, out, "connection lost, trying to restore it")
	require.NotContains(t, out, "event validation check")
}
