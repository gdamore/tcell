# Layouts

This directory contains packages which implement keyboard layouts for the emulator
subsystem.

The emulator subsystem is designed to support writing your own terminal emulators,
including the simulation (mock) terminal emulator used to test TCell.

If you're not implementing a terminal emulator, you don't need this.

## Designing a New Layout

The layout subsystem is intended to be extensible, where layouts can inherit from
other layouts, so that you only need to implement the differences.

A good place to start is by looking at the us/ package, and the US International Layout
in particular (which extends on the vanilla ANSI layout.  The ANSI layout is always available.)

## Finding Layout Information

A good source for information about keyboard layouts is [this site](https://kbdlayout.info).
Pay particular attention to shift states for a given layout, and the dead keys.

## Limitations

As of this writing, support for Ligature Keys is absent.
