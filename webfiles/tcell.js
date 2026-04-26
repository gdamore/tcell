// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

const wasmFilePath = "main.wasm";
const container = document.getElementById("terminal");
const beepAudio = new Audio("beep.wav");
const textDecoder = new TextDecoder();
const ghosttyWebURL = new URL("./ghostty-web/ghostty-web.js", import.meta.url);
const ghosttyWasmURL = new URL("./ghostty-web/ghostty-vt.wasm", import.meta.url);

window.addEventListener("error", (event) => {
  showFatal(event.error ?? event.message);
});
window.addEventListener("unhandledrejection", (event) => {
  showFatal(event.reason);
});

const { FitAddon, Ghostty, Terminal } = await importGhosttyWeb();
const modes = {
  mouseButton: false,
  mouseDrag: false,
  mouseMotion: false,
  mouseSgr: false,
  focus: false,
  paste: false,
  resize: false,
  win32Input: false,
};
const pressedButtons = new Set();
const modifierState = {
  ShiftLeft: false,
  ShiftRight: false,
  ControlLeft: false,
  ControlRight: false,
  AltLeft: false,
  AltRight: false,
  MetaLeft: false,
  MetaRight: false,
};

const initialSize = configuredSize();
const ghostty = await Ghostty.load(ghosttyWasmURL.href);

const term = new Terminal({
  ghostty,
  cols: initialSize.cols,
  rows: initialSize.rows,
  cursorBlink: true,
  fontFamily: '"Menlo", "Andale Mono", "Courier New", monospace',
  fontSize: 14,
  theme: {
    background: "#000000",
    foreground: "#e5e5e5",
  },
});
const fit = new FitAddon();
term.loadAddon(fit);
term.open(container);
fitTerminal();

term.onData((data) => {
  if (modes.win32Input) {
    return;
  }
  globalThis.tcellRead?.(data);
});
term.onResize(() => {
  globalThis.tcellResize?.();
});
term.onBell(() => {
  beep();
});

globalThis.tcellWrite = (data) => {
  const text = data instanceof Uint8Array ? textDecoder.decode(data, { stream: true }) : data;
  handleHostOutput(text);
  term.write(text);
};

globalThis.tcellWindowSize = () => {
  const metrics = term.renderer?.getMetrics?.();
  return {
    cols: term.cols,
    rows: term.rows,
    pixelWidth: metrics ? Math.floor(metrics.width * term.cols) : 0,
    pixelHeight: metrics ? Math.floor(metrics.height * term.rows) : 0,
  };
};

globalThis.beep = beep;

function beep() {
  beepAudio.currentTime = 0;
  beepAudio.play()?.catch(() => {});
}

async function importGhosttyWeb() {
  try {
    return await import(/* @vite-ignore */ ghosttyWebURL.href);
  } catch (error) {
    throw new Error(
      "Failed to load ./ghostty-web/ghostty-web.js. Copy webfiles/ghostty-web next to tcell.js.",
      { cause: error }
    );
  }
}

function showFatal(error) {
  const message = error instanceof Error ? error.message : String(error);
  const cause = error instanceof Error && error.cause instanceof Error ? `\n${error.cause.message}` : "";
  container.textContent = `Tcell WASM startup failed:\n${message}${cause}`;
  container.style.whiteSpace = "pre-wrap";
  container.style.color = "#ff6b6b";
  container.style.backgroundColor = "#111";
  container.style.padding = "1rem";
  console.error(error);
}

function configuredSize() {
  return {
    cols: positiveInt(container.dataset.cols) ?? 80,
    rows: positiveInt(container.dataset.rows) ?? 24,
  };
}

function explicitSize() {
  return {
    cols: positiveInt(container.dataset.cols),
    rows: positiveInt(container.dataset.rows),
  };
}

function positiveInt(value) {
  const n = Number.parseInt(value ?? "", 10);
  return Number.isFinite(n) && n > 0 ? n : undefined;
}

function fitTerminal() {
  const size = explicitSize();
  if (size.cols != null && size.rows != null) {
    if (term.cols !== size.cols || term.rows !== size.rows) {
      term.resize(size.cols, size.rows);
      reportResize();
    }
    return;
  }

  fit.fit();
  if (size.cols != null || size.rows != null) {
    const cols = size.cols ?? term.cols;
    const rows = size.rows ?? term.rows;
    if (term.cols !== cols || term.rows !== rows) {
      term.resize(cols, rows);
    }
  }
  reportResize();
}

