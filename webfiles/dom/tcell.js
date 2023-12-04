// Copyright 2023 The TCell Authors
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

const wasmFilePath = "main.wasm"
const fitwindow = true
const term = document.getElementById("terminal")
var width = 80; var height = 24
const beepAudio = new Audio("beep.wav");

var cx = -1; var cy = -1
var cursorClass = "cursor-blinking-block"

function initialize() {
    resize(width, height) // initialize content
}

function resize(w, h) {
    width = w
    height = h
    clearScreen()
}

function clearScreen(fg, bg) {
    if (fg) { term.style.color = intToHex(fg) }
    if (bg) { term.style.backgroundColor = intToHex(bg) }

    term.innerHTML = ""
    for (let i = 0; i < height; i++) {
        row = document.createElement("span")
        for (let j = 0; j < width; j++) {
            row.appendChild(document.createTextNode(" "))
        }
        row.appendChild(document.createTextNode("\n"))
        term.appendChild(row)
    }
}

function drawCell(x, y, mainc, combc, fg, bg, attrs) {
    var combString = String.fromCharCode(mainc)
    combc.forEach(char => { combString += String.fromCharCode(char) });

    var span = document.createElement("span")
    var use = false

    if ((attrs & (1 << 2)) != 0) { // reverse video
        var temp = bg
        bg = fg
        fg = temp
        use = true
    }
    if (fg != -1) { span.style.color = intToHex(fg); use = true }
    if (bg != -1) { span.style.backgroundColor = intToHex(bg); use = true }

    if ((x == cx) && (y == cy) && ((attrs & (1 << 8)) == 0)) {
        cx = -1
        cy = -1
    }

    if (attrs != 0) {
        use = true
        if ((attrs & 1) != 0) { span.classList.add("bold") }
        if ((attrs & (1 << 1)) != 0) { span.classList.add("blink") }
        if ((attrs & (1 << 3)) != 0) { span.classList.add("underline") }
        if ((attrs & (1 << 4)) != 0) { span.classList.add("dim") }
        if ((attrs & (1 << 5)) != 0) { span.classList.add("italic") }
        if ((attrs & (1 << 6)) != 0) { span.classList.add("strikethrough") }
        if ((attrs & (1 << 8)) != 0) { span.classList.add(cursorClass); cx = x; cy = y }
    }

    var textnode = document.createTextNode(combString)
    span.appendChild(textnode)

    term.childNodes[y].childNodes[x].replaceWith(use ? span : textnode)
}

function setCursorStyle(newClass) {
    if (newClass == cursorClass) {
        return
    }

    if (!(cx < 0 || cy < 0)) {
        content.dirty = true
        content.data[cy].previous = null

        if (!term.childNodes[cy].childNodes[cx].classList) {
            var span = document.createElement("span")
            span.appendChild(term.childNodes[cy].childNodes[cx])
            term.childNodes[cy].childNodes[cx].replaceWith(span)
        }
        term.childNodes[cy].childNodes[cx].classList.remove(cursorClass)
        term.childNodes[cy].childNodes[cx].classList.add(newClass)
    }

    cursorClass = newClass
}

function beep() {
    beepAudio.currentTime = 0;
    beepAudio.play();
}

function intToHex(n) {
    return "#" + n.toString(16).padStart(6, '0')
}

initialize()

let fontwidth = term.clientWidth / width
let fontheight = term.clientHeight / height

if (fitwindow) {
    document.documentElement.style.overflow = 'hidden';
    resize(Math.floor(document.documentElement.clientWidth / fontwidth), Math.floor(document.documentElement.clientHeight / fontheight))
}

document.addEventListener("keydown", e => {
    onKeyEvent(e.key, e.shiftKey, e.altKey, e.ctrlKey, e.metaKey)
})

term.addEventListener("click", e => {
    onMouseClick(Math.min((e.offsetX / fontwidth) | 0, width - 1), Math.min((e.offsetY / fontheight) | 0, height - 1), e.which, e.shiftKey, e.altKey, e.ctrlKey)
})

term.addEventListener("mousemove", e => {
    onMouseMove(Math.min((e.offsetX / fontwidth) | 0, width - 1), Math.min((e.offsetY / fontheight) | 0, height - 1), e.which, e.shiftKey, e.altKey, e.ctrlKey)
})

term.addEventListener("focus", e => {
    onFocus(true)
})

term.addEventListener("blur", e => {
    onFocus(false)
})
term.tabIndex = 0


document.addEventListener("paste", e => {
    e.preventDefault();
    var text = (e.originalEvent || e).clipboardData.getData('text/plain');
    onPaste(true)
    for (let i = 0; i < text.length; i++) {
        onKeyEvent(text.charAt(i), false, false, false, false)
    }
    onPaste(false)
});

if (fitwindow) {
    document.defaultView.addEventListener("resize", e => {
        const charWidth = Math.floor(document.documentElement.clientWidth / fontwidth)
        const charHeight = Math.floor(document.documentElement.clientHeight / fontheight)
        onResizeEvent(charWidth, charHeight)
    })
}

const go = new Go();
go.env.LINES = height.toString();
go.env.COLUMNS = width.toString();
WebAssembly.instantiateStreaming(fetch(wasmFilePath), go.importObject).then((result) => {
    go.run(result.instance);
});
