package golog

import (
	"runtime"
	"strings"
	"testing"
)

func TestSystem(t *testing.T) {
	var s System
	s.Init()
	t.Logf("%+v", s)
	t.Logf("%+v", runtime.Version())
	t.Logf("%+v", strings.TrimPrefix(runtime.Version(), "go"))
}
