// -*- tab-width: 4; -*-

// Package main for the shapepuzzle program.
package main

import (
	"fmt"
	"log"
	"runtime"
	"shapepuzzle/board"
	"shapepuzzle/shape"
)

// There are 8 permutations for every shape: flipped over or not, then
// four rotations for each side up.

// getShapes generates the shapes for a specific shape puzzle to be solved.
//
func getShapes() []shape.Shape {
	grids := [][][]int{{
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
	return shape.MakeShapes(grids)
}

type _NullWriter struct {
}

func (nw _NullWriter) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func main() {

	logenable := false

	runtime.GOMAXPROCS(8)
	log.SetFlags(0)

	if !logenable {
		log.SetOutput(_NullWriter{})
	}

	b := board.NewBoard(8, 8)
	shapes := getShapes()

	nshapes := len(shapes)
	for i := 0; i < nshapes; i++ {
		s := &shapes[i]
		perms := s.Permutations()
		for _, p := range perms {
			log.Printf(p.String())
		}
	}

	fmt.Printf("Initial board:\n%v", b)

	bc := b.Solve(shapes)
	nfound := 0
	for b := range bc {
		fmt.Printf("Solution found.\n")
		fmt.Printf("%s\n", b)
		nfound++
	}
	if nfound == 0 {
		fmt.Printf("No solution found.\n")
	} else if nfound == 1 {
		fmt.Printf("One solution found.\n")
	} else {
		fmt.Printf("%d solutions found.\n", nfound)
	}
}
