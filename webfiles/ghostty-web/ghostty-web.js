var m = /* @__PURE__ */ ((E) => (E[E.CURSOR_KEY_APPLICATION = 0] = "CURSOR_KEY_APPLICATION", E[E.KEYPAD_KEY_APPLICATION = 1] = "KEYPAD_KEY_APPLICATION", E[E.IGNORE_KEYPAD_WITH_NUMLOCK = 2] = "IGNORE_KEYPAD_WITH_NUMLOCK", E[E.ALT_ESC_PREFIX = 3] = "ALT_ESC_PREFIX", E[E.MODIFY_OTHER_KEYS_STATE_2 = 4] = "MODIFY_OTHER_KEYS_STATE_2", E[E.KITTY_KEYBOARD_FLAGS = 5] = "KITTY_KEYBOARD_FLAGS", E))(m || {}), b = /* @__PURE__ */ ((E) => (E[E.RELEASE = 0] = "RELEASE", E[E.PRESS = 1] = "PRESS", E[E.REPEAT = 2] = "REPEAT", E))(b || {}), D = /* @__PURE__ */ ((E) => (E[E.A = 4] = "A", E[E.B = 5] = "B", E[E.C = 6] = "C", E[E.D = 7] = "D", E[E.E = 8] = "E", E[E.F = 9] = "F", E[E.G = 10] = "G", E[E.H = 11] = "H", E[E.I = 12] = "I", E[E.J = 13] = "J", E[E.K = 14] = "K", E[E.L = 15] = "L", E[E.M = 16] = "M", E[E.N = 17] = "N", E[E.O = 18] = "O", E[E.P = 19] = "P", E[E.Q = 20] = "Q", E[E.R = 21] = "R", E[E.S = 22] = "S", E[E.T = 23] = "T", E[E.U = 24] = "U", E[E.V = 25] = "V", E[E.W = 26] = "W", E[E.X = 27] = "X", E[E.Y = 28] = "Y", E[E.Z = 29] = "Z", E[E.ONE = 30] = "ONE", E[E.TWO = 31] = "TWO", E[E.THREE = 32] = "THREE", E[E.FOUR = 33] = "FOUR", E[E.FIVE = 34] = "FIVE", E[E.SIX = 35] = "SIX", E[E.SEVEN = 36] = "SEVEN", E[E.EIGHT = 37] = "EIGHT", E[E.NINE = 38] = "NINE", E[E.ZERO = 39] = "ZERO", E[E.ENTER = 40] = "ENTER", E[E.ESCAPE = 41] = "ESCAPE", E[E.BACKSPACE = 42] = "BACKSPACE", E[E.TAB = 43] = "TAB", E[E.SPACE = 44] = "SPACE", E[E.MINUS = 45] = "MINUS", E[E.EQUAL = 46] = "EQUAL", E[E.BRACKET_LEFT = 47] = "BRACKET_LEFT", E[E.BRACKET_RIGHT = 48] = "BRACKET_RIGHT", E[E.BACKSLASH = 49] = "BACKSLASH", E[E.SEMICOLON = 51] = "SEMICOLON", E[E.QUOTE = 52] = "QUOTE", E[E.GRAVE = 53] = "GRAVE", E[E.COMMA = 54] = "COMMA", E[E.PERIOD = 55] = "PERIOD", E[E.SLASH = 56] = "SLASH", E[E.CAPS_LOCK = 57] = "CAPS_LOCK", E[E.F1 = 58] = "F1", E[E.F2 = 59] = "F2", E[E.F3 = 60] = "F3", E[E.F4 = 61] = "F4", E[E.F5 = 62] = "F5", E[E.F6 = 63] = "F6", E[E.F7 = 64] = "F7", E[E.F8 = 65] = "F8", E[E.F9 = 66] = "F9", E[E.F10 = 67] = "F10", E[E.F11 = 68] = "F11", E[E.F12 = 69] = "F12", E[E.PRINT_SCREEN = 70] = "PRINT_SCREEN", E[E.SCROLL_LOCK = 71] = "SCROLL_LOCK", E[E.PAUSE = 72] = "PAUSE", E[E.INSERT = 73] = "INSERT", E[E.HOME = 74] = "HOME", E[E.PAGE_UP = 75] = "PAGE_UP", E[E.DELETE = 76] = "DELETE", E[E.END = 77] = "END", E[E.PAGE_DOWN = 78] = "PAGE_DOWN", E[E.RIGHT = 79] = "RIGHT", E[E.LEFT = 80] = "LEFT", E[E.DOWN = 81] = "DOWN", E[E.UP = 82] = "UP", E[E.NUM_LOCK = 83] = "NUM_LOCK", E[E.KP_DIVIDE = 84] = "KP_DIVIDE", E[E.KP_MULTIPLY = 85] = "KP_MULTIPLY", E[E.KP_MINUS = 86] = "KP_MINUS", E[E.KP_PLUS = 87] = "KP_PLUS", E[E.KP_ENTER = 88] = "KP_ENTER", E[E.KP_1 = 89] = "KP_1", E[E.KP_2 = 90] = "KP_2", E[E.KP_3 = 91] = "KP_3", E[E.KP_4 = 92] = "KP_4", E[E.KP_5 = 93] = "KP_5", E[E.KP_6 = 94] = "KP_6", E[E.KP_7 = 95] = "KP_7", E[E.KP_8 = 96] = "KP_8", E[E.KP_9 = 97] = "KP_9", E[E.KP_0 = 98] = "KP_0", E[E.KP_PERIOD = 99] = "KP_PERIOD", E[E.NON_US_BACKSLASH = 100] = "NON_US_BACKSLASH", E[E.APPLICATION = 101] = "APPLICATION", E[E.F13 = 104] = "F13", E[E.F14 = 105] = "F14", E[E.F15 = 106] = "F15", E[E.F16 = 107] = "F16", E[E.F17 = 108] = "F17", E[E.F18 = 109] = "F18", E[E.F19 = 110] = "F19", E[E.F20 = 111] = "F20", E[E.F21 = 112] = "F21", E[E.F22 = 113] = "F22", E[E.F23 = 114] = "F23", E[E.F24 = 115] = "F24", E))(D || {}), q = /* @__PURE__ */ ((E) => (E[E.NONE = 0] = "NONE", E[E.SHIFT = 1] = "SHIFT", E[E.CTRL = 2] = "CTRL", E[E.ALT = 4] = "ALT", E[E.SUPER = 8] = "SUPER", E[E.CAPSLOCK = 16] = "CAPSLOCK", E[E.NUMLOCK = 32] = "NUMLOCK", E))(q || {});
const d = 80;
var U = /* @__PURE__ */ ((E) => (E[E.BOLD = 1] = "BOLD", E[E.ITALIC = 2] = "ITALIC", E[E.UNDERLINE = 4] = "UNDERLINE", E[E.STRIKETHROUGH = 8] = "STRIKETHROUGH", E[E.INVERSE = 16] = "INVERSE", E[E.INVISIBLE = 32] = "INVISIBLE", E[E.BLINK = 64] = "BLINK", E[E.FAINT = 128] = "FAINT", E))(U || {});
class x {
  constructor(A) {
    this.exports = A.exports, this.memory = this.exports.memory;
  }
  /**
   * Get current memory buffer (may change when memory grows)
   */
  getBuffer() {
    return this.memory.buffer;
  }
  /**
   * Create a key encoder instance
   */
  createKeyEncoder() {
    return new z(this.exports);
  }
  /**
   * Create a terminal emulator instance
   */
  createTerminal(A = 80, B = 24, Q) {
    return new j(this.exports, this.memory, A, B, Q);
  }
  /**
   * Load Ghostty WASM from URL or file path
   * If no path is provided, attempts to load from common default locations
   */
  static async load(A) {
    const B = [
      // Inline base64 wasm defaults removed (tcell de-inline): the wasm ships as the
      // separate ghostty-vt.wasm and is loaded by URL. Callers should pass an explicit
      // URL to Ghostty.load(); the relative paths below remain as no-argument fallbacks.
      // When used from CDN or local dev
      "./ghostty-vt.wasm",
      "/ghostty-vt.wasm"
    ], Q = A ? [A] : B;
    let g = null;
    for (const C of Q)
      try {
        let I;
        const i = await fetch(C);
        if (!i.ok)
          throw new Error(`Failed to fetch WASM: ${i.status} ${i.statusText}`);
        if (I = await i.arrayBuffer(), I.byteLength === 0)
          throw new Error(`WASM file is empty (0 bytes). Check path: ${C}`);
        const w = await WebAssembly.instantiate(I, {
          env: {
            log: (o, i) => {
              const M = w.instance, k = new Uint8Array(M.exports.memory.buffer, o, i), G = new TextDecoder().decode(k);
              console.log("[ghostty-wasm]", G);
            }
          }
        });
        return new x(w.instance);
      } catch (I) {
        g = I instanceof Error ? I : new Error(String(I));
      }
    throw new Error(
      `Failed to load ghostty-vt.wasm. Tried paths: ${Q.join(", ")}. Last error: ${g == null ? void 0 : g.message}. You can specify a custom path with: new Terminal({ wasmPath: './path/to/ghostty-vt.wasm' })`
    );
  }
}
class z {
  constructor(A) {
    this.encoder = 0, this.exports = A;
    const B = this.exports.ghostty_wasm_alloc_opaque(), Q = this.exports.ghostty_key_encoder_new(0, B);
    if (Q !== 0)
      throw new Error(`Failed to create key encoder: ${Q}`);
    const g = new DataView(this.exports.memory.buffer);
    this.encoder = g.getUint32(B, !0), this.exports.ghostty_wasm_free_opaque(B);
  }
  /**
   * Set an encoder option
   */
  setOption(A, B) {
    const Q = this.exports.ghostty_wasm_alloc_u8(), g = new DataView(this.exports.memory.buffer);
    typeof B == "boolean" ? g.setUint8(Q, B ? 1 : 0) : g.setUint8(Q, B);
    const C = this.exports.ghostty_key_encoder_setopt(this.encoder, A, Q);
    if (this.exports.ghostty_wasm_free_u8(Q), C !== void 0 && C !== 0)
      throw new Error(`Failed to set encoder option: ${C}`);
  }
  /**
   * Enable Kitty keyboard protocol with specified flags
   */
  setKittyFlags(A) {
    this.setOption(m.KITTY_KEYBOARD_FLAGS, A);
  }
  /**
   * Encode a key event to escape sequence
   */
  encode(A) {
    const B = this.exports.ghostty_wasm_alloc_opaque(), Q = this.exports.ghostty_key_event_new(0, B);
    if (Q !== 0)
      throw new Error(`Failed to create key event: ${Q}`);
    const g = new DataView(this.exports.memory.buffer), C = g.getUint32(B, !0);
    if (this.exports.ghostty_wasm_free_opaque(B), this.exports.ghostty_key_event_set_action(C, A.action), this.exports.ghostty_key_event_set_key(C, A.key), this.exports.ghostty_key_event_set_mods(C, A.mods), A.utf8) {
      const J = new TextEncoder().encode(A.utf8), s = this.exports.ghostty_wasm_alloc_u8_array(J.length);
      new Uint8Array(this.exports.memory.buffer).set(J, s), this.exports.ghostty_key_event_set_utf8(C, s, J.length), this.exports.ghostty_wasm_free_u8_array(s, J.length);
    }
    const I = 32, w = this.exports.ghostty_wasm_alloc_u8_array(I), o = this.exports.ghostty_wasm_alloc_usize(), i = this.exports.ghostty_key_encoder_encode(
      this.encoder,
      C,
      w,
      I,
      o
    );
    if (i !== 0)
      throw this.exports.ghostty_wasm_free_u8_array(w, I), this.exports.ghostty_wasm_free_usize(o), this.exports.ghostty_key_event_free(C), new Error(`Failed to encode key: ${i}`);
    const M = g.getUint32(o, !0), k = new Uint8Array(this.exports.memory.buffer, w, M).slice();
    return this.exports.ghostty_wasm_free_u8_array(w, I), this.exports.ghostty_wasm_free_usize(o), this.exports.ghostty_key_event_free(C), k;
  }
  /**
   * Free encoder resources
   */
  dispose() {
    this.encoder && (this.exports.ghostty_key_encoder_free(this.encoder), this.encoder = 0);
  }
}
const f = class Y {
  /**
   * Create a new terminal.
   *
   * @param exports WASM exports
   * @param memory WASM memory
   * @param cols Number of columns (default: 80)
   * @param rows Number of rows (default: 24)
   * @param config Optional terminal configuration (colors, scrollback)
   * @throws Error if allocation fails
   */
  constructor(A, B, Q = 80, g = 24, C) {
    var w;
    this.exports = A, this.memory = B, this._cols = Q, this._rows = g;
    let I;
    if (C) {
      const o = this.exports.ghostty_wasm_alloc_u8_array(d);
      if (o === 0)
        throw new Error("Failed to allocate config (out of memory)");
      try {
        const i = new DataView(this.memory.buffer);
        let M = o;
        i.setUint32(M, C.scrollbackLimit ?? 1e4, !0), M += 4, i.setUint32(M, C.fgColor ?? 0, !0), M += 4, i.setUint32(M, C.bgColor ?? 0, !0), M += 4, i.setUint32(M, C.cursorColor ?? 0, !0), M += 4;
        for (let k = 0; k < 16; k++) {
          const G = ((w = C.palette) == null ? void 0 : w[k]) ?? 0;
          i.setUint32(M, G, !0), M += 4;
        }
        I = this.exports.ghostty_terminal_new_with_config(Q, g, o);
      } finally {
        this.exports.ghostty_wasm_free_u8_array(o, d);
      }
    } else
      I = this.exports.ghostty_terminal_new(Q, g);
    if (I === 0)
      throw new Error("Failed to allocate terminal (out of memory)");
    this.handle = I;
  }
  /**
   * Free the terminal. Must be called to prevent memory leaks.
   */
  free() {
    this.handle !== 0 && (this.exports.ghostty_terminal_free(this.handle), this.handle = 0);
  }
  /**
   * Write data to terminal (parses VT sequences and updates screen).
   *
   * @param data UTF-8 string or Uint8Array
   *
   * @example
   * ```typescript
   * term.write('Hello, World!\n');
   * term.write('\x1b[1;31mBold Red\x1b[0m\n');
   * term.write(new Uint8Array([0x1b, 0x5b, 0x41])); // Up arrow
   * ```
   */
  write(A) {
    const B = typeof A == "string" ? new TextEncoder().encode(A) : A;
    if (B.length === 0)
      return;
    const Q = this.exports.ghostty_wasm_alloc_u8_array(B.length);
    new Uint8Array(this.memory.buffer).set(B, Q);
    try {
      this.exports.ghostty_terminal_write(this.handle, Q, B.length);
    } finally {
      this.exports.ghostty_wasm_free_u8_array(Q, B.length);
    }
  }
  /**
   * Resize the terminal.
   *
   * @param cols New column count
   * @param rows New row count
   */
  resize(A, B) {
    this.exports.ghostty_terminal_resize(this.handle, A, B), this._cols = A, this._rows = B;
  }
  /**
   * Get terminal dimensions.
   */
  get cols() {
    return this._cols;
  }
  get rows() {
    return this._rows;
  }
  /**
   * Get terminal dimensions (for IRenderable compatibility)
   */
  getDimensions() {
    return { cols: this._cols, rows: this._rows };
  }
  /**
   * Get cursor position and visibility.
   */
  getCursor() {
    return {
      x: this.exports.ghostty_terminal_get_cursor_x(this.handle),
      y: this.exports.ghostty_terminal_get_cursor_y(this.handle),
      visible: this.exports.ghostty_terminal_get_cursor_visible(this.handle)
    };
  }
  /**
   * Get scrollback length (number of lines in history).
   */
  getScrollbackLength() {
    return this.exports.ghostty_terminal_get_scrollback_length(this.handle);
  }
  /**
   * Check if terminal is in alternate screen buffer mode.
   *
   * The alternate screen is used by vim, less, htop, etc.
   * When active, normal buffer is preserved and restored when the app exits.
   *
   * @returns true if in alternate screen, false if in normal screen
   *
   * @example
   * ```typescript
   * // Detect if vim is running
   * if (term.isAlternateScreen()) {
   *   console.log('Full-screen app is active');
   * }
   * ```
   */
  isAlternateScreen() {
    return !!this.exports.ghostty_terminal_is_alternate_screen(this.handle);
  }
  /**
   * Check if a row is wrapped from the previous row.
   *
   * Wrapped rows are continuations of long lines that exceeded terminal width.
   * Used for text selection to treat wrapped lines as single logical lines.
   *
   * @param row Row index (0 = top visible line)
   * @returns true if row continues from previous line, false otherwise
   *
   * @example
   * ```typescript
   * // Get full logical line including wraps
   * let text = '';
   * for (let row = 0; row < term.rows; row++) {
   *   const line = term.getLine(row);
   *   text += lineToString(line);
   *
   *   // Only add newline if NOT wrapped
   *   if (!term.isRowWrapped(row + 1)) {
   *     text += '\n';
   *   }
   * }
   * ```
   */
  isRowWrapped(A) {
    return A < 0 || A >= this._rows ? !1 : !!this.exports.ghostty_terminal_is_row_wrapped(this.handle, A);
  }
  /**
   * Get a line of cells from the visible screen.
   *
   * @param y Line number (0 = top visible line)
   * @returns Array of cells, or null if y is out of bounds
   *
   * @example
   * ```typescript
   * const cells = term.getLine(0);
   * if (cells) {
   *   for (const cell of cells) {
   *     const char = String.fromCodePoint(cell.codepoint);
   *     const isBold = (cell.flags & CellFlags.BOLD) !== 0;
   *     console.log(`"${char}" ${isBold ? 'bold' : 'normal'}`);
   *   }
   * }
   * ```
   */
  getLine(A) {
    if (A < 0 || A >= this._rows)
      return null;
    const B = this._cols * Y.CELL_SIZE, Q = this.exports.ghostty_wasm_alloc_u8_array(B);
    try {
      const g = this.exports.ghostty_terminal_get_line(this.handle, A, Q, this._cols);
      if (g < 0)
        return null;
      const C = [], I = new DataView(this.memory.buffer, Q, B);
      for (let w = 0; w < g; w++) {
        const o = w * Y.CELL_SIZE;
        C.push({
          codepoint: I.getUint32(o, !0),
          fg_r: I.getUint8(o + 4),
          fg_g: I.getUint8(o + 5),
          fg_b: I.getUint8(o + 6),
          bg_r: I.getUint8(o + 7),
          bg_g: I.getUint8(o + 8),
          bg_b: I.getUint8(o + 9),
          flags: I.getUint8(o + 10),
          width: I.getUint8(o + 11),
          hyperlink_id: I.getUint16(o + 12, !0)
        });
      }
      return C;
    } finally {
      this.exports.ghostty_wasm_free_u8_array(Q, B);
    }
  }
  /**
   * Get a line from scrollback history.
   *
   * @param offset Line offset from top of scrollback (0 = oldest line)
   * @returns Array of cells, or null if not available
   */
  getScrollbackLine(A) {
    const B = this.getScrollbackLength();
    if (A < 0 || A >= B)
      return null;
    const Q = this._cols * Y.CELL_SIZE, g = this.exports.ghostty_wasm_alloc_u8_array(Q);
    try {
      const C = this.exports.ghostty_terminal_get_scrollback_line(
        this.handle,
        A,
        g,
        this._cols
      );
      if (C < 0)
        return null;
      const I = [], w = new DataView(this.memory.buffer, g, Q);
      for (let o = 0; o < C; o++) {
        const i = o * Y.CELL_SIZE;
        I.push({
          codepoint: w.getUint32(i, !0),
          fg_r: w.getUint8(i + 4),
          fg_g: w.getUint8(i + 5),
          fg_b: w.getUint8(i + 6),
          bg_r: w.getUint8(i + 7),
          bg_g: w.getUint8(i + 8),
          bg_b: w.getUint8(i + 9),
          flags: w.getUint8(i + 10),
          width: w.getUint8(i + 11),
          hyperlink_id: w.getUint16(i + 12, !0)
        });
      }
      return I;
    } finally {
      this.exports.ghostty_wasm_free_u8_array(g, Q);
    }
  }
  /**
   * Check if any part of the screen is dirty.
   */
  isDirty() {
    return this.exports.ghostty_terminal_is_dirty(this.handle);
  }
  /**
   * Check if a specific row is dirty.
   */
  isRowDirty(A) {
    return A < 0 || A >= this._rows ? !1 : this.exports.ghostty_terminal_is_row_dirty(this.handle, A);
  }
  /**
   * Clear all dirty flags (call after rendering).
   */
  clearDirty() {
    this.exports.ghostty_terminal_clear_dirty(this.handle);
  }
  /**
   * Get all visible lines at once (convenience method).
   *
   * @returns Array of line arrays, or empty array on error
   */
  getAllLines() {
    const A = [];
    for (let B = 0; B < this._rows; B++) {
      const Q = this.getLine(B);
      Q && A.push(Q);
    }
    return A;
  }
  /**
   * Get only the dirty lines (for optimized rendering).
   *
   * @returns Map of row number to cell array
   */
  getDirtyLines() {
    const A = /* @__PURE__ */ new Map();
    for (let B = 0; B < this._rows; B++)
      if (this.isRowDirty(B)) {
        const Q = this.getLine(B);
        Q && A.set(B, Q);
      }
    return A;
  }
  /**
   * Get hyperlink URI by ID
   *
   * @param hyperlinkId Hyperlink ID from a GhosttyCell (0 = no link)
   * @returns URI string or null if ID is invalid/not found
   */
  getHyperlinkUri(A) {
    if (A === 0)
      return null;
    const B = 2048, Q = this.exports.ghostty_wasm_alloc_u8_array(B);
    try {
      const g = this.exports.ghostty_terminal_get_hyperlink_uri(
        this.handle,
        A,
        Q,
        B
      );
      if (g === 0)
        return null;
      const C = new Uint8Array(this.memory.buffer, Q, g);
      return new TextDecoder().decode(C);
    } finally {
      this.exports.ghostty_wasm_free_u8_array(Q, B);
    }
  }
  // ============================================================================
  // Terminal Modes
  // ============================================================================
  /**
   * Query terminal mode state
   */
  getMode(A, B = !1) {
    return this.exports.ghostty_terminal_get_mode(this.handle, A, B ? 1 : 0) !== 0;
  }
  /**
   * Check if bracketed paste mode is enabled
   */
  hasBracketedPaste() {
    return this.exports.ghostty_terminal_has_bracketed_paste(this.handle) !== 0;
  }
  /**
   * Check if focus event reporting is enabled
   */
  hasFocusEvents() {
    return this.exports.ghostty_terminal_has_focus_events(this.handle) !== 0;
  }
  /**
   * Check if mouse tracking is enabled
   */
  hasMouseTracking() {
    return this.exports.ghostty_terminal_has_mouse_tracking(this.handle) !== 0;
  }
};
f.CELL_SIZE = 16;
let j = f;
class K {
  constructor() {
    this.listeners = [], this.event = (A) => (this.listeners.push(A), {
      dispose: () => {
        const B = this.listeners.indexOf(A);
        B >= 0 && this.listeners.splice(B, 1);
      }
    });
  }
  fire(A) {
    for (const B of this.listeners)
      B(A);
  }
  dispose() {
    this.listeners = [];
  }
}
class V {
  constructor(A) {
    this.bufferChangeEmitter = new K(), this.terminal = A;
  }
  get active() {
    const A = this.terminal.wasmTerm;
    return A ? A.isAlternateScreen() ? this.alternate : this.normal : this.normal;
  }
  get normal() {
    return this._normalBuffer || (this._normalBuffer = new n(this.terminal, "normal")), this._normalBuffer;
  }
  get alternate() {
    return this._alternateBuffer || (this._alternateBuffer = new n(this.terminal, "alternate")), this._alternateBuffer;
  }
  get onBufferChange() {
    return this.bufferChangeEmitter.event;
  }
  /**
   * Internal: Fire buffer change event when screen switches
   * Should be called by Terminal when detecting screen change
   */
  _fireBufferChange(A) {
    this.bufferChangeEmitter.fire(A);
  }
}
class n {
  constructor(A, B) {
    this.terminal = A, this.bufferType = B;
    const Q = {
      codepoint: 0,
      fg_r: 204,
      fg_g: 204,
      fg_b: 204,
      bg_r: 0,
      bg_g: 0,
      bg_b: 0,
      flags: 0,
      width: 1,
      hyperlink_id: 0
    };
    this.nullCell = new S(Q, 0);
  }
  get type() {
    return this.bufferType;
  }
  get cursorX() {
    const A = this.getWasmTerm();
    return A ? A.getCursor().x : 0;
  }
  get cursorY() {
    const A = this.getWasmTerm();
    return A ? A.getCursor().y : 0;
  }
  get viewportY() {
    return 0;
  }
  get baseY() {
    return 0;
  }
  get length() {
    const A = this.getWasmTerm();
    return A ? this.bufferType === "alternate" ? A.rows : A.getScrollbackLength() + A.rows : 0;
  }
  getLine(A) {
    const B = this.getWasmTerm();
    if (!B || A < 0 || A >= this.length)
      return;
    const Q = B.getScrollbackLength();
    let g, C, I;
    if (this.bufferType === "normal" && A < Q) {
      const w = A;
      g = B.getScrollbackLine(w), I = !1;
    } else
      C = this.bufferType === "normal" ? A - Q : A, g = B.getLine(C), I = B.isRowWrapped(C);
    if (g)
      return new W(g, I, B.cols);
  }
  getNullCell() {
    return this.nullCell;
  }
  getWasmTerm() {
    return this.terminal.wasmTerm;
  }
}
class W {
  constructor(A, B, Q) {
    this.cells = A, this._isWrapped = B, this._length = Q;
  }
  get length() {
    return this._length;
  }
  get isWrapped() {
    return this._isWrapped;
  }
  getCell(A) {
    if (!(A < 0 || A >= this._length))
      return A >= this.cells.length ? new S(
        {
          codepoint: 0,
          fg_r: 204,
          fg_g: 204,
          fg_b: 204,
          bg_r: 0,
          bg_g: 0,
          bg_b: 0,
          flags: 0,
          width: 1,
          hyperlink_id: 0
        },
        A
      ) : new S(this.cells[A], A);
  }
  translateToString(A = !1, B = 0, Q = this._length) {
    const g = Math.max(0, Math.min(B, this._length)), C = Math.max(g, Math.min(Q, this._length));
    let I = "";
    for (let w = g; w < C; w++) {
      const o = this.getCell(w);
      if (o) {
        const i = o.getChars();
        I += i;
      }
    }
    return A && (I = I.trimEnd()), I;
  }
}
class S {
  constructor(A, B) {
    this.cell = A, this.x = B;
  }
  getChars() {
    const A = this.cell.codepoint;
    return A === 0 ? "" : A < 0 || A > 1114111 || A >= 55296 && A <= 57343 ? "�" : String.fromCodePoint(A);
  }
  getCode() {
    return this.cell.codepoint;
  }
  getWidth() {
    return this.cell.width;
  }
  getFgColorMode() {
    return -1;
  }
  getBgColorMode() {
    return -1;
  }
  getFgColor() {
    return this.cell.fg_r << 16 | this.cell.fg_g << 8 | this.cell.fg_b;
  }
  getBgColor() {
    return this.cell.bg_r << 16 | this.cell.bg_g << 8 | this.cell.bg_b;
  }
  isBold() {
    return this.cell.flags & U.BOLD ? 1 : 0;
  }
  isItalic() {
    return this.cell.flags & U.ITALIC ? 1 : 0;
  }
  isUnderline() {
    return this.cell.flags & U.UNDERLINE ? 1 : 0;
  }
  isStrikethrough() {
    return this.cell.flags & U.STRIKETHROUGH ? 1 : 0;
  }
  isBlink() {
    return this.cell.flags & U.BLINK ? 1 : 0;
  }
  isInverse() {
    return this.cell.flags & U.INVERSE ? 1 : 0;
  }
  isInvisible() {
    return this.cell.flags & U.INVISIBLE ? 1 : 0;
  }
  isFaint() {
    return this.cell.flags & U.FAINT ? 1 : 0;
  }
  /**
   * Get hyperlink ID for this cell (0 = no link)
   * Used by link detection system
   */
  getHyperlinkId() {
    return this.cell.hyperlink_id;
  }
  /**
   * Get the Unicode codepoint for this cell
   * Used by link detection system
   */
  getCodepoint() {
    return this.cell.codepoint;
  }
  /**
   * Check if cell has dim/faint attribute
   * Added for IBufferCell compatibility
   */
  isDim() {
    return (this.cell.flags & U.FAINT) !== 0;
  }
}
const Z = {
  // Letters
  KeyA: D.A,
  KeyB: D.B,
  KeyC: D.C,
  KeyD: D.D,
  KeyE: D.E,
  KeyF: D.F,
  KeyG: D.G,
  KeyH: D.H,
  KeyI: D.I,
  KeyJ: D.J,
  KeyK: D.K,
  KeyL: D.L,
  KeyM: D.M,
  KeyN: D.N,
  KeyO: D.O,
  KeyP: D.P,
  KeyQ: D.Q,
  KeyR: D.R,
  KeyS: D.S,
  KeyT: D.T,
  KeyU: D.U,
  KeyV: D.V,
  KeyW: D.W,
  KeyX: D.X,
  KeyY: D.Y,
  KeyZ: D.Z,
  // Numbers
  Digit1: D.ONE,
  Digit2: D.TWO,
  Digit3: D.THREE,
  Digit4: D.FOUR,
  Digit5: D.FIVE,
  Digit6: D.SIX,
  Digit7: D.SEVEN,
  Digit8: D.EIGHT,
  Digit9: D.NINE,
  Digit0: D.ZERO,
  // Special keys
  Enter: D.ENTER,
  Escape: D.ESCAPE,
  Backspace: D.BACKSPACE,
  Tab: D.TAB,
  Space: D.SPACE,
  // Punctuation
  Minus: D.MINUS,
  Equal: D.EQUAL,
  BracketLeft: D.BRACKET_LEFT,
  BracketRight: D.BRACKET_RIGHT,
  Backslash: D.BACKSLASH,
  Semicolon: D.SEMICOLON,
  Quote: D.QUOTE,
  Backquote: D.GRAVE,
  Comma: D.COMMA,
  Period: D.PERIOD,
  Slash: D.SLASH,
  // Function keys
  CapsLock: D.CAPS_LOCK,
  F1: D.F1,
  F2: D.F2,
  F3: D.F3,
  F4: D.F4,
  F5: D.F5,
  F6: D.F6,
  F7: D.F7,
  F8: D.F8,
  F9: D.F9,
  F10: D.F10,
  F11: D.F11,
  F12: D.F12,
  // Special function keys
  PrintScreen: D.PRINT_SCREEN,
  ScrollLock: D.SCROLL_LOCK,
  Pause: D.PAUSE,
  Insert: D.INSERT,
  Home: D.HOME,
  PageUp: D.PAGE_UP,
  Delete: D.DELETE,
  End: D.END,
  PageDown: D.PAGE_DOWN,
  // Arrow keys
  ArrowRight: D.RIGHT,
  ArrowLeft: D.LEFT,
  ArrowDown: D.DOWN,
  ArrowUp: D.UP,
  // Keypad
  NumLock: D.NUM_LOCK,
  NumpadDivide: D.KP_DIVIDE,
  NumpadMultiply: D.KP_MULTIPLY,
  NumpadSubtract: D.KP_MINUS,
  NumpadAdd: D.KP_PLUS,
  NumpadEnter: D.KP_ENTER,
  Numpad1: D.KP_1,
  Numpad2: D.KP_2,
  Numpad3: D.KP_3,
  Numpad4: D.KP_4,
  Numpad5: D.KP_5,
  Numpad6: D.KP_6,
  Numpad7: D.KP_7,
  Numpad8: D.KP_8,
  Numpad9: D.KP_9,
  Numpad0: D.KP_0,
  NumpadDecimal: D.KP_PERIOD,
  // International
  IntlBackslash: D.NON_US_BACKSLASH,
  ContextMenu: D.APPLICATION,
  // Additional function keys
  F13: D.F13,
  F14: D.F14,
  F15: D.F15,
  F16: D.F16,
  F17: D.F17,
  F18: D.F18,
  F19: D.F19,
  F20: D.F20,
  F21: D.F21,
  F22: D.F22,
  F23: D.F23,
  F24: D.F24
};
class X {
  /**
   * Create a new InputHandler
   * @param ghostty - Ghostty instance (for creating KeyEncoder)
   * @param container - DOM element to attach listeners to
   * @param onData - Callback for terminal data (escape sequences to send to PTY)
   * @param onBell - Callback for bell/beep event
   * @param onKey - Optional callback for raw key events
   * @param customKeyEventHandler - Optional custom key event handler
   */
  constructor(A, B, Q, g, C, I) {
    this.keydownListener = null, this.keypressListener = null, this.pasteListener = null, this.isDisposed = !1, this.encoder = A.createKeyEncoder(), this.container = B, this.onDataCallback = Q, this.onBellCallback = g, this.onKeyCallback = C, this.customKeyEventHandler = I, this.attach();
  }
  /**
   * Set custom key event handler (for runtime updates)
   */
  setCustomKeyEventHandler(A) {
    this.customKeyEventHandler = A;
  }
  /**
   * Attach keyboard event listeners to container
   */
  attach() {
    typeof this.container.hasAttribute == "function" && typeof this.container.setAttribute == "function" && (this.container.hasAttribute("tabindex") || this.container.setAttribute("tabindex", "0"), this.container.style && (this.container.style.outline = "none")), this.keydownListener = this.handleKeyDown.bind(this), this.container.addEventListener("keydown", this.keydownListener), this.pasteListener = this.handlePaste.bind(this), this.container.addEventListener("paste", this.pasteListener);
  }
  /**
   * Map KeyboardEvent.code to USB HID Key enum value
   * @param code - KeyboardEvent.code value
   * @returns Key enum value or null if unmapped
   */
  mapKeyCode(A) {
    return Z[A] ?? null;
  }
  /**
   * Extract modifier flags from KeyboardEvent
   * @param event - KeyboardEvent
   * @returns Mods flags
   */
  extractModifiers(A) {
    let B = q.NONE;
    return A.shiftKey && (B |= q.SHIFT), A.ctrlKey && (B |= q.CTRL), A.altKey && (B |= q.ALT), A.metaKey && (B |= q.SUPER), B;
  }
  /**
   * Check if this is a printable character with no special modifiers
   * @param event - KeyboardEvent
   * @returns true if printable character
   */
  isPrintableCharacter(A) {
    return A.ctrlKey && !A.altKey || A.altKey && !A.ctrlKey || A.metaKey ? !1 : A.key.length === 1;
  }
  /**
   * Handle keydown event
   * @param event - KeyboardEvent
   */
  handleKeyDown(A) {
    if (this.isDisposed)
      return;
    if (this.onKeyCallback && this.onKeyCallback({ key: A.key, domEvent: A }), this.customKeyEventHandler && this.customKeyEventHandler(A)) {
      A.preventDefault();
      return;
    }
    if ((A.ctrlKey || A.metaKey) && A.code === "KeyV" || A.metaKey && A.code === "KeyC")
      return;
    if (this.isPrintableCharacter(A)) {
      A.preventDefault(), this.onDataCallback(A.key);
      return;
    }
    const B = this.mapKeyCode(A.code);
    if (B === null)
      return;
    const Q = this.extractModifiers(A);
    if (Q === q.NONE || Q === q.SHIFT) {
      let C = null;
      switch (B) {
        case D.ENTER:
          C = "\r";
          break;
        case D.TAB:
          C = "	";
          break;
        case D.BACKSPACE:
          C = "";
          break;
        case D.ESCAPE:
          C = "\x1B";
          break;
        case D.UP:
          C = "\x1B[A";
          break;
        case D.DOWN:
          C = "\x1B[B";
          break;
        case D.RIGHT:
          C = "\x1B[C";
          break;
        case D.LEFT:
          C = "\x1B[D";
          break;
        case D.HOME:
          C = "\x1B[H";
          break;
        case D.END:
          C = "\x1B[F";
          break;
        case D.INSERT:
          C = "\x1B[2~";
          break;
        case D.DELETE:
          C = "\x1B[3~";
          break;
        case D.PAGE_UP:
          C = "\x1B[5~";
          break;
        case D.PAGE_DOWN:
          C = "\x1B[6~";
          break;
        case D.F1:
          C = "\x1BOP";
          break;
        case D.F2:
          C = "\x1BOQ";
          break;
        case D.F3:
          C = "\x1BOR";
          break;
        case D.F4:
          C = "\x1BOS";
          break;
        case D.F5:
          C = "\x1B[15~";
          break;
        case D.F6:
          C = "\x1B[17~";
          break;
        case D.F7:
          C = "\x1B[18~";
          break;
        case D.F8:
          C = "\x1B[19~";
          break;
        case D.F9:
          C = "\x1B[20~";
          break;
        case D.F10:
          C = "\x1B[21~";
          break;
        case D.F11:
          C = "\x1B[23~";
          break;
        case D.F12:
          C = "\x1B[24~";
          break;
      }
      if (C !== null) {
        A.preventDefault(), this.onDataCallback(C);
        return;
      }
    }
    const g = b.PRESS;
    try {
      const C = A.key.length === 1 && A.key.charCodeAt(0) < 128 ? A.key.toLowerCase() : void 0, I = this.encoder.encode({
        action: g,
        key: B,
        mods: Q,
        utf8: C
      }), o = new TextDecoder().decode(I);
      A.preventDefault(), A.stopPropagation(), o.length > 0 && this.onDataCallback(o);
    } catch (C) {
      console.warn("Failed to encode key:", A.code, C);
    }
  }
  /**
   * Handle paste event from clipboard
   * @param event - ClipboardEvent
   */
  handlePaste(A) {
    if (this.isDisposed)
      return;
    A.preventDefault(), A.stopPropagation();
    const B = A.clipboardData;
    if (!B) {
      console.warn("No clipboard data available");
      return;
    }
    const Q = B.getData("text/plain");
    if (!Q) {
      console.warn("No text in clipboard");
      return;
    }
    this.onDataCallback(Q);
  }
  /**
   * Dispose the InputHandler and remove event listeners
   */
  dispose() {
    this.isDisposed || (this.keydownListener && (this.container.removeEventListener("keydown", this.keydownListener), this.keydownListener = null), this.keypressListener && (this.container.removeEventListener("keypress", this.keypressListener), this.keypressListener = null), this.pasteListener && (this.container.removeEventListener("paste", this.pasteListener), this.pasteListener = null), this.isDisposed = !0);
  }
  /**
   * Check if handler is disposed
   */
  isActive() {
    return !this.isDisposed;
  }
}
class P {
  // Terminal instance for buffer access
  constructor(A) {
    this.terminal = A, this.providers = [], this.linkCache = /* @__PURE__ */ new Map(), this.scannedRows = /* @__PURE__ */ new Set();
  }
  /**
   * Register a link provider
   */
  registerProvider(A) {
    this.providers.push(A), this.invalidateCache();
  }
  /**
   * Get link at the specified buffer position
   * @param col Column (0-based)
   * @param row Absolute row in buffer (0-based)
   * @returns Link at position, or undefined if none
   */
  async getLinkAt(A, B) {
    const Q = this.terminal.buffer.active.getLine(B);
    if (!Q || A < 0 || A >= Q.length)
      return;
    const g = Q.getCell(A);
    if (!g)
      return;
    const C = g.getHyperlinkId();
    if (C > 0) {
      const I = `h${C}`;
      if (this.linkCache.has(I))
        return this.linkCache.get(I);
    }
    if (this.scannedRows.has(B) || await this.scanRow(B), C > 0) {
      const I = `h${C}`, w = this.linkCache.get(I);
      if (w)
        return w;
    }
    for (const I of this.linkCache.values())
      if (this.isPositionInLink(A, B, I))
        return I;
  }
  /**
   * Scan a row for links using all registered providers
   */
  async scanRow(A) {
    this.scannedRows.add(A);
    const B = [];
    for (const Q of this.providers) {
      const g = await new Promise((C) => {
        Q.provideLinks(A, C);
      });
      g && B.push(...g);
    }
    for (const Q of B)
      this.cacheLink(Q);
  }
  /**
   * Cache a link for fast lookup
   */
  cacheLink(A) {
    const { start: B } = A.range, Q = this.terminal.buffer.active.getLine(B.y);
    if (Q) {
      const w = Q.getCell(B.x);
      if (!w) {
        const { start: i, end: M } = A.range, k = `r${i.y}:${i.x}-${M.x}`;
        this.linkCache.set(k, A);
        return;
      }
      const o = w.getHyperlinkId();
      if (o > 0) {
        this.linkCache.set(`h${o}`, A);
        return;
      }
    }
    const { start: g, end: C } = A.range, I = `r${g.y}:${g.x}-${C.x}`;
    this.linkCache.set(I, A);
  }
  /**
   * Check if a position is within a link's range
   */
  isPositionInLink(A, B, Q) {
    const { start: g, end: C } = Q.range;
    return B < g.y || B > C.y ? !1 : g.y === C.y ? A >= g.x && A <= C.x : B === g.y ? A >= g.x : B === C.y ? A <= C.x : !0;
  }
  /**
   * Invalidate cache when terminal content changes
   * Should be called on terminal write, resize, or clear
   */
  invalidateCache() {
    this.linkCache.clear(), this.scannedRows.clear();
  }
  /**
   * Invalidate cache for specific rows
   * Used when only part of the terminal changed
   */
  invalidateRows(A, B) {
    for (let g = A; g <= B; g++)
      this.scannedRows.delete(g);
    const Q = [];
    for (const [g, C] of this.linkCache.entries()) {
      const { start: I, end: w } = C.range;
      (I.y >= A && I.y <= B || w.y >= A && w.y <= B || I.y < A && w.y > B) && Q.push(g);
    }
    for (const g of Q)
      this.linkCache.delete(g);
  }
  /**
   * Dispose and cleanup
   */
  dispose() {
    var A;
    this.linkCache.clear(), this.scannedRows.clear();
    for (const B of this.providers)
      (A = B.dispose) == null || A.call(B);
    this.providers = [];
  }
}
class v {
  constructor(A) {
    this.terminal = A;
  }
  /**
   * Provide all OSC 8 links on the given row
   * Note: This may return links that span multiple rows
   */
  provideLinks(A, B) {
    const Q = [], g = /* @__PURE__ */ new Set(), C = this.terminal.buffer.active.getLine(A);
    if (!C) {
      B(void 0);
      return;
    }
    for (let I = 0; I < C.length; I++) {
      const w = C.getCell(I);
      if (!w)
        continue;
      const o = w.getHyperlinkId();
      if (o === 0 || g.has(o))
        continue;
      g.add(o);
      const i = this.findLinkRange(o, A, I);
      if (!this.terminal.wasmTerm)
        continue;
      const M = this.terminal.wasmTerm.getHyperlinkUri(o);
      M && Q.push({
        text: M,
        range: i,
        activate: (k) => {
          (k.ctrlKey || k.metaKey) && window.open(M, "_blank", "noopener,noreferrer");
        }
      });
    }
    B(Q.length > 0 ? Q : void 0);
  }
  /**
   * Find the full extent of a link by scanning for contiguous cells
   * with the same hyperlink_id. Handles multi-line links.
   */
  findLinkRange(A, B, Q) {
    const g = this.terminal.buffer.active;
    let C = B, I = Q;
    for (; I > 0; ) {
      const M = g.getLine(C);
      if (!M)
        break;
      const k = M.getCell(I - 1);
      if (!k || k.getHyperlinkId() !== A)
        break;
      I--;
    }
    if (I === 0 && C > 0) {
      let M = C - 1;
      for (; M >= 0; ) {
        const k = g.getLine(M);
        if (!k || k.length === 0)
          break;
        const G = k.getCell(k.length - 1);
        if (!G || G.getHyperlinkId() !== A)
          break;
        C = M, I = 0;
        for (let J = k.length - 1; J >= 0; J--) {
          const s = k.getCell(J);
          if (!s || s.getHyperlinkId() !== A) {
            I = J + 1;
            break;
          }
        }
        if (I === 0)
          M--;
        else
          break;
      }
    }
    let w = B, o = Q;
    const i = g.getLine(w);
    if (i) {
      for (; o < i.length - 1; ) {
        const M = i.getCell(o + 1);
        if (!M || M.getHyperlinkId() !== A)
          break;
        o++;
      }
      if (o === i.length - 1) {
        let M = w + 1;
        const k = g.length;
        for (; M < k; ) {
          const G = g.getLine(M);
          if (!G || G.length === 0)
            break;
          const J = G.getCell(0);
          if (!J || J.getHyperlinkId() !== A)
            break;
          w = M, o = 0;
          for (let s = 0; s < G.length; s++) {
            const F = G.getCell(s);
            if (!F)
              break;
            if (F.getHyperlinkId() !== A) {
              o = s - 1;
              break;
            }
            o = s;
          }
          if (o === G.length - 1)
            M++;
          else
            break;
        }
      }
    }
    return {
      start: { x: I, y: C },
      end: { x: o, y: w }
    };
  }
  dispose() {
  }
}
const e = class H {
  constructor(A) {
    this.terminal = A;
  }
  /**
   * Provide all regex-detected URLs on the given row
   */
  provideLinks(A, B) {
    const Q = [], g = this.terminal.buffer.active.getLine(A);
    if (!g) {
      B(void 0);
      return;
    }
    const C = this.lineToText(g);
    H.URL_REGEX.lastIndex = 0;
    let I = H.URL_REGEX.exec(C);
    for (; I !== null; ) {
      let w = I[0];
      const o = I.index;
      let i = I.index + w.length - 1;
      const M = w.replace(H.TRAILING_PUNCTUATION, "");
      M.length < w.length && (w = M, i = o + w.length - 1), w.length > 8 && Q.push({
        text: w,
        range: {
          start: { x: o, y: A },
          end: { x: i, y: A }
        },
        activate: (k) => {
          (k.ctrlKey || k.metaKey) && window.open(w, "_blank", "noopener,noreferrer");
        }
      }), I = H.URL_REGEX.exec(C);
    }
    B(Q.length > 0 ? Q : void 0);
  }
  /**
   * Convert a buffer line to plain text string
   */
  lineToText(A) {
    const B = [];
    for (let Q = 0; Q < A.length; Q++) {
      const g = A.getCell(Q);
      if (!g) {
        B.push(" ");
        continue;
      }
      const C = g.getCodepoint();
      C === 0 || C < 32 ? B.push(" ") : B.push(String.fromCodePoint(C));
    }
    return B.join("");
  }
  dispose() {
  }
};
e.URL_REGEX = /(?:https?:\/\/|mailto:|ftp:\/\/|ssh:\/\/|git:\/\/|tel:|magnet:|gemini:\/\/|gopher:\/\/|news:)[\w\-.~:\/?#@!$&*+,;=%]+/gi;
e.TRAILING_PUNCTUATION = /[.,;!?)\]]+$/;
let u = e;
const r = {
  foreground: "#d4d4d4",
  background: "#1e1e1e",
  cursor: "#ffffff",
  cursorAccent: "#1e1e1e",
  selectionBackground: "rgba(255, 255, 255, 0.3)",
  selectionForeground: "#d4d4d4",
  black: "#000000",
  red: "#cd3131",
  green: "#0dbc79",
  yellow: "#e5e510",
  blue: "#2472c8",
  magenta: "#bc3fbc",
  cyan: "#11a8cd",
  white: "#e5e5e5",
  brightBlack: "#666666",
  brightRed: "#f14c4c",
  brightGreen: "#23d18b",
  brightYellow: "#f5f543",
  brightBlue: "#3b8eea",
  brightMagenta: "#d670d6",
  brightCyan: "#29b8db",
  brightWhite: "#ffffff"
};
class _ {
  constructor(A, B = {}) {
    this.cursorVisible = !0, this.lastCursorPosition = { x: 0, y: 0 }, this.lastViewportY = 0, this.hoveredHyperlinkId = 0, this.previousHoveredHyperlinkId = 0, this.hoveredLinkRange = null, this.previousHoveredLinkRange = null, this.canvas = A;
    const Q = A.getContext("2d", { alpha: !1 });
    if (!Q)
      throw new Error("Failed to get 2D rendering context");
    this.ctx = Q, this.fontSize = B.fontSize ?? 15, this.fontFamily = B.fontFamily ?? "monospace", this.cursorStyle = B.cursorStyle ?? "block", this.cursorBlink = B.cursorBlink ?? !1, this.theme = { ...r, ...B.theme }, this.devicePixelRatio = B.devicePixelRatio ?? window.devicePixelRatio ?? 1, this.palette = [
      this.theme.black,
      this.theme.red,
      this.theme.green,
      this.theme.yellow,
      this.theme.blue,
      this.theme.magenta,
      this.theme.cyan,
      this.theme.white,
      this.theme.brightBlack,
      this.theme.brightRed,
      this.theme.brightGreen,
      this.theme.brightYellow,
      this.theme.brightBlue,
      this.theme.brightMagenta,
      this.theme.brightCyan,
      this.theme.brightWhite
    ], this.metrics = this.measureFont(), this.cursorBlink && this.startCursorBlink();
  }
  // ==========================================================================
  // Font Metrics Measurement
  // ==========================================================================
  measureFont() {
    const B = document.createElement("canvas").getContext("2d");
    B.font = `${this.fontSize}px ${this.fontFamily}`;
    const Q = B.measureText("M"), g = Math.ceil(Q.width), C = Q.actualBoundingBoxAscent || this.fontSize * 0.8, I = Q.actualBoundingBoxDescent || this.fontSize * 0.2, w = Math.ceil(C + I) + 2, o = Math.ceil(C) + 1;
    return { width: g, height: w, baseline: o };
  }
  /**
   * Remeasure font metrics (call after font loads or changes)
   */
  remeasureFont() {
    this.metrics = this.measureFont();
  }
  // ==========================================================================
  // Color Conversion
  // ==========================================================================
  rgbToCSS(A, B, Q) {
    return `rgb(${A}, ${B}, ${Q})`;
  }
  // ==========================================================================
  // Canvas Sizing
  // ==========================================================================
  /**
   * Resize canvas to fit terminal dimensions
   */
  resize(A, B) {
    const Q = A * this.metrics.width, g = B * this.metrics.height;
    this.canvas.style.width = `${Q}px`, this.canvas.style.height = `${g}px`, this.canvas.width = Q * this.devicePixelRatio, this.canvas.height = g * this.devicePixelRatio, this.ctx.scale(this.devicePixelRatio, this.devicePixelRatio), this.ctx.textBaseline = "alphabetic", this.ctx.textAlign = "left", this.ctx.fillStyle = this.theme.background, this.ctx.fillRect(0, 0, Q, g);
  }
  // ==========================================================================
  // Main Rendering
  // ==========================================================================
  /**
   * Render the terminal buffer to canvas
   */
  render(A, B = !1, Q = 0, g, C = 1) {
    const I = A.getCursor(), w = A.getDimensions(), o = g ? g.getScrollbackLength() : 0;
    (this.canvas.width !== w.cols * this.metrics.width * this.devicePixelRatio || this.canvas.height !== w.rows * this.metrics.height * this.devicePixelRatio) && (this.resize(w.cols, w.rows), B = !0), Q !== this.lastViewportY && (B = !0, this.lastViewportY = Q);
    const M = I.x !== this.lastCursorPosition.x || I.y !== this.lastCursorPosition.y;
    if (M || this.cursorBlink) {
      if (!B && !A.isRowDirty(I.y)) {
        const N = A.getLine(I.y);
        N && this.renderLine(N, I.y, w.cols);
      }
      if (M && this.lastCursorPosition.y !== I.y && !B && !A.isRowDirty(this.lastCursorPosition.y)) {
        const N = A.getLine(this.lastCursorPosition.y);
        N && this.renderLine(N, this.lastCursorPosition.y, w.cols);
      }
    }
    const k = this.selectionManager && this.selectionManager.hasSelection(), G = /* @__PURE__ */ new Set();
    if (k) {
      const N = this.selectionManager.getSelectionCoords();
      if (N)
        for (let c = N.startRow; c <= N.endRow; c++)
          G.add(c);
    }
    if (this.selectionManager) {
      const N = this.selectionManager.getDirtySelectionRows();
      if (N.size > 0) {
        for (const c of N)
          G.add(c);
        this.selectionManager.clearDirtySelectionRows();
      }
    }
    const J = /* @__PURE__ */ new Set(), s = this.hoveredHyperlinkId !== this.previousHoveredHyperlinkId, F = JSON.stringify(this.hoveredLinkRange) !== JSON.stringify(this.previousHoveredLinkRange);
    if (s) {
      for (let N = 0; N < w.rows; N++) {
        let c = null;
        if (Q > 0)
          if (N < Q && g) {
            const h = o - Math.floor(Q) + N;
            c = g.getScrollbackLine(h);
          } else {
            const h = N - Math.floor(Q);
            c = A.getLine(h);
          }
        else
          c = A.getLine(N);
        if (c) {
          for (const h of c)
            if (h.hyperlink_id === this.hoveredHyperlinkId || h.hyperlink_id === this.previousHoveredHyperlinkId) {
              J.add(N);
              break;
            }
        }
      }
      this.previousHoveredHyperlinkId = this.hoveredHyperlinkId;
    }
    if (F) {
      if (this.previousHoveredLinkRange)
        for (let N = this.previousHoveredLinkRange.startY; N <= this.previousHoveredLinkRange.endY; N++)
          J.add(N);
      if (this.hoveredLinkRange)
        for (let N = this.hoveredLinkRange.startY; N <= this.hoveredLinkRange.endY; N++)
          J.add(N);
      this.previousHoveredLinkRange = this.hoveredLinkRange;
    }
    let a = !1;
    for (let N = 0; N < w.rows; N++) {
      if (!(Q > 0 ? !0 : B || A.isRowDirty(N) || G.has(N) || J.has(N)))
        continue;
      a = !0;
      let h = null;
      if (Q > 0)
        if (N < Q && g) {
          const y = o - Math.floor(Q) + N;
          h = g.getScrollbackLine(y);
        } else {
          const y = Q > 0 ? N - Math.floor(Q) : N;
          h = A.getLine(y);
        }
      else
        h = A.getLine(N);
      h && this.renderLine(h, N, w.cols);
    }
    k && a && this.renderSelection(w.cols), Q === 0 && I.visible && this.cursorVisible && this.renderCursor(I.x, I.y), g && C > 0 && this.renderScrollbar(Q, o, w.rows, C), this.lastCursorPosition = { x: I.x, y: I.y }, B || A.clearDirty();
  }
  /**
   * Render a single line
   */
  renderLine(A, B, Q) {
    const g = B * this.metrics.height;
    this.ctx.fillStyle = this.theme.background, this.ctx.fillRect(0, g, Q * this.metrics.width, this.metrics.height);
    for (let C = 0; C < A.length; C++) {
      const I = A[C];
      I.width !== 0 && this.renderCell(I, C, B);
    }
  }
  /**
   * Render a single cell
   */
  renderCell(A, B, Q) {
    const g = B * this.metrics.width, C = Q * this.metrics.height, I = this.metrics.width * A.width;
    let w = A.fg_r, o = A.fg_g, i = A.fg_b, M = A.bg_r, k = A.bg_g, G = A.bg_b;
    if (A.flags & U.INVERSE && ([w, o, i, M, k, G] = [M, k, G, w, o, i]), this.ctx.fillStyle = this.rgbToCSS(M, k, G), this.ctx.fillRect(g, C, I, this.metrics.height), A.flags & U.INVISIBLE)
      return;
    let J = "";
    A.flags & U.ITALIC && (J += "italic "), A.flags & U.BOLD && (J += "bold "), this.ctx.font = `${J}${this.fontSize}px ${this.fontFamily}`, this.ctx.fillStyle = this.rgbToCSS(w, o, i), A.flags & U.FAINT && (this.ctx.globalAlpha = 0.5);
    const s = g, F = C + this.metrics.baseline, a = String.fromCodePoint(A.codepoint || 32);
    if (this.ctx.fillText(a, s, F), A.flags & U.FAINT && (this.ctx.globalAlpha = 1), A.flags & U.UNDERLINE) {
      const N = C + this.metrics.baseline + 2;
      this.ctx.strokeStyle = this.ctx.fillStyle, this.ctx.lineWidth = 1, this.ctx.beginPath(), this.ctx.moveTo(g, N), this.ctx.lineTo(g + I, N), this.ctx.stroke();
    }
    if (A.flags & U.STRIKETHROUGH) {
      const N = C + this.metrics.height / 2;
      this.ctx.strokeStyle = this.ctx.fillStyle, this.ctx.lineWidth = 1, this.ctx.beginPath(), this.ctx.moveTo(g, N), this.ctx.lineTo(g + I, N), this.ctx.stroke();
    }
    if (A.hyperlink_id > 0 && A.hyperlink_id === this.hoveredHyperlinkId) {
      const c = C + this.metrics.baseline + 2;
      this.ctx.strokeStyle = "#4A90E2", this.ctx.lineWidth = 1, this.ctx.beginPath(), this.ctx.moveTo(g, c), this.ctx.lineTo(g + I, c), this.ctx.stroke();
    }
    if (this.hoveredLinkRange) {
      const N = this.hoveredLinkRange;
      if (Q === N.startY && B >= N.startX && (Q < N.endY || B <= N.endX) || Q > N.startY && Q < N.endY || Q === N.endY && B <= N.endX && (Q > N.startY || B >= N.startX)) {
        const h = C + this.metrics.baseline + 2;
        this.ctx.strokeStyle = "#4A90E2", this.ctx.lineWidth = 1, this.ctx.beginPath(), this.ctx.moveTo(g, h), this.ctx.lineTo(g + I, h), this.ctx.stroke();
      }
    }
  }
  /**
   * Render cursor
   */
  renderCursor(A, B) {
    const Q = A * this.metrics.width, g = B * this.metrics.height;
    switch (this.ctx.fillStyle = this.theme.cursor, this.cursorStyle) {
      case "block":
        this.ctx.fillRect(Q, g, this.metrics.width, this.metrics.height);
        break;
      case "underline":
        const C = Math.max(2, Math.floor(this.metrics.height * 0.15));
        this.ctx.fillRect(
          Q,
          g + this.metrics.height - C,
          this.metrics.width,
          C
        );
        break;
      case "bar":
        const I = Math.max(2, Math.floor(this.metrics.width * 0.15));
        this.ctx.fillRect(Q, g, I, this.metrics.height);
        break;
    }
  }
  // ==========================================================================
  // Cursor Blinking
  // ==========================================================================
  startCursorBlink() {
    this.cursorBlinkInterval = window.setInterval(() => {
      this.cursorVisible = !this.cursorVisible;
    }, 530);
  }
  stopCursorBlink() {
    this.cursorBlinkInterval !== void 0 && (clearInterval(this.cursorBlinkInterval), this.cursorBlinkInterval = void 0), this.cursorVisible = !0;
  }
  // ==========================================================================
  // Public API
  // ==========================================================================
  /**
   * Update theme colors
   */
  setTheme(A) {
    this.theme = { ...r, ...A }, this.palette = [
      this.theme.black,
      this.theme.red,
      this.theme.green,
      this.theme.yellow,
      this.theme.blue,
      this.theme.magenta,
      this.theme.cyan,
      this.theme.white,
      this.theme.brightBlack,
      this.theme.brightRed,
      this.theme.brightGreen,
      this.theme.brightYellow,
      this.theme.brightBlue,
      this.theme.brightMagenta,
      this.theme.brightCyan,
      this.theme.brightWhite
    ];
  }
  /**
   * Update font size
   */
  setFontSize(A) {
    this.fontSize = A, this.metrics = this.measureFont();
  }
  /**
   * Update font family
   */
  setFontFamily(A) {
    this.fontFamily = A, this.metrics = this.measureFont();
  }
  /**
   * Update cursor style
   */
  setCursorStyle(A) {
    this.cursorStyle = A;
  }
  /**
   * Enable/disable cursor blinking
   */
  setCursorBlink(A) {
    A && !this.cursorBlink ? (this.cursorBlink = !0, this.startCursorBlink()) : !A && this.cursorBlink && (this.cursorBlink = !1, this.stopCursorBlink());
  }
  /**
   * Get current font metrics
   */
  /**
   * Render scrollbar (Phase 2)
   * Shows scroll position and allows click/drag interaction
   * @param opacity Opacity level (0-1) for fade in/out effect
   */
  renderScrollbar(A, B, Q, g = 1) {
    if (g <= 0 || B === 0)
      return;
    const C = this.ctx, I = this.canvas.height / this.devicePixelRatio, w = this.canvas.width / this.devicePixelRatio, o = 8, i = w - o - 4, M = 4, k = I - M * 2, G = B + Q, J = Math.max(20, Q / G * k), s = A / B, F = M + (k - J) * (1 - s);
    C.fillStyle = `rgba(128, 128, 128, ${0.1 * g})`, C.fillRect(i, M, o, k);
    const N = A > 0 ? 0.5 : 0.3;
    C.fillStyle = `rgba(128, 128, 128, ${N * g})`, C.fillRect(i, F, o, J);
  }
  getMetrics() {
    return { ...this.metrics };
  }
  /**
   * Get canvas element (needed by SelectionManager)
   */
  getCanvas() {
    return this.canvas;
  }
  /**
   * Set selection manager (for rendering selection overlay)
   */
  setSelectionManager(A) {
    this.selectionManager = A;
  }
  /**
   * Set the currently hovered hyperlink ID for rendering underlines
   */
  setHoveredHyperlinkId(A) {
    this.hoveredHyperlinkId = A;
  }
  /**
   * Set the currently hovered link range for rendering underlines (for regex-detected URLs)
   * Pass null to clear the hover state
   */
  setHoveredLinkRange(A) {
    this.hoveredLinkRange = A;
  }
  /**
   * Get character cell width (for coordinate conversion)
   */
  get charWidth() {
    return this.metrics.width;
  }
  /**
   * Get character cell height (for coordinate conversion)
   */
  get charHeight() {
    return this.metrics.height;
  }
  /**
   * Clear entire canvas
   */
  clear() {
    this.ctx.fillStyle = this.theme.background, this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
  }
  /**
   * Render selection overlay
   */
  renderSelection(A) {
    const B = this.selectionManager.getSelectionCoords();
    if (!B)
      return;
    const { startCol: Q, startRow: g, endCol: C, endRow: I } = B;
    this.ctx.save(), this.ctx.fillStyle = this.theme.selectionBackground, this.ctx.globalAlpha = 0.5;
    for (let w = g; w <= I; w++) {
      const o = w === g ? Q : 0, i = w === I ? C : A - 1, M = o * this.metrics.width, k = w * this.metrics.height, G = (i - o + 1) * this.metrics.width, J = this.metrics.height;
      this.ctx.fillRect(M, k, G, J);
    }
    this.ctx.restore();
  }
  /**
   * Cleanup resources
   */
  dispose() {
    this.stopCursorBlink();
  }
}
const L = class R {
  // ms between scroll steps
  constructor(A, B, Q, g) {
    this.selectionStart = null, this.selectionEnd = null, this.isSelecting = !1, this.mouseDownTarget = null, this.dirtySelectionRows = /* @__PURE__ */ new Set(), this.selectionChangedEmitter = new K(), this.boundMouseUpHandler = null, this.boundContextMenuHandler = null, this.boundClickHandler = null, this.boundDocumentMouseMoveHandler = null, this.autoScrollInterval = null, this.autoScrollDirection = 0, this.terminal = A, this.renderer = B, this.wasmTerm = Q, this.textarea = g, this.attachEventListeners();
  }
  // pixels from edge to trigger scroll
  /**
   * Get current viewport Y position (how many lines scrolled into history)
   */
  getViewportY() {
    const A = typeof this.terminal.getViewportY == "function" ? this.terminal.getViewportY() : this.terminal.viewportY || 0;
    return Math.max(0, Math.floor(A));
  }
  /**
   * Convert viewport row to absolute buffer row
   * Absolute row is an index into combined buffer: scrollback (0 to len-1) + screen (len to len+rows-1)
   */
  viewportRowToAbsolute(A) {
    const B = this.wasmTerm.getScrollbackLength(), Q = this.getViewportY();
    return B + A - Q;
  }
  /**
   * Convert absolute buffer row to viewport row (may be outside visible range)
   */
  absoluteRowToViewport(A) {
    const B = this.wasmTerm.getScrollbackLength(), Q = this.getViewportY();
    return A - B + Q;
  }
  // ==========================================================================
  // Public API
  // ==========================================================================
  /**
   * Get the selected text as a string
   */
  getSelection() {
    if (!this.selectionStart || !this.selectionEnd)
      return "";
    let { col: A, absoluteRow: B } = this.selectionStart, { col: Q, absoluteRow: g } = this.selectionEnd;
    (B > g || B === g && A > Q) && ([A, Q] = [Q, A], [B, g] = [g, B]);
    const C = this.wasmTerm.getScrollbackLength();
    let I = "";
    for (let w = B; w <= g; w++) {
      let o = null;
      if (w < C)
        o = this.wasmTerm.getScrollbackLine(w);
      else {
        const J = w - C;
        o = this.wasmTerm.getLine(J);
      }
      if (!o)
        continue;
      let i = -1;
      const M = w === B ? A : 0, k = w === g ? Q : o.length - 1;
      let G = "";
      for (let J = M; J <= k; J++) {
        const s = o[J];
        if (s && s.codepoint !== 0) {
          const F = String.fromCodePoint(s.codepoint);
          G += F, F.trim() && (i = G.length);
        } else
          G += " ";
      }
      i >= 0 ? G = G.substring(0, i) : G = "", I += G, w < g && (I += `
`);
    }
    return I;
  }
  /**
   * Check if there's an active selection
   */
  hasSelection() {
    return !this.selectionStart || !this.selectionEnd ? !1 : !(this.selectionStart.col === this.selectionEnd.col && this.selectionStart.absoluteRow === this.selectionEnd.absoluteRow);
  }
  /**
   * Clear the selection
   */
  clearSelection() {
    if (!this.hasSelection())
      return;
    const A = this.normalizeSelection();
    if (A)
      for (let B = A.startRow; B <= A.endRow; B++)
        this.dirtySelectionRows.add(B);
    this.selectionStart = null, this.selectionEnd = null, this.isSelecting = !1, this.requestRender();
  }
  /**
   * Select all text in the terminal
   */
  selectAll() {
    const A = this.wasmTerm.getDimensions(), B = this.getViewportY();
    this.selectionStart = { col: 0, absoluteRow: B }, this.selectionEnd = { col: A.cols - 1, absoluteRow: B + A.rows - 1 }, this.requestRender(), this.selectionChangedEmitter.fire();
  }
  /**
   * Select text at specific column and row with length
   * xterm.js compatible API
   */
  select(A, B, Q) {
    const g = this.wasmTerm.getDimensions();
    B = Math.max(0, Math.min(B, g.rows - 1)), A = Math.max(0, Math.min(A, g.cols - 1));
    let C = B, I = A + Q - 1;
    for (; I >= g.cols; )
      I -= g.cols, C++;
    C = Math.min(C, g.rows - 1);
    const w = this.getViewportY();
    this.selectionStart = { col: A, absoluteRow: w + B }, this.selectionEnd = { col: I, absoluteRow: w + C }, this.requestRender(), this.selectionChangedEmitter.fire();
  }
  /**
   * Select entire lines from start to end
   * xterm.js compatible API
   */
  selectLines(A, B) {
    const Q = this.wasmTerm.getDimensions();
    A = Math.max(0, Math.min(A, Q.rows - 1)), B = Math.max(0, Math.min(B, Q.rows - 1)), A > B && ([A, B] = [B, A]);
    const g = this.getViewportY();
    this.selectionStart = { col: 0, absoluteRow: g + A }, this.selectionEnd = { col: Q.cols - 1, absoluteRow: g + B }, this.requestRender(), this.selectionChangedEmitter.fire();
  }
  /**
   * Get selection position as buffer range
   * xterm.js compatible API
   */
  getSelectionPosition() {
    const A = this.normalizeSelection();
    if (A)
      return {
        start: { x: A.startCol, y: A.startRow },
        end: { x: A.endCol, y: A.endRow }
      };
  }
  /**
   * Deselect all text
   * xterm.js compatible API
   */
  deselect() {
    this.clearSelection(), this.selectionChangedEmitter.fire();
  }
  /**
   * Focus the terminal (make it receive keyboard input)
   */
  focus() {
    const A = this.renderer.getCanvas();
    A.parentElement && A.parentElement.focus();
  }
  /**
   * Get current selection coordinates (for rendering)
   */
  getSelectionCoords() {
    return this.normalizeSelection();
  }
  /**
   * Get dirty selection rows that need redraw (for clearing old highlight)
   */
  getDirtySelectionRows() {
    return this.dirtySelectionRows;
  }
  /**
   * Clear the dirty selection rows tracking (after redraw)
   */
  clearDirtySelectionRows() {
    this.dirtySelectionRows.clear();
  }
  /**
   * Get selection change event accessor
   */
  get onSelectionChange() {
    return this.selectionChangedEmitter.event;
  }
  /**
   * Cleanup resources
   */
  dispose() {
    this.selectionChangedEmitter.dispose(), this.stopAutoScroll(), this.boundMouseUpHandler && (document.removeEventListener("mouseup", this.boundMouseUpHandler), this.boundMouseUpHandler = null), this.boundDocumentMouseMoveHandler && (document.removeEventListener("mousemove", this.boundDocumentMouseMoveHandler), this.boundDocumentMouseMoveHandler = null), this.boundContextMenuHandler && (this.renderer.getCanvas().removeEventListener("contextmenu", this.boundContextMenuHandler), this.boundContextMenuHandler = null), this.boundClickHandler && (document.removeEventListener("click", this.boundClickHandler), this.boundClickHandler = null);
  }
  // ==========================================================================
  // Private Methods
  // ==========================================================================
  /**
   * Attach mouse event listeners to canvas
   */
  attachEventListeners() {
    const A = this.renderer.getCanvas();
    A.addEventListener("mousedown", (B) => {
      if (B.button === 0) {
        A.parentElement && A.parentElement.focus();
        const Q = this.pixelToCell(B.offsetX, B.offsetY);
        this.hasSelection() && this.clearSelection();
        const C = this.viewportRowToAbsolute(Q.row);
        this.selectionStart = { col: Q.col, absoluteRow: C }, this.selectionEnd = { col: Q.col, absoluteRow: C }, this.isSelecting = !0;
      }
    }), A.addEventListener("mousemove", (B) => {
      if (this.isSelecting) {
        this.markCurrentSelectionDirty();
        const Q = this.pixelToCell(B.offsetX, B.offsetY), g = this.viewportRowToAbsolute(Q.row);
        this.selectionEnd = { col: Q.col, absoluteRow: g }, this.requestRender(), this.updateAutoScroll(B.offsetY, A.clientHeight);
      }
    }), A.addEventListener("mouseleave", (B) => {
      if (this.isSelecting) {
        const Q = A.getBoundingClientRect();
        B.clientY < Q.top ? this.startAutoScroll(-1) : B.clientY > Q.bottom && this.startAutoScroll(1);
      }
    }), A.addEventListener("mouseenter", () => {
      this.isSelecting && this.stopAutoScroll();
    }), this.boundDocumentMouseMoveHandler = (B) => {
      if (this.isSelecting) {
        const Q = A.getBoundingClientRect(), g = Math.max(Q.left, Math.min(B.clientX, Q.right)), C = Math.max(Q.top, Math.min(B.clientY, Q.bottom)), I = g - Q.left, w = C - Q.top;
        if ((B.clientX < Q.left || B.clientX > Q.right || B.clientY < Q.top || B.clientY > Q.bottom) && (B.clientY < Q.top ? this.startAutoScroll(-1) : B.clientY > Q.bottom ? this.startAutoScroll(1) : this.stopAutoScroll(), this.autoScrollDirection === 0)) {
          this.markCurrentSelectionDirty();
          const o = this.pixelToCell(I, w), i = this.viewportRowToAbsolute(o.row);
          this.selectionEnd = { col: o.col, absoluteRow: i }, this.requestRender();
        }
      }
    }, document.addEventListener("mousemove", this.boundDocumentMouseMoveHandler), document.addEventListener("mousedown", (B) => {
      this.mouseDownTarget = B.target;
    }), this.boundMouseUpHandler = (B) => {
      if (this.isSelecting) {
        this.isSelecting = !1, this.stopAutoScroll();
        const Q = this.getSelection();
        Q && (this.copyToClipboard(Q), this.selectionChangedEmitter.fire());
      }
    }, document.addEventListener("mouseup", this.boundMouseUpHandler), A.addEventListener("dblclick", (B) => {
      const Q = this.pixelToCell(B.offsetX, B.offsetY), g = this.getWordAtCell(Q.col, Q.row);
      if (g) {
        const C = this.viewportRowToAbsolute(Q.row);
        this.selectionStart = { col: g.startCol, absoluteRow: C }, this.selectionEnd = { col: g.endCol, absoluteRow: C }, this.requestRender();
        const I = this.getSelection();
        I && (this.copyToClipboard(I), this.selectionChangedEmitter.fire());
      }
    }), this.boundContextMenuHandler = (B) => {
      if (this.renderer.getCanvas().getBoundingClientRect(), this.textarea.style.position = "fixed", this.textarea.style.left = `${B.clientX}px`, this.textarea.style.top = `${B.clientY}px`, this.textarea.style.width = "1px", this.textarea.style.height = "1px", this.textarea.style.zIndex = "1000", this.textarea.style.opacity = "0", this.textarea.style.pointerEvents = "auto", this.hasSelection()) {
        const g = this.getSelection();
        this.textarea.value = g, this.textarea.select(), this.textarea.setSelectionRange(0, g.length);
      } else
        this.textarea.value = "";
      this.textarea.focus(), setTimeout(() => {
        const g = () => {
          this.textarea.style.pointerEvents = "none", this.textarea.style.zIndex = "-10", this.textarea.style.width = "0", this.textarea.style.height = "0", this.textarea.style.left = "0", this.textarea.style.top = "0", this.textarea.value = "", document.removeEventListener("click", g), document.removeEventListener("contextmenu", g), this.textarea.removeEventListener("blur", g);
        };
        document.addEventListener("click", g, { once: !0 }), document.addEventListener("contextmenu", g, { once: !0 }), this.textarea.addEventListener("blur", g, { once: !0 });
      }, 10);
    }, A.addEventListener("contextmenu", this.boundContextMenuHandler), this.boundClickHandler = (B) => {
      if (this.isSelecting || this.mouseDownTarget && A.contains(this.mouseDownTarget))
        return;
      const g = B.target;
      A.contains(g) || this.hasSelection() && this.clearSelection();
    }, document.addEventListener("click", this.boundClickHandler);
  }
  /**
   * Mark current selection rows as dirty for redraw
   */
  markCurrentSelectionDirty() {
    const A = this.normalizeSelection();
    if (A)
      for (let B = A.startRow; B <= A.endRow; B++)
        this.dirtySelectionRows.add(B);
  }
  /**
   * Update auto-scroll based on mouse Y position within canvas
   */
  updateAutoScroll(A, B) {
    const Q = R.AUTO_SCROLL_EDGE_SIZE;
    A < Q ? this.startAutoScroll(-1) : A > B - Q ? this.startAutoScroll(1) : this.stopAutoScroll();
  }
  /**
   * Start auto-scrolling in the given direction
   */
  startAutoScroll(A) {
    this.autoScrollInterval !== null && this.autoScrollDirection === A || (this.stopAutoScroll(), this.autoScrollDirection = A, this.autoScrollInterval = setInterval(() => {
      if (!this.isSelecting) {
        this.stopAutoScroll();
        return;
      }
      const B = R.AUTO_SCROLL_SPEED * this.autoScrollDirection;
      if (this.terminal.scrollLines(B), this.selectionEnd) {
        const Q = this.wasmTerm.getDimensions(), g = this.getViewportY();
        if (this.autoScrollDirection < 0) {
          const C = g;
          C < this.selectionEnd.absoluteRow && (this.selectionEnd = { col: 0, absoluteRow: C });
        } else {
          const C = g + Q.rows - 1;
          C > this.selectionEnd.absoluteRow && (this.selectionEnd = { col: Q.cols - 1, absoluteRow: C });
        }
      }
      this.requestRender();
    }, R.AUTO_SCROLL_INTERVAL));
  }
  /**
   * Stop auto-scrolling
   */
  stopAutoScroll() {
    this.autoScrollInterval !== null && (clearInterval(this.autoScrollInterval), this.autoScrollInterval = null), this.autoScrollDirection = 0;
  }
  /**
   * Convert pixel coordinates to terminal cell coordinates
   */
  pixelToCell(A, B) {
    const Q = this.renderer.getMetrics(), g = Math.floor(A / Q.width), C = Math.floor(B / Q.height);
    return {
      col: Math.max(0, Math.min(g, this.terminal.cols - 1)),
      row: Math.max(0, Math.min(C, this.terminal.rows - 1))
    };
  }
  /**
   * Normalize selection coordinates (handle backward selection)
   * Returns coordinates in VIEWPORT space for rendering, clamped to visible area
   */
  normalizeSelection() {
    if (!this.selectionStart || !this.selectionEnd)
      return null;
    let { col: A, absoluteRow: B } = this.selectionStart, { col: Q, absoluteRow: g } = this.selectionEnd;
    (B > g || B === g && A > Q) && ([A, Q] = [Q, A], [B, g] = [g, B]);
    let C = this.absoluteRowToViewport(B), I = this.absoluteRowToViewport(g);
    const w = this.wasmTerm.getDimensions(), o = w.rows - 1;
    return I < 0 || C > o ? null : (C < 0 && (C = 0, A = 0), I > o && (I = o, Q = w.cols - 1), { startCol: A, startRow: C, endCol: Q, endRow: I });
  }
  /**
   * Get word boundaries at a cell position
   */
  getWordAtCell(A, B) {
    const Q = this.wasmTerm.getLine(B);
    if (!Q)
      return null;
    const g = (w) => {
      if (!w || w.codepoint === 0)
        return !1;
      const o = String.fromCodePoint(w.codepoint);
      return /[\w-]/.test(o);
    };
    if (!g(Q[A]))
      return null;
    let C = A;
    for (; C > 0 && g(Q[C - 1]); )
      C--;
    let I = A;
    for (; I < Q.length - 1 && g(Q[I + 1]); )
      I++;
    return { startCol: C, endCol: I };
  }
  /**
   * Copy text to clipboard
   */
  async copyToClipboard(A) {
    if (navigator.clipboard && navigator.clipboard.writeText)
      try {
        await navigator.clipboard.writeText(A);
        return;
      } catch {
      }
    const B = document.activeElement;
    try {
      const Q = this.textarea;
      Q.value = A, Q.style.position = "fixed", Q.style.left = "-9999px", Q.style.top = "0", Q.style.width = "1px", Q.style.height = "1px", Q.style.opacity = "0", Q.focus(), Q.select(), Q.setSelectionRange(0, A.length);
      const g = document.execCommand("copy");
      B && B.focus(), g || console.error("❌ execCommand copy failed");
    } catch (Q) {
      console.error("❌ Fallback copy failed:", Q), B && B.focus();
    }
  }
  /**
   * Request a render update (triggers selection overlay redraw)
   */
  requestRender() {
  }
};
L.AUTO_SCROLL_EDGE_SIZE = 30;
L.AUTO_SCROLL_SPEED = 3;
L.AUTO_SCROLL_INTERVAL = 50;
let $ = L;
class CA {
  // 200ms fade animation
  constructor(A = {}) {
    this.unicode = {
      get activeVersion() {
        return "15.1";
      }
    }, this.dataEmitter = new K(), this.resizeEmitter = new K(), this.bellEmitter = new K(), this.selectionChangeEmitter = new K(), this.keyEmitter = new K(), this.titleChangeEmitter = new K(), this.scrollEmitter = new K(), this.renderEmitter = new K(), this.cursorMoveEmitter = new K(), this.onData = this.dataEmitter.event, this.onResize = this.resizeEmitter.event, this.onBell = this.bellEmitter.event, this.onSelectionChange = this.selectionChangeEmitter.event, this.onKey = this.keyEmitter.event, this.onTitleChange = this.titleChangeEmitter.event, this.onScroll = this.scrollEmitter.event, this.onRender = this.renderEmitter.event, this.onCursorMove = this.cursorMoveEmitter.event, this.isOpen = !1, this.isDisposed = !1, this.addons = [], this.currentTitle = "", this.viewportY = 0, this.targetViewportY = 0, this.lastCursorY = 0, this.isDraggingScrollbar = !1, this.scrollbarDragStart = null, this.scrollbarDragStartViewportY = 0, this.scrollbarVisible = !1, this.scrollbarOpacity = 0, this.SCROLLBAR_HIDE_DELAY_MS = 1500, this.SCROLLBAR_FADE_DURATION_MS = 200, this.animateScroll = () => {
      if (!this.wasmTerm || this.scrollAnimationStartTime === void 0)
        return;
      const Q = this.options.smoothScrollDuration ?? 100, g = this.targetViewportY - this.viewportY;
      if (Math.abs(g) < 0.01) {
        this.viewportY = this.targetViewportY, this.scrollEmitter.fire(Math.floor(this.viewportY)), this.getScrollbackLength() > 0 && this.showScrollbar(), this.scrollAnimationFrame = void 0, this.scrollAnimationStartTime = void 0, this.scrollAnimationStartY = void 0;
        return;
      }
      const w = 1 - (1 / (Q / 1e3 * 60)) ** 2;
      this.viewportY += g * w;
      const o = Math.floor(this.viewportY);
      this.scrollEmitter.fire(o), this.getScrollbackLength() > 0 && this.showScrollbar(), this.scrollAnimationFrame = requestAnimationFrame(this.animateScroll);
    }, this.handleMouseMove = (Q) => {
      if (!(!this.canvas || !this.renderer || !this.wasmTerm)) {
        if (this.isDraggingScrollbar) {
          this.processScrollbarDrag(Q);
          return;
        }
        if (this.linkDetector) {
          if (this.mouseMoveThrottleTimeout) {
            this.pendingMouseMove = Q;
            return;
          }
          this.processMouseMove(Q), this.mouseMoveThrottleTimeout = window.setTimeout(() => {
            if (this.mouseMoveThrottleTimeout = void 0, this.pendingMouseMove) {
              const g = this.pendingMouseMove;
              this.pendingMouseMove = void 0, this.processMouseMove(g);
            }
          }, 16);
        }
      }
    }, this.handleMouseLeave = () => {
      var Q, g;
      this.renderer && this.wasmTerm && ((this.renderer.hoveredHyperlinkId || 0) > 0 && this.renderer.setHoveredHyperlinkId(0), this.renderer.setHoveredLinkRange(null)), this.currentHoveredLink && ((g = (Q = this.currentHoveredLink).hover) == null || g.call(Q, !1), this.currentHoveredLink = void 0, this.element && (this.element.style.cursor = "text"));
    }, this.handleClick = async (Q) => {
      if (!this.canvas || !this.renderer || !this.linkDetector || !this.wasmTerm)
        return;
      const g = this.canvas.getBoundingClientRect(), C = Math.floor((Q.clientX - g.left) / this.renderer.charWidth), w = Math.floor((Q.clientY - g.top) / this.renderer.charHeight), o = this.wasmTerm.getScrollbackLength();
      let i;
      const M = this.getViewportY(), k = Math.max(0, Math.floor(M));
      if (k > 0)
        if (w < k)
          i = o - k + w;
        else {
          const J = w - k;
          i = o + J;
        }
      else
        i = o + w;
      const G = await this.linkDetector.getLinkAt(C, i);
      G && (G.activate(Q), (Q.ctrlKey || Q.metaKey) && Q.preventDefault());
    }, this.handleWheel = (Q) => {
      var C, I, w;
      if (Q.preventDefault(), Q.stopPropagation(), this.customWheelEventHandler && this.customWheelEventHandler(Q))
        return;
      if (((C = this.wasmTerm) == null ? void 0 : C.isAlternateScreen()) ?? !1) {
        const o = Q.deltaY > 0 ? "down" : "up", i = Math.min(Math.abs(Math.round(Q.deltaY / 33)), 5);
        for (let M = 0; M < i; M++)
          o === "up" ? this.dataEmitter.fire("\x1B[A") : this.dataEmitter.fire("\x1B[B");
      } else {
        let o;
        if (Q.deltaMode === WheelEvent.DOM_DELTA_PIXEL) {
          const i = ((w = (I = this.renderer) == null ? void 0 : I.getMetrics()) == null ? void 0 : w.height) ?? 20;
          o = Q.deltaY / i;
        } else
          Q.deltaMode === WheelEvent.DOM_DELTA_LINE ? o = Q.deltaY : Q.deltaMode === WheelEvent.DOM_DELTA_PAGE ? o = Q.deltaY * this.rows : o = Q.deltaY / 33;
        if (o !== 0) {
          const i = this.viewportY - o;
          this.smoothScrollTo(i);
        }
      }
    }, this.handleMouseDown = (Q) => {
      if (!this.canvas || !this.renderer || !this.wasmTerm)
        return;
      const g = this.wasmTerm.getScrollbackLength();
      if (g === 0)
        return;
      const C = this.canvas.getBoundingClientRect(), I = Q.clientX - C.left, w = Q.clientY - C.top, o = C.width, i = C.height, M = 8, k = o - M - 4, G = 4;
      if (I >= k && I <= k + M) {
        Q.preventDefault(), Q.stopPropagation(), Q.stopImmediatePropagation();
        const J = i - G * 2, s = this.rows, F = g + s, a = Math.max(20, s / F * J), N = this.viewportY / g, c = G + (J - a) * (1 - N);
        if (w >= c && w <= c + a)
          this.isDraggingScrollbar = !0, this.scrollbarDragStart = w, this.scrollbarDragStartViewportY = this.viewportY, this.canvas && (this.canvas.style.userSelect = "none", this.canvas.style.webkitUserSelect = "none");
        else {
          const y = 1 - (w - G) / J, O = Math.round(y * g);
          this.scrollToLine(Math.max(0, Math.min(g, O)));
        }
      }
    }, this.handleMouseUp = () => {
      this.isDraggingScrollbar && (this.isDraggingScrollbar = !1, this.scrollbarDragStart = null, this.canvas && (this.canvas.style.userSelect = "", this.canvas.style.webkitUserSelect = ""), this.scrollbarVisible && this.getScrollbackLength() > 0 && this.showScrollbar());
    }, this.ghostty = A.ghostty ?? gA();
    const B = {
      cols: A.cols ?? 80,
      rows: A.rows ?? 24,
      cursorBlink: A.cursorBlink ?? !1,
      cursorStyle: A.cursorStyle ?? "block",
      theme: A.theme ?? {},
      scrollback: A.scrollback ?? 1e3,
      fontSize: A.fontSize ?? 15,
      fontFamily: A.fontFamily ?? "monospace",
      allowTransparency: A.allowTransparency ?? !1,
      convertEol: A.convertEol ?? !1,
      disableStdin: A.disableStdin ?? !1,
      smoothScrollDuration: A.smoothScrollDuration ?? 100
      // Default: 100ms smooth scroll
    };
    this.options = new Proxy(B, {
      set: (Q, g, C) => {
        const I = Q[g];
        return Q[g] = C, this.isOpen && this.handleOptionChange(g, C, I), !0;
      }
    }), this.cols = this.options.cols, this.rows = this.options.rows, this.buffer = new V(this);
  }
  // ==========================================================================
  // Theme to WASM Config Conversion
  // ==========================================================================
  /**
   * Parse a CSS color string to 0xRRGGBB format.
   * Returns 0 if the color is undefined or invalid.
   */
  parseColorToHex(A) {
    if (!A)
      return 0;
    if (A.startsWith("#")) {
      let Q = A.slice(1);
      Q.length === 3 && (Q = Q[0] + Q[0] + Q[1] + Q[1] + Q[2] + Q[2]);
      const g = Number.parseInt(Q, 16);
      return Number.isNaN(g) ? 0 : g;
    }
    const B = A.match(/rgb\((\d+),\s*(\d+),\s*(\d+)\)/);
    if (B) {
      const Q = Number.parseInt(B[1], 10), g = Number.parseInt(B[2], 10), C = Number.parseInt(B[3], 10);
      return Q << 16 | g << 8 | C;
    }
    return 0;
  }
  /**
   * Convert terminal options to WASM terminal config.
   */
  buildWasmConfig() {
    const A = this.options.theme, B = this.options.scrollback;
    if (!A && B === 1e3)
      return;
    const Q = [
      this.parseColorToHex(A == null ? void 0 : A.black),
      this.parseColorToHex(A == null ? void 0 : A.red),
      this.parseColorToHex(A == null ? void 0 : A.green),
      this.parseColorToHex(A == null ? void 0 : A.yellow),
      this.parseColorToHex(A == null ? void 0 : A.blue),
      this.parseColorToHex(A == null ? void 0 : A.magenta),
      this.parseColorToHex(A == null ? void 0 : A.cyan),
      this.parseColorToHex(A == null ? void 0 : A.white),
      this.parseColorToHex(A == null ? void 0 : A.brightBlack),
      this.parseColorToHex(A == null ? void 0 : A.brightRed),
      this.parseColorToHex(A == null ? void 0 : A.brightGreen),
      this.parseColorToHex(A == null ? void 0 : A.brightYellow),
      this.parseColorToHex(A == null ? void 0 : A.brightBlue),
      this.parseColorToHex(A == null ? void 0 : A.brightMagenta),
      this.parseColorToHex(A == null ? void 0 : A.brightCyan),
      this.parseColorToHex(A == null ? void 0 : A.brightWhite)
    ];
    return {
      scrollbackLimit: B,
      fgColor: this.parseColorToHex(A == null ? void 0 : A.foreground),
      bgColor: this.parseColorToHex(A == null ? void 0 : A.background),
      cursorColor: this.parseColorToHex(A == null ? void 0 : A.cursor),
      palette: Q
    };
  }
  // ==========================================================================
  // Option Change Handling (for mutable options)
  // ==========================================================================
  /**
   * Handle runtime option changes (called when options are modified after terminal is open)
   * This enables xterm.js compatibility where options can be changed at runtime
   */
  handleOptionChange(A, B, Q) {
    if (B !== Q)
      switch (A) {
        case "disableStdin":
          break;
        case "cursorBlink":
        case "cursorStyle":
          this.renderer && (this.renderer.setCursorStyle(this.options.cursorStyle), this.renderer.setCursorBlink(this.options.cursorBlink));
          break;
        case "theme":
          this.renderer && console.warn("ghostty-web: theme changes after open() are not yet fully supported");
          break;
        case "fontSize":
        case "fontFamily":
          this.renderer && console.warn("ghostty-web: font changes after open() are not yet fully supported");
          break;
        case "cols":
        case "rows":
          this.resize(this.options.cols, this.options.rows);
          break;
      }
  }
  // ==========================================================================
  // Lifecycle Methods
  // ==========================================================================
  /**
   * Open terminal in a parent element
   *
   * Initializes all components and starts rendering.
   * Requires a pre-loaded Ghostty instance passed to the constructor.
   */
  open(A) {
    if (this.isOpen)
      throw new Error("Terminal is already open");
    if (this.isDisposed)
      throw new Error("Terminal has been disposed");
    this.element = A, this.isOpen = !0;
    try {
      A.hasAttribute("tabindex") || A.setAttribute("tabindex", "0");
      const B = this.buildWasmConfig();
      this.wasmTerm = this.ghostty.createTerminal(this.cols, this.rows, B), this.canvas = document.createElement("canvas"), this.canvas.style.display = "block", A.appendChild(this.canvas), this.textarea = document.createElement("textarea"), this.textarea.setAttribute("autocorrect", "off"), this.textarea.setAttribute("autocapitalize", "off"), this.textarea.setAttribute("spellcheck", "false"), this.textarea.setAttribute("tabindex", "-1"), this.textarea.setAttribute("aria-label", "Terminal input"), this.textarea.style.position = "absolute", this.textarea.style.left = "0", this.textarea.style.top = "0", this.textarea.style.width = "0", this.textarea.style.height = "0", this.textarea.style.zIndex = "-10", this.textarea.style.opacity = "0", this.textarea.style.overflow = "hidden", this.textarea.style.pointerEvents = "none", this.textarea.style.resize = "none", this.textarea.style.border = "none", this.textarea.style.outline = "none", A.appendChild(this.textarea), this.renderer = new _(this.canvas, {
        fontSize: this.options.fontSize,
        fontFamily: this.options.fontFamily,
        cursorStyle: this.options.cursorStyle,
        cursorBlink: this.options.cursorBlink,
        theme: this.options.theme
      }), this.renderer.resize(this.cols, this.rows), this.inputHandler = new X(
        this.ghostty,
        A,
        (Q) => {
          this.options.disableStdin || this.dataEmitter.fire(Q);
        },
        () => {
          this.bellEmitter.fire();
        },
        (Q) => {
          this.keyEmitter.fire(Q);
        },
        this.customKeyEventHandler
      ), this.selectionManager = new $(
        this,
        this.renderer,
        this.wasmTerm,
        this.textarea
      ), this.renderer.setSelectionManager(this.selectionManager), this.selectionManager.onSelectionChange(() => {
        this.selectionChangeEmitter.fire();
      }), this.textarea.addEventListener("paste", (Q) => {
        var C;
        Q.preventDefault(), Q.stopPropagation();
        const g = (C = Q.clipboardData) == null ? void 0 : C.getData("text");
        g && this.paste(g);
      }), this.linkDetector = new P(this), this.linkDetector.registerProvider(new v(this)), this.linkDetector.registerProvider(new u(this)), A.addEventListener("mousedown", this.handleMouseDown, { capture: !0 }), A.addEventListener("mousemove", this.handleMouseMove), A.addEventListener("mouseleave", this.handleMouseLeave), A.addEventListener("click", this.handleClick), document.addEventListener("mouseup", this.handleMouseUp), A.addEventListener("wheel", this.handleWheel, { passive: !1, capture: !0 }), this.renderer.render(this.wasmTerm, !0, this.viewportY, this, this.scrollbarOpacity), this.startRenderLoop(), this.focus();
    } catch (B) {
      throw this.isOpen = !1, this.cleanupComponents(), new Error(`Failed to open terminal: ${B}`);
    }
  }
  /**
   * Write data to terminal
   */
  write(A, B) {
    this.assertOpen(), this.options.convertEol && typeof A == "string" && (A = A.replace(/\n/g, `\r
`)), this.writeInternal(A, B);
  }
  /**
   * Internal write implementation (extracted from write())
   */
  writeInternal(A, B) {
    var Q;
    this.wasmTerm.write(A), typeof A == "string" && A.includes("\x07") ? this.bellEmitter.fire() : A instanceof Uint8Array && A.includes(7) && this.bellEmitter.fire(), (Q = this.linkDetector) == null || Q.invalidateCache(), this.viewportY !== 0 && this.scrollToBottom(), typeof A == "string" && A.includes("\x1B]") && this.checkForTitleChange(A), B && requestAnimationFrame(B);
  }
  /**
   * Write data with newline
   */
  writeln(A, B) {
    if (typeof A == "string")
      this.write(A + `\r
`, B);
    else {
      const Q = new Uint8Array(A.length + 2);
      Q.set(A), Q[A.length] = 13, Q[A.length + 1] = 10, this.write(Q, B);
    }
  }
  /**
   * Paste text into terminal (triggers bracketed paste if supported)
   */
  paste(A) {
    this.assertOpen(), !this.options.disableStdin && (this.wasmTerm.hasBracketedPaste() ? this.dataEmitter.fire("\x1B[200~" + A + "\x1B[201~") : this.dataEmitter.fire(A));
  }
  /**
   * Input data into terminal (as if typed by user)
   *
   * @param data - Data to input
   * @param wasUserInput - If true, triggers onData event (default: false for compat with some apps)
   */
  input(A, B = !1) {
    this.assertOpen(), !this.options.disableStdin && (B ? this.dataEmitter.fire(A) : this.write(A));
  }
  /**
   * Resize terminal
   */
  resize(A, B) {
    if (this.assertOpen(), A === this.cols && B === this.rows)
      return;
    this.cols = A, this.rows = B, this.wasmTerm.resize(A, B), this.renderer.resize(A, B);
    const Q = this.renderer.getMetrics();
    this.canvas.width = Q.width * A, this.canvas.height = Q.height * B, this.canvas.style.width = `${Q.width * A}px`, this.canvas.style.height = `${Q.height * B}px`, this.resizeEmitter.fire({ cols: A, rows: B }), this.renderer.render(this.wasmTerm, !0, this.viewportY, this);
  }
  /**
   * Clear terminal screen
   */
  clear() {
    this.assertOpen(), this.wasmTerm.write("\x1B[2J\x1B[H");
  }
  /**
   * Reset terminal state
   */
  reset() {
    this.assertOpen(), this.wasmTerm && this.wasmTerm.free();
    const A = this.buildWasmConfig();
    this.wasmTerm = this.ghostty.createTerminal(this.cols, this.rows, A), this.renderer.clear(), this.currentTitle = "";
  }
  /**
   * Focus terminal input
   */
  focus() {
    this.isOpen && this.element && (this.element.focus(), setTimeout(() => {
      var A;
      (A = this.element) == null || A.focus();
    }, 0));
  }
  /**
   * Blur terminal (remove focus)
   */
  blur() {
    this.isOpen && this.element && this.element.blur();
  }
  /**
   * Load an addon
   */
  loadAddon(A) {
    A.activate(this), this.addons.push(A);
  }
  // ==========================================================================
  // Selection API (xterm.js compatible)
  // ==========================================================================
  /**
   * Get the selected text as a string
   */
  getSelection() {
    var A;
    return ((A = this.selectionManager) == null ? void 0 : A.getSelection()) || "";
  }
  /**
   * Check if there's an active selection
   */
  hasSelection() {
    var A;
    return ((A = this.selectionManager) == null ? void 0 : A.hasSelection()) || !1;
  }
  /**
   * Clear the current selection
   */
  clearSelection() {
    var A;
    (A = this.selectionManager) == null || A.clearSelection();
  }
  /**
   * Select all text in the terminal
   */
  selectAll() {
    var A;
    (A = this.selectionManager) == null || A.selectAll();
  }
  /**
   * Select text at specific column and row with length
   */
  select(A, B, Q) {
    var g;
    (g = this.selectionManager) == null || g.select(A, B, Q);
  }
  /**
   * Select entire lines from start to end
   */
  selectLines(A, B) {
    var Q;
    (Q = this.selectionManager) == null || Q.selectLines(A, B);
  }
  /**
   * Get selection position as buffer range
   */
  /**
   * Get the current viewport Y position.
   *
   * This is the number of lines scrolled back from the bottom of the
   * scrollback buffer. It may be fractional during smooth scrolling.
   */
  getViewportY() {
    return this.viewportY;
  }
  getSelectionPosition() {
    var A;
    return (A = this.selectionManager) == null ? void 0 : A.getSelectionPosition();
  }
  // ==========================================================================
  // Phase 1: Custom Event Handlers
  // ==========================================================================
  /**
   * Attach a custom keyboard event handler
   * Returns true to prevent default handling
   */
  attachCustomKeyEventHandler(A) {
    this.customKeyEventHandler = A, this.inputHandler && this.inputHandler.setCustomKeyEventHandler(A);
  }
  /**
   * Attach a custom wheel event handler (Phase 2)
   * Returns true to prevent default handling
   */
  attachCustomWheelEventHandler(A) {
    this.customWheelEventHandler = A;
  }
  // ==========================================================================
  // Link Detection Methods
  // ==========================================================================
  /**
   * Register a custom link provider
   * Multiple providers can be registered to detect different types of links
   *
   * @example
   * ```typescript
   * term.registerLinkProvider({
   *   provideLinks(y, callback) {
   *     // Detect URLs, file paths, etc.
   *     callback(detectedLinks);
   *   }
   * });
   * ```
   */
  registerLinkProvider(A) {
    if (!this.linkDetector)
      throw new Error("Terminal must be opened before registering link providers");
    this.linkDetector.registerProvider(A);
  }
  // ==========================================================================
  // Phase 2: Scrolling Methods
  // ==========================================================================
  /**
   * Scroll viewport by a number of lines
   * @param amount Number of lines to scroll (positive = down, negative = up)
   */
  scrollLines(A) {
    if (!this.wasmTerm)
      throw new Error("Terminal not open");
    const B = this.getScrollbackLength(), g = Math.max(0, Math.min(B, this.viewportY - A));
    g !== this.viewportY && (this.viewportY = g, this.scrollEmitter.fire(this.viewportY), B > 0 && this.showScrollbar());
  }
  /**
   * Scroll viewport by a number of pages
   * @param amount Number of pages to scroll (positive = down, negative = up)
   */
  scrollPages(A) {
    this.scrollLines(A * this.rows);
  }
  /**
   * Scroll viewport to the top of the scrollback buffer
   */
  scrollToTop() {
    const A = this.getScrollbackLength();
    A > 0 && this.viewportY !== A && (this.viewportY = A, this.scrollEmitter.fire(this.viewportY), this.showScrollbar());
  }
  /**
   * Scroll viewport to the bottom (current output)
   */
  scrollToBottom() {
    this.viewportY !== 0 && (this.viewportY = 0, this.scrollEmitter.fire(this.viewportY), this.getScrollbackLength() > 0 && this.showScrollbar());
  }
  /**
   * Scroll viewport to a specific line in the buffer
   * @param line Line number (0 = top of scrollback, scrollbackLength = bottom)
   */
  scrollToLine(A) {
    const B = this.getScrollbackLength(), Q = Math.max(0, Math.min(B, A));
    Q !== this.viewportY && (this.viewportY = Q, this.scrollEmitter.fire(this.viewportY), B > 0 && this.showScrollbar());
  }
  /**
   * Smoothly scroll to a target viewport position
   * @param targetY Target viewport Y position (in lines, can be fractional)
   */
  smoothScrollTo(A) {
    if (!this.wasmTerm)
      return;
    const B = this.getScrollbackLength(), g = Math.max(0, Math.min(B, A));
    if ((this.options.smoothScrollDuration ?? 100) === 0) {
      this.viewportY = g, this.targetViewportY = g, this.scrollEmitter.fire(Math.floor(this.viewportY)), B > 0 && this.showScrollbar();
      return;
    }
    this.targetViewportY = g, !this.scrollAnimationFrame && (this.scrollAnimationStartTime = Date.now(), this.scrollAnimationStartY = this.viewportY, this.animateScroll());
  }
  // ==========================================================================
  // Lifecycle
  // ==========================================================================
  /**
   * Dispose terminal and clean up resources
   */
  dispose() {
    if (!this.isDisposed) {
      this.isDisposed = !0, this.isOpen = !1, this.animationFrameId && (cancelAnimationFrame(this.animationFrameId), this.animationFrameId = void 0), this.scrollAnimationFrame && (cancelAnimationFrame(this.scrollAnimationFrame), this.scrollAnimationFrame = void 0), this.mouseMoveThrottleTimeout && (clearTimeout(this.mouseMoveThrottleTimeout), this.mouseMoveThrottleTimeout = void 0), this.pendingMouseMove = void 0;
      for (const A of this.addons)
        A.dispose();
      this.addons = [], this.cleanupComponents(), this.dataEmitter.dispose(), this.resizeEmitter.dispose(), this.bellEmitter.dispose(), this.selectionChangeEmitter.dispose(), this.keyEmitter.dispose(), this.titleChangeEmitter.dispose(), this.scrollEmitter.dispose(), this.renderEmitter.dispose(), this.cursorMoveEmitter.dispose();
    }
  }
  // ==========================================================================
  // Private Methods
  // ==========================================================================
  /**
   * Start the render loop
   */
  startRenderLoop() {
    const A = () => {
      if (!this.isDisposed && this.isOpen) {
        const B = this.wasmTerm.getCursor();
        B.y !== this.lastCursorY && (this.lastCursorY = B.y, this.cursorMoveEmitter.fire()), this.renderer.render(this.wasmTerm, !1, this.viewportY, this, this.scrollbarOpacity), this.animationFrameId = requestAnimationFrame(A);
      }
    };
    A();
  }
  /**
   * Get a line from native WASM scrollback buffer
   * Implements IScrollbackProvider
   */
  getScrollbackLine(A) {
    return this.wasmTerm ? this.wasmTerm.getScrollbackLine(A) : null;
  }
  /**
   * Get scrollback length from native WASM
   * Implements IScrollbackProvider
   */
  getScrollbackLength() {
    return this.wasmTerm ? this.wasmTerm.getScrollbackLength() : 0;
  }
  /**
   * Clean up components (called on dispose or error)
   */
  cleanupComponents() {
    this.selectionManager && (this.selectionManager.dispose(), this.selectionManager = void 0), this.inputHandler && (this.inputHandler.dispose(), this.inputHandler = void 0), this.renderer && (this.renderer.dispose(), this.renderer = void 0), this.canvas && this.canvas.parentNode && (this.canvas.parentNode.removeChild(this.canvas), this.canvas = void 0), this.textarea && this.textarea.parentNode && (this.textarea.parentNode.removeChild(this.textarea), this.textarea = void 0), this.element && (this.element.removeEventListener("wheel", this.handleWheel), this.element.removeEventListener("mousedown", this.handleMouseDown, { capture: !0 }), this.element.removeEventListener("mousemove", this.handleMouseMove), this.element.removeEventListener("mouseleave", this.handleMouseLeave), this.element.removeEventListener("click", this.handleClick)), this.isOpen && typeof document < "u" && document.removeEventListener("mouseup", this.handleMouseUp), this.scrollbarHideTimeout && (window.clearTimeout(this.scrollbarHideTimeout), this.scrollbarHideTimeout = void 0), this.linkDetector && (this.linkDetector.dispose(), this.linkDetector = void 0), this.wasmTerm && (this.wasmTerm.free(), this.wasmTerm = void 0), this.ghostty = void 0, this.element = void 0, this.textarea = void 0;
  }
  /**
   * Assert terminal is open (throw if not)
   */
  assertOpen() {
    if (this.isDisposed)
      throw new Error("Terminal has been disposed");
    if (!this.isOpen)
      throw new Error("Terminal must be opened before use. Call terminal.open(parent) first.");
  }
  /**
   * Process mouse move for link detection (internal, called by throttled handler)
   */
  processMouseMove(A) {
    if (!this.canvas || !this.renderer || !this.linkDetector || !this.wasmTerm)
      return;
    const B = this.canvas.getBoundingClientRect(), Q = Math.floor((A.clientX - B.left) / this.renderer.charWidth), C = Math.floor((A.clientY - B.top) / this.renderer.charHeight);
    let I = 0, w = null;
    const o = this.getViewportY(), i = Math.max(0, Math.floor(o));
    if (i > 0) {
      const F = this.wasmTerm.getScrollbackLength();
      if (C < i) {
        const a = F - i + C;
        w = this.wasmTerm.getScrollbackLine(a);
      } else {
        const a = C - i;
        w = this.wasmTerm.getLine(a);
      }
    } else
      w = this.wasmTerm.getLine(C);
    w && Q >= 0 && Q < w.length && (I = w[Q].hyperlink_id);
    const M = this.renderer.hoveredHyperlinkId || 0;
    I !== M && this.renderer.setHoveredHyperlinkId(I);
    const k = this.wasmTerm.getScrollbackLength();
    let G;
    const J = this.getViewportY(), s = Math.max(0, Math.floor(J));
    if (s > 0)
      if (C < s)
        G = k - s + C;
      else {
        const F = C - s;
        G = k + F;
      }
    else
      G = k + C;
    this.linkDetector.getLinkAt(Q, G).then((F) => {
      var a, N, c, h;
      if (F !== this.currentHoveredLink && ((N = (a = this.currentHoveredLink) == null ? void 0 : a.hover) == null || N.call(a, !1), this.currentHoveredLink = F, (c = F == null ? void 0 : F.hover) == null || c.call(F, !0), this.element && (this.element.style.cursor = F ? "pointer" : "text"), this.renderer))
        if (F) {
          const y = ((h = this.wasmTerm) == null ? void 0 : h.getScrollbackLength()) || 0, O = this.getViewportY(), p = Math.max(0, Math.floor(O)), T = F.range.start.y - y + p, l = F.range.end.y - y + p;
          T < this.rows && l >= 0 ? this.renderer.setHoveredLinkRange({
            startX: F.range.start.x,
            startY: Math.max(0, T),
            endX: F.range.end.x,
            endY: Math.min(this.rows - 1, l)
          }) : this.renderer.setHoveredLinkRange(null);
        } else
          this.renderer.setHoveredLinkRange(null);
    }).catch((F) => {
      console.warn("Link detection error:", F);
    });
  }
  /**
   * Process scrollbar drag movement
   */
  processScrollbarDrag(A) {
    if (!this.canvas || !this.renderer || !this.wasmTerm || this.scrollbarDragStart === null)
      return;
    const B = this.wasmTerm.getScrollbackLength();
    if (B === 0)
      return;
    const Q = this.canvas.getBoundingClientRect(), C = A.clientY - Q.top - this.scrollbarDragStart, o = Q.height - 4 * 2, i = this.rows, M = B + i, k = Math.max(20, i / M * o), G = -C / (o - k), J = Math.round(G * B), s = this.scrollbarDragStartViewportY + J;
    this.scrollToLine(Math.max(0, Math.min(B, s)));
  }
  /**
   * Show scrollbar with fade-in and schedule auto-hide
   */
  showScrollbar() {
    this.scrollbarHideTimeout && (window.clearTimeout(this.scrollbarHideTimeout), this.scrollbarHideTimeout = void 0), this.scrollbarVisible ? this.scrollbarOpacity = 1 : (this.scrollbarVisible = !0, this.scrollbarOpacity = 0, this.fadeInScrollbar()), this.isDraggingScrollbar || (this.scrollbarHideTimeout = window.setTimeout(() => {
      this.hideScrollbar();
    }, this.SCROLLBAR_HIDE_DELAY_MS));
  }
  /**
   * Hide scrollbar with fade-out
   */
  hideScrollbar() {
    this.scrollbarHideTimeout && (window.clearTimeout(this.scrollbarHideTimeout), this.scrollbarHideTimeout = void 0), this.scrollbarVisible && this.fadeOutScrollbar();
  }
  /**
   * Fade in scrollbar
   */
  fadeInScrollbar() {
    const A = Date.now(), B = () => {
      const Q = Date.now() - A, g = Math.min(Q / this.SCROLLBAR_FADE_DURATION_MS, 1);
      this.scrollbarOpacity = g, g < 1 && requestAnimationFrame(B);
    };
    B();
  }
  /**
   * Fade out scrollbar
   */
  fadeOutScrollbar() {
    const A = Date.now(), B = this.scrollbarOpacity, Q = () => {
      const g = Date.now() - A, C = Math.min(g / this.SCROLLBAR_FADE_DURATION_MS, 1);
      this.scrollbarOpacity = B * (1 - C), C < 1 ? requestAnimationFrame(Q) : (this.scrollbarVisible = !1, this.scrollbarOpacity = 0);
    };
    Q();
  }
  /**
   * Check for title changes in written data (OSC sequences)
   * Simplified implementation - looks for OSC 0, 1, 2
   */
  checkForTitleChange(A) {
    const B = /\x1b\]([012]);([^\x07\x1b]*?)(?:\x07|\x1b\\)/g;
    let Q = null;
    for (; (Q = B.exec(A)) !== null; ) {
      const g = Q[1], C = Q[2];
      (g === "0" || g === "2") && C !== this.currentTitle && (this.currentTitle = C, this.titleChangeEmitter.fire(C));
    }
  }
  // ============================================================================
  // Terminal Modes
  // ============================================================================
  /**
   * Query terminal mode state
   *
   * @param mode Mode number (e.g., 2004 for bracketed paste)
   * @param isAnsi True for ANSI modes, false for DEC modes (default: false)
   * @returns true if mode is enabled
   */
  getMode(A, B = !1) {
    return this.assertOpen(), this.wasmTerm.getMode(A, B);
  }
  /**
   * Check if bracketed paste mode is enabled
   */
  hasBracketedPaste() {
    return this.assertOpen(), this.wasmTerm.hasBracketedPaste();
  }
  /**
   * Check if focus event reporting is enabled
   */
  hasFocusEvents() {
    return this.assertOpen(), this.wasmTerm.hasFocusEvents();
  }
  /**
   * Check if mouse tracking is enabled
   */
  hasMouseTracking() {
    return this.assertOpen(), this.wasmTerm.hasMouseTracking();
  }
}
const AA = 2, QA = 1, BA = 15, EA = 100;
class IA {
  constructor() {
    this._isResizing = !1;
  }
  /**
   * Activate the addon (called by Terminal.loadAddon)
   */
  activate(A) {
    this._terminal = A;
  }
  /**
   * Dispose the addon and clean up resources
   */
  dispose() {
    this._resizeObserver && (this._resizeObserver.disconnect(), this._resizeObserver = void 0), this._resizeDebounceTimer && (clearTimeout(this._resizeDebounceTimer), this._resizeDebounceTimer = void 0), this._lastCols = void 0, this._lastRows = void 0, this._terminal = void 0;
  }
  /**
   * Fit the terminal to its container
   *
   * Calculates optimal dimensions and resizes the terminal.
   * Does nothing if dimensions cannot be calculated or haven't changed.
   */
  fit() {
    if (this._isResizing)
      return;
    const A = this.proposeDimensions();
    if (!A || !this._terminal)
      return;
    const B = this._terminal, Q = B.cols, g = B.rows;
    if (!(A.cols === this._lastCols && A.rows === this._lastRows || A.cols === Q && A.rows === g)) {
      this._lastCols = A.cols, this._lastRows = A.rows, this._isResizing = !0;
      try {
        B.resize && typeof B.resize == "function" && B.resize(A.cols, A.rows);
      } finally {
        setTimeout(() => {
          this._isResizing = !1;
        }, 50);
      }
    }
  }
  /**
   * Propose dimensions to fit the terminal to its container
   *
   * Calculates cols and rows based on:
   * - Terminal container element dimensions (clientWidth/Height)
   * - Terminal element padding
   * - Font metrics (character cell size)
   * - Scrollbar width reservation
   *
   * @returns Proposed dimensions or undefined if cannot calculate
   */
  proposeDimensions() {
    var a;
    if (!((a = this._terminal) != null && a.element))
      return;
    const B = this._terminal.renderer;
    if (!B || typeof B.getMetrics != "function")
      return;
    const Q = B.getMetrics();
    if (!Q || Q.width === 0 || Q.height === 0)
      return;
    const g = this._terminal.element;
    if (typeof g.clientWidth > "u")
      return;
    const C = window.getComputedStyle(g), I = Number.parseInt(C.getPropertyValue("padding-top")) || 0, w = Number.parseInt(C.getPropertyValue("padding-bottom")) || 0, o = Number.parseInt(C.getPropertyValue("padding-left")) || 0, i = Number.parseInt(C.getPropertyValue("padding-right")) || 0, M = g.clientWidth, k = g.clientHeight;
    if (M === 0 || k === 0)
      return;
    const G = M - o - i - BA, J = k - I - w, s = Math.max(AA, Math.floor(G / Q.width)), F = Math.max(QA, Math.floor(J / Q.height));
    return { cols: s, rows: F };
  }
  /**
   * Observe the terminal's container for resize events
   *
   * Sets up a ResizeObserver to automatically call fit() when the
   * container size changes. Resize events are debounced to avoid
   * excessive calls during window drag operations.
   *
   * Call dispose() to stop observing.
   */
  observeResize() {
    var A;
    (A = this._terminal) != null && A.element && (this._resizeObserver || (this._resizeObserver = new ResizeObserver((B) => {
      this._isResizing || !B[0] || (this._resizeDebounceTimer && clearTimeout(this._resizeDebounceTimer), this._resizeDebounceTimer = setTimeout(() => {
        this.fit();
      }, EA));
    }), this._resizeObserver.observe(this._terminal.element)));
  }
}
let t = null;
async function DA() {
  t || (t = await x.load());
}
function gA() {
  if (!t)
    throw new Error(
      `ghostty-web not initialized. Call init() before creating Terminal instances.
Example:
  import { init, Terminal } from "ghostty-web";
  await init();
  const term = new Terminal();

For tests, pass a Ghostty instance directly:
  import { Ghostty, Terminal } from "ghostty-web";
  const ghostty = await Ghostty.load();
  const term = new Terminal({ ghostty });`
    );
  return t;
}
export {
  _ as CanvasRenderer,
  U as CellFlags,
  K as EventEmitter,
  IA as FitAddon,
  x as Ghostty,
  j as GhosttyTerminal,
  X as InputHandler,
  z as KeyEncoder,
  m as KeyEncoderOption,
  P as LinkDetector,
  v as OSC8LinkProvider,
  $ as SelectionManager,
  CA as Terminal,
  u as UrlRegexProvider,
  gA as getGhostty,
  DA as init
};
