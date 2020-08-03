// +build ignore

package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// $ unzip-server -h 0.0.0.0:8100
var (
	host  = flag.String("h", "0.0.0.0:8100", "host address")
	zipRc *zip.ReadCloser
)

// public.zip
func publicZip(writer http.ResponseWriter, req *http.Request) {
	name := req.RequestURI[1:]
	if len(name) == 0 {
		name = "index.html"
	}

	b := false
	for _, file := range zipRc.File {
		fmt.Printf("file.Name = %s\n", file.Name)
		if file.Name == name {
			send(writer, file, name)
			b = true
			break
		}
	}

	// 404
	if b == false {
		fmt.Printf("404 of %s\n", name)
		writer.WriteHeader(404)
	}
}

func send(writer http.ResponseWriter, file *zip.File, name string) {
	fmt.Printf("Contents %s\n", name)
	writer.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(name)))

	rc, _ := file.Open()
	defer rc.Close()

	if _, err := io.Copy(writer, rc); err != nil {
		log.Panic(err)
	}
}

//
func main() {
	flag.Parse()
	fmt.Fprintf(os.Stderr, "usage: unzip-server -h %v\n", *host)

	if _, err := os.Stat("public.zip"); os.IsNotExist(err) {
		log.Fatal(err)
	} else {
		zipRc, _ = zip.OpenReader("public.zip")
	}
	defer zipRc.Close()

	http.HandleFunc("/", publicZip)

	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal(err)
	}
}
