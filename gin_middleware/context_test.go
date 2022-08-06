package gin_middleware

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {

	ctx := context.Background()
	t.Log(ctx)

	ctx = context.WithValue(ctx, "k1", "测试")

	t.Log()
}
