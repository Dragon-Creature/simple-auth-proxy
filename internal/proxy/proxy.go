package proxy

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func PostAuth(c echo.Context) error {
	setCookie(c)
	err := c.Redirect(http.StatusMovedPermanently, "/")
	if err != nil {
		return err
	}
	return nil
}

func setCookie(c echo.Context) {
	expire := time.Now().Add(20 * time.Minute) // Expires in 20 minutes
	cookie := http.Cookie{Name: "token", Value: "username:password", Path: "/", Expires: expire, MaxAge: 86400}
	c.SetCookie(&cookie)
	cookie = http.Cookie{Name: "token", Value: "username:password", Path: "/", Expires: expire, MaxAge: 86400, HttpOnly: true, Secure: true}
	c.SetCookie(&cookie)
}
