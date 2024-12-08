package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	dataset := flag.String("dataset", "", "If specified, this dataset will be loaded")
	runUI := flag.Bool("tui", false, "Should we run the TUI")
	tuiIndex := flag.Int("index", 0, "The index of puzzle to load into the TUI from the dataset")
	outputDir := flag.String("output", "./sudoku-output", "The output directory to save files to")
	flag.Parse()

	datasetLines := []string{}

	if *dataset != "" {
		datasetLines = loadDataset(*dataset)
	}

	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	if *runUI {
		mainUI(datasetLines, *tuiIndex)
	} else {
		if len(datasetLines) == 0 {
			log.Fatalf("Must specify a valid dataset for computation in cli mode\n")
		}
		mainSolver(datasetLines, *outputDir)
	}

}

func loadDataset(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open dataset at '%s': %v\n", path, err)
	}
	defer f.Close()
	datasetBytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read dataset at '%s': %v\n", path, err)
	}
	lines := make([]string, 0)
	for _, l := range strings.Split(string(datasetBytes), "\n") {
		l = strings.Trim(l, "\n\r \t")
		if len(l) == 0 {
			continue
		}
		lines = append(lines, l)
	}
	return lines
}

func mainSolver(boardStrings []string, outputDir string) {
	fmt.Println("Loaded", len(boardStrings), "boards")

	solved := 0
	solvedPuzzles := make([]string, 0)
	unsolvedPuzzles := make([]string, 0)
	solutions := make([]string, 0, len(boardStrings))

	tstart := time.Now()
	for bsi, bs := range boardStrings {
		b := NewBoard([]ConsistencyRule{
			&UniqueGroupRule{RowGroup},
			&UniqueGroupRule{ColGroup},
			&UniqueGroupRule{SquareGroup},
			&CountPossibilityGroupRule{RowGroup, nil, nil},
			&CountPossibilityGroupRule{ColGroup, nil, nil},
			&CountPossibilityGroupRule{SquareGroup, nil, nil},
		})
		err := b.LoadString(bs)
		if err != nil {
			panic(err)
		}
		if b.Solved() {
			solved++
			solvedPuzzles = append(solvedPuzzles, boardStrings[bsi])
		} else {
			unsolvedPuzzles = append(unsolvedPuzzles, boardStrings[bsi])
		}
		solutions = append(solutions, b.ExportString())
		if bsi%100 == 0 {
			fmt.Printf("\r%d (%.2f%%)   ", bsi, float64(bsi)/float64(len(boardStrings))*100)
		}
	}
	fmt.Println()
	tdone := time.Since(tstart)
	fmt.Printf("Solved %d of %d (%.2f%%) leaving %d unsolved [%v]\n", solved, len(boardStrings), (float64(solved)/float64(len(boardStrings)))*100, len(boardStrings)-solved, tdone)

	for _, pair := range []pair[string, []string]{
		{"solvable_sudoku.txt", solvedPuzzles},
		{"unsolvable_sudoku.txt", unsolvedPuzzles},
		{"results_sudoku.txt", solutions},
		{"input_sudoku.txt", boardStrings},
	} {
		func() {
			f, err := os.Create(path.Join(outputDir, pair.A))
			if err != nil {
				panic(err)
			}
			defer f.Close()
			f.WriteString(strings.Join(pair.B, "\n"))
		}()
	}
}
