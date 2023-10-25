package dag

import (
	"context"
	"sync"
)

type Action interface {
	Run(ctx context.Context) (err error)
}

// Dag interface definition
type Dag interface {
	Next() []Vertex
	HasFailed() bool
	HasSucceeded() bool
	HasFinished() bool
}

type DAG struct {
	mu       sync.Mutex
	vertices map[string]*dagVertex
	started  bool
}

func NewDAG() *DAG {
	return &DAG{vertices: make(map[string]*dagVertex)}
}

func (d *DAG) AddVertex(id string, allowFail bool, action Action) error {
	if d.started {
		panic(ErrCanNotModify)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.vertices[id]; exists {
		panic(ErrVertexExists{id}.Error())
	}

	d.vertices[id] = newVertex(d, id, allowFail, action)

	return nil
}

func (d *DAG) AddEdge(fromID, toID string) error {
	if d.started {
		panic(ErrCanNotModify)
	}

	fromVertex, ok1 := d.vertices[fromID]
	if !ok1 {
		panic(ErrVertexInvalid{id: fromID}.Error())
	}

	toVertex, ok2 := d.vertices[toID]
	if !ok2 {
		panic(ErrVertexInvalid{id: toID}.Error())
	}

	// check if the edge already exists
	for _, v := range fromVertex.outEdges {
		if v == toVertex {
			return nil // edge already exists, do nothing
		}
	}

	// deep check for cyclic relations
	if toVertex.hasDescendant(fromVertex) {
		panic(ErrCyclicRelation{fromID, toID}.Error())
	}

	fromVertex.outEdges = append(fromVertex.outEdges, toVertex)
	toVertex.inEdges = append(toVertex.inEdges, fromVertex)

	return nil
}

func (d *DAG) Next() []Vertex {
	if d.HasFinished() {
		panic(ErrDagHasFinished)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.started = true

	// find all ready vertices (vertices with no inEdges or all inEdges are processed)
	var readyVertices []Vertex
	for _, v := range d.vertices {
		if v.state == Unprocessed {
			if len(v.inEdges) == 0 || v.allPredecessorsProcessed() {
				readyVertices = append(readyVertices, v)
			}
		}
	}
	return readyVertices
}

func (d *DAG) HasFailed() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, v := range d.vertices {
		if v.state == Failed && !v.allowFail {
			return true
		}
	}
	return false
}

func (d *DAG) HasSucceeded() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, v := range d.vertices {
		if v.state != Passed || (v.state == Failed && v.allowFail) {
			return false
		}
	}
	return true
}

// HasFinished returns true if all vertices have been processed.
// Needs to be improved, because if DAG grows large, this will be slow. Maybe cache not processed vertices?
func (d *DAG) HasFinished() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, v := range d.vertices {
		if v.state == Unprocessed {
			return false
		}
	}
	return true
}
