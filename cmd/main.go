package main

import (
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/login"
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/proxy"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/*", login.GetLoginPage)
	e.POST("/auth", proxy.PostAuth)
	e.Logger.Fatal(e.Start(":8080"))
}
