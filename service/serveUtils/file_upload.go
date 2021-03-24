package serveUtils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
)

func FileUploadService(r *http.Request) (string, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", err
	}

	file, handler, err := r.FormFile("image")

	if err != nil {
		return "", err
	}

	defer file.Close()

	mimeType := handler.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg", "image/jpg":
		name, err := saveFile(file, mimeType)
		if err != nil {
			return "", err
		}
		return name, nil
	case "image/png":
		name, err := saveFile(file, mimeType)
		if err != nil {
			return "", err
		}
		return name, nil
	default:
		return "", errors.New("unsupported file type")
	}
}

func saveFile(file multipart.File, mime string) (string, error) {
	name, err := getFileName(mime)
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile("./static/"+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	defer f.Close()
	io.Copy(f, file)

	return name, nil
}

func getFileName(mimes string) (string, error) {
	fileName := randToken(16)

	fileEndings, err := mime.ExtensionsByType(mimes)
	if err != nil {
		return "", err
	}

	name := fileName + fileEndings[0]

	return name, nil
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
