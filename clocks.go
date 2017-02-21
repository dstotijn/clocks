package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	gfx "github.com/veandco/go-sdl2/sdl_gfx"
)

type clock struct {
	minute int
	hour   int
}

type grid [][]clock

const (
	clockRadius = 32
	gutter      = 6
)

func main() {
	g := make(grid, 16)
	for i := range g {
		g[i] = make([]clock, 9)
	}
	if err := run(&g); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
}

func run(g *grid) error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(1280, 720, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()
	defer r.Destroy()

	if err := r.Clear(); err != nil {
		return fmt.Errorf("could not clear renderer: %v", err)
	}

	for dx := range *g {
		for dy := range (*g)[dx] {
			if err := (*g)[dx][dy].draw(r, dx, dy); err != nil {
				return fmt.Errorf("could not draw clock: %v", err)
			}
		}
	}

	r.Present()
	time.Sleep(time.Second * 5)

	return nil
}

func (c *clock) draw(r *sdl.Renderer, dx, dy int) error {
	x := dx*2*(clockRadius+gutter) + clockRadius*2
	y := dy*2*(clockRadius+gutter) + clockRadius*2 - gutter
	color := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	if result := gfx.CircleColor(r, x, y, clockRadius, color); result == false {
		return fmt.Errorf("could not draw circle: %v", sdl.GetError())
	}

	return nil
}
