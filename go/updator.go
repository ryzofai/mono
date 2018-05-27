package main

import (
"fmt"
"io"
"net/http"
"os"
"crypto/md5"
"strings"
"encoding/base64"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle all errors via http response

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	username, password, authOK := r.BasicAuth()
	if authOK == false {
		http.Error(w, "Not authorized", 401)
		return
	}
	if username != "secret user" || password != "secret password" {
		http.Error(w, "Not authorized", 401)
		return
	}

	savfile := r.Header.Get("file")
	file, err := os.Create(savfile)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(err.Error())))
	}
	fmt.Println(r.Header.Get("file"))

	info, err := os.Stat(r.Header.Get("file"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf(err.Error())))
	}
	fmt.Println(info.ModTime())

	n, err := io.Copy(file, r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(err.Error())))
	}

	if b, err := ComputeMd5(savfile); err != nil {
		w.Write([]byte(fmt.Sprintf(err.Error())))
	} else {
		w.Write([]byte(fmt.Sprintf("md5 checksum is: %x", b)))
		w.Write([]byte(fmt.Sprintf("\n%d bytes are recieved\n", n)))
		w.Write([]byte(fmt.Sprintf("file timestamp:", info.ModTime())))
	}
	file.Close()
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("C:\\Projects\\Go"))
	http.Handle("/static", http.StripPrefix("/static/", fs))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/st", auth(staticHandler))
	fs := http.FileServer(http.Dir("C:\\Projects\\Go"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	//http.HandleFunc("/st", basicAuth(staticHandler))
	http.ListenAndServe(":5050", nil)
}

func ComputeMd5(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}
	return hash.Sum(result), nil
}

func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if user !="test" && pass != "pass" {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		fn(w, r)
	}
}

type handler func(w http.ResponseWriter, r *http.Request)

func basicAuth(pass handler) handler {

    return func(w http.ResponseWriter, r *http.Request) {
        
        auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

        if len(auth) != 2 || auth[0] != "Basic" {
            http.Error(w, "authorization failed", http.StatusUnauthorized)
            return
        }

        payload, _ := base64.StdEncoding.DecodeString(auth[1])
        pair := strings.SplitN(string(payload), ":", 2)

        if len(pair) != 2 || !validate(pair[0], pair[1]) {
            http.Error(w, "authorization failed", http.StatusUnauthorized)
            return
        }

        pass(w, r)
    }
}

func validate(username, password string) bool {
    if username == "test" && password == "test" {
        return true
    }
    return false
}

func GetOnly(h handler) handler {

    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            h(w, r)
            return
        }
        http.Error(w, "get only", http.StatusMethodNotAllowed)
    }
}
