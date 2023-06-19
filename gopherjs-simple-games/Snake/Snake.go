// FIXME: add hunger and length lost feature
// FIXME: add text output feature ('yam-yam', 'what a tasty banana!')
// FIXME: monomial (5 banana in a row) score bonus
// FIXME: ticks should be every 20 ms. (for a better keyboard response)

package main

import (
	"strconv"
	"github.com/gopherjs/gopherjs/js"
)

const (
	CELL_SIZE = 30 // px
	BASE_TICK_DURATION = 250 // ms
	MAX_LAST_FRUITS_LEN = 5
)

var (
	g_width, g_height int
	g_width_px, g_height_px int
	g_canvas *Canvas
	g_document *js.Object
	g_math *js.Object
	g_snake *Snake
	g_level int
	g_fruits [3]*Fruit
	g_game_over bool
	g_paused bool
	g_game_over_status string
	g_last_fruits = make([]FruitType, 0, MAX_LAST_FRUITS_LEN+1)
	g_last_fruits_count int
)

//
// Canvas
//

type Canvas struct {
	obj *js.Object
	ctx *js.Object
	//width, height int
}

func NewCanvasById(id string) *Canvas {
	canvas := &Canvas{}
	canvas.obj = g_document.Call("getElementById", id)
	canvas.ctx = canvas.obj.Call("getContext", "2d")
	return canvas
}

func (c *Canvas) fillRect(x, y, w, h int) *Canvas {
	c.ctx.Call("fillRect", x, y, w, h)
	return c
}

func (c *Canvas) strokeRect(x, y, w, h int) *Canvas {
	c.ctx.Call("strokeRect", x, y, w, h)
	return c
}

func (c *Canvas) strokeStyle(style string) *Canvas {
	c.ctx.Set("strokeStyle", style)
	return c
}

func (c *Canvas) fillStyle(style string) *Canvas {
	c.ctx.Set("fillStyle", style)
	return c
}

func (c *Canvas) setWidth(w int) *Canvas {
	c.obj.Set("width", w)
	return c
}

func (c *Canvas) setHeight(h int) *Canvas {
	c.obj.Set("height", h)
	return c
}

//
// Cell
//

type Cell struct {
	x, y int
}

func (c *Cell) draw() {
	g_canvas.fillRect(c.x*CELL_SIZE+1, c.y*CELL_SIZE+1, CELL_SIZE-2, CELL_SIZE-2)
}

//
// Fruit
//

type FruitType int

const (
	Apple FruitType = iota
	Banana
	Orange
	Watermelon
	Melon
	LAST_FRUIT
)

func (f FruitType) Color() string {
	switch (f) {
	case Apple:
		return "green"
	case Banana:
		return "yellow"
	case Orange:
		return "orange"
	case Watermelon:
		return "pink"
	case Melon:
		return "lightyellow"
	default:
		return "<UNKNOWN COLOR>"
	}
}

func (f FruitType) String() string {
	switch (f) {
	case Apple:
		return "apple"
	case Banana:
		return "banana"
	case Orange:
		return "orange"
	case Watermelon:
		return "watermelon"
	case Melon:
		return "melon"
	default:
		return "<UNKNOWN COLOR>"
	}
}

type Fruit struct {
	x, y int
	duration int // ticks
	t FruitType // type
}

func GenerateNewFruit() *Fruit {
	f := &Fruit{}

	f.t = FruitType(random(int(LAST_FRUIT)))
	f.duration = 100

	for {
		f.x = random(g_width)
		f.y = random(g_height)
		invalid := false
		for _,c := range g_snake.body {
			if f.x == c.x && f.y == c.y {
				invalid = true
				break
			}
		}
		if !invalid {
			break
		}
	}

	return f
}

func (f *Fruit) draw() {
	g_canvas.fillStyle(f.t.Color())
	g_canvas.fillRect(f.x*CELL_SIZE+1, f.y*CELL_SIZE+1, CELL_SIZE-2, CELL_SIZE-2)
}

//
// Snake
//

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

type Snake struct {
	body []Cell
	score int
	direction Direction
	fruits_eaten int
}

func NewSnake() *Snake {
	s := &Snake{}
	s.body = make([]Cell, 3, 100)
	s.body[0] = Cell{g_width / 2 - 1, g_height / 2}
	s.body[1] = Cell{g_width / 2, g_height / 2}
	s.body[2] = Cell{g_width / 2 + 1, g_height / 2}
	s.score = 0
	s.direction = EAST
	return s
}

func (s *Snake) draw() {
	g_canvas.fillStyle("lightblue")
	for _, b := range s.body {
		b.draw()
	}
	g_canvas.fillStyle("blue")
	s.body[len(s.body)-1].draw()
}

func (s *Snake) head() Cell {
	return s.body[len(s.body)-1]
}

