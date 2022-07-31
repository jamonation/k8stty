const fitAddon = new FitAddon.FitAddon();
const utf8decoder = new TextDecoder('utf-8');

// var now = new Date();
// var end = new Date();
// var interval = null;

// function timer() {
//   if ((Math.floor(end - now) / 1000) <= 10) {
//     //window.parent.postMessage("Expired", document.referrer);
//     alert("Closing terminal");
//     window.clearInterval(interval);
//   } else {
//     now = now + 1000;
//     console.log(Date(now));
//     console.log(end - now);
//   }

// }

async function requestTerminal() {
    let image = new URLSearchParams(location.search).get("image") || "ubuntu:jammy";
    console.log(image);
    // add some animation() function
    const headers = [
        ['Content-Type', 'application/json']
    ];
    response = await fetch(`/api/v1/terminal/create?image=${image}`, { headers });
    let msg = await response.json();
    return msg;
}

async function attachTerminal() {
    var loadoverlay = document.querySelector("#loadoverlay");
    loadoverlay.className = 'loading';


    let resp = await requestTerminal();
    if (resp.status_code != 201) {
        return;
    }
    var terminal_id = resp.msg.id;

    let term = new Terminal({
        cursorBlink: true,
        cursorStyle: "block",
        fontFamily: "monospace",
        fontSize: "18",
        rendererType: "dom",
        theme: {
            background: "white",
            cursor: "green",
            cursorAccent: "white",
            selection: "green",
            foreground: "black"
        },
    });

    term.active = false;
    term.loadAddon(fitAddon);
    term.clear();

    term.onData(function (data, event) {
        sendData(ws, data);
    });

    term.onResize(function (event) {
        sendResize(ws, event);
    });

    if (window.location.protocol === 'https:') {
        ws_protocol = 'wss';
    } else {
        ws_protocol = 'ws';
    }

    let command = new URLSearchParams(location.search).get("command") || "/bin/sh";

    let ws = new WebSocket(
        `${ws_protocol}://${window.location.host}/api/v1/terminal/attach/${terminal_id}?command=${command}`);
    ws.binaryType = 'arraybuffer';

    ws.onopen = function (evt) {
        term.open(document.getElementById('k8stty'));
        fitAddon.fit();
    };

    ws.onmessage = evt => {
        if (!evt.data instanceof ArrayBuffer) {
            console.log('expecting array buffer: ', evt.data);
        }
        loadoverlay.className = 'hidden';
        term.write(utf8decoder.decode(evt.data).replace('\0', ''));
    };

    ws.onclose = function (evt) {
        term.write(`websocket closed`);
    };

    ws.onerror = function (evt) {
        console.log(`websocket error: `, evt);
    };
}

function sendData(ws, data) {
    // Send \x00 first byte to indicate this is a data channel message
    ws.send(new TextEncoder().encode('\x00' + data));
}

function sendResize(ws, e) {
    message = JSON.stringify({
        "event": "resize",
        "size": { cols: e.cols, rows: e.rows }
    });
    // Send \x01 first byte to indicate this is a control channel message
    ws.send(new TextEncoder().encode('\x01' + message));
}

document.addEventListener("DOMContentLoaded", () => {
    window.addEventListener("resize", () => {
        // fitAddon triggers a term.onResize event
        // which calls sendResize with rows & cols
        fitAddon.fit();
    });
});

attachTerminal();