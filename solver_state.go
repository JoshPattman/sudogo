package main

import "fmt"

func NewCellState() *CellState {
	return &CellState{
		mask:          [9]bool{true, true, true, true, true, true, true, true, true},
		possibleCount: 9,
		resolved:      false,
		resolvedTo:    0,
		wasOriginal:   false,
	}
}

type CellState struct {
	mask          [9]bool
	possibleCount int
	resolved      bool
	resolvedTo    int
	wasOriginal   bool
}

func (c *CellState) validateI(i int) {
	if i < 1 || i > 9 {
		panic(fmt.Sprintf("tried to ask if cell state %v was possible", i))
	}
}

// Returns if index i is a possibility for this state
func (c *CellState) Possible(i int) bool {
	c.validateI(i)
	return c.mask[i-1]
}

func (c *CellState) PossibleCount() int {
	return c.possibleCount
}

func (c *CellState) Possibilities() []int {
	ps := make([]int, 0)
	for i, ok := range c.mask {
		if ok {
			ps = append(ps, i+1)
		}
	}
	return ps
}

func (c *CellState) RemoveExcept(i int) (bool, error) {
	c.validateI(i)
	if c.mask[i-1] {
		changed := false
		for j := 1; j <= 9; j++ {
			if i == j {
				continue
			}
			if c.mask[j-1] {
				c.mask[j-1] = false
				changed = true
			}
		}
		c.possibleCount = 1
		c.resolved = true
		c.resolvedTo = i
		return changed, nil
	} else {
		return false, fmt.Errorf("impossible, reached count of %d", 0)
	}
}

func (c *CellState) Remove(i int) (bool, error) {
	c.validateI(i)
	if c.mask[i-1] {
		c.mask[i-1] = false
		c.possibleCount--
		if c.possibleCount < 1 {
			return false, fmt.Errorf("impossible, reached count of %d", c.possibleCount)
		}
		c.resolved = c.possibleCount == 1
		if c.resolved {
			for i := range 9 {
				if c.mask[i] {
					c.resolvedTo = i + 1
					break
				}
			}
		}
		return true, nil
	}
	return false, nil
}

func (c *CellState) Resolved() (int, bool) {
	if c.resolved {
		return c.resolvedTo, true
	}
	return 0, false
}
