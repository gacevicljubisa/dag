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
func newVertex(id string, allowFail bool, action Action) *dagVertex {
	return &dagVertex{
		id:        id,
		allowFail: allowFail,
		action:    action,
	}
}

type dagVertex struct {
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
	v.state = Passed
}

func (v *dagVertex) SetFail() {
	if v.allowFail {
		v.state = Passed
	} else {
		v.state = Failed
	}
}

func (v *dagVertex) State() vertexState {
	return v.state
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

// allPredecessorsProcessed check if all its predecessors are processed with success
func (v *dagVertex) allPredecessorsProcessed() bool {
	for _, pred := range v.inEdges {
		if pred.state == Unprocessed { //|| (pred.state == Failed && !pred.allowFail) //check if this is the correct behavior
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
