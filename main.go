package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	path         string
	port         string
	nameAuth     string
	passwordAuth string
)

func cleanOldFiles() {
	fmt.Println("cleanOldFiles start")
	// 调用find命令删除30天以上的文件
	cmd := exec.Command("bash", "clean.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func cronCleanOldFiles() {
	cleanOldFiles()

	fmt.Println("cronCleanOldFiles start")
	ticker := time.NewTicker(12 * time.Hour)
	for range ticker.C {
		cleanOldFiles()
	}
}

func main() {
	flag.StringVar(&port, "port", "8000", "Port number")
	flag.StringVar(&path, "path", "nutshell", "File server path")
	flag.StringVar(&nameAuth, "name", "", "auth")
	flag.StringVar(&passwordAuth, "password", "", "auth")

	flag.Parse()

	http.HandleFunc("/", detector)

	fmt.Println("use: http://localhost:8000/index to upload file")
	fmt.Println("use: http://localhost:8000/index/ to view file")

	go cronCleanOldFiles()

	_ = http.ListenAndServe(":"+port, nil)
}

func detector(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.RequestURI, "upload") {
		uploadHandler(w, r)
		return
	}
	if nameAuth != "" {
		name, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, "need log in "))
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			if !(name == nameAuth && password == passwordAuth) {
				w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, "need log in "))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}

	// print logs
	ip := strings.Split(r.RemoteAddr, ":")[0]
	log.Println(ip, r.RequestURI, "visited")

	if strings.HasSuffix(r.RequestURI, "index") {
		uploadPageHandler(w, r)
		return
	}
	http.FileServer(http.Dir(path)).ServeHTTP(w, r)
}
