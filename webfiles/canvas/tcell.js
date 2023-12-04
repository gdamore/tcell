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
const fontclass = "16px \"Menlo\", \"Andale Mono\", \"Courier New\", Monospace"
const fontcolor = "#FFFFFF"
const fontback = "#000000"
const dimcolor = "#7F7F7F"
const dimback = "#000000"
const fitwindow = true

/** @type HTMLCanvasElement */
const termElement = document.getElementById("terminal")
const term = termElement.getContext("2d")
term.font = fontclass
const fontSize = term.measureText("X")
// Results are less than ideal if these are not integers.
const fontwidth = Math.round(fontSize.width)
const fontheight = Math.round(fontSize.fontBoundingBoxAscent + fontSize.fontBoundingBoxDescent)
const vertoffset = Math.round(fontSize.fontBoundingBoxAscent)

var width = 80; var height = 24
if (fitwindow) {
    document.documentElement.style.overflow = 'hidden';
    width = Math.floor(document.documentElement.clientWidth / fontwidth)
    height = Math.floor(document.documentElement.clientHeight / fontheight)
}

const beepAudio = new Audio("beep.wav");

var cursorClass = "cursor-blinking-block"

var blinkState = false
var blinkCells = []

var blinkId = setInterval(function () {
    blinkState = !blinkState
    blinkCells.forEach(cell => {
        drawCell(cell.x, cell.y, cell.mainc, cell.combc, cell.fg, cell.bg, cell.attrs, true)
    })
}, 500);

function initialize() {
    resize(width, height) // initialize content
}

function resize(w, h) {
    width = w
    height = h
    termElement.width = fontwidth * width
    termElement.height = fontheight * height
    clearScreen()
}

function clearScreen(fg, bg) {
    if (fg) { fontcolor = intToHex(fg) }
    if (bg) { fontback = intToHex(bg) }

    term.fillStyle = fontback
    term.fillRect(0, 0, termElement.width, termElement.height)
    blinkCells = []
}

function invert(data) {
    for (var idx = 0; idx < data.data.byteLength; idx++) {
        if (idx % 4 != 3) { // don't invert the alpha channel
            data.data[idx] = 255 - data.data[idx]
        }
    }
}

function drawCell(x, y, mainc, combc, fg, bg, attrs, blink) {
    if (!blink) {
        if ((attrs & ((1 << 1) | (1 << 8))) != 0) {
            // include x and y so we don't have to deconvert the id
            blinkCells[y * width + x] = { x: x, y: y, mainc: mainc, combc: combc, fg: fg, bg: bg, attrs: attrs }
        } else {
            delete blinkCells[y * width + x]
        }
    }

    var combString = String.fromCharCode(mainc)
    combc.forEach(char => { combString += String.fromCharCode(char) });

    let fgc = fontcolor
    let bgc = fontback
    if ((attrs & (1 << 4)) == 0) {
        if (fg != -1) { fgc = intToHex(fg) }
        if (bg != -1) { bgc = intToHex(bg) }
    } else { // dim
        fgc = dimcolor
        if (fg != -1) { fgc = intToHex((fg & 0xFEFEFE) >> 1) }
        bgc = dimback
        if (bg != -1) { bgc = intToHex((bg & 0xFEFEFE) >> 1) }
    }

    if (((attrs & (1 << 1)) != 0) && !blinkState) { // blink off.  Just blank the cell and return.
        term.fillStyle = bgc
        term.fillRect(x * fontwidth, y * fontheight, fontwidth, fontheight)
        return
    }

    if ((attrs & (1 << 2)) != 0) { // reverse video
        var temp = bgc
        bgc = fgc
        fgc = temp
    }

    term.fillStyle = bgc
    term.fillRect(x * fontwidth, y * fontheight, fontwidth, fontheight)

    term.fillStyle = fgc

    let fc = fontclass

    if (attrs != 0) {
        if ((attrs & 1) != 0) { fc = "bold " + fc }
        if ((attrs & (1 << 3)) != 0) { // underscore
            term.fillRect(x * fontwidth, y * fontheight + fontheight - 2, fontwidth, 2)
        }
        if ((attrs & (1 << 5)) != 0) { fc = "italic " + fc }
        if ((attrs & (1 << 6)) != 0) { // strikethrough
            term.fillRect(x * fontwidth, y * fontheight + fontheight / 2 - 1, fontwidth, 2)
        }
    }

    term.font = fc
    term.fillText(combString, x * fontwidth, y * fontheight + vertoffset)

    if ((attrs & (1 << 8)) != 0) { // cursor: invert the cursor area
        var data
        switch (cursorClass) {
            case "cursor-blinking-block":
                if (!blinkState) break;
            case "cursor-steady-block":
                data = term.getImageData(x * fontwidth, y * fontheight, fontwidth, fontheight)
                invert(data)
                term.putImageData(data, x * fontwidth, y * fontheight)
            case "cursor-blinking-underline":
                if (!blinkState) break;
            case "cursor-steady-underline":
                data = term.getImageData(x * fontwidth, y * fontheight + fontheight - 2, fontwidth, 2)
                invert(data)
                term.putImageData(data, x * fontwidth, y * fontheight + fontheight - 2)
                break;
            case "cursor-blinking-bar":
                if (!blinkState) break;
            case "cursor-steady-bar":
                data = term.getImageData(x * fontwidth, y * fontheight, 2, fontheight)
                invert(data)
                term.putImageData(data, x * fontwidth, y * fontheight)
                break;
        }
    }
}

function setCursorStyle(newClass) {
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

document.addEventListener("keydown", e => {
    onKeyEvent(e.key, e.shiftKey, e.altKey, e.ctrlKey, e.metaKey)
})

termElement.addEventListener("click", e => {
    onMouseClick(Math.min((e.offsetX / fontwidth) | 0, width - 1), Math.min((e.offsetY / fontheight) | 0, height - 1), e.which, e.shiftKey, e.altKey, e.ctrlKey)
})

termElement.addEventListener("mousemove", e => {
    onMouseMove(Math.min((e.offsetX / fontwidth) | 0, width - 1), Math.min((e.offsetY / fontheight) | 0, height - 1), e.which, e.shiftKey, e.altKey, e.ctrlKey)
})

termElement.addEventListener("focus", e => {
    onFocus(true)
})

termElement.addEventListener("blur", e => {
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
