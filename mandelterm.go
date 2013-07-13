// mandelterm - a terminal based Mandelbrot set viewer
package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
	"math/cmplx"
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
	Print(5, 4, termbox.ColorRed, termbox.ColorWhite, "Terminal Mandlebrot by Nick Craig-Wood")
	Print(5, 6, termbox.ColorBlack, termbox.ColorWhite, "* Arrow keys to Move")
	Print(5, 7, termbox.ColorBlack, termbox.ColorWhite, "* PgUp/PdDown/+/- to zoom")
	Print(5, 8, termbox.ColorBlack, termbox.ColorWhite, "* [/] to change depth")
	Print(5, 9, termbox.ColorBlack, termbox.ColorWhite, "* h toggle help on and off")
	Print(5, 10, termbox.ColorBlack, termbox.ColorWhite, "* i toggle info on and off")
	Print(5, 11, termbox.ColorBlack, termbox.ColorWhite, "* q/ESC/c-C to quit")
	Print(5, 12, termbox.ColorBlack, termbox.ColorWhite, "* r to reset to start")
}

// Info
func Info(dt time.Duration) {
	Printf(0, 0, termbox.ColorBlack, termbox.ColorWhite, "c = %g, r = %g, Depth %d, rendered in %s", center, radius, depth, dt)
}

// Draw the mandelbrot set
func Draw() {
	start := time.Now()
	w, h := termbox.Size()
	// Choose shortest direction for reference
	var dx, dy float64
	if float64(h) > float64(w)/aspect {
		dx = 2 * radius / float64(w)
		dy = 2 * radius / float64(w) * aspect
	} else {
		dx = 2 * radius / float64(h) / aspect
		dy = 2 * radius / float64(h)
	}

	// FIXME aspect ratio
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	fy := imag(center) + dy*float64(-h/2)
	for y := 0; y < h; y++ {
		fx := real(center) + dx*float64(-w/2)
		for x := 0; x < w; x++ {
			z := complex(0, 0)
			c := complex(fx, fy)
			var i int
			for i = 0; i < depth; i++ {
				if cmplx.Abs(z) >= 2 {
					break
				}
				z = z*z + c
			}
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.Attribute(i%8)+1)
			fx += dx
		}
		fy += dy
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
