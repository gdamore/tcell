// Copyright 2024 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

const wasmFilePath = "main.wasm";
const term = document.getElementById("terminal");
var width = 80;
var height = 24;
const beepAudio = new Audio("beep.wav");

var cx = -1;
var cy = -1;
var cursorClass = "cursor-blinking-block";
var cursorColor = "";

var content; // {data: row[height], dirty: bool}
// row = {data: element[width], previous: span}
// dirty/[previous being null] indicates if previous (or entire terminal) needs to be recalculated.
// dirty is true/null if terminal/previous need to be re-calculated/shown

function initialize() {
  resize(width, height); // initialize content
  show(); // then show the screen
}

function resize(w, h) {
  width = w;
  height = h;
  content = { data: new Array(height), dirty: true };
  for (let i = 0; i < height; i++) {
    content.data[i] = { data: new Array(width), previous: null };
  }

  clearScreen();
}

function clearScreen(fg, bg) {
  if (fg) {
    term.style.color = intToHex(fg);
  }
  if (bg) {
    term.style.backgroundColor = intToHex(bg);
  }

  content.dirty = true;
  for (let i = 0; i < height; i++) {
    content.data[i].previous = null; // we set the row to be recalculated later
    for (let j = 0; j < width; j++) {
      content.data[i].data[j] = document.createTextNode(" "); // set the entire row to spaces.
    }
  }
}

function drawCell(x, y, s, fg, bg, attrs, us, uc) {
  var span = document.createElement("span");
  var use = false;

  if ((attrs & (1 << 2)) != 0) {
    // reverse video
    var temp = bg;
    bg = fg;
    fg = temp;
    use = true;
  }
  if (fg != -1) {
    span.style.color = intToHex(fg);
    use = true;
  }
  if (bg != -1) {
    span.style.backgroundColor = intToHex(bg);
    use = true;
  }

  // NB: these has to be updated if Attrs.go changes
  if (attrs != 0) {
    use = true;
    if ((attrs & 1) != 0) {
      span.classList.add("bold");
    }
    if ((attrs & (1 << 4)) != 0) {
      span.classList.add("dim");
    }
    if ((attrs & (1 << 5)) != 0) {
      span.classList.add("italic");
    }
    if ((attrs & (1 << 6)) != 0) {
      span.classList.add("strikethrough");
    }
  }
  if (us != 0) {
    use = true;
    if (us == 1) {
      span.classList.add("underline");
    } else if (us == 2) {
      span.classList.add("double_underline");
    } else if (us == 3) {
      span.classList.add("curly_underline");
    } else if (us == 4) {
      span.classList.add("dotted_underline");
    } else if (us == 5) {
      span.classList.add("dashed_underline");
    }
    if (uc != -1) {
      span.style.textDecorationColor = intToHex(uc);
    }
  }

  if ((attrs & (1 << 1)) != 0) {
    var blink = document.createElement("span");
    blink.classList.add("blink");
    var textnode = document.createTextNode(s);
    blink.appendChild(textnode);
    span.appendChild(blink);
  } else {
    var textnode = document.createTextNode(s);
    span.appendChild(textnode);
  }

  content.dirty = true; // invalidate terminal- new cell
  content.data[y].previous = null; // invalidate row- new row
  content.data[y].data[x] = use ? span : textnode;
}

function show() {
  if (!content.dirty) {
    return; // no new draws; no need to update
  }

  displayCursor();

  term.innerHTML = "";
  content.data.forEach((row) => {
    if (row.previous == null) {
      row.previous = document.createElement("span");
      row.data.forEach((c) => {
        row.previous.appendChild(c);
      });
      row.previous.appendChild(document.createTextNode("\n"));
    }
    term.appendChild(row.previous);
  });

  content.dirty = false;
}

function showCursor(x, y) {
  content.dirty = true;

  if (!(cx < 0 || cy < 0)) {
    // if original position is a valid cursor position
    content.data[cy].previous = null;
    if (content.data[cy].data[cx].classList) {
      content.data[cy].data[cx].classList.remove(cursorClass);
    }
  }

  cx = x;
  cy = y;
}

function displayCursor() {
  content.dirty = true;

  if (!(cx < 0 || cy < 0)) {
    // if new position is a valid cursor position
    content.data[cy].previous = null;

    if (!content.data[cy].data[cx].classList) {
      var span = document.createElement("span");
      span.appendChild(content.data[cy].data[cx]);
      content.data[cy].data[cx] = span;
    }

    if (cursorColor != "") {
      term.style.setProperty("--cursor-color", cursorColor);
    } else {
      term.style.setProperty("--cursor-color", "lightgrey");
    }

    content.data[cy].data[cx].classList.add(cursorClass);
  }
}

function setCursorStyle(newClass, newColor) {
  if (newClass == cursorClass && newColor == cursorColor) {
    return;
  }

  if (!(cx < 0 || cy < 0)) {
    // mark cursor row as dirty; new class has been applied to (cx, cy)
    content.dirty = true;
    content.data[cy].previous = null;

    if (content.data[cy].data[cx].classList) {
      content.data[cy].data[cx].classList.remove(cursorClass);
    }

    // adding the new class will be dealt with when displayCursor() is called
  }

  cursorClass = newClass;
  cursorColor = newColor;
}

function beep() {
  beepAudio.currentTime = 0;
  beepAudio.play();
}

function setTitle(title) {
  document.title = title;
}

function intToHex(n) {
  return "#" + n.toString(16).padStart(6, "0");
}

initialize();

let fontwidth = term.clientWidth / width;
let fontheight = term.clientHeight / height;

document.addEventListener("keydown", (e) => {
  onKeyEvent(e.key, e.shiftKey, e.altKey, e.ctrlKey, e.metaKey);
});

term.addEventListener("click", (e) => {
  onMouseClick(
    Math.min((e.offsetX / fontwidth) | 0, width - 1),
    Math.min((e.offsetY / fontheight) | 0, height - 1),
    e.which,
    e.shiftKey,
    e.altKey,
    e.ctrlKey
  );
});

term.addEventListener("mousemove", (e) => {
  onMouseMove(
    Math.min((e.offsetX / fontwidth) | 0, width - 1),
    Math.min((e.offsetY / fontheight) | 0, height - 1),
    e.which,
    e.shiftKey,
    e.altKey,
    e.ctrlKey
  );
});

term.addEventListener("focus", (e) => {
  onFocus(true);
});

term.addEventListener("blur", (e) => {
  onFocus(false);
});
term.tabIndex = 0;

document.addEventListener("paste", (e) => {
  e.preventDefault();
  var text = (e.originalEvent || e).clipboardData.getData("text/plain");
  onPaste(true);
  for (let i = 0; i < text.length; i++) {
    onKeyEvent(text.charAt(i), false, false, false, false);
  }
  onPaste(false);
});

const go = new Go();
WebAssembly.instantiateStreaming(fetch(wasmFilePath), go.importObject).then(
  (result) => {
    go.run(result.instance);
  }
);
