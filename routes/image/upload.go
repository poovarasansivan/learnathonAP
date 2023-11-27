package image

import (
	"io"
	"net/http"
	"os"
)

func Upload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("UploadFiles")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close() 

	// var response map[string]interface{}

	outputFile, err := os.Create("./questionImage/" + handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()
	io.Copy(outputFile, file)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Image uploaded successfully"))
}
