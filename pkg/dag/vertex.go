package dag

// Vertex interface definition
type Vertex interface {
	Id() string
	SetPass()
	SetFail()
	State() vertexState
}

type vertexState int

const (
	Unprocessed vertexState = iota
	Passed
	Failed
)

type dagVertex struct {
	id        string
	state     vertexState
	allowFail bool
	inEdges   []*dagVertex
	outEdges  []*dagVertex
	Action    Action
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

func (v *dagVertex) hasDescendant(target *dagVertex) bool {
	for _, descendant := range v.outEdges {
		if descendant == target || descendant.hasDescendant(target) {
			return true
		}
	}
	return false
}
