// mandelterm - a terminal based Mandelbrot set viewer
package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
	"math/cmplx"
	"runtime"
	"sync"
	"time"
)

// Constants
const (
	// On my terminal 80 x 24 screen is 800 x 504 pixels
	// So each character is 10 x 21 pixels
	// aspect ratio of terminal characters
	aspect = 2.1

	// Fraction of the radius we pan on each keypress
	pan = 0.2

	// Factor we zoom in on each keypress
	zoom = 2
)

// Globals
var (
	showHelp = true
	showInfo = true
	center   complex128
	radius   float64
	depth    int
)

// Reset to the start position
func Reset() {
	center = complex(0, 0)
	radius = 2.0
	depth = 256
}

// Print a string
func Print(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

// Printf a string
func Printf(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	Print(x, y, fg, bg, s)
}

// Help
func Help() {
	Print(0, 0, termbox.ColorRed, termbox.ColorWhite, "Terminal Mandlebrot by Nick Craig-Wood")
	Print(0, 1, termbox.ColorBlack, termbox.ColorWhite, "* Arrow keys to Move")
	Print(0, 2, termbox.ColorBlack, termbox.ColorWhite, "* PgUp/PdDown/+/- to zoom")
	Print(0, 3, termbox.ColorBlack, termbox.ColorWhite, "* [/] to change depth")
	Print(0, 4, termbox.ColorBlack, termbox.ColorWhite, "* h toggle help on and off")
	Print(0, 5, termbox.ColorBlack, termbox.ColorWhite, "* i toggle info on and off")
	Print(0, 6, termbox.ColorBlack, termbox.ColorWhite, "* q/ESC/c-C to quit")
	Print(0, 7, termbox.ColorBlack, termbox.ColorWhite, "* r to reset to start")
}

// Info
func Info(dt time.Duration) {
	_, h := termbox.Size()
	Printf(0, h-1, termbox.ColorBlack, termbox.ColorWhite, "c = %g, r = %g, Depth %d, rendered in %s", center, radius, depth, dt)
}

// Plot a horizontal line from the mandelbrot set
func calculateLine(fx, fy, dx float64, line []termbox.Attribute, wg *sync.WaitGroup) {
	defer wg.Done()
	for x := range line {
		z := complex(0, 0)
		c := complex(fx, fy)
		var i int
		for i = 0; i < depth; i++ {
			if cmplx.Abs(z) >= 2 {
				break
			}
			z = z*z + c
		}
		line[x] = termbox.Attribute(i%8) + 1
		fx += dx
	}
}

// Draw the mandelbrot set
func Draw() {
	start := time.Now()
	w, h := termbox.Size()

	// Choose shortest direction for radius
	var dx, dy float64
	if float64(h) > float64(w)/aspect {
		dx = 2 * radius / float64(w)
		dy = 2 * radius / float64(w) * aspect
	} else {
		dx = 2 * radius / float64(h) / aspect
		dy = 2 * radius / float64(h)
	}

	// Calculate mandelbrot into screen
	screen := make([][]termbox.Attribute, h)
	fy := imag(center) + dy*float64(-h/2)
	var wg sync.WaitGroup
	for y := range screen {
		fx := real(center) + dx*float64(-w/2)
		screen[y] = make([]termbox.Attribute, w)
		wg.Add(1)
		go calculateLine(fx, fy, dx, screen[y], &wg)
		fy += dy

	}
	wg.Wait()

	// Plot
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y := range screen {
		for x, col := range screen[y] {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, col)
		}
	}
	if showInfo {
		Info(time.Now().Sub(start))
	}
	if showHelp {
		Help()
	}
	termbox.Flush()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()
	Reset()
	Draw()
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Key + termbox.Key(ev.Ch) {
			case termbox.KeyEsc, termbox.KeyCtrlC, 'q':
				return
			case termbox.KeyArrowUp:
				center += complex(0.0, -radius*pan)
			case termbox.KeyArrowDown:
				center += complex(0.0, radius*pan)
			case termbox.KeyArrowLeft:
				center += complex(-radius*pan, 0.0)
			case termbox.KeyArrowRight:
				center += complex(radius*pan, 0.0)
			case termbox.KeyPgup, '=', '+':
				radius /= zoom
			case termbox.KeyPgdn, '-', '_':
				radius *= zoom
			case ']':
				depth *= 2
			case '[':
				depth /= 2
				if depth < 64 {
					depth = 64
				}
			case 'h':
				showHelp = !showHelp
			case 'i':
				showInfo = !showInfo
			case 'r':
				Reset()
			}
		}
		Draw()
	}
}
