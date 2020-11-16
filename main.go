package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dubyte/uploader/upload"
)

func main() {
	host := flag.String("host", "localhost", "the host where the web server will listen")
	port := flag.String("port", "8080", "sets the port to listen")
	route := flag.String("route", "/upload", "the route to upload files")
	multiple := flag.Bool("multiple", false, "select file allows multiple files")
	pathToSaveFiles := flag.String("save-path", "/tmp", "the path to store the uploaded files")
	perms := flag.Int("perms", 0644, "permissions of the files uploaded")

	flag.Parse()

	log.Print("uploader: starting server...")

	u := upload.Uploader{
		Multiple:             *multiple,
		PathToSaveFiles:      *pathToSaveFiles,
		Permission:           *perms,
		RouteToHandleUploads: *route,
	}

	http.HandleFunc(*route, u.Handler)

	addr := fmt.Sprintf("%s:%s", *host, *port)

	log.Printf("uploader: listening on: %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
