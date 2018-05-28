package main

import (
"fmt"
"io"
"net/http"
"os"
"crypto/md5"
"io/ioutil"
"log"
"github.com/abbot/go-http-auth"
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
	listdir()
	savfile := r.Header.Get("filename")
	file, err := os.Create(savfile)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(err.Error())))
	}
	fmt.Println(r.Header.Get("file"))

	info, err := os.Stat(r.Header.Get("filename"))
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
	authenticator := auth.NewBasicAuthenticator("localhost", Secret)
	http.HandleFunc("/static/", auth.JustCheck(authenticator, handleFileServer("D:\\Golanger\\update.exe\\", "/static/")))
	//fs := http.FileServer(http.Dir("D:\\Golanger\\update.exe\\"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))
	//http.HandleFunc("/listdir", auth(listdHandler))
	//http.HandleFunc("/stream", auth(streamHandler))
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

/*func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if user !="test" && pass != "pass" {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		fn(w, r)
	}
}*/

func listdir() {
    files, err := ioutil.ReadDir("./")
    if err != nil {
        fmt.Println(err)
    }

    for _, f := range files {
            fmt.Println(f.Name())
    }
}

func handleFileServer(dir, prefix string) http.HandlerFunc {
    fs := http.FileServer(http.Dir(dir))
    realHandler := http.StripPrefix(prefix, fs).ServeHTTP
    return func(w http.ResponseWriter, req *http.Request) {
        log.Println(req.URL)
        realHandler(w, req)
    }
}

func Secret(user, realm string) string {
    users := map[string]string{
        "john": "$2a$14$qtJ1USY91vKe2fOWXD2piuzLaraB1H8EiajP0C5RsicEZz.4Epgnu", //hello
    }

    if a, ok := users[user]; ok {
        return a
    }
    return ""
}
