package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/gacevicljubisa/dag/pkg/dag"
	"golang.org/x/sync/errgroup"
)

func main() {
	d := dag.NewDAG()
	action := &SimpleAction{
		Func: func(ctx context.Context) error {
			return errors.New("error")
		},
	}

	ctx := context.Background()

	var err error

	d.AddVertex("A", false, getUsersParalel())
	d.AddVertex("B", false, getUsers())
	d.AddVertex("C", true, action)
	d.AddVertex("D", true, NewLoop(6))
	d.AddEdge("A", "B")
	d.AddEdge("B", "C")
	d.AddEdge("A", "C")
	// d.AddEdge("C", "A") // panic: cyclic relation

	err = nextRecursive(ctx, d)
	if err != nil {
		fmt.Println("Error occured: ", err)
	}

	// d.Next() // panic: DAG has finished

	fmt.Println("Is failed: ", d.HasFailed())
	fmt.Println("Is succeded: ", d.HasSucceeded())
	fmt.Println("Is finished: ", d.HasFinished())
}

// nextRecursive processes the DAG's vertices in parallel, marking each as 'Pass' or 'Fail' based on execution.
// The function recurses until all vertices are processed.
func nextRecursive(ctx context.Context, d *dag.DAG) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, v := range d.Next() {
		currentVertex := v
		g.Go(func() error {
			errExec := currentVertex.Execute(ctx)
			if errExec != nil {
				currentVertex.SetFail()
				return errExec
			}
			currentVertex.SetPass()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	if !d.HasFinished() {
		return nextRecursive(ctx, d)
	}

	fmt.Println("Is finished: ", d.HasFinished())

	return nil
}
