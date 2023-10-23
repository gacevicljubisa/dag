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
	return d.addVertex(id, allowFail, action)
}

func (d *DAG) addVertex(id string, allowFail bool, action Action) error {
	if d.started {
		return ErrCanNotModify
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.vertices[id]; exists {
		return ErrVertexExists{id}
	}

	d.vertices[id] = &dagVertex{
		id:        id,
		allowFail: allowFail,
		Action:    action,
	}

	return nil
}

func (d *DAG) AddEdge(fromID, toID string) error {
	if d.started {
		return ErrCanNotModify
	}

	fromVertex, ok1 := d.vertices[fromID]
	if !ok1 {
		return ErrVertexInvalid{id: fromID}
	}

	toVertex, ok2 := d.vertices[toID]
	if !ok2 {
		return ErrVertexInvalid{id: toID}
	}

	// check if the edge already exists
	for _, v := range fromVertex.outEdges {
		if v == toVertex {
			return nil // edge already exists, do nothing
		}
	}

	// deep check for cyclic relations
	if toVertex.hasDescendant(fromVertex) {
		return ErrCyclicRelation
	}

	fromVertex.outEdges = append(fromVertex.outEdges, toVertex)
	toVertex.inEdges = append(toVertex.inEdges, fromVertex)

	return nil
}

func (d *DAG) Next() []Vertex {
	if d.HasFinished() {
		panic("DAG has finished")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.started = true

	// TODO: check this if ok
	var readyVertices []Vertex
	for _, v := range d.vertices {
		if v.state == Unprocessed && len(v.inEdges) == 0 {
			readyVertices = append(readyVertices, v)
		}
	}
	return readyVertices
}

func (d *DAG) HasFailed() bool {
	for _, v := range d.vertices {
		if v.state == Failed {
			return true
		}
	}
	return false
}

func (d *DAG) HasSucceeded() bool {
	for _, v := range d.vertices {
		if v.state != Passed {
			return false
		}
	}
	return true
}

// HasFinished returns true if all vertices have been processed.
// Needs to be improved, because if DAG grows large, this will be slow. Maybe cache not processed vertices?
func (d *DAG) HasFinished() bool {
	for _, v := range d.vertices {
		if v.state == Unprocessed {
			return false
		}
	}
	return true
}
