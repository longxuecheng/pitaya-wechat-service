package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	l := New()
	l.Debug("Debug message ...")
	l.Info("Info message ...")
	l.Warn("Warn message ...")
	l.Error("Error message ...")
	l.Fatal("Fatal message ...")
}