function reportResize() {
  globalThis.tcellResize?.();
  if (modes.resize) {
    sendInput(`\x1b[4;${term.rows};${term.cols}t`);
  }
}

const resizeObserver = new ResizeObserver(fitTerminal);
resizeObserver.observe(container);

window.addEventListener("resize", fitTerminal);

container.addEventListener("mousedown", (event) => {
  if (!mouseReportingEnabled()) {
    return;
  }
  captureTerminalEvent(event);
  focusTerminal();
  pressedButtons.add(event.button);
  sendMouse(event, event.button, false, true);
}, true);

container.addEventListener("mouseup", (event) => {
  if (!mouseReportingEnabled()) {
    return;
  }
  captureTerminalEvent(event);
  sendMouse(event, event.button, false, false);
  pressedButtons.delete(event.button);
}, true);

container.addEventListener("mousemove", (event) => {
  if (!mouseReportingEnabled()) {
    return;
  }
  if (!modes.mouseMotion && !(modes.mouseDrag && pressedButtons.size > 0)) {
    return;
  }
  captureTerminalEvent(event);
  sendMouse(event, firstPressedButton(), true, true);
}, true);

container.addEventListener("contextmenu", (event) => {
  if (modes.mouseButton) {
    captureTerminalEvent(event);
  }
}, true);

container.addEventListener("auxclick", (event) => {
  if (modes.mouseButton) {
    captureTerminalEvent(event);
  }
}, true);

container.addEventListener("selectstart", (event) => {
  if (modes.mouseButton) {
    captureTerminalEvent(event);
  }
}, true);

container.addEventListener(
  "wheel",
  handleWheel,
  { passive: false }
);

window.addEventListener("wheel", (event) => {
  if (container.contains(event.target)) {
    handleWheel(event);
  }
}, { capture: true, passive: false });

container.addEventListener("focusin", () => {
  if (modes.focus) {
    sendInput("\x1b[I");
  }
});

container.addEventListener("focusout", () => {
  if (modes.focus) {
    sendInput("\x1b[O");
  }
});

window.addEventListener("focus", () => {
  if (modes.focus) {
    sendInput("\x1b[I");
  }
});

window.addEventListener("blur", () => {
  clearModifierState();
  pressedButtons.clear();
  if (modes.focus) {
    sendInput("\x1b[O");
  }
});

container.addEventListener(
  "paste",
  (event) => {
    if (!modes.paste) {
      return;
    }
    event.preventDefault();
    event.stopPropagation();
    const text = event.clipboardData?.getData("text/plain") ?? "";
    sendInput(`\x1b[200~${text}\x1b[201~`);
  },
  true
);

window.addEventListener("keydown", (event) => {
  if (!modes.win32Input) {
    return;
  }
  captureTerminalEvent(event);
  updateModifierState(event, true);
  sendWinKey(event, true);
}, true);

window.addEventListener("keyup", (event) => {
  if (!modes.win32Input) {
    return;
  }
  captureTerminalEvent(event);
  updateModifierState(event, false);
  sendWinKey(event, false);
}, true);

function handleHostOutput(text) {
  replyToQueries(text);
  updateModes(text);
}

function replyToQueries(text) {
  if (text.includes("\x1b[c")) {
    // VT420-ish primary DA with color and clipboard capability.
    sendInput("\x1b[?64;22;52c");
  }
  if (text.includes("\x1b[>q")) {
    sendInput("\x1bP>|ghostty-web tcell\x1b\\");
  }
  const modeStates = {
    1000: modes.mouseButton,
    1002: modes.mouseDrag,
    1003: modes.mouseMotion,
    1004: modes.focus,
    1006: modes.mouseSgr,
    2004: modes.paste,
    2048: modes.resize,
    9001: modes.win32Input,
  };
  for (const [mode, enabled] of Object.entries(modeStates)) {
    if (text.includes(`\x1b[?${mode}$p`)) {
      sendInput(`\x1b[?${mode};${enabled ? 1 : 2}$y`);
    }
  }
  // Deliberately do not advertise Kitty/XTerm keyboard modes; the bridge uses
  // tcell's existing Win32 input parser for browser advanced-key reporting.
}

