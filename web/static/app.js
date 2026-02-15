(function () {
  const termEl = document.getElementById("terminal");
  const term = new Terminal({
    cursorBlink: true,
    convertEol: true,
    scrollback: 0,
    fontFamily: "ui-monospace, Menlo, Monaco, Consolas, monospace",
    fontSize: 15,
    theme: {
      background: "#111822",
    },
  });

  const fitAddon = new FitAddon.FitAddon();
  term.loadAddon(fitAddon);
  term.open(termEl);
  fitAddon.fit();
  term.focus();

  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const ws = new WebSocket(
    `${protocol}://${window.location.host}/ws${window.location.search}`
  );
  ws.binaryType = "arraybuffer";

  function sendResize() {
    if (ws.readyState !== WebSocket.OPEN) {
      return;
    }
    ws.send(
      JSON.stringify({
        type: "resize",
        cols: term.cols,
        rows: term.rows,
      })
    );
  }

  ws.onopen = function () {
    sendResize();
  };

  ws.onmessage = function (event) {
    if (event.data instanceof ArrayBuffer) {
      term.write(new Uint8Array(event.data));
      return;
    }

    if (typeof event.data === "string") {
      term.write(event.data);
    }
  };

  ws.onerror = function () {
    term.write("\r\n\x1b[31mWebSocket error\x1b[0m\r\n");
  };

  ws.onclose = function () {
    term.write("\r\n\x1b[33mDisconnected\x1b[0m\r\n");
  };

  term.onData(function (data) {
    if (ws.readyState !== WebSocket.OPEN) {
      return;
    }
    ws.send(new TextEncoder().encode(data));
  });

  window.addEventListener("resize", function () {
    fitAddon.fit();
    sendResize();
  });
})();
