## tcell

[![Linux Status](https://img.shields.io/travis/gdamore/tcell.svg?label=linux)](https://travis-ci.org/gdamore/tcell)
[![Windows Status](https://img.shields.io/appveyor/ci/gdamore/tcell.svg?label=windows)](https://ci.appveyor.com/project/gdamore/tcell)
[![GitHub License](https://img.shields.io/github/license/gdamore/tcell.svg)](https://github.com/gdamore/tcell/blob/master/LICENSE)
[![Issues](https://img.shields.io/github/issues/gdamore/tcell.svg)](https://github.com/gdamore/tcell/issues)
[![Gitter](https://img.shields.io/badge/gitter-join-brightgreen.svg)](https://gitter.im/gdamore/tcell)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/gdamore/tcell)

> _Tcell is a work in progress (Beta).
> Please use with caution; interfaces may change in before final release._

Package tcell provides a cell based view for text terminals, like xterm.
It was inspired by termbox, but differs from termbox in some important
ways.

First, it includes a full parser and expander for terminfo capability strings,
so that it can avoid hard coding escape strings for formatting.  It also favors
portability, and includes support for all POSIX systems, at the slight expense
of needing cgo support for terminal initializations.  (This will be corrected
when Go provides standard support for terminal handling via termio ioctls on
all POSIX platforms.)

Also, this code is able to operate without requiring
SIGIO signals, or asynchronous I/O, and can instead use standard Go file
objects and Go routines.

It also includes enhanced support for Unicode, include wide characters and
combining characters.  It also has richer support for a larger number of
special keys that some terminals can send.

It will respect your terminal's color space as specified within your terminfo
entries, so that for example attempts to emit color sequences on VT100 terminals
won't result in unintended consequences.

I started this project originally by submitting patches to the author of
go-termbox, but due to some fundamental differences of opinion, I thought
it might be simpler just to start from scratch.

## Termbox compatibility 

A compatibility layer for termbox is provided in the compat
directory.  To use it, try importing "github.com/gdamore/tcell/termbox"
instead.  Most termbox-go programs will probably work without further
modification.

## Working with Unicode

This version of the tcells expects that your terminal can support Unicode
on output.  That is, if you submit Unicode sequences to it, it will attempt
send Unicode to the terminal.  This works for modern xterm and other emulators,
but legacy systems may have poor results.  I'm interested to hear reports from
folks who need support for other character sets.

## Wide & Combining Characters

The Setcell() API takes a sequence of runes; exactly least one of them should
be a non-zero width.  Combining runes may follow.  If any of the runes
is a wide (East Asian) rune occupying two cells, then the library will skip
output from the following cell, but care must be taken in the application to
avoid explicitly attempting to set content in the next cell, otherwise the
results are undefined.  (Normally the wide character will not be displayed.)

## Colors

We assume the ANSI/XTerm color model, including the 256 color map that
XTerm uses when it supports 256 colors.  The terminfo guidance will be
honored, with respect to the number of colors supported.  Also, only
terminals which expose ANSI style setaf and setab will support color;
if you have a color terminal that only has setf and setb, please let me
know; it wouldn't be hard to add that if there is need.

## Performance

Reasonable attempts have been made to minimize sending data to terminals,
avoiding repeated sequences or drawing the same cell on refresh updates.

## Terminfo

The Terminfo implementation operates with two forms of database.  The first
is the database.go file, which contains a number of real database entries
that are compiled into the program directly.  This should minimize calling
out to database file searches.

The second is a JSON file, that contains the same information, which can
be located either by the $TCELLDB environment file, or is located in the
Go source directory.

These files (both the Go database.go and the database.json) file can be
generated using the mkinfo.go program.  If you need to regnerate the
entire set for some reason, run the mkdatabase.sh file.  The generation
uses the terminfo routines on the system to populate the data files.

The mkinfo.go program can also be used to generate specific database entries
for named terminals, in case your favorite terminal is missing.  (If you
find that this is the case, please let me know and I'll try to add it!)

Tcell requires that the terminal support the 'cup' mode of cursor addressing.
Terminals without absolute cursor addressability are not supported.
This is unlikely to be a problem; such terminals have not been mass produced
since the early 1970s.

## Mouse Support

Mouse support is detected via the "kmous" terminfo variable, however,
enablement/disablement and decoding mouse events is done using hard coded
sequences based on the XTerm X11 model.  As of this writing all popular
terminals with mouse tracking support this model.  (Full terminfo support
is not possible as terminfo sequences are not defined.)

## Platforms

This system requires a POSIX termios implementation with /dev/tty to run.
It also requires functional cgo to run.  As of this writing, Cgo is available
on all POSIX Go 1.5 platforms.

Windows console support is forthcoming, cygwin should work now.
