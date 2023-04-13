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
	"path/filepath"
	"strings"
	"time"
)

type Proxy struct {
	TargetProtocol string
	TargetURL      string
	HtpasswdFile   string
	CookieMaxAge   int
}

func (p *Proxy) HandleTraffic(c echo.Context) error {
	username, password := getToken(c)
	if p.checkToken(Credential{Username: username, Password: password}) {
		if c.IsWebSocket() {
			err := p.handleWebSocket(c)
			if err != nil {
				c.Logger().Error(err)
			}
			return err
		}
		err := p.handleHttpProxy(c)
		if err != nil {
			err = errors.WithStack(err)
			c.Logger().Error(err)
		}
		return nil
	}
	err := sendLoginFiles(c)
	if err != nil {
		c.Logger().Error(errors.WithStack(err))
		return errors.WithStack(err)
	}
	return nil
}

func (p *Proxy) handleWebSocket(c echo.Context) error {
	path := c.Request().RequestURI
	w := ws.CreateClient(p.TargetURL, path)
	err := w.HandleWebsocket(c)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (p *Proxy) handleHttpProxy(c echo.Context) error {
	client := http.Client{Timeout: time.Second * 5}
	request, err := http.NewRequest(c.Request().Method, fmt.Sprintf("%s%s%s", p.TargetProtocol, p.TargetURL, c.Request().RequestURI), nil)
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
	ext := filepath.Ext(filename)

	fileExtToContentType := map[string]string{
		"html": "text/html",
		"css":  "text/css",
		"json": "application/json",
		"txt":  "text/plain",
		"png":  "image/png",
		"js":   "application/javascript",
		"ico":  "image/x-icon",
	}
	err = c.Blob(http.StatusOK, fileExtToContentType[ext], data)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (p *Proxy) PostAuth(c echo.Context) error {
	bodyData, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return errors.WithStack(err)
	}
	credential := Credential{}
	err = json.Unmarshal(bodyData, &credential)
	if err != nil {
		err = errors.WithStack(err)
		c.Logger().Error(err)
		return err
	}
	hashCredential, auth := p.checkCredentials(credential)
	if auth {
		p.setCookie(c, *hashCredential)
		err = c.Redirect(http.StatusMovedPermanently, "/")
		if err != nil {
			err = errors.WithStack(err)
			c.Logger().Error(err)
		}
		return nil
	} else {
		err = c.String(http.StatusUnauthorized, "invalid credentials")
		if err != nil {
			err = errors.WithStack(err)
			c.Logger().Error(err)
			return err
		}
	}

	return nil
}

func (p *Proxy) setCookie(c echo.Context, credential Credential) {
	expire := time.Now().Add(time.Duration(p.CookieMaxAge) * time.Second)
	cookie := http.Cookie{Name: "token", Value: fmt.Sprintf("%s:%s", credential.Username, credential.Password), Path: "/", Expires: expire, MaxAge: p.CookieMaxAge}
	c.SetCookie(&cookie)
	cookie = http.Cookie{Name: "token", Value: fmt.Sprintf("%s:%s", credential.Username, credential.Password), Path: "/", Expires: expire, MaxAge: p.CookieMaxAge, HttpOnly: true, Secure: true}
	c.SetCookie(&cookie)
}

func (p *Proxy) checkToken(credential Credential) bool {
	credentials, err := p.getPasswdFile()
	if err != nil {
		fmt.Println(errors.WithStack(err))
		return false
	}
	for _, c := range credentials {
		if c.Username == credential.Username && c.Password == credential.Password {
			return true
		}
	}
	return false
}

func (p *Proxy) checkCredentials(credential Credential) (*Credential, bool) {
	credentials, err := p.getPasswdFile()
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

func (p *Proxy) getPasswdFile() ([]Credential, error) {
	data, err := os.ReadFile(p.HtpasswdFile)
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