function updateModes(text) {
  const re = /\x1b\[\?([0-9;]+)([hl])/g;
  let match;
  while ((match = re.exec(text)) !== null) {
    const enabled = match[2] === "h";
    for (const rawMode of match[1].split(";")) {
      setMode(Number(rawMode), enabled);
    }
  }
}

function setMode(mode, enabled) {
  switch (mode) {
    case 1000:
      modes.mouseButton = enabled;
      break;
    case 1002:
      modes.mouseDrag = enabled;
      break;
    case 1003:
      modes.mouseMotion = enabled;
      break;
    case 1006:
      modes.mouseSgr = enabled;
      break;
    case 1004:
      modes.focus = enabled;
      if (enabled) {
        sendInput(document.hasFocus() ? "\x1b[I" : "\x1b[O");
      }
      break;
    case 2004:
      modes.paste = enabled;
      break;
    case 2048:
      modes.resize = enabled;
      break;
    case 9001:
      modes.win32Input = enabled;
      break;
  }
}

function sendMouse(event, button, motion, down, wheel = false) {
  const pos = mouseCell(event);
  if (pos == null) {
    return;
  }
  let cb;
  if (wheel) {
    cb = 64 + button;
  } else {
    cb = buttonCode(button);
    if (motion) {
      cb += 32;
    }
  }
  if (event.shiftKey) {
    cb += 4;
  }
  if (event.altKey || event.metaKey) {
    cb += 8;
  }
  if (event.ctrlKey) {
    cb += 16;
  }
  if (modes.mouseSgr) {
    sendInput(`\x1b[<${cb};${pos.x};${pos.y}${down ? "M" : "m"}`);
    return;
  }

  // Legacy mouse reports are byte-oriented, not UTF-8 text. Use Uint8Array so
  // coordinates above 95 columns do not get corrupted by string encoding.
  if (!down) {
    cb = 3;
    if (event.shiftKey) {
      cb += 4;
    }
    if (event.altKey || event.metaKey) {
      cb += 8;
    }
    if (event.ctrlKey) {
      cb += 16;
    }
  }
  sendInput(new Uint8Array([
    0x1b,
    0x5b,
    0x4d,
    cb + 32,
    Math.min(pos.x, 223) + 32,
    Math.min(pos.y, 223) + 32,
  ]));
}

function handleWheel(event) {
  if (!mouseReportingEnabled()) {
    return;
  }
  captureTerminalEvent(event);
  const horizontal = Math.abs(event.deltaX) > Math.abs(event.deltaY);
  const button = horizontal
    ? event.deltaX < 0
      ? 2
      : 3
    : event.deltaY < 0
      ? 0
      : 1;
  sendMouse(event, button, false, true, true);
}

function mouseCell(event) {
  const metrics = term.renderer?.getMetrics?.();
  const rect = container.getBoundingClientRect();
  const cellWidth = metrics?.width || rect.width / term.cols;
  const cellHeight = metrics?.height || rect.height / term.rows;
  if (!cellWidth || !cellHeight) {
    return null;
  }
  return {
    x: Math.max(1, Math.min(term.cols, Math.floor((event.clientX - rect.left) / cellWidth) + 1)),
    y: Math.max(1, Math.min(term.rows, Math.floor((event.clientY - rect.top) / cellHeight) + 1)),
  };
}

function buttonCode(button) {
  switch (button) {
    case 1:
      return 1;
    case 2:
      return 2;
    default:
      return 0;
  }
}

function firstPressedButton() {
  return pressedButtons.size === 0 ? 0 : pressedButtons.values().next().value;
}

function mouseReportingEnabled() {
  return modes.mouseButton || modes.mouseDrag || modes.mouseMotion;
}

function sendInput(data) {
  globalThis.tcellRead?.(data);
}

function sendWinKey(event, down) {
  const vk = virtualKey(event);
  const uc = unicodeChar(event);
  const cs = controlKeyState(event);
  sendInput(`\x1b[${vk};0;${uc};${down ? 1 : 0};${cs};1_`);
}

function virtualKey(event) {
  const modifierVK = modifierVirtualKey(event);
  if (modifierVK != null) {
    return modifierVK;
  }
  if (event.keyCode) {
    return event.keyCode;
  }
  if (/^Key[A-Z]$/.test(event.code)) {
    return event.code.charCodeAt(3);
  }
  if (/^Digit[0-9]$/.test(event.code)) {
    return event.code.charCodeAt(5);
  }
  if (/^F([1-9]|1[0-9]|2[0-4])$/.test(event.key)) {
    return 0x70 + Number(event.key.slice(1)) - 1;
  }
  return keyToVK[event.key] ?? 0;
}

function unicodeChar(event) {
  if (event.key.length === 1) {
    return event.key.codePointAt(0);
  }
  return 0;
}

function controlKeyState(event) {
  let state = 0;
  if (modifierState.ShiftLeft || modifierState.ShiftRight || event.shiftKey) {
    state |= 0x10;
  }
  if (modifierState.ControlLeft) {
    state |= 0x08;
  }
  if (modifierState.ControlRight) {
    state |= 0x04;
  }
  if (event.ctrlKey && !modifierState.ControlLeft && !modifierState.ControlRight) {
    state |= event.location === KeyboardEvent.DOM_KEY_LOCATION_RIGHT ? 0x04 : 0x08;
  }
  if (modifierState.AltLeft) {
    state |= 0x02;
  }
  if (modifierState.AltRight) {
    state |= 0x01;
  }
  if (event.altKey && !modifierState.AltLeft && !modifierState.AltRight) {
    state |= event.location === KeyboardEvent.DOM_KEY_LOCATION_RIGHT ? 0x01 : 0x02;
  }
  if (modifierState.MetaLeft) {
    state |= 0x200;
  }
  if (modifierState.MetaRight) {
    state |= 0x400;
  }
  if (event.metaKey && !modifierState.MetaLeft && !modifierState.MetaRight) {
    state |= event.location === KeyboardEvent.DOM_KEY_LOCATION_RIGHT ? 0x400 : 0x200;
  }
  return state;
}

function updateModifierState(event, down) {
  if (event.code in modifierState) {
    modifierState[event.code] = down;
    return;
  }
  switch (event.key) {
    case "Shift":
      modifierState[event.location === KeyboardEvent.DOM_KEY_LOCATION_RIGHT ? "ShiftRight" : "ShiftLeft"] = down;
      break;
    case "Control":
      modifierState[event.location === KeyboardEvent.DOM_KEY_LOCATION_RIGHT ? "ControlRight" : "ControlLeft"] = down;
      break;
    case "Alt":
      modifierState[event.location === KeyboardEvent.DOM_KEY_LOCATION_RIGHT ? "AltRight" : "AltLeft"] = down;
      break;
    case "Meta":
      modifierState[event.location === KeyboardEvent.DOM_KEY_LOCATION_RIGHT ? "MetaRight" : "MetaLeft"] = down;
      break;
  }
}

function clearModifierState() {
  for (const key of Object.keys(modifierState)) {
    modifierState[key] = false;
  }
}

function modifierVirtualKey(event) {
  switch (event.code) {
    case "ShiftLeft":
      return 0xa0;
    case "ShiftRight":
      return 0xa1;
    case "ControlLeft":
      return 0xa2;
    case "ControlRight":
      return 0xa3;
    case "AltLeft":
      return 0xa4;
    case "AltRight":
      return 0xa5;
    case "MetaLeft":
      return 0x5b;
    case "MetaRight":
      return 0x5c;
  }
  return null;
}

function captureTerminalEvent(event) {
  event.preventDefault();
  event.stopPropagation();
  event.stopImmediatePropagation();
}

function focusTerminal() {
  term.focus?.();
  container.focus?.();
}

const keyToVK = {
  Backspace: 0x08,
  Tab: 0x09,
  Enter: 0x0d,
  Pause: 0x13,
  Escape: 0x1b,
  PageUp: 0x21,
  PageDown: 0x22,
  End: 0x23,
  Home: 0x24,
  ArrowLeft: 0x25,
  ArrowUp: 0x26,
  ArrowRight: 0x27,
  ArrowDown: 0x28,
  PrintScreen: 0x2c,
  Insert: 0x2d,
  Delete: 0x2e,
  Help: 0x2f,
  Shift: 0x10,
  Control: 0x11,
  Alt: 0x12,
  Meta: 0x5b,
  CapsLock: 0x14,
};

const go = new Go();
const result = await WebAssembly.instantiateStreaming(
  fetch(wasmFilePath),
  go.importObject
);
go.run(result.instance);
