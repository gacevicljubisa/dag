package main

import (
	"context"
	"fmt"

	"github.com/gacevicljubisa/dag/pkg/dag"
)

func main() {
	d := dag.NewDAG()
	action := &SimpleAction{
		Func: func(ctx context.Context) error {
			fmt.Println("Hello World 1")
			return nil
		},
	}

	d.AddVertex("A", false, action)
	d.AddVertex("B", false, action)
	d.AddVertex("C", true, action)
	d.AddEdge("A", "B")

	fmt.Println("Initial Next:", d.Next())

	v := d.Next()[0]
	v.SetPass()
	fmt.Println("Next after A:", d.Next())

	v = d.Next()[0]
	v.SetFail()
	fmt.Println("HasFailed:", d.HasFailed())
	fmt.Println("HasSucceeded:", d.HasSucceeded())
	fmt.Println("HasFinished:", d.HasFinished())
}
