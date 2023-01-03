//go:build js && wasm
// +build js,wasm

package tcell

import (
	"errors"
	"strings"
	"sync"
	"syscall/js"
	"unicode/utf8"
)

func NewTerminfoScreen() (Screen, error) {
	t := &tScreen{}
	t.fallback = make(map[rune]string)

	return t, nil
}

type tScreen struct {
	w, h  int
	style Style
	cells CellBuffer

	running      bool
	fini         bool // dummy var; html doesn't need to get restored or "shut down"
	clear        bool
	flagsPresent bool
	pasteEnabled bool
	mouseFlags   MouseFlags

	cursorStyle CursorStyle
	cx, cy      int // dummies so web tScreen can use generic tScreen.Sync

	quit     chan struct{}
	evch     chan Event
	fallback map[rune]string

	sync.Mutex
}

func (t *tScreen) Init() error {
	t.w, t.h = 80, 24 // default for html as of now
	t.evch = make(chan Event, 10)
	t.quit = make(chan struct{})

	t.Lock()
	t.running = true
	t.cx, t.cy = -1, -1
	t.style = StyleDefault
	t.cells.Resize(t.w, t.h)
	t.Unlock()

	js.Global().Set("onKeyEvent", js.FuncOf(t.onKeyEvent))

	return nil
}

func (t *tScreen) Fini() {
	close(t.quit)
}

func (t *tScreen) drawCell(x, y int) int {
	mainc, combc, style, width := t.cells.GetContent(x, y)

	if !t.cells.Dirty(x, y) {
		return width
	}

	if style == StyleDefault {
		style = t.style
	}

	fg, bg := style.fg.Hex(), style.bg.Hex()

	var combcarr []interface{} = make([]interface{}, len(combc))
	for i, c := range combc {
		combcarr[i] = c
	}

	t.cells.SetDirty(x, y, false)
	js.Global().Call("drawCell", x, y, mainc, combcarr, fg, bg, int(style.attrs))

	return width
}

func (t *tScreen) ShowCursor(x, y int) {
	t.Lock()
	js.Global().Call("showCursor", x, y)
	t.Unlock()
}

func (t *tScreen) SetCursorStyle(cs CursorStyle) {
	t.Lock()
	js.Global().Call("setCursorStyle", curStyleClasses[cs])
	t.Unlock()
}

func (t *tScreen) clearScreen() {
	js.Global().Call("clearScreen", t.style.fg.Hex(), t.style.bg.Hex())
	t.clear = false
}

func (t *tScreen) draw() {
	if t.clear {
		t.clearScreen()
	}

	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			width := t.drawCell(x, y)
			x += width - 1
		}
	}

	js.Global().Call("show")
}

func (t *tScreen) enableMouse(f MouseFlags) {
	if f&MouseButtonEvents != 0 {
		js.Global().Set("onMouseClick", js.FuncOf(t.onMouseEvent))
	} else {
		js.Global().Set("onMouseClick", js.FuncOf(t.unset))
	}

	if f&MouseDragEvents != 0 || f&MouseMotionEvents != 0 {
		js.Global().Set("onMouseMove", js.FuncOf(t.onMouseEvent))
	} else {
		js.Global().Set("onMouseMove", js.FuncOf(t.unset))
	}
}

func (t *tScreen) enablePasting(on bool) {
	if on {
		js.Global().Set("onPaste", js.FuncOf(t.onPaste))
	} else {
		js.Global().Set("onPaste", js.FuncOf(t.unset))
	}
}

// resize does nothing, as asking the web window to resize
// without a specified width or height will cause no change.
func (t *tScreen) resize() {}

func (t *tScreen) Colors() int {
	return 16777216 // 256 ^ 3
}

func (t *tScreen) onMouseEvent(this js.Value, args []js.Value) interface{} {
	mod := ModNone
	button := ButtonNone

	switch args[2].Int() {
	case 0:
		if t.mouseFlags&MouseMotionEvents == 0 {
			// don't want this event! is a mouse motion event, but user has asked not.
			return nil
		}
		button = ButtonNone
	case 1:
		button = Button1
	case 2:
		button = Button3 // Note we prefer to treat right as button 2
	case 3:
		button = Button2 // And the middle button as button 3
	}

	if args[3].Bool() { // mod shift
		mod |= ModShift
	}

	if args[4].Bool() { // mod alt
		mod |= ModAlt
	}

	if args[5].Bool() { // mod ctrl
		mod |= ModCtrl
	}

	t.PostEventWait(NewEventMouse(args[0].Int(), args[1].Int(), button, mod))
	return nil
}

