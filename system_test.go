package golog

import (
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestSystem(t *testing.T) {
	var s System
	s.Init()
	t.Logf("%+v", s)
	t.Logf("%+v", runtime.Version())
	goVersion, _ := strconv.ParseFloat(strings.TrimPrefix(runtime.Version(), "go"), 64)
	t.Logf("%+v", goVersion)
}
