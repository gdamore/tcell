## Breaking Changes in _Tcell_ v3

> [!NOTE]
> _Tcell_ v3 is currently in development, and these details are subject to change
> before we release.

### Termbox Compatibility Removed

The `termbox` compatibility package is removed. Few applications were using it,
and the compatibility was imperfect. Also the package had limited support for many
newer features. Further, _Termbox_ itself is no longer being maintained.
Applications that still need this should keep using _Tcell_ v2.

### Cell and Contents APIs

In order to improve support for multi-rune grapheme clusters, and to provide an
experience that reduces friction when using it, some APIs have been removed, and
newer APIs exist in their place.

- `SetCell` is removed.  Use `Put` instead.
- `SetContents` is deprecated and may be removed before release.  Use `Put` instead.
- `GetContents` is removed. Use `Get` instead.

### Support for Grapheme Clusters in EventKey

`EventKey` now carries a string for `KeyRune` instead of a single rune.
As a result the old `Rune` method for `EventKey` is replaced by `Str`.
The main difference for most users will be that `Str` returns a string, and most
of the time that string will consist of only a single rune. However, it is possible
now to inject synthetic key strokes consisting of multi-rune grapheme clusters.

### Terminfo Redesign

The Terminfo subsystem is being replaced entirely.  This is currently still a work in progress.
Essentially the old terminfo based design has long proved to be inferior for modern terminal
applications, and has not kept up with newer terminal features such as 24-bit color,
different mouse reporting modes, bracketed paste, advanced text styling, and so forth.

As part of this, we're removing the parsed terminfo logic entirely.  Most of the terminal
descriptions are being consolidated into a set that resemble either _xterm_ or legacy Digital
VT100 and successors. There may be some outliers remaining like _aixterm_ and certain FOSS
consoles, but we hope even those will be consolidated into a some basic few based on ECMA-48.

A consequence of this is that the Terminfo libraries and descriptions are subject to removal entirely.
Do not depend on Terminfo.

A further consequence of this is that support for some legacy terminals that are either functionally
extinct (such as _hpterm_) or unlikely to be found outside of a museum (such as VT52, Wyse50, or
anything produced more than 40 years ago.)

Note that VT100 and later will work in emulation, and VT220 and later physical terminals should still work. 
VT100 physical terminals may or may not work, as we are removing the special padding delays
that hurt emulations that do not need them, and existed only to accommodate limitations found on the
physical hardware from the 1970s.

Note that we still examine `$TERM` when appropriate, but if the value is not one we recognize,
then we will assume something reasonably capable and compatible at some level with _xterm_ or
at least ECMA-48.

### Color Bit Size

The `Color` type is now only 32-bits, which should save some memory on large terminal windows.

### Underline

`AttrUnderline` is gone.  It was not sufficient to describe styled and colored underlines.

### Blocking PostEvent

The `PostEventWait` function is gone. Use `PostEvent` instead, and check the error status
to determine if the event queue was too full.  (Unless your application is buggy or you
are overwhelming it with events, the event queue should never fill up.)

### Removed Capability Queries

Deprecated APIs `HasKey` and `CanDisplay` are removed.
These functions weren't reliable and served no useful purpose.

### Windows Console API

`NewConsoleScreen` is removed as is support for Windows console mode.

Instead this uses the more modern Windows VT modes.
As a consequence, this means that _Tcell_ on Windows requires at least Winows 10 build 1703 (the Creators Update).
If you are using a version of Windows 10 older than that, you should really upgrade for _many_ reasons, not just
because _Tcell_ doesn't support it anymore.
