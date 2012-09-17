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
	return shape.MakeShapes(grids)
}

type NullWriter struct {
}

func (nw NullWriter) Write(b []byte) (n int, err error) {
    return len(b), nil
}

func main() {

	logenable := true

	runtime.GOMAXPROCS(8)
	log.SetFlags(0)
	
	if ! logenable {
		log.SetOutput(NullWriter{})
	}

	b := board.NewBoard(8, 8)
	shapes := getShapes()

	nshapes := len(shapes)
	for i := 0; i < nshapes; i += 1 {
		s := &shapes[i]
		perms := s.Permutations()
		for _, p := range perms {
			log.Printf(p.String())
		}
	}

	fmt.Printf("Initial board:\n%v", b)

	b, ok := b.Solve(shapes)
	if ok {
		fmt.Println("Solution found.")
		fmt.Println(b)
	} else {
		fmt.Println("No solution found.")
	}

}
