package main

import "fmt"

type PositionalError interface {
	ErrorPosition() BoardPos
}

type ImpossibleStateError struct {
	Err      error
	Position BoardPos
}

func (i ImpossibleStateError) Error() string {
	return fmt.Sprintf("Cell %v errored with %v", i.Position, i.Err)
}

func (i ImpossibleStateError) ErrorPosition() BoardPos {
	return i.Position
}
