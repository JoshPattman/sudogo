package main

import "iter"

func RowGroup(of BoardPos) iter.Seq[BoardPos] {
	return func(yield func(BoardPos) bool) {
		for i := range 9 {
			if !yield(BoardPos{of.Row, i}) {
				return
			}
		}
	}
}
func ColGroup(of BoardPos) iter.Seq[BoardPos] {
	return func(yield func(BoardPos) bool) {
		for i := range 9 {
			if !yield(BoardPos{i, of.Col}) {
				return
			}
		}
	}
}
func SquareGroup(of BoardPos) iter.Seq[BoardPos] {
	return func(yield func(BoardPos) bool) {
		sRow := (of.Row / 3) * 3
		sCol := (of.Col / 3) * 3
		for ro := range 3 {
			for co := range 3 {
				if !yield(BoardPos{sRow + ro, sCol + co}) {
					return
				}
			}
		}
	}
}
