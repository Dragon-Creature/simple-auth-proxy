package main

import (
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/login"
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/proxy"
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/ws"
	"github.com/labstack/echo/v4"
)

func main() {
	w1 := ws.CreateClient("/ws")
	w2 := ws.CreateClient("/live/Front_Deck")

	e := echo.New()
	e.GET("/*", login.GetLoginPage)
	e.POST("/auth", proxy.PostAuth)
	e.GET("/ws", w1.EstablishConnection)
	e.GET("/live/Front_Deck", w2.EstablishConnection)
	e.Logger.Fatal(e.Start(":8080"))
}
