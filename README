
Shapepuzzle is a Go program for solving a puzzle where interlocking shapes
must be placed to fill a square board.  The shapes are defined as 2D
bitmask arrays which can be rotated and flipped and placed on the board by
OR'ing the bits.  When all the shapes have been placed without colliding
with other shape placements, the puzzle is solved.

The placement of each piece is a step in the solution search space and runs
in its own goroutine.  The first goroutine places the first piece onto the
open spots on the board for each of the piece's permutations, then pushes
that board to the next goroutine.  If a piece cannot be placed successfully
on a board, or if the board contains gaps into which no other piece can be
placed, then the board is rejected.  Otherwise the new board state with the
new shape placement is pushed to the succeeding placement goroutine.  The
very last goroutine places the last piece, and if a placement succeeds the
board is pushed to the last channel as a solution.

