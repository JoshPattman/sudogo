package main

import "fmt"

type BoardPos struct {
	Row int
	Col int
}

type Board struct {
	cells [9][9]*CellState
	rules []ConsistencyRule
}

func NewBoard(rules []ConsistencyRule) *Board {
	b := &Board{
		rules: rules,
	}
	for row := range 9 {
		for col := range 9 {
			b.cells[row][col] = NewCellState()
		}
	}
	return b
}

func NewStaticBoard() *Board {
	return NewBoard([]ConsistencyRule{})
}

func NewAutosolveBoard() *Board {
	return NewBoard([]ConsistencyRule{
		&UniqueGroupRule{RowGroup},
		&UniqueGroupRule{ColGroup},
		&UniqueGroupRule{SquareGroup},
		&CountPossibilityGroupRule{RowGroup, nil, nil},
		&CountPossibilityGroupRule{ColGroup, nil, nil},
		&CountPossibilityGroupRule{SquareGroup, nil, nil},
	})
}

func (b *Board) LoadString(s string) error {
	if len(s) != 9*9 {
		return fmt.Errorf("incorrect length of string board: %v", s)
	}
	for row := range 9 {
		for col := range 9 {
			val := int(s[col+row*9] - '0')
			if val >= 1 && val <= 9 {
				err := b.SetCertain(BoardPos{row, col}, val)
				if err != nil {
					return err
				}
				b.At(BoardPos{row, col}).wasOriginal = true
			}
		}
	}
	return nil
}

func (b *Board) ExportString() string {
	s := ""
	for row := range 9 {
		for col := range 9 {
			v := b.At(BoardPos{row, col})
			if resVal, resOK := v.Resolved(); resOK {
				s += fmt.Sprint(resVal)
			} else {
				s += "."
			}
		}
	}
	return s
}

func (b *Board) FocussedString(bp BoardPos, errorPos BoardPos) string {
	s := ""
	for row := range 9 {
		for col := range 9 {
			rval, resolved := b.cells[row][col].Resolved()
			space := " "
			bpt := BoardPos{row, col}
			if bpt == bp {
				space = "<"
			}
			cell := ""
			if resolved {
				cell = fmt.Sprintf("%v", rval)
			} else {
				cell = "-"
			}
			cell += space
			if errorPos == bpt {
				cell = fmt.Sprintf("\033[31m%s\033[0m", cell)
			} else if b.At(bpt).wasOriginal {
				cell = fmt.Sprintf("\033[32m%s\033[0m", cell)
			}
			s += cell
			if col == 2 || col == 5 {
				s += "| "
			}
		}
		s += "\n"
		if row == 2 || row == 5 {
			s += "----------------------\n"
		}
	}
	return s
}

func (b *Board) String() string {
	return b.FocussedString(BoardPos{-1, -1}, BoardPos{-1, -1})
}

func (b *Board) SetCertain(bp BoardPos, to int) error {
	state := b.At(bp)
	if _, err := state.RemoveExcept(to); err != nil {
		return ImpossibleStateError{err, bp}
	}
	return b.maintainConsistency(bp)
}

func (b *Board) maintainConsistency(pos BoardPos) error {
	stack := make(map[BoardPos]struct{})
	for _, rule := range b.rules {
		changed, err := rule.Apply(b, pos)
		if err != nil {
			return err
		}
		for _, c := range changed {
			stack[c] = struct{}{}
		}
	}
	for bp := range stack {
		if err := b.maintainConsistency(bp); err != nil {
			return err
		}
	}
	return nil
}

func (b *Board) At(bp BoardPos) *CellState {
	return b.cells[bp.Row][bp.Col]
}

func (b *Board) Solved() bool {
	for row := range 9 {
		for col := range 9 {
			if _, ok := b.At(BoardPos{row, col}).Resolved(); !ok {
				return false
			}
		}
	}
	return true
}
