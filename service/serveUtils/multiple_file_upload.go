package serveUtils

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func MultipleFileUploadService(r *http.Request) ([]string, error) {
	err := r.ParseMultipartForm(200000)
	if err != nil {
		return nil, err
	}

	formdata := r.MultipartForm
	files := formdata.File["image"]

	var fileName []string

	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		if err != nil {
			return nil, err
		}

		defer file.Close()

		mimeType := files[i].Header.Get("Content-Type")
		switch mimeType {
		case "image/jpeg", "image/jpg":
			name, err := getFileName(mimeType)
			if err != nil {
				return nil, err
			}

			fileName = append(fileName, name)
			_, _ = saveFiles(file, name)
		case "image/png":
			name, err := getFileName(mimeType)
			if err != nil {
				return nil, err
			}

			fileName = append(fileName, name)
			_, _ = saveFiles(file, name)
		default:
			return nil, errors.New("unsupported file type")
		}
	}
	return fileName, nil
}

func saveFiles(file multipart.File, name string) (string, error) {
	f, err := os.OpenFile("./static/"+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}

	return name, nil
}
