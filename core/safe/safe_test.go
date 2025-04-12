package safe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	testFn := func() {
		defer Recover()

		panic("recover panic")
	}
	runFn := func() (err error) {
		defer RecoverToError(&err)

		testFn()

		return err
	}

	assert.NoError(t, runFn())
}

func TestGo(t *testing.T) {
	runFn := func() (err error) {
		defer RecoverToError(&err)

		done := make(chan struct{})

		Go(func() {
			close(done)
			panic("go panic")
		})

		<-done

		return err
	}

	assert.NoError(t, runFn())
}

func TestFn(t *testing.T) {
	got := Fn(func() error {
		panic("fn panic")
	})

	assert.Error(t, got())
}
