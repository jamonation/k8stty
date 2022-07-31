package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingWait = (pongWait * 9) / 10
)

type WS struct {
	rChan resizeChan
	rbuf  []byte
	cond  *sync.Cond
	conn  *websocket.Conn
	sync.Mutex
}

type resizeChan chan remotecommand.TerminalSize

type ttySize struct {
	Rows uint16
	Cols uint16
	X    uint16
	Y    uint16
}

// The json field names here are only used as annotations to
// describe the json keys of incoming messages. They aren't sent to
// clients in this version of the terminal.
type controlMessage struct {
	Event string  `json:"event"`
	Data  string  `json:"data"`
	Size  ttySize `json:"size"`
}

// Run reads from client (this runs in a goroutine)
func (w *WS) Run() {
	for {
		messageType, r, err := w.conn.ReadMessage()
		if err != nil {
			log.Printf("error reading websocket: %s", err)
			return
		}

		// TODO: turn this into a control channel for the socket,
		// perhaps in conjunction with h.Conn.ReadJSON
		if messageType == websocket.TextMessage {
			log.Printf("expected BinaryMessage\n")
			continue
		}

		// First byte is either \x00 or \x01 indicating message type
		typeByteMarker, err := bytes.NewReader(r).ReadByte()
		if err != nil {
			log.Printf("error reading type byte marker: %v\n", err)
			return
		}

		switch typeByteMarker {
		// this is a data channel message prepended with \x00 on the client
		case 0:
			{
				w.Lock()
				w.rbuf = append(w.rbuf, r[1:]...) // skip first byte control character
				w.cond.Signal()
				w.Unlock()
			}
		// this is a control channel message prepended with \x01 on the client
		case 1:
			{
				msg := controlMessage{}
				err := json.Unmarshal(r[1:], &msg) // skip first byte control character
				if err != nil {
					log.Printf("error decoding control message: %v\n", err)
				}

				switch msg.Event {
				case "resize":
					// log.Printf("received resize message: %v\n", msg)
					size := remotecommand.TerminalSize{
						Width:  msg.Size.Cols,
						Height: msg.Size.Rows}
					w.rChan <- size
				}
			}
		default:
			continue
		}
	}
}

func (w *WS) Next() *remotecommand.TerminalSize {
	termSize := <-w.rChan
	size := &termSize
	return size
}

func (w *WS) Read(b []byte) (int, error) {
	w.Lock()
	for len(w.rbuf) == 0 {
		w.cond.Wait()
	}
	size := copy(b, w.rbuf)
	w.rbuf = w.rbuf[size:]
	w.Unlock()
	return size, nil
}

func (w *WS) Write(b []byte) (int, error) {
	size := len(b)
	err := w.conn.WriteMessage(websocket.BinaryMessage, b)
	return size, err
}
