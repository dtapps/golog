package golog

import "testing"

var a = App{}

func TestApp(t *testing.T) {
	a.Pgsql = nil
}
