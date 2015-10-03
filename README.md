## tcell

[![Linux Status](https://img.shields.io/travis/gdamore/tcell.svg?label=linux)](https://travis-ci.org/gdamore/tcell)
[![Windows Status](https://img.shields.io/appveyor/ci/gdamore/tcell.svg?label=windows)](https://ci.appveyor.com/project/gdamore/tcell)
[![GitHub License](https://img.shields.io/github/license/gdamore/tcell.svg)](https://github.com/gdamore/tcell/blob/master/LICENSE)
[![Issues](https://img.shields.io/github/issues/gdamore/tcell.svg)](https://github.com/gdamore/tcell/issues)
[![Gitter](https://img.shields.io/badge/gitter-join-brightgreen.svg)](https://gitter.im/gdamore/tcell)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/gdamore/tcell)

> _Tcell is a work in progress (Beta).
> Please use with caution; interfaces may change in before final release.
> That said, our confidence in Tcell's stability is increasing.  If you
> would like to use it in your own application, it is recommended that
> you drop a message to garrett@damore.org before commitment._

Package tcell provides a cell based view for text terminals, like xterm.
It was inspired by termbox, but differs from termbox in some important
ways.  It also adds signficant functionality beyond termbox.

## Pure Go Terminfo Database

First, it includes a full parser and expander for terminfo capability strings,
so that it can avoid hard coding escape strings for formatting.  It also favors
portability, and includes support for all POSIX systems, at the slight expense
of needing cgo support for terminal initializations.  (This may be corrected
when Go provides standard support for terminal handling via termio ioctls on
all POSIX platforms.)  The database itself, while built using CGO, as well
as the parser for it, is implemented in Pure Go.

The database is also flexible & extensibel, and can modified by either running a
program to build the database, or hand editing of simple JSON files.

## More Portable

Tcell is portable to a wider variety of systems.  It relies on standard
POSIX supported function calls (on POSIX platforms) for setting terminal
modes, which leads to improved support for a broader array of platforms.
This does come at the cost of requiring your code to be able to use CGO, but
we believe that the vastly improved portability justifies this
requirement.  Note that the functions called are part of the standard C
library, so there shouldn't be any additional external requirements beyond
that required for every POSIX program.

## No async IO

Termbox code is able to operate without requiring
SIGIO signals, or asynchronous I/O, and can instead use standard Go file
objects and Go routines.  This means it should be safe, especially for
use with programs that use exec, or otherwise need to manipulate the
tty streams.  This model is also much closer to idiomatic Go, leading
to fewer surprises.

## Richer Unicode support

Tcell includes enhanced support for Unicode, include wide characters and
combining characters, provided your terminal can support them.  Note that
Windows terminals generally don't support the full Unicode repertoire.

## More Function Keys

It also has richer support for a larger number of
special keys that some terminals can send.

## Better color handling

 Tcell will respect your terminal's color space as specified within your terminfo
entries, so that for example attempts to emit color sequences on VT100 terminals
won't result in unintended consequences.

In Windows mode, we support 16 colors, underline, bold, dim, and reverse,
instead of just termbox's 8 colors with reverse.  (Note that there is some
conflation with bold/dim and colors.)

## Better mouse support

It supports enhanced mouse tracking mode, so your application can receive
regular mouse motion events, and wheel events, if your terminal supports it.

## Why not just patch termbox-go?

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

Experience has shown that the vanilla Windows 8 console application does not
support any of these characters properly, but at least some options like
ConEmu do support Wide characters at least.  Combining characters are
disabled for Windows in the release.

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

Windows still needs some work here, as it can make numerous system calls
a bit less efficiently than it could.

## Terminfo

(Not relevent for Windows users.)

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

On Windows, the mouse works normally.

Mouse wheel buttons on various terminals are known to work, but the support
in terminal emulators, as well as support for various buttons and
live mouse tracking, varies widely.  As a particular datum, MacOS X Terminal
does not support Mouse events at all (as of MacOS 10.10, aka Yosemite.)  The
excellent iTerm application is fully supported, as is vanilla XTerm.

Mouse tracking with live tracking also varies widely.  Current XTerm
implementations, as well as current Screen and iTerm2, and Windows
consoles, all support this quite nicely.  On other platforms you might
find that only mouse click and release events are reported, with
no intervening motion events.  It really depends on your terminal.

## Platforms

On POSIX systems, a POSIX termios implementation with /dev/tty is required.
It also requires functional cgo to run.  As of this writing, Cgo is available
on all POSIX Go 1.5 platforms.

Windows console mode applications are supported.  Unfortunately mintty
and other cygwin style applications are not supported.

Modern console applications like ConEmu support all the good features
(resize, mouse tracking, etc.)

I haven't figured out how to cleanly resolve the dichotomy between cygwin
style termios and the Windows Console API; it seems that perhaps nobody else
has either.  If anyone has suggestions, let me know!  Really, if you're
using a Windows application, you should use the native Windows console or a
fully compatible consule implementation.  We expect that Windows 10
ships with a less crippled implementation than prior releases -- we haven't
tested that, lacking Windows 10 ourselves.

