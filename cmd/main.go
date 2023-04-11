package main

import (
	"fmt"
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/login"
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/proxy"
	"github.com/labstack/echo/v4"
	"github.com/rgamba/evtwebsocket"
	"golang.org/x/net/websocket"
	"log"
	"os"
)

type ws struct {
	con evtwebsocket.Conn
	msg chan []byte
}

func (w *ws) hello(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		go func() {
			for {
				msg := <-w.msg
				fmt.Printf("server: %s\n", string(msg))
				err := websocket.Message.Send(ws, msg)
				if err != nil {
					log.Fatal(err)
				}
			}
		}()

		for {
			message := ""
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("client: %s\n", message)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	message1 := make(chan []byte, 10)

	conn1 := evtwebsocket.Conn{
		// Fires when the connection is established
		OnConnected: func(w *evtwebsocket.Conn) {
			fmt.Println("Connected!")
		},
		// Fires when a new message arrives from the server
		OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
			message1 <- msg
		},
		// Fires when an error occurs and connection is closed
		OnError: func(err error) {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		},
		// Ping interval in secs (optional)
		PingIntervalSecs: 5,
		// Ping message to send (optional)
		PingMsg: []byte("PING"),
	}

	err := conn1.Dial("ws://localhost:30085/ws", "")
	if err != nil {
		log.Fatal(err)
	}

	w1 := ws{
		con: conn1,
		msg: message1,
	}

	message2 := make(chan []byte, 10)

	conn2 := evtwebsocket.Conn{
		// Fires when the connection is established
		OnConnected: func(w *evtwebsocket.Conn) {
			fmt.Println("Connected!")
		},
		// Fires when a new message arrives from the server
		OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
			message2 <- msg
		},
		// Fires when an error occurs and connection is closed
		OnError: func(err error) {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		},
		// Ping interval in secs (optional)
		PingIntervalSecs: 5,
		// Ping message to send (optional)
		PingMsg: []byte("PING"),
	}

	err = conn2.Dial("ws://localhost:30085/live/Front_Deck", "")
	if err != nil {
		log.Fatal(err)
	}

	w2 := ws{
		con: conn2,
		msg: message2,
	}

	e := echo.New()
	e.GET("/*", login.GetLoginPage)
	e.POST("/auth", proxy.PostAuth)
	e.GET("/ws", w1.hello)
	e.GET("/live/Front_Deck", w2.hello)
	e.Logger.Fatal(e.Start(":8080"))

	//r := chi.NewRouter()
	//r.Use(middleware.Logger)
	//wsServer := websocket.Start(context.Background())
	//
	//r.HandleFunc("/ws", wsServer.Handler)
	//
	//r.HandleFunc("/*", login.GetLoginPage)
	//r.Post("/auth", proxy.PostAuth)
	//
	//wsServer.On("Front_Deck", func(c *websocket.Conn, msg *websocket.Message) {
	//	conn := evtwebsocket.Conn{
	//		// Fires when the connection is established
	//		OnConnected: func(w *evtwebsocket.Conn) {
	//			fmt.Println("Connected!")
	//		},
	//		// Fires when a new message arrives from the server
	//		OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
	//			fmt.Printf("New message: %s\n", msg)
	//		},
	//		// Fires when an error occurs and connection is closed
	//		OnError: func(err error) {
	//			fmt.Printf("Error: %s\n", err.Error())
	//			os.Exit(1)
	//		},
	//		// Ping interval in secs (optional)
	//		PingIntervalSecs: 5,
	//		// Ping message to send (optional)
	//		PingMsg: []byte("PING"),
	//	}
	//
	//	err := conn.Dial("http://localhost:30085/ws", "")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	_ = c.Emit("echo", msg.Data)
	//})
	//err := http.ListenAndServe(":8080", r)
	//if err != nil {
	//	panic(err)
	//}
}
