package time

import (
	"context"
	"testing"
	"time"
)

func TestShrink(t *testing.T) {
	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	if err != nil {
		t.Fatalf("TestShrink:  d.UnmarshalText failed!err:=%v", err)
	}
	c := context.Background()
	to, ctx, cancel := d.Shrink(c)
	defer cancel()
	if time.Duration(to) != time.Second {
		t.Fatalf("new timeout must be equal 1 second")
	}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > time.Second || time.Until(deadline) < time.Millisecond*500 {
		t.Fatalf("ctx deadline must be less than 1s and greater than 500ms")
	}
}
