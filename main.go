package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gacevicljubisa/dag/pkg/dag"
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
		panic(err)
	}

	// d.Next() // panic: DAG has finished

	fmt.Println("Is succeded: ", d.HasSucceeded())
}

// nextRecursive processes the DAG's vertices in parallel, marking each as 'Pass' or 'Fail' based on execution.
// The function recurses until all vertices are processed.
func nextRecursive(ctx context.Context, d *dag.DAG) error {
	var wg sync.WaitGroup
	for _, v := range d.Next() {
		wg.Add(1)
		go func(ctx context.Context, currentVertex dag.Vertex) {
			defer wg.Done()
			errExec := currentVertex.Execute(ctx)
			if errExec != nil {
				currentVertex.SetFail()
			} else {
				currentVertex.SetPass()
			}
		}(ctx, v)
	}

	wg.Wait()

	if !d.HasFinished() {
		return nextRecursive(ctx, d)
	}

	fmt.Println("Is finished: ", d.HasFinished())

	return nil
}
