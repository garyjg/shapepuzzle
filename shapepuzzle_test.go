// -*- tab-width: 4; -*-

package main

import (
	"fmt"
	"log"
	"shapepuzzle/board"
	"shapepuzzle/shape"
	"testing"
)

func testShapes() []shape.Shape {
	grids := [][][]int{{
		{1, 1, 1}, {1, 0, 0}, {1, 0, 0}, {1, 0, 0}}, {
		{1, 1, 0}, {1, 1, 1}}, {
		{1, 1, 1}, {0, 1, 0}}, {
		{0, 0, 1, 1}, {1, 1, 1, 0}}, {
		{1, 0, 1}, {1, 1, 1}}}
	shapes := []shape.Shape{}
	for id, grid := range grids {
		s := shape.NewShape(id+1, grid)
		shapes = append(shapes, s)
	}
	return shapes
}

func TestSolution(t *testing.T) {
	b := board.NewBoard(5, 5)
	shapes := testShapes()

	bc := b.Solve(shapes)
	b, ok := <-bc
	if !ok {
		t.Errorf("No solution found!")
	}
	fmt.Println(b)
}

func TestPlacements(t *testing.T) {

	tb := board.NewBoard(5, 5)
	shapes := testShapes()
	log.Printf("Empty board (mask=%v):\n%v", tb.Mask(), tb)
	tb = tb.Place(shapes[0].Translate(0, 0))
	tb = tb.Place(shapes[1].Translate(1, 1))
	log.Println(tb)
}
