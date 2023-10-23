package dag

import (
	"errors"
	"fmt"
)

var (
	ErrCanNotModify   = errors.New("can't modify DAG after calling Next")
	ErrCyclicRelation = errors.New("cyclic relation detected")
)

type ErrVertexExists struct {
	id string
}

func (e ErrVertexExists) Error() string {
	return fmt.Sprintf("vertex with ID '%s' already exists", e.id)
}

type ErrVertexInvalid struct {
	id string
}

func (e ErrVertexInvalid) Error() string {
	return fmt.Sprintf("vertex with ID '%s' doesn't exist", e.id)
}