func (s *Snake) move() {
	h := s.head()
	switch (s.direction) {
	case NORTH:
		h.y--
	case EAST:
		h.x++
	case SOUTH:
		h.y++
	case WEST:
		h.x--
	}

	for h.x < 0 {
		h.x += g_width
	}
	for h.x >= g_width {
		h.x -= g_width
	}

	for h.y < 0 {
		h.y += g_height
	}
	for h.y >= g_height {
		h.y -= g_height
	}

	for i, v := range g_snake.body {
		for ii, vv := range g_snake.body {
			if i != ii {
				if v.x == vv.x && v.y == vv.y {
					g_game_over = true
					g_game_over_status = "<b>GAME OVER<br />" +
						"YOU LOSE<br />" +
						"The snake is not allowed to eat itself.</b>"
					return
				}
			}
		}
	}

	expand := 1

	for i, f := range g_fruits {
		if f.x == h.x && f.y == h.y {
			g_snake.score += 10
			g_snake.fruits_eaten += 1

			g_last_fruits_count++
			g_last_fruits = append(g_last_fruits, f.t)
			if len(g_last_fruits) >= MAX_LAST_FRUITS_LEN {
				copy(g_last_fruits, g_last_fruits[1:6])
				g_last_fruits = g_last_fruits[0:5]
			}

			g_fruits[i] = GenerateNewFruit()
			expand = 0
		}
	}

	s.body = append(s.body[expand:], h)
}

func (s *Snake) turnCW() {
	if s.direction == WEST {
		s.direction = NORTH
	} else {
		s.direction++
	}
}

func (s *Snake) turnCCW() {
	if s.direction == NORTH {
		s.direction = WEST
	} else {
		s.direction--
	}
}

//
//
//

func onKeyPress(e *js.Object) {
	if g_game_over {
		return
	}

	ck := e.Get("charCode").Int()
	//println("onKeyPress: ck:", ck)

	switch (ck) {
	case int(' '):
		g_paused = !g_paused
		if g_paused {
			println("PAUSED")
		} else {
			tick_duration := BASE_TICK_DURATION - g_snake.score / 10;
			js.Global.Call("setTimeout", oneIteration, tick_duration)
			println("UNPAUSED")
		}
	case int('j'):
		if !g_paused {
			g_snake.turnCCW()
			println("turn CCW")
		}
	case int('l'):
		if !g_paused {
			g_snake.turnCW()
			println("turn CW")
		}
	}
}

func init() {
	println("init() started")
	g_document = js.Global.Get("document")
	g_math = js.Global.Get("Math")

	g_width_px = 1280
	g_height_px = 1024

	g_width = g_width_px / CELL_SIZE
	g_height = g_height_px / CELL_SIZE

	// correct g_{width,height}_px
	g_width_px = g_width * CELL_SIZE
	g_height_px = g_height * CELL_SIZE

	g_canvas = NewCanvasById("Canvas")
	g_canvas.setWidth(g_width_px)
	g_canvas.setHeight(g_height_px)

	g_snake = NewSnake()

	for i := 0; i < len(g_fruits); i++ {
		g_fruits[i] = GenerateNewFruit()
	}

	g_paused = true

	body := g_document.Call("getElementsByTagName", "body").Index(0)
	body.Call("addEventListener", "keypress", onKeyPress, true)

	println("init() finished")
}

//
//
//

func random(n int) int {
	return int(g_math.Call("random").Float() * float64(n))
}

func draw_board() {
	for y := 0; y < g_height; y++ {
		for x := 0; x < g_width; x++ {
			g_canvas.strokeRect(x*CELL_SIZE, y*CELL_SIZE, CELL_SIZE, CELL_SIZE)
		}
	}
}

func draw_all() {
	g_canvas.strokeStyle("rgba(242, 198, 65, 0.1)")
	g_canvas.fillStyle("rgb(38, 38, 38)")

	g_canvas.fillRect(0, 0, g_width_px, g_height_px)

	draw_board()
	g_snake.draw()

	for _, f := range g_fruits {
		f.draw()
	}
}

func show_game_status() {
	score := g_document.Call("getElementById", "Score")

	if g_game_over {
		score.Set("innerHTML", g_game_over_status)
		return
	}


	score_string := "<b>Score</b>: " + strconv.Itoa(g_snake.score)

	pause_status := "<b>UNPAUSED</b>"
	if g_paused {
		pause_status = "<b>PAUSED</b>"
	}

	tick_duration_status := "<b>Tick duration</b>: " +
		strconv.Itoa(calculate_tick_duration())

	all := score_string + "<br />" + pause_status + "<br />" + tick_duration_status

	score.Set("innerHTML", all)
}

func show_last_fruits() {
	all := ""
	for _, f := range g_last_fruits {
		all += "<span class=\"" + f.String() + "\">" + f.String() + "</span>   "
	}

	fruits := g_document.Call("getElementById", "Fruits")
	fruits.Set("innerHTML", all)
}

func oneIteration() {
	draw_all()
	if !g_paused && !g_game_over {
		g_snake.move()
		js.Global.Call("setTimeout", oneIteration, calculate_tick_duration())
	}
	show_game_status()
	if !g_game_over {
		show_last_fruits()
	}
}

func calculate_tick_duration() (tick_duration int) {
	tick_duration = BASE_TICK_DURATION - g_snake.fruits_eaten * 5
	if tick_duration < 50 {
		tick_duration = 50
	}
	return tick_duration
}

func main() {
	js.Global.Call("setTimeout", oneIteration, BASE_TICK_DURATION)
}
