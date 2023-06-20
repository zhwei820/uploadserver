package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	path         string
	port         string
	nameAuth     string
	passwordAuth string
)

func main() {
	flag.StringVar(&port, "port", "8000", "Port number")
	flag.StringVar(&path, "path", "nutshell", "File server path")
	flag.StringVar(&nameAuth, "name", "", "auth")
	flag.StringVar(&passwordAuth, "password", "", "auth")

	flag.Parse()

	http.HandleFunc("/", detector)

	fmt.Println("use: http://localhost:8000/index to upload file")
	fmt.Println("use: http://localhost:8000/index/ to view file")
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
