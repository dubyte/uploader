package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dubyte/uploader/upload"
)

func main() {
	port := flag.String("port", "8080", "set the port to listen")
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

	log.Printf("uploader: listening on port %s", *port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
