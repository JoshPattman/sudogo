package main

import (
	"iter"
)

type ConsistencyRule interface {
	Apply(board *Board, bp BoardPos) ([]BoardPos, error)
}

type UniqueGroupRule struct {
	groupFunc func(BoardPos) iter.Seq[BoardPos]
}

func (g *UniqueGroupRule) Apply(board *Board, bp BoardPos) ([]BoardPos, error) {
	res, isRes := board.At(bp).Resolved()
	if !isRes {
		return nil, nil
	}
	changed := make([]BoardPos, 0)
	for tbp := range g.groupFunc(bp) {
		if bp == tbp {
			continue // Do not apply rule to own cell
		}
		if removed, err := board.At(tbp).Remove(res); err != nil {
			return nil, ImpossibleStateError{err, bp}
		} else if removed {
			changed = append(changed, tbp)
		}
	}
	return changed, nil
}

type CountPossibilityGroupRule struct {
	groupFunc     func(BoardPos) iter.Seq[BoardPos]
	countsBuffer  []pair[BoardPos, int]
	changedBuffer []BoardPos
}

type pair[T, U any] struct {
	A T
	B U
}

func (g *CountPossibilityGroupRule) Apply(board *Board, bp BoardPos) ([]BoardPos, error) {
	if g.countsBuffer == nil {
		g.countsBuffer = make([]pair[BoardPos, int], 9)
	} else {
		for i := range g.countsBuffer {
			g.countsBuffer[i] = pair[BoardPos, int]{}
		}
	}
	for cellBP := range g.groupFunc(bp) {
		cell := board.At(cellBP)
		for posI, isPos := range cell.mask {
			if !isPos {
				continue
			}
			g.countsBuffer[posI].A = cellBP
			g.countsBuffer[posI].B += 1
		}
	}
	if g.changedBuffer == nil {
		g.changedBuffer = make([]BoardPos, 0, 9)
	} else {
		g.changedBuffer = g.changedBuffer[:0]
	}
	for mustBeMinusOne, atCount := range g.countsBuffer {
		if atCount.B == 1 {
			changed, err := board.At(atCount.A).RemoveExcept(mustBeMinusOne + 1)
			if err != nil {
				return nil, ImpossibleStateError{err, atCount.A}
			} else if changed {
				g.changedBuffer = append(g.changedBuffer, atCount.A)
			}
		}
	}
	return g.changedBuffer, nil
}