func (t *tScreen) onKeyEvent(this js.Value, args []js.Value) interface{} {
	key := args[0].String()

	// don't accept any modifier keys as their own
	if key == "Control" || key == "Alt" || key == "Meta" || key == "Shift" {
		return nil
	}

	mod := ModNone
	if args[1].Bool() { // mod shift
		mod |= ModShift
	}

	if args[2].Bool() { // mod alt
		mod |= ModAlt
	}

	if args[3].Bool() { // mod ctrl
		mod |= ModCtrl
	}

	if args[4].Bool() { // mod meta
		mod |= ModMeta
	}

	// check for special case of Ctrl + key
	if mod == ModCtrl {
		if k, ok := WebKeyNames["Ctrl-"+strings.ToLower(key)]; ok {
			t.PostEventWait(NewEventKey(k, 0, mod))
			return nil
		}
	}

	// next try function keys
	if k, ok := WebKeyNames[key]; ok {
		t.PostEventWait(NewEventKey(k, 0, mod))
		return nil
	}

	// finally try normal, printable chars
	r, _ := utf8.DecodeRuneInString(key)
	t.PostEventWait(NewEventKey(KeyRune, r, mod))
	return nil
}

func (t *tScreen) onPaste(this js.Value, args []js.Value) interface{} {
	t.PostEventWait(NewEventPaste(args[0].Bool()))
	return nil
}

// unset is a dummy function for js when we want nothing to
// happen when javascript calls a function (for example, when
// mouse input is disabled, when onMouseEvent() is called from
// js, it redirects here and does nothing).
func (t *tScreen) unset(this js.Value, args []js.Value) interface{} {
	return nil
}

func (t *tScreen) CharacterSet() string {
	return "UTF-8"
}

func (t *tScreen) CanDisplay(r rune, checkFallbacks bool) bool {
	if utf8.ValidRune(r) {
		return true
	}
	if !checkFallbacks {
		return false
	}
	if _, ok := t.fallback[r]; ok {
		return true
	}
	return false
}

func (t *tScreen) HasMouse() bool {
	return true
}

func (t *tScreen) HasKey(k Key) bool {
	return true
}

func (t *tScreen) SetSize(w, h int) {
	if w == t.w && h == t.h {
		return
	}

	t.cells.Invalidate()
	t.cells.Resize(w, h)
	js.Global().Call("resize", w, h)
	t.w, t.h = w, h
	t.PostEvent(NewEventResize(w, h))
}

// Suspend simply pauses all input and output, and clears the screen.
// There isn't a "default terminal" to go back to.
func (t *tScreen) Suspend() error {
	t.Lock()
	if !t.running {
		t.Unlock()
		return nil
	}
	t.running = false
	t.clearScreen()
	t.enableMouse(0)
	t.enablePasting(false)
	js.Global().Set("onKeyEvent", js.FuncOf(t.unset)) // stop keypresses
	return nil
}

func (t *tScreen) Resume() error {
	t.Lock()

	if t.running {
		return errors.New("already engaged")
	}
	t.running = true

	t.enableMouse(t.mouseFlags)
	t.enablePasting(t.pasteEnabled)

	js.Global().Set("onKeyEvent", js.FuncOf(t.onKeyEvent))

	t.Unlock()
	return nil
}

func (t *tScreen) Beep() error {
	js.Global().Call("beep")
	return nil
}

