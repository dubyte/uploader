package upload

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Uploader struct {
	Multiple             bool
	RouteToHandleUploads string
	PathToSaveFiles      string
	Permission           int
}

var FormTemplate = defaultHTML
var Log = flogf

var defaultHTML = `<!DOCTYPE html>
<html>
	<head>
	<title>File Uploader</title>
	</head>
	<body>
		<form name="uploader" action="{{.RouteToHandleUploads}}", method="POST", enctype="multipart/form-data">
			<label for="myfile">Select a file:</label><br><br>
			<input type="file" name="toUpload" {{if .Multiple}} multiple {{end}}><br><br>
			<input type="submit">
		</form>
	</body>
</html>`

func (u Uploader) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := u.uploadFile(w, r)
		if err != nil {
			Log(w, "handler: uploadFile error: %v", err)
		}
	case http.MethodGet:
		err := u.showUploadForm(w, r)
		if err != nil {
			Log(w, "handler: showUploadForm error: %v", err)
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		Log(ioutil.Discard, "unsupported method: '%v' was called", r.Method)
	}
}

func (u Uploader) showUploadForm(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("content-type", "text/html")

	t := template.Must(template.New("DefaultHTML").Parse(FormTemplate))
	err := t.Execute(w, u)
	if err != nil {
		return fmt.Errorf("tempalte execute: %w", err)
	}

	return nil
}

func (u Uploader) uploadFile(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		return fmt.Errorf("parseMultipartForm: %w", err)
	}

	Log(w, "Starting upload ...\n")
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fHeader := range fileHeaders {
			filename := fHeader.Filename
			Log(w, "filename: %s\n", filename)
			Log(w, "\tfile size: %d\n", fHeader.Size)

			f, err := fHeader.Open()
			if err != nil {
				Log(w, "\topen error: %v\n", err)
				continue
			}

			fContent, err := ioutil.ReadAll(f)
			if err != nil {
				Log(w, "\tread file error: %v\n", err)
				continue
			}

			err = f.Close()
			if err != nil {
				Log(w, "\tclose file error: %v\n", err)
				continue
			}

			fpath := filepath.Join(u.PathToSaveFiles, filename)

			err = fileNotExists(fpath)
			if err != nil {
				Log(w, "\tfileNotExists: %v\n", err)
				continue
			}
			err = ioutil.WriteFile(fpath, fContent, os.FileMode(u.Permission))
			if err != nil {
				Log(w, "\twrite file error: %v\n", err)
				continue
			}

			Log(w, "\tfile uploaded successfully\n")
		}
	}
	return nil
}

func flogf(w io.Writer, format string, v ...interface{}) {
	log.Printf(format, v...)
	fmt.Fprintf(w, format, v...)
}

func fileNotExists(filename string) error {
	_, err := os.Stat(filename)
	if err == nil {
		return errors.New("file exists")
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("os.Stat: unexpected error: %w", err)
	} else {
		return nil
	}
}
