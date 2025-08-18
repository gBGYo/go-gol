package terminal

import (
	"fmt"
	"math/rand"
	"time"
)

const ROWS int64 = 20
const COLS int64 = 40

type Grid [ROWS][COLS]int64

func mod(a, b int64) (c int64) {
	c = ((a % b) + b) % b
	return
}

func (g *Grid) randInit(threshold float64) {
	for y := range ROWS {
		for x := range COLS {
			if rand.Float64() < threshold {
				g[y][x] = 1
			} else {
				g[y][x] = 0
			}
		}
	}
}

func (g Grid) display() {
	for y := range ROWS {
		for x := range COLS {
			if g[y][x] == 1 {
				fmt.Print("\u25CF")
			} else {
				// fmt.Print("\u25CB")
				fmt.Print("\u25CC")
			}
		}
		fmt.Println()
	}
}

func (g Grid) update() Grid {
	var newGrid Grid
	for y := range ROWS {
		for x := range COLS {
			var neighbors int64

			neighbors += g[mod(y-1, ROWS)][mod(x-1, COLS)]
			neighbors += g[mod(y-1, ROWS)][mod(x, COLS)]
			neighbors += g[mod(y-1, ROWS)][mod(x+1, COLS)]

			neighbors += g[mod(y, ROWS)][mod(x-1, COLS)]
			neighbors += g[mod(y, ROWS)][mod(x+1, COLS)]

			neighbors += g[mod(y+1, ROWS)][mod(x-1, COLS)]
			neighbors += g[mod(y+1, ROWS)][mod(x, COLS)]
			neighbors += g[mod(y+1, ROWS)][mod(x+1, COLS)]

			// fmt.Print(neighbors)

			if neighbors < 2 || neighbors > 3 {
				newGrid[y][x] = 0
			}
			if (neighbors == 2 || neighbors == 3) && g[y][x] == 1 {
				newGrid[y][x] = 1
			}
			if neighbors == 3 && g[y][x] == 0 {
				newGrid[y][x] = 1
			}
		}
		// fmt.Println()
	}
	return newGrid
}

func Run() {
	var grid Grid

	// Random Grid
	grid.randInit(0.2)

	// Glider
	// offset := 5
	// grid[0+offset][1+offset] = 1
	// grid[1+offset][2+offset] = 1
	// grid[2+offset][0+offset] = 1
	// grid[2+offset][1+offset] = 1
	// grid[2+offset][2+offset] = 1

	for {
		fmt.Print("\033[H")
		grid.display()
		grid = grid.update()
		time.Sleep(100 * time.Millisecond)
	}
}
