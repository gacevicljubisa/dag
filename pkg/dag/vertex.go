package dag

import (
	"context"
	"fmt"
)

// Vertex interface definition
type Vertex interface {
	Id() string
	SetPass()
	SetFail()
	State() vertexState
	Execute(ctx context.Context) error
}

type vertexState int

const (
	Unprocessed vertexState = iota
	Passed
	Failed
)

// newVertex creates a new Vertex
func newVertex(dag *DAG, id string, allowFail bool, action Action) *dagVertex {
	return &dagVertex{
		dag:       dag,
		id:        id,
		allowFail: allowFail,
		action:    action,
	}
}

type dagVertex struct {
	dag       *DAG
	id        string
	state     vertexState
	allowFail bool
	inEdges   []*dagVertex
	outEdges  []*dagVertex
	action    Action
}

func (v *dagVertex) Id() string {
	return v.id
}

func (v *dagVertex) SetPass() {
	v.setState(Passed)
}

func (v *dagVertex) SetFail() {
	v.setState(Failed)
}

func (v *dagVertex) State() vertexState {
	return v.state
}

func (v *dagVertex) setState(state vertexState) {
	if v.dag.HasFailed() {
		panic(ErrDagAlreadyFailed)
	}
	v.state = state
}

// hasDescendant check if the vertex has the target vertex as a descendant
func (v *dagVertex) hasDescendant(target *dagVertex) bool {
	for _, descendant := range v.outEdges {
		if descendant == target || descendant.hasDescendant(target) {
			return true
		}
	}
	return false
}

// allPredecessorsProcessed check if all its predecessors are processed
func (v *dagVertex) allPredecessorsProcessed() bool {
	for _, pred := range v.inEdges {
		if pred.state == Unprocessed {
			return false
		}
	}
	return true
}

// Execute executes the action of the vertex
func (v *dagVertex) Execute(ctx context.Context) error {
	fmt.Printf("current vertex: %s\n", v.Id())
	return v.action.Run(ctx)
}
