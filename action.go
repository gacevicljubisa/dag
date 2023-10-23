package main

import (
	"context"

	"github.com/gacevicljubisa/dag/pkg/dag"
	"golang.org/x/sync/errgroup"
)

// SimpleAction is an action that runs a function
type SimpleAction struct {
	Func func(ctx context.Context) error
}

// Run runs the function
func (sa *SimpleAction) Run(ctx context.Context) error {
	return sa.Func(ctx)
}

// CompositeAction is an action that runs multiple actions in sequence
type CompositeAction struct {
	Actions []dag.Action
}

// Run runs all actions in sequence
func (ca *CompositeAction) Run(ctx context.Context) error {
	for _, action := range ca.Actions {
		if err := action.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}

// ParallelActions is an action that runs multiple actions in parallel
type ParallelActions struct {
	Actions []dag.Action
}

// Run runs all actions in parallel
func (pa *ParallelActions) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, action := range pa.Actions {
		currAction := action

		g.Go(func() error {
			return currAction.Run(ctx)
		})
	}

	return g.Wait()
}
