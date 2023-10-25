package dag

import (
	"errors"
	"fmt"
)

var (
	ErrCanNotModify     = errors.New("can't modify DAG after calling Next function")
	ErrDagHasFinished   = errors.New("DAG has finished")
	ErrDagAlreadyFailed = errors.New("can't set vertex state because DAG has already failed")
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

type ErrCyclicRelation struct {
	fromID, toID string
}

func (e ErrCyclicRelation) Error() string {
	return fmt.Sprintf("cyclic relation detected between vertex '%s' and '%s'", e.fromID, e.toID)
}
