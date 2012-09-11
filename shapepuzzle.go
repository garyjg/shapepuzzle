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

// Given a set of shapes, a board 8x8, and their permutations and
// translations on the board, find the combination which covers every spot
// on the board.  The board state is represented as a 64-bit mask, and
// shape placements can be converted to a mask on that board to quickly
// determine if that placement is valid, ie, it does not collied any other
// placements.
// 
// To search the solution space, each piece gets its own goroutine in which
// it generates all the available moves on the current board state.  It
// puts each valid move and the resulting board state onto its output
// channel.  The goroutines for each piece can be chained together so that
// if the final piece in the chain puts a new board state on its output
// channel, then that must be a solution to the puzzle.
// 
// As soon as a board state is found on which a piece cannot be placed,
// that board state is discarded and not propagated.  Therefore not all
// pieces will need to be attempted before a placement sequence is
// recognized as a dead end.
// 
// Because the first placement is symmetrical across the four quadrants of
// the board, it only needs to generate placements in one quadrant.  That
// at least reduces the search space by close to a fourth.

// This might be optimized by immediately pruning board states which cannot
// yield solutions because there are gaps which are too small for any shape.
// These gaps can be found by matching against a set of masks:
//
//  A square with one space inside, shifted to all positions on the board.
//  A space on the edge surrounded by filled positions.
//  A rectangle with two spaces inside.
//  A rectangle on the edge.
//  Corner masks which only fit one, two, three, or four spaces.
// 
// [ [ 0, 1, 0 ], [ 1, 0, 1 ], [ 0, 1, 0 ] ]
// 
// The corners do not need to be filled in, only the edges.  The point is
// that empty spaces are surrounded on all sides by either filled spaces or
// the border.
// 
// Basically, given a rectangle with 1, 2, 3, or 4 empty spaces inside, AND
// that filled shape with the board, then OR the mask with the inside
// spaces removed.  If the result is not equal to the filled mask, then
// there are gaps which cannot fit a shape.

// Also, the masks for a shape placement do not change depending upon the
// incoming board, so really the placements can all be computed ahead of
// time and then the masks scanned for one that fits on the incoming board.


type NullWriter struct {
}

func (nw NullWriter) Write(b []byte) (n int, err error) {
    return len(b), nil
}

func main() {

	testboard := true
	logenable := true

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
		s.ComputeMask(b.NumRows(), b.NumCols())
		perms := s.Permutations()
		for _, p := range perms {
			(&p).ComputeMask(b.NumRows(), b.NumCols())
			log.Printf(p.String())
		}
	}

	if testboard {
		tb := b
		log.Printf("Empty board (mask=%v):\n%v", tb.Mask(), tb)
		tb = tb.Place(board.NewPlacement(&shapes[0], &tb, 0, 0))
		tb = tb.Place(board.NewPlacement(&shapes[1], &tb, 1, 1))
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
