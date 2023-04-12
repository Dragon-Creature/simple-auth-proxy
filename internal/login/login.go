package login

import (
	"fmt"
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/ws"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetLoginPage(c echo.Context) error {
	username, password := getToken(c)
	if username == "username" && password == "password" {
		if c.IsWebSocket() {
			path := c.Request().RequestURI
			w := ws.CreateClient(path)
			err := w.HandleWebsocket(c)
			if err != nil {
				return err
			}
			return nil
		}
		client := http.Client{Timeout: time.Second * 5}
		request, err := http.NewRequest(c.Request().Method, fmt.Sprintf("http://localhost:30085%s", c.Request().RequestURI), nil)
		if err != nil {
			panic(err)
		}
		for key, value := range c.Request().Header {
			request.Header.Add(key, value[0])
		}
		response, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		data, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		for key, value := range response.Header {
			c.Response().Header().Set(key, value[0])
		}
		err = c.String(http.StatusOK, string(data))
		if err != nil {
			panic(err)
		}
		return nil
	}
	sendFiles(c)
	return nil
}

func getToken(c echo.Context) (string, string) {
	for _, cookie := range c.Cookies() {
		if cookie.Name == "token" {
			values := strings.Split(cookie.Value, ":")
			username := values[0]
			password := values[1]
			return username, password
		}
	}
	return "", ""
}

func sendFiles(c echo.Context) {
	filename := c.Request().RequestURI
	if filename == "/" {
		filename = "index.html"
	}
	filename = fmt.Sprintf("build/%s", filename)
	file, err := os.Open(filename)
	if err != nil {
		//c.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.String(http.StatusOK, string(data))
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
