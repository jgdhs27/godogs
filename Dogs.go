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

func (fg *FilledGrid) IsSquareOccupied(x, y int) bool {
	return fg.rows[y][x] != " "
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

func (b *Board) Insert(other *FilledGrid, x, y int, charToInsert string) (*Board, error) {
	newState := NewBoard("\n")
	newState.rows = make([][]string, len(b.rows))
	for i, row := range b.rows {
		newState.rows[i] = make([]string, len(row))
		copy(newState.rows[i], row)
	}

	for y_offset, row := range other.rows {
		for x_offset, _ := range row {
			if other.IsSquareOccupied(x_offset, y_offset) {
				if b.IsSquareOccupied(x+x_offset, y+y_offset) {
					return nil, &DoesNotFitError{}
				}
				newState.Set(x+x_offset, y+y_offset, charToInsert)
			}
		}
	}
	return newState, nil
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

	const (
		VIVI = `
X
`
ROCCO = `
X
X
`
YASKA = `
XXX
`
EIKA = `
 X
XX
`
GARIBALDI = `
XX
XX
`
MARBLE = `
XXX
 X 
`
OAKLEY = `
X 
XX
 X
`
STRAWBERRY = `
X 
X 
XX
`
MOE = `
 X
 X
XX
`
BELLE = `
XX 
XXX
`
SIR_ALFIE = `
X 
XX
XX
`
ZOE = `
XXX
X X
`
MUCKI = `
XX   
 XXXX
`
MIKKO = `
  XX
XXXX
`
BORKO = `
 X 
XXX
XX 
`
POLKA = `
XX 
XX 
XXX
`
ELLIE = `
 XX
 XX
XXX
`
DUKE = `
XX
 X
 X
 X
XX
`
LULU = `
XX 
XXX
XXX
`
ABBY = `
  XX
  XX
XXXX
`
REX = `
XXX
XXX
XXX
`
WANNI = `
XX  
XXXX
XXXX
`
RUST = `
XXX
 XX
 XX
XXX
`
KAFKA = `
X     
X     
XX    
XX    
XXXX  
XXXXXX
`
MAX = `
XX  X
XX  X
XXXXX
XXXXX
XX XX
`
GOLDIE = `
   X
XXXX
XXXX
   X
`
ROMY = `
XXXXXXXX
XXXXXXXX
  X   XX
      XX
`
KORRA = `
XX  
XX  
XXX 
XXXX
XXXX
`
	)

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
				func() {
					defer func() {
						if err := recover(); err != nil {
							// "slice bounds out of range" error is caught here
						}
					}()

					newBoard, err := board.Insert(dogOrientation, x, y, fmt.Sprint(len(dogs)-1))
					if err == nil {
						solution, _ := solve(otherDogs, newBoard)
						if solution != nil {
							panic("solution found")
						}
					}
				}()
			}
		}
	}
	return nil, &DoesNotFitError{}
}