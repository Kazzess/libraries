package safe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Kazzess/libraries/core/safe"

	"github.com/stretchr/testify/assert"
)

func TestGroup_Go(t *testing.T) {
	ctx := context.Background()
	want := errors.New("test err")

	eg, egCtx := safe.WithContext(ctx)

	eg.Go(func() error {
		<-egCtx.Done()

		return nil
	})

	eg.Go(func() error {
		<-egCtx.Done()

		return nil
	})

	eg.Run(func(ctx context.Context) error {
		<-ctx.Done()

		return nil
	})

	eg.Run(func(ctx context.Context) error {
		return want
	})

	got := eg.Wait()
	assert.Equal(t, want, got)
}

func TestGroup_Run(t *testing.T) {
	ctx := context.Background()

	eg, _ := safe.WithContext(ctx)

	eg.Go(func() error {
		panic("test panic")

		return nil //nolint:govet
	})

	got := eg.Wait()
	assert.Error(t, got)
}
