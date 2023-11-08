package image

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func ServeImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageFileName := vars["filename"]

	imageDir := "./questionImage/" // Specify the directory where your images are stored

	imagePath := filepath.Join(imageDir, imageFileName)

	imageFile, err := os.Open(imagePath)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	defer imageFile.Close()

	contentType := "image/" + getFileExtension(imageFileName)
	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, imageFile)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getFileExtension(fileName string) string {
	return filepath.Ext(fileName)[1:]
}
