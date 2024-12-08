package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func mainUI(dataset []string, datasetIndex int) {
	var bs string
	if len(dataset) > datasetIndex && datasetIndex >= 0 {
		bs = dataset[datasetIndex]
	} else {
		bs = strings.Repeat(".", 9*9)
	}
	p := tea.NewProgram(HomeScreenModel{loadedString: bs})
	go p.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

var _ tea.Model = HomeScreenModel{}

type HomeScreenModel struct {
	loadedBoard  *Board
	loadedString string
	pos          BoardPos
	err          error
}

// Init implements tea.Model.
func (h HomeScreenModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (h HomeScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return h, tea.Quit
		case tea.KeyLeft.String():
			if h.pos.Col > 0 {
				h.pos.Col -= 1
			}
		case tea.KeyRight.String():
			if h.pos.Col < 8 {
				h.pos.Col += 1
			}
		case tea.KeyUp.String():
			if h.pos.Row > 0 {
				h.pos.Row -= 1
			}
		case tea.KeyDown.String():
			if h.pos.Row < 8 {
				h.pos.Row += 1
			}
		case "r":
			h.loadedBoard = NewStaticBoard()
			h.err = h.loadedBoard.LoadString(h.loadedString)
		case "s":
			h.loadedBoard = NewAutosolveBoard()
			h.err = h.loadedBoard.LoadString(h.loadedString)
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if h.err != nil {
				break
			}
			n := int(msg.String()[0] - '0')
			if h.loadedBoard != nil {
				h.err = h.loadedBoard.SetCertain(h.pos, n)
			}
		}
	}
	return h, nil
}

// View implements tea.Model.
func (h HomeScreenModel) View() string {
	titleString := strings.Repeat("=", 10) + " Sudoku Master " + strings.Repeat("=", 10)
	var lb *Board
	if h.loadedBoard == nil {
		lb = NewStaticBoard()
	} else {
		lb = h.loadedBoard
	}
	errPos := BoardPos{-1, -1}
	if err, ok := h.err.(PositionalError); ok {
		errPos = err.ErrorPosition()
	}
	bString := Pad(lb.FocussedString(h.pos, errPos), 7)
	instructions := "     Use arrow keys to move around"
	currentState := lb.At(h.pos)
	currentStateText := fmt.Sprintf("      Possibilities: %v", currentState.Possibilities())
	var errString string
	if h.err != nil {
		errString = fmt.Sprintf("   \033[31mError: %s\033[0m", h.err.Error())
	} else {
		errString = "             \033[32mConsistent\033[0m"
	}
	return titleString + "\n\n" + bString + "\n" + instructions + "\n" + currentStateText + "\n" + errString + "\n"
}
