package auth

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetLoginPage(w http.ResponseWriter, r *http.Request) {
	filename := r.RequestURI
	if filename == "/" {
		filename = "index.html"
	}
	filename = fmt.Sprintf("build/%s", filename)
	file, err := os.Open(filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
