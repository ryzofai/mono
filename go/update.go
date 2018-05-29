package main

import (
"fmt"
"io"
"net/http"
"os"
"crypto/md5"
"io/ioutil"
"time"

"github.com/abbot/go-http-auth"
"github.com/tkanos/gonfig"
)

type Configuration struct {
	Port			string
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	listdir()
	savfile := r.Header.Get("filename")
	file, err := os.Create(savfile)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(err.Error())))
	}
	fmt.Println(r.Header.Get("filename"))

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
		w.Write([]byte(fmt.Sprintf("file modified date:", info.ModTime())))
	}
	file.Close()
}

func main() {
	//http.HandleFunc("/upload", uploadHandler)

	authenticator := auth.NewBasicAuthenticator("localhost", Secret)
	http.HandleFunc("/hb/", auth.JustCheck(authenticator,heartbeat))
	http.HandleFunc("/ud/", auth.JustCheck(authenticator, uploadHandler))
	http.HandleFunc("/md/", auth.JustCheck(authenticator, md5Handler))
	http.HandleFunc("/ls/", auth.JustCheck(authenticator, handleFileServer("./", "/ls/")))

	configuration := Configuration{}
	err := gonfig.GetConf("update.cfg", &configuration)
	if err != nil {  
		fmt.Println(err)
	}
	http.ListenAndServe(":" + configuration.Port, nil)
}

func md5Handler(w http.ResponseWriter, r *http.Request) {
	file, ok := r.URL.Query()["file"]
	f := file[0]
	
    
    if !ok || len(file) < 1 {
        fmt.Println("\nUrl Param 'file' is missing")
        w.Write([]byte(fmt.Sprintf("\nUrl Param 'file' is missing")))
    } else {
    	if b, err := ComputeMd5(f); err != nil {
    		w.Write([]byte(fmt.Sprintf(err.Error())))
    	} else {
    		w.Write([]byte(fmt.Sprintf("md5 checksum is: %x", b)))
    		info, err := os.Stat(f)
    		if err != nil {
    			w.Write([]byte(fmt.Sprintf(err.Error())))
    		}
    		w.Write([]byte(fmt.Sprintf("\nfile modified date:", info.ModTime())))
    	}
    
    w.Write([]byte(fmt.Sprintf("\nUrl Param 'file' is: " + string(f))))
    }
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
		fmt.Println(req.URL)
		realHandler(w, req)
	}
}

func Secret(user, realm string) string {
	users := map[string]string{
		"secret_admin": "$2a$14$Rxe926OsV6RaWT1RZZWzceFzQaY28XIFoM2s9.ICVhUJcC3ZHtgL2",
	}

	if a, ok := users[user]; ok {
		return a
	}
	return ""
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("status: up \n" + "time now: " + time.Now().String()  )))
}
