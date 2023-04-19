package ws

import (
	"testing"
)

func TestHandleWebsocket(t *testing.T) {
	//e := echo.New()
	//
	////todo handle this entirely in the test
	//client, err := CreateClient("localhost:30085", "ws")
	//require.NoError(t, err)
	//
	//e.GET("/*", client.HandleWebsocket)
	//
	//go func() {
	//	e.Logger.Fatal(e.Start(":8080"))
	//}()
	//
	//time.Sleep(2 * time.Second)
	//
	//u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}
	//
	//// connect to the WebSocket server
	//conn, _, err := gorilla.DefaultDialer.Dial(u.String(), nil)
	//require.NoError(t, err)
	//
	//err = conn.WriteMessage(gorilla.TextMessage, []byte("test"))
	//require.NoError(t, err)
}
