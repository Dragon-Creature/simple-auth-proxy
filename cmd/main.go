package main

import (
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/env"
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/proxy"
	"github.com/labstack/echo/v4"
	"strconv"
)

func main() {
	e := echo.New()
	targetProtocol := env.GetEnvOrDefault("TARGET_PROTOCOL", "http://")
	targetURL := env.GetEnvOrDefault("TARGET_URL", "localhost:30085")
	htpasswdFile := env.GetEnvOrDefault("HTPASSWD_FILE", "htpasswd")
	cookieMaxAge := env.GetEnvOrDefault("COOKIE_MAX_AGE", "86400") //default is 1 day
	cookieMaxAgeInt, err := strconv.Atoi(cookieMaxAge)
	if err != nil {
		e.Logger.Fatal(err)
	}

	p := proxy.Proxy{
		TargetProtocol: targetProtocol,
		TargetURL:      targetURL,
		HtpasswdFile:   htpasswdFile,
		CookieMaxAge:   cookieMaxAgeInt,
	}

	e.GET("/*", p.HandleTraffic)
	e.POST("/auth", p.PostAuth)
	e.Logger.Fatal(e.Start(":8080"))
}
