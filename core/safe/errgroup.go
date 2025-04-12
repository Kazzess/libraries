package safe

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Group represent a wrapper for *errgroup.Group.
type Group struct {
	// The *errgroup.Group instance.
	eg *errgroup.Group
	// Context for Run.
	ctx context.Context
}

// Go calls fn with error recovery.
func (g *Group) Go(fn func() error) {
	g.eg.Go(func() (err error) {
		defer RecoverToError(&err)

		return fn()
	})
}

// Run calls fn with group context with error recovery.
func (g *Group) Run(fn func(ctx context.Context) error) {
	g.eg.Go(func() (err error) {
		defer RecoverToError(&err)

		return fn(g.ctx)
	})
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	return g.eg.Wait()
}

// WithContext returns a new *Group with associated context.
func WithContext(ctx context.Context) (*Group, context.Context) {
	group, ctx := errgroup.WithContext(ctx)

	return &Group{
		eg:  group,
		ctx: ctx,
	}, ctx
}
