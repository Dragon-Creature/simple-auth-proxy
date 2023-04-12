package ws

import (
	gorilla "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"log"
	"net/url"
)

type WS struct {
	send    chan []byte
	receive chan string
}

func (w *WS) HandleWebsocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		go func() {
			for {
				msg := <-w.send
				err := websocket.Message.Send(ws, msg)
				if err != nil {
					log.Println(err)
				}
			}
		}()

		for {
			//message := ""
			//err := websocket.Message.Receive(ws, &message)
			//if err != nil {
			//	c.Logger().Error(err)
			//}
			//w.receive <- message
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func CreateClient(targetURL string, path string) WS {
	send := make(chan []byte, 10)
	receive := make(chan string, 10)

	u := url.URL{Scheme: "ws", Host: targetURL, Path: path}
	log.Printf("connecting to %s", u.String())

	// connect to the WebSocket server
	conn, _, err := gorilla.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}

	// start a goroutine to read messages from the WebSocket server
	go func() {
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			send <- msg
		}
	}()

	go func() {
		for {
			msg := <-receive

			err = conn.WriteMessage(gorilla.TextMessage, []byte(msg))
			if err != nil {
				log.Println("read error:", err)
				return
			}
		}
	}()
	return WS{
		send:    send,
		receive: receive,
	}
}
