package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_gfx"
)

type clock struct {
	minute float64
	hour   float64
}

type grid [][]clock

const (
	clockRadius = 32
	gutter      = 6
	frameRate   = 60
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

	go func() {
		for range time.Tick(time.Second) {
			for dx := range *g {
				for dy := range (*g)[dx] {
					now := time.Now()
					(*g)[dx][dy].minute = (float64(now.Second()) / 60) * 2 * math.Pi
				}
			}
		}
	}()

	return loop(r, g)
}

func loop(r *sdl.Renderer, g *grid) error {
	for {
		r.SetDrawColor(0, 0, 0, 255)
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
		sdl.Delay(1000 / frameRate)
	}
}

func (c *clock) draw(r *sdl.Renderer, dx, dy int) error {
	x := dx*2*(clockRadius+gutter) + clockRadius*2
	y := dy*2*(clockRadius+gutter) + clockRadius*2 - gutter
	color := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	if result := gfx.AACircleColor(r, x, y, clockRadius, color); result == false {
		return fmt.Errorf("could not draw circle: %v", sdl.GetError())
	}

	mx := x + int(clockRadius*0.8*math.Cos(3*math.Pi/2+c.minute)-0.5)
	my := y + int(clockRadius*0.8*math.Sin(3*math.Pi/2+c.minute)-0.5)
	if result := gfx.AALineColor(r, x, y, mx, my, color); result == false {
		return fmt.Errorf("could not draw minute hand: %v", sdl.GetError())
	}

	hx := x + int(clockRadius*0.6*math.Cos(3*math.Pi/2+c.hour)-0.5)
	hy := y + int(clockRadius*0.6*math.Sin(3*math.Pi/2+c.hour)-0.5)
	if result := gfx.AALineColor(r, x, y, hx, hy, color); result == false {
		return fmt.Errorf("could not draw minute hand: %v", sdl.GetError())
	}

	return nil
}
