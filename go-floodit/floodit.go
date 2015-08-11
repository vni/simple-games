/* TODO: Finish Legend */
/* TODO: FIX score accounting if using the same color */
/* TODO: Clean the code */

package main

import "github.com/nsf/termbox-go"
import "fmt"
import "time"
import "math/rand"

var grid [25][25]termbox.Attribute
var nmoves int

func init() {
	rand.Seed(time.Now().UnixNano())

	for x:=0; x<25; x++ {
		for y:=0; y<25; y++ {
			grid[x][y] = termbox.ColorRed + termbox.Attribute(rand.Int() % 6)
		}
	}
}

func draw_cell(x, y int, color termbox.Attribute) {
	termbox.SetCell(4*x+0, 2*y, ' ', termbox.ColorDefault, color)
	termbox.SetCell(4*x+1, 2*y, ' ', termbox.ColorDefault, color)
	termbox.SetCell(4*x+2, 2*y, ' ', termbox.ColorDefault, color)
	termbox.SetCell(4*x+3, 2*y, ' ', termbox.ColorDefault, color)

	termbox.SetCell(4*x+0, 2*y+1, ' ', termbox.ColorDefault, color)
	termbox.SetCell(4*x+1, 2*y+1, ' ', termbox.ColorDefault, color)
	termbox.SetCell(4*x+2, 2*y+1, ' ', termbox.ColorDefault, color)
	termbox.SetCell(4*x+3, 2*y+1, ' ', termbox.ColorDefault, color)
}

func draw() {
	for x:=0; x<len(grid); x++ {
		for y:=0; y<len(grid[x]); y++ {
			draw_cell(x, y, grid[x][y])
		}
	}

	draw_cell(26, 1, termbox.ColorRed+termbox.Attribute(0))
	draw_cell(28, 1, termbox.ColorRed+termbox.Attribute(1))
	draw_cell(30, 1, termbox.ColorRed+termbox.Attribute(2))
	draw_cell(26, 3, termbox.ColorRed+termbox.Attribute(3))
	draw_cell(28, 3, termbox.ColorRed+termbox.Attribute(4))
	draw_cell(30, 3, termbox.ColorRed+termbox.Attribute(5))

	termbox.Flush()
}

func flood(x, y int, orig, n termbox.Attribute) {
	if x < 0 || x >= 25 || y < 0 || y >= 25 || grid[x][y] != orig {
		return
	}
	grid[x][y] = n
	flood(x+1, y, orig, n)
	flood(x, y+1, orig, n)
}

func finished() bool {
	for x:=0; x<len(grid); x++ {
		for y:=0; y<len(grid); y++ {
			if grid[0][0] != grid[x][y] {
				return false
			}
		}
	}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	tb_print(2, 1, termbox.ColorDefault, termbox.ColorDefault, "CONGRATULATIONS! YOU WIN!")
	tb_printf(2, 2, termbox.ColorDefault, termbox.ColorDefault, "with %d moves", nmoves)
	termbox.Flush()
	time.Sleep(3*time.Second)

	return true
}

func tb_print(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func tb_printf(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	tb_print(x, y, fg, bg, s)
}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err);
	}
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	loop: for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		draw()
		termbox.Flush()

		select {
			case ev := <-event_queue:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc {
					break loop
				} else if ev.Ch == 'u' {
					flood(0, 0, grid[0][0], termbox.ColorRed)
				} else if ev.Ch == 'i' {
					flood(0, 0, grid[0][0], termbox.ColorGreen)
				} else if ev.Ch == 'o' {
					flood(0, 0, grid[0][0], termbox.ColorYellow)
				} else if ev.Ch == 'j' {
					flood(0, 0, grid[0][0], termbox.ColorBlue)
				} else if ev.Ch == 'k' {
					flood(0, 0, grid[0][0], termbox.ColorMagenta)
				} else if ev.Ch == 'l' {
					flood(0, 0, grid[0][0], termbox.ColorCyan)
				}

				if finished() {
					break loop
				}

				nmoves++
			}
		}
	}
}
