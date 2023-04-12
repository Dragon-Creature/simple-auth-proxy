package proxy

import (
	"encoding/json"
	"fmt"
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/ws"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func HandleTraffic(c echo.Context) error {
	username, password := getToken(c)
	if checkToken(Credential{Username: username, Password: password}) {
		if c.IsWebSocket() {
			return handleWebSocket(c)
		}
		return handleHttpProxy(c)
	}
	err := sendLoginFiles(c)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func handleWebSocket(c echo.Context) error {
	path := c.Request().RequestURI
	w := ws.CreateClient(path)
	err := w.HandleWebsocket(c)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func handleHttpProxy(c echo.Context) error {
	client := http.Client{Timeout: time.Second * 5}
	request, err := http.NewRequest(c.Request().Method, fmt.Sprintf("http://localhost:30085%s", c.Request().RequestURI), nil)
	if err != nil {
		return errors.WithStack(err)
	}
	for key, value := range c.Request().Header {
		request.Header.Add(key, value[0])
	}
	response, err := client.Do(request)
	if err != nil {
		return errors.WithStack(err)
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
		return errors.WithStack(err)
	}
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

func sendLoginFiles(c echo.Context) error {
	filename := c.Request().RequestURI
	if filename == "/" {
		filename = "index.html"
	}
	filename = fmt.Sprintf("build/%s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return errors.WithStack(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return errors.WithStack(err)
	}
	err = c.HTML(http.StatusOK, string(data))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func PostAuth(c echo.Context) error {
	bodyData, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return errors.WithStack(err)
	}
	credential := Credential{}
	err = json.Unmarshal(bodyData, &credential)
	if err != nil {
		return errors.WithStack(err)
	}
	hashCredential, auth := checkCredentials(credential)
	if auth {
		setCookie(c, *hashCredential)
		err = c.Redirect(http.StatusMovedPermanently, "/")
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		err = c.String(http.StatusUnauthorized, "invalid credentials")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func setCookie(c echo.Context, credential Credential) {
	expire := time.Now().Add(20 * time.Minute) // Expires in 20 minutes
	cookie := http.Cookie{Name: "token", Value: fmt.Sprintf("%s:%s", credential.Username, credential.Password), Path: "/", Expires: expire, MaxAge: 86400}
	c.SetCookie(&cookie)
	cookie = http.Cookie{Name: "token", Value: fmt.Sprintf("%s:%s", credential.Username, credential.Password), Path: "/", Expires: expire, MaxAge: 86400, HttpOnly: true, Secure: true}
	c.SetCookie(&cookie)
}

func checkToken(credential Credential) bool {
	credentials, err := getPasswdFile()
	if err != nil {
		return false
	}
	for _, c := range credentials {
		if c.Username == credential.Username && c.Password == credential.Password {
			return true
		}
	}
	return false
}

func checkCredentials(credential Credential) (*Credential, bool) {
	credentials, err := getPasswdFile()
	if err != nil {
		return nil, false
	}
	for _, c := range credentials {
		if c.Username == credential.Username {
			err = checkPassword(c.Password, credential.Password)
			if err != nil {
				return nil, false
			}
			return &c, true
		}
	}
	return nil, false
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func getPasswdFile() ([]Credential, error) {
	data, err := os.ReadFile("htpasswd")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	lines := strings.Split(string(data), "\n")
	credentials := []Credential{}
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		credentials = append(credentials, Credential{
			Username: fields[0],
			Password: fields[1],
		})
	}
	return credentials, nil
}

func checkPassword(hashed string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
