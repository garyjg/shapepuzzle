# Shapepuzzle

Shapepuzzle is a [Go](http://golang.org) program for solving a puzzle where
interlocking 2D shapes must be placed to fill a square board.

## Data model

The shapes are defined as arrays of bitmasks which can be rotated and flipped
and placed on the board by OR'ing the bits.  Each byte is one row on the
board, and each bit is a column.  So any board size up to 8x8 can fit into an
unsigned long integer, but the board size for this particular puzzle is 8x8.

For example, this grid:

```go
        {1, 1, 0}, {1, 1, 1}
```

describes a shape like this:

<table style='width: 120p; height: 80p;'>
<tr><td bgcolor='black'/><td bgcolor='black'/><td bgcolor='white'/></tr>
<tr><td bgcolor='black'/><td bgcolor='black'/><td bgcolor='black'/></tr>
</table>

## Algorithm

The placement of each piece is a step in the solution search space and runs
in its own goroutine.  The first goroutine places the first piece onto the
open spots on the board for each of the piece's permutations, then pushes
that board to the next goroutine.  If a piece cannot be placed successfully
on a board, or if the board contains gaps into which no other piece can be
placed, then the board is rejected.  Otherwise the new board state with the
new shape placement is pushed to the succeeding placement goroutine.  The
very last goroutine places the last piece, and if a placement succeeds the
board is pushed to the last channel as a solution.

When all the shapes have been placed on the board without colliding with
other shape placements, the puzzle is solved.

This division of the search space is obviously not optimal, since the work
done by each goroutine, ie, the space searched by each goroutine, gets
smaller as more pieces are placed on the board.

Also, as pieces are placed, there are fewer spots to fit more pieces, so it
might be possible to limit the shape permutations which are tried by first
looking for the open spots.

Since the larger pieces generate fewer possible boards, placing the largest
(or most irregular?) pieces first should reduce the search spaces passed to
the subsequent goroutines.