// WebKeyNames maps string names reported from HTML
// (KeyboardEvent.key) to tcell accepted keys.
var WebKeyNames = map[string]Key{
	"Enter":      KeyEnter,
	"Backspace":  KeyBackspace,
	"Tab":        KeyTab,
	"Backtab":    KeyBacktab,
	"Escape":     KeyEsc,
	"Backspace2": KeyBackspace2,
	"Delete":     KeyDelete,
	"Insert":     KeyInsert,
	"ArrowUp":    KeyUp,
	"ArrowDown":  KeyDown,
	"ArrowLeft":  KeyLeft,
	"ArrowRight": KeyRight,
	"Home":       KeyHome,
	"End":        KeyEnd,
	"UpLeft":     KeyUpLeft,    // not supported by HTML
	"UpRight":    KeyUpRight,   // not supported by HTML
	"DownLeft":   KeyDownLeft,  // not supported by HTML
	"DownRight":  KeyDownRight, // not supported by HTML
	"Center":     KeyCenter,
	"PgDn":       KeyPgDn,
	"PgUp":       KeyPgUp,
	"Clear":      KeyClear,
	"Exit":       KeyExit,
	"Cancel":     KeyCancel,
	"Pause":      KeyPause,
	"Print":      KeyPrint,
	"F1":         KeyF1,
	"F2":         KeyF2,
	"F3":         KeyF3,
	"F4":         KeyF4,
	"F5":         KeyF5,
	"F6":         KeyF6,
	"F7":         KeyF7,
	"F8":         KeyF8,
	"F9":         KeyF9,
	"F10":        KeyF10,
	"F11":        KeyF11,
	"F12":        KeyF12,
	"F13":        KeyF13,
	"F14":        KeyF14,
	"F15":        KeyF15,
	"F16":        KeyF16,
	"F17":        KeyF17,
	"F18":        KeyF18,
	"F19":        KeyF19,
	"F20":        KeyF20,
	"F21":        KeyF21,
	"F22":        KeyF22,
	"F23":        KeyF23,
	"F24":        KeyF24,
	"F25":        KeyF25,
	"F26":        KeyF26,
	"F27":        KeyF27,
	"F28":        KeyF28,
	"F29":        KeyF29,
	"F30":        KeyF30,
	"F31":        KeyF31,
	"F32":        KeyF32,
	"F33":        KeyF33,
	"F34":        KeyF34,
	"F35":        KeyF35,
	"F36":        KeyF36,
	"F37":        KeyF37,
	"F38":        KeyF38,
	"F39":        KeyF39,
	"F40":        KeyF40,
	"F41":        KeyF41,
	"F42":        KeyF42,
	"F43":        KeyF43,
	"F44":        KeyF44,
	"F45":        KeyF45,
	"F46":        KeyF46,
	"F47":        KeyF47,
	"F48":        KeyF48,
	"F49":        KeyF49,
	"F50":        KeyF50,
	"F51":        KeyF51,
	"F52":        KeyF52,
	"F53":        KeyF53,
	"F54":        KeyF54,
	"F55":        KeyF55,
	"F56":        KeyF56,
	"F57":        KeyF57,
	"F58":        KeyF58,
	"F59":        KeyF59,
	"F60":        KeyF60,
	"F61":        KeyF61,
	"F62":        KeyF62,
	"F63":        KeyF63,
	"F64":        KeyF64,
	"Ctrl-a":     KeyCtrlA,          // not reported by HTML- need to do special check
	"Ctrl-b":     KeyCtrlB,          // not reported by HTML- need to do special check
	"Ctrl-c":     KeyCtrlC,          // not reported by HTML- need to do special check
	"Ctrl-d":     KeyCtrlD,          // not reported by HTML- need to do special check
	"Ctrl-e":     KeyCtrlE,          // not reported by HTML- need to do special check
	"Ctrl-f":     KeyCtrlF,          // not reported by HTML- need to do special check
	"Ctrl-g":     KeyCtrlG,          // not reported by HTML- need to do special check
	"Ctrl-j":     KeyCtrlJ,          // not reported by HTML- need to do special check
	"Ctrl-k":     KeyCtrlK,          // not reported by HTML- need to do special check
	"Ctrl-l":     KeyCtrlL,          // not reported by HTML- need to do special check
	"Ctrl-n":     KeyCtrlN,          // not reported by HTML- need to do special check
	"Ctrl-o":     KeyCtrlO,          // not reported by HTML- need to do special check
	"Ctrl-p":     KeyCtrlP,          // not reported by HTML- need to do special check
	"Ctrl-q":     KeyCtrlQ,          // not reported by HTML- need to do special check
	"Ctrl-r":     KeyCtrlR,          // not reported by HTML- need to do special check
	"Ctrl-s":     KeyCtrlS,          // not reported by HTML- need to do special check
	"Ctrl-t":     KeyCtrlT,          // not reported by HTML- need to do special check
	"Ctrl-u":     KeyCtrlU,          // not reported by HTML- need to do special check
	"Ctrl-v":     KeyCtrlV,          // not reported by HTML- need to do special check
	"Ctrl-w":     KeyCtrlW,          // not reported by HTML- need to do special check
	"Ctrl-x":     KeyCtrlX,          // not reported by HTML- need to do special check
	"Ctrl-y":     KeyCtrlY,          // not reported by HTML- need to do special check
	"Ctrl-z":     KeyCtrlZ,          // not reported by HTML- need to do special check
	"Ctrl- ":     KeyCtrlSpace,      // not reported by HTML- need to do special check
	"Ctrl-_":     KeyCtrlUnderscore, // not reported by HTML- need to do special check
	"Ctrl-]":     KeyCtrlRightSq,    // not reported by HTML- need to do special check
	"Ctrl-\\":    KeyCtrlBackslash,  // not reported by HTML- need to do special check
	"Ctrl-^":     KeyCtrlCarat,      // not reported by HTML- need to do special check
}

var curStyleClasses = map[CursorStyle]string{
	CursorStyleDefault:           "cursor-blinking-block",
	CursorStyleBlinkingBlock:     "cursor-blinking-block",
	CursorStyleSteadyBlock:       "cursor-steady-block",
	CursorStyleBlinkingUnderline: "cursor-blinking-underline",
	CursorStyleSteadyUnderline:   "cursor-steady-underline",
	CursorStyleBlinkingBar:       "cursor-blinking-bar",
	CursorStyleSteadyBar:         "cursor-steady-bar",
}
