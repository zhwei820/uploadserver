package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	ospath "path"
	"regexp"
	"strings"

	"github.com/chinglinwen/log"
)

// concern for name collision, filename need to be unique
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	relFilePath := r.FormValue("file_path")
	ip := strings.Split(r.RemoteAddr, ":")[0]
	uri := strings.NewReplacer("/uploadapi", "").Replace(r.RequestURI)

	fileUrl2 := ospath.Join(path, uri, relFilePath)

	var validPath = regexp.MustCompile(`.*\.\.\/.*`)
	if validPath.MatchString(fileUrl2) {
		fmt.Fprintf(w, "file path should not contain the two dot\n")
		return
	}

	if ip == "[" {
		ip = "127.0.0.1"
	}

	// no bigger than 10G
	err := r.ParseMultipartForm(10000000000)
	if err != nil {
		log.Println(w, err)
		return
	}

	formdata := r.MultipartForm

	// multipart parameter "file" need to specified when upload
	files := formdata.File["files"]

	if len(files) == 0 {
		fmt.Fprintf(w, "need to provide file(multipart form) or data\n")
		return
	}

	for i := range files {
		file, err := files[i].Open()
		if err != nil {
			log.Println(w, err)
			continue
		}
		defer file.Close()

		f := files[i].Filename

		fileUrl := ospath.Join(path, uri, relFilePath, f)
		fmt.Println("fileUrl", fileUrl)
		fmt.Println("uri", uri)

		dir := ospath.Dir(fileUrl)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Println(w, "mkdir error %v\n", err)
			continue
		}

		out, err := os.OpenFile(fileUrl, mode(false), 0644)
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing")
			continue
		}
		defer out.Close()

		n, err := io.Copy(out, file)
		if err != nil {
			log.Println(w, err)
			continue
		}

		note := "(replace old file if file exist)"

		fmt.Fprintf(w, "Files uploaded successfully : %v %v bytes %v\n", fileUrl, n, note)
		log.Println(ip, uri, f, "created")
	}
}

func mode(append bool) int {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if append {
		flags = flags | os.O_APPEND
		flags -= os.O_TRUNC
	}
	return flags
}
