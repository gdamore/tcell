package tcell

import (
	"reflect"
	"testing"

	runewidth "github.com/mattn/go-runewidth"
)

// SetContent sets the contents (primary rune, combining runes,
// and style) for a cell at a given location.  If the background or
// foreground of the style is set to ColorNone, then the respective
// color is left un changed.
func (cb *CellBuffer) SetContentOld(x int, y int,
	mainc rune, combc []rune, style Style,
) {
	if x >= 0 && y >= 0 && x < cb.w && y < cb.h {
		c := &cb.cells[(y*cb.w)+x]

		// Wide characters: we want to mark the "wide" cells
		// dirty as well as the base cell, to make sure we consider
		// both cells as dirty together.  We only need to do this
		// if we're changing content
		if (c.width > 0) && (mainc != c.currMain || len(combc) != len(c.currComb) || (len(combc) > 0 && !reflect.DeepEqual(combc, c.currComb))) {
			for i := 0; i < c.width; i++ {
				cb.SetDirty(x+i, y, true)
			}
		}

		c.currComb = append([]rune{}, combc...)

		if c.currMain != mainc {
			c.width = runewidth.RuneWidth(mainc)
		}
		c.currMain = mainc
		if style.fg == ColorNone {
			style.fg = c.currStyle.fg
		}
		if style.bg == ColorNone {
			style.bg = c.currStyle.bg
		}
		c.currStyle = style
	}
}

func Benchmark_SetContentOld_ascii(b *testing.B) {
	buffer := &CellBuffer{}
	buffer.Resize(100, 100)
	for i := 0; i < b.N; i++ {
		for w := 0; w < 100; w++ {
			for h := 0; h < 100; h++ {
				buffer.SetContentOld(w, h, 'a', nil, StyleDefault)
			}
		}
		for w := 0; w < 100; w++ {
			for h := 0; h < 100; h++ {
				buffer.SetContentOld(w, h, 'b', nil, StyleDefault)
			}
		}
	}
}

func Benchmark_SetContent_ascii(b *testing.B) {
	buffer := &CellBuffer{}
	buffer.Resize(100, 100)
	for i := 0; i < b.N; i++ {
		for w := 0; w < 100; w++ {
			for h := 0; h < 100; h++ {
				buffer.SetContent(w, h, 'a', nil, StyleDefault)
			}
		}
		for w := 0; w < 100; w++ {
			for h := 0; h < 100; h++ {
				buffer.SetContent(w, h, 'b', nil, StyleDefault)
			}
		}
	}
}

func Benchmark_SetContentOld(b *testing.B) {
	buffer := &CellBuffer{}
	buffer.Resize(100, 100)
	flag := []rune("🇦🇺")
	for i := 0; i < b.N; i++ {
		for w := 0; w < 100; w++ {
			for h := 0; h < 100; h++ {
				buffer.SetContentOld(w, h, 'a', nil, StyleDefault)
			}
		}
		for w := 0; w < 100; w++ {
			for h := 0; h < 100; h++ {
				buffer.SetContentOld(w, h, flag[0], flag[1:], StyleDefault)
			}
		}
	}
}

func Benchmark_SetContent(b *testing.B) {
	buffer := &CellBuffer{}
	buffer.Resize(100, 100)
	flag := []rune("🇦🇺")
	for i := 0; i < b.N; i++ {
		for w := 0; w < 100; w++ {
			for h := 0; h < 100; h++ {
				buffer.SetContent(w, h, 'a', nil, StyleDefault)
			}
		}
		for w := 0; w < 100; w++ {
			for h := 0; h < 100; h++ {
				buffer.SetContent(w, h, flag[0], flag[1:], StyleDefault)
			}
		}
	}
}
