package main

import (
	"fmt"
	"sort"
	"strings"
)

type DoesNotFitError struct{}

func (e *DoesNotFitError) Error() string {
	return "Does Not Fit"
}

type FilledGrid struct {
	rows [][]string
}

func (fg *FilledGrid) String() string {
	var s strings.Builder
	for _, row := range fg.rows {
		s.WriteString(strings.Join(row, "") + "\n")
	}
	return s.String()
}

func NewFilledGrid(inputString string) *FilledGrid {
	fg := &FilledGrid{}
	splitInputString := strings.Split(inputString, "\n")
	lines := splitInputString[1 : len(splitInputString)-1]
	for _, line := range lines {
		fg.rows = append(fg.rows, strings.Split(line, ""))
	}
	return fg
}

func (fg *FilledGrid) Set(x, y int, char string) {
	fg.rows[y][x] = char
}

func (fg *FilledGrid) canInsert(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}
	if y >= len(fg.rows) || x >= len(fg.rows[0]) {
		return false
	}
	return fg.rows[y][x] == " "
}

func (fg *FilledGrid) Rotate() *FilledGrid {
	newRows := make([][]string, len(fg.rows[0]))
	for i := 0; i < len(fg.rows[0]); i++ {
		newRow := make([]string, len(fg.rows))
		for j := 0; j < len(fg.rows); j++ {
			newRow[j] = fg.rows[j][len(fg.rows[0])-i-1]
		}
		newRows[i] = newRow
	}

	newDog := &FilledGrid{rows: newRows}
	return newDog
}

func (fg *FilledGrid) Size() int {
	size := 0
	for _, row := range fg.rows {
		for _, val := range row {
			if val != " " {
				size++
			}
		}
	}
	return size
}

type Board struct {
	*FilledGrid
}

func NewBoard(inputString string) *Board {
	return &Board{FilledGrid: NewFilledGrid(inputString)}
}

func (b *Board) Insert(other *FilledGrid, x, y int, charToInsert string) error {
	for y_offset, row := range other.rows {
		for x_offset := range row {
			if !other.canInsert(x_offset, y_offset) {
				if !b.canInsert(x+x_offset, y+y_offset) {
					return &DoesNotFitError{}
				}				
			}
		}
	}
	for y_offset, row := range other.rows {
		for x_offset := range row {
			if !other.canInsert(x_offset, y_offset) {
				b.Set(x+x_offset, y+y_offset, charToInsert)
			}
		}
	}
	return nil
}

func (b *Board) Remove(charToRemove string) {
	for y, row := range b.rows {
		for x := range row {
			if b.rows[y][x] == charToRemove {
				b.rows[y][x] = " "
			}
		}
	}
}

type Dog struct {
	orientations []*FilledGrid
}

func NewDog(inputString string) *Dog {
	d := &Dog{}
	orientation := NewFilledGrid(inputString)
	for i := 0; i < 4; i++ {
		d.orientations = append(d.orientations, orientation)
		orientation = orientation.Rotate()
	}
	return d
}

func (d *Dog) Size() int {
	return d.orientations[0].Size()
}

func main() {
	board := NewBoard(`
■     ■■
       ■
        
        
        
        
■       
■■     ■
`)

	dogs := []*Dog{
		NewDog(EIKA),
		NewDog(WANNI),
		NewDog(VIVI),
		NewDog(BORKO),
		NewDog(RUST),
		NewDog(OAKLEY),
		NewDog(MARBLE),
		NewDog(MIKKO),
		NewDog(BORKO),
		NewDog(MUCKI),
	}

	// Sort dogs based on size in descending order
	sortDogsBySize(dogs)

	boardSizeX := len(board.rows[0])
	boardSizeY := len(board.rows)

	dogsTotal := sumDogSizes(dogs)
	boardTotal := boardSizeX*boardSizeY - board.Size()
	fmt.Println("Total dog size:", dogsTotal)
	fmt.Println("Total board size:", boardTotal)
	if dogsTotal != boardTotal {
		fmt.Println("Error: Total dog size and board size do not match.")
		return
	}

	solution, err := solve(dogs, board)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	println("Found solution:")
	fmt.Println(solution)
}

func sumDogSizes(dogs []*Dog) int {
	total := 0
	for _, dog := range dogs {
		total += dog.Size()
	}
	return total
}

func sortDogsBySize(dogs []*Dog) {
	sort.Slice(dogs, func(i, j int) bool {
		return dogs[i].Size() > dogs[j].Size()
	})
}

func solve(dogs []*Dog, board *Board) (*Board, error) {
	if len(dogs) == 0 {
		return board, nil
	}

	if len(dogs) < 3 {
		fmt.Println(board)
	}

	dog := dogs[0]
	otherDogs := dogs[1:]
	for _, dogOrientation := range dog.orientations {
		for x := 0; x < len(board.rows[0]); x++ {
			for y := 0; y < len(board.rows); y++ {
				err := board.Insert(dogOrientation, x, y, fmt.Sprint(len(dogs)-1))
				if err == nil {
					solution, _ := solve(otherDogs, board)
					if solution != nil {
						return solution, nil
					}
					board.Remove(fmt.Sprint(len(dogs)-1)) // Backtrack: remove the dog from the board
				}
			}
		}
	}
	return nil, &DoesNotFitError{}
}