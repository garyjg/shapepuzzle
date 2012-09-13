// -*- tab-width: 4; -*-

package main

import (
	"fmt"
	"shapepuzzle/shape"
	"shapepuzzle/board"
	"log"
	"runtime"
)


// There are 8 permutations for every shape: flipped over or not, then
// four rotations for each side up.

func getShapes() []shape.Shape {
	grids := [][][]int { {
			{1, 1, 0}, {1, 1, 1}}, {
			{1, 0, 1}, {1, 1, 1}}, {
			{1, 0, 0, 0, 0}, {1, 1, 1, 1, 1}}, {
			{1, 1, 1, 1}, {1, 0, 0, 1}}, {
			{1, 1, 1}, {1, 1, 1}, {0, 1, 1}}, {
			{0, 1, 0}, {1, 1, 1}, {0, 1, 0}}, {
			{0, 1, 0}, {0, 1, 0}, {1, 1, 1}}, {
			{0, 0, 1, 1}, {1, 1, 1, 1}}, {
			{0, 1, 1}, {1, 1, 0}, {1, 0, 0}}, {
			{1, 0, 0, 0}, {1, 0, 0, 0}, {1, 0, 0, 0}, {1, 1, 1, 1}}, {
			{1, 0, 0, 0}, {1, 1, 1, 1}, {1, 0, 0, 0}}}
	shapes := []shape.Shape{}
	for id, grid := range grids {
		s := shape.NewShape(id+1, grid)
		shapes = append(shapes, s)
	}
	return shapes
}

func testShapes() []shape.Shape {
	grids := [][][]int { {
			{1, 1, 1}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0}}, {
			{1, 1, 0}, {1, 1, 1}}, {
			{1, 1, 1}, {0, 1, 0}}, {
			{0,0,1,1}, {1,1,1,0}}, {
			{1, 0, 1}, {1, 1, 1}}}
	shapes := []shape.Shape{}
	for id, grid := range grids {
		s := shape.NewShape(id+1, grid)
		shapes = append(shapes, s)
	}
	return shapes
}

// To search the solution space, each piece gets its own goroutine in which
// it generates all the available moves on the current board state.  It
// puts each valid move and the resulting board state onto its output
// channel.  The goroutines for each piece can be chained together so that
// if the final piece in the chain puts a new board state on its output
// channel, then that must be a solution to the puzzle.
// 


type NullWriter struct {
}

func (nw NullWriter) Write(b []byte) (n int, err error) {
    return len(b), nil
}

func main() {

	testboard := false
	logenable := false

	runtime.GOMAXPROCS(8)
	log.SetFlags(0)
	
	if ! logenable {
		log.SetOutput(NullWriter{})
	}

	b := board.NewBoard(8, 8)
	shapes := getShapes()
	if testboard {
		b = board.NewBoard(5, 5)
		shapes = testShapes()
	}

	nshapes := len(shapes)
	for i := 0; i < nshapes; i += 1 {
		s := &shapes[i]
		perms := s.Permutations()
		for _, p := range perms {
			log.Printf(p.String())
		}
	}

	if testboard {
		tb := b
		log.Printf("Empty board (mask=%v):\n%v", tb.Mask(), tb)
		tb = tb.Place(board.NewPlacement(shapes[0], 0, 0))
		tb = tb.Place(board.NewPlacement(shapes[1], 1, 1))
		log.Println(tb)
	}

	fmt.Printf("Initial board:\n%v", b)

	// Set up a channel for each shape to be placed.
	channels := make ([]board.BoardChannel, nshapes)
	for i := 0; i < nshapes; i +=1 {
		channels[i] = make(board.BoardChannel, 10000)
	}
	
	// Chain the channels.  Generate first placements for the first shape,
	// and tell it to put those new boards on its channel.

	go board.FirstPlacements(shapes[0], b, channels[0])

	if false {
	for nextboard := range channels[0] {
		fmt.Println(nextboard)
	} 
	return
	}

	for i := 1; i < len(shapes); i += 1 {
		go board.NextPlacements(shapes[i], b, channels[i-1], channels[i])
    }

	if false {
	for nextboard := range channels[1] {
		fmt.Println(nextboard)
	} 
	return
	}

	// Finally listen for a solution (or not) to be pushed to the last
	// channel.

	b, ok := <-channels[nshapes-1]
	if ok {
		fmt.Println("Solution found.")
		fmt.Println(b)
	} else {
		fmt.Println("No solution found.")
	}

}
