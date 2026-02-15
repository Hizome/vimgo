package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

const defaultAddr = ":8080"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsControlMessage struct {
	Type string `json:"type"`
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

func main() {
	bin := "./vimgo"
	staticDir := filepath.Join("web", "static")

	http.Handle("/", http.FileServer(http.Dir(staticDir)))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		size := parseBoardSize(r.URL.Query().Get("size"))
		handleWS(w, r, bin, size)
	})

	log.Printf("web ui listening on http://localhost%s", defaultAddr)
	log.Fatal(http.ListenAndServe(defaultAddr, nil))
}

func handleWS(w http.ResponseWriter, r *http.Request, bin string, size int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	cmd := exec.Command(bin, "-size", strconv.Itoa(size))
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("failed to start pty process: %v", err)
		writeErr(conn, "failed to start vimgo process: "+err.Error())
		return
	}
	_ = pty.Setsize(ptmx, &pty.Winsize{Cols: 120, Rows: 40})

	defer func() {
		_ = ptmx.Close()
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_, _ = cmd.Process.Wait()
		}
	}()

	done := make(chan struct{})

	go func() {
		defer close(done)
		buf := make([]byte, 8192)
		for {
			n, readErr := ptmx.Read(buf)
			if n > 0 {
				if writeErr := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); writeErr != nil {
					return
				}
			}
			if readErr != nil {
				if readErr != io.EOF && !errors.Is(readErr, syscall.EIO) {
					log.Printf("pty read error: %v", readErr)
				}
				return
			}
		}
	}()

	for {
		msgType, payload, err := conn.ReadMessage()
		if err != nil {
			break
		}

		switch msgType {
		case websocket.BinaryMessage:
			if _, err := ptmx.Write(payload); err != nil {
				log.Printf("pty write error: %v", err)
				return
			}
		case websocket.TextMessage:
			var ctl wsControlMessage
			if err := json.Unmarshal(payload, &ctl); err != nil {
				continue
			}
			if ctl.Type == "resize" && ctl.Cols > 0 && ctl.Rows > 0 {
				if err := pty.Setsize(ptmx, &pty.Winsize{Cols: ctl.Cols, Rows: ctl.Rows}); err != nil {
					log.Printf("pty resize error: %v", err)
				}
			}
		}
	}

	<-done
}

func parseBoardSize(raw string) int {
	if raw == "" {
		return 19
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 19
	}
	if v == 9 || v == 13 || v == 19 {
		return v
	}
	log.Printf("invalid board size %q, falling back to 19", raw)
	return 19
}

func writeErr(conn *websocket.Conn, msg string) {
	_ = conn.WriteMessage(websocket.TextMessage, []byte("\r\n"+fmt.Sprintf("[vimgo-web] %s", msg)+"\r\n"))
}
