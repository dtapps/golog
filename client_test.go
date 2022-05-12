package golog

import "testing"

var c = &Client{}

func TestClient(t *testing.T) {
	c = NewClientGin(nil, "")
	c = NewClientApi(nil, "")
}
