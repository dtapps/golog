package golog

import (
	"testing"
)

func TestSystem(t *testing.T) {
	var s System
	s.Init()
	t.Logf("%+v", s)
}
