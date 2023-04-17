package ws

import (
	gorilla "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"log"
	"net/url"
)

type WS struct {
	send    chan []byte
	receive chan []byte
}

func (w *WS) HandleWebsocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				err = errors.WithStack(err)
				log.Println(err)
				return
			}
		}()

		go func() {
			for {
				msg := <-w.send
				err := websocket.Message.Send(ws, msg)
				if err != nil {
					err = errors.WithStack(err)
					log.Println(err)
					return
				}
			}
		}()

		for {
			message := []byte{}
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				err = errors.WithStack(err)
				log.Println(err)
				return
			}
			w.receive <- message
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func CreateClient(targetURL string, path string) (*WS, error) {
	send := make(chan []byte, 10)
	receive := make(chan []byte, 10)

	u := url.URL{Scheme: "ws", Host: targetURL, Path: path}
	log.Printf("connecting to %s", u.String())

	// connect to the WebSocket server
	conn, _, err := gorilla.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// start a goroutine to read messages from the WebSocket server
	go func() {
		defer func() {
			err = conn.Close()
			err = errors.WithStack(err)
			log.Println(err)
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				err = errors.WithStack(err)
				log.Println(err)
				return
			}
			send <- msg
		}
	}()

	go func() {
		for {
			msg := <-receive

			//todo handle different message flags
			err = conn.WriteMessage(gorilla.TextMessage, msg)
			if err != nil {
				err = errors.WithStack(err)
				log.Println(err)
				return
			}
		}
	}()
	return &WS{
		send:    send,
		receive: receive,
	}, nil
}
