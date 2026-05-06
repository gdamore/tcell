// benchmarks for SimulationScreen. run with:
//
//	go test -run '^$' -bench '^BenchmarkSimulation' -benchmem .
//
// to attribute allocations:
//
//	go test -run '^$' -bench '^BenchmarkSimulationShow_FullDirty$' \
//	    -benchmem -memprofile=mem.pprof .
//	go tool pprof -alloc_objects -top -cum mem.pprof

package tcell

import (
	"fmt"
	"testing"
)

// fillScreen writes a deterministic pattern into every cell. used by
// the benchmarks to mark every cell dirty before the next Show.
func fillScreen(s SimulationScreen, w, h int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := rune('a' + ((x + y) % 26))
			s.SetContent(x, y, r, nil, StyleDefault)
		}
	}
}

// BenchmarkSimulationShow_FullDirty exercises Show with every cell
// dirty since the last Show, so drawCell runs for every cell. surfaces
// the per-cell allocs in drawCell -- atm append([]rune{mainc}, combc...)
// plus the two make([]byte, 12) encode buffers.
func BenchmarkSimulationShow_FullDirty(b *testing.B) {
	for _, dim := range []struct{ w, h int }{
		{80, 24},
		{200, 60},
	} {
		name := fmt.Sprintf("%dx%d", dim.w, dim.h)
		b.Run(name, func(b *testing.B) {
			s := NewSimulationScreen("UTF-8")
			if err := s.Init(); err != nil {
				b.Fatalf("init: %v", err)
			}
			defer s.Fini()
			s.SetSize(dim.w, dim.h)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; b.Loop(); i++ {
				// rewrite every cell so they're all dirty for this Show.
				fillScreen(s, dim.w, dim.h)
				s.Show()
				_ = i
			}
		})
	}
}

// BenchmarkSimulationShow_Clean exercises Show when nothing's changed
// since the previous Show. drawCell short-circuits via the !Dirty check,
// so this isolates the per-cell loop overhead from the per-cell allocs.
func BenchmarkSimulationShow_Clean(b *testing.B) {
	for _, dim := range []struct{ w, h int }{
		{80, 24},
		{200, 60},
	} {
		name := fmt.Sprintf("%dx%d", dim.w, dim.h)
		b.Run(name, func(b *testing.B) {
			s := NewSimulationScreen("UTF-8")
			if err := s.Init(); err != nil {
				b.Fatalf("init: %v", err)
			}
			defer s.Fini()
			s.SetSize(dim.w, dim.h)
			fillScreen(s, dim.w, dim.h)
			s.Show() // prime: clear all dirty bits
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				s.Show()
			}
		})
	}
}

// BenchmarkSimulationClear exercises the Clear+Show path. Clear marks
// every cell dirty so Show then redraws all cells. mirrors the per-frame
// pattern in apps that issue a Clear at the start of each render tick.
func BenchmarkSimulationClear(b *testing.B) {
	for _, dim := range []struct{ w, h int }{
		{80, 24},
		{200, 60},
	} {
		name := fmt.Sprintf("%dx%d", dim.w, dim.h)
		b.Run(name, func(b *testing.B) {
			s := NewSimulationScreen("UTF-8")
			if err := s.Init(); err != nil {
				b.Fatalf("init: %v", err)
			}
			defer s.Fini()
			s.SetSize(dim.w, dim.h)
			fillScreen(s, dim.w, dim.h)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				s.Clear()
				s.Show()
			}
		})
	}
}

// BenchmarkSimulationSetContent isolates the SetContent path itself
// (no Show) so changes there can be tracked independently of drawCell.
func BenchmarkSimulationSetContent(b *testing.B) {
	const w, h = 200, 60
	s := NewSimulationScreen("UTF-8")
	if err := s.Init(); err != nil {
		b.Fatalf("init: %v", err)
	}
	defer s.Fini()
	s.SetSize(w, h)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		x := i % w
		y := (i / w) % h
		s.SetContent(x, y, 'x', nil, StyleDefault)
	}
}
