package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	l "github.com/Guploader/log"
	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"uuid"`
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {

	l.Tlog.Info("checking Method of request")
	l.Tlog.Warn("Method can only be GET ")
	// filter the type of the request that the server can process
	if r.Method != http.MethodGet {
		l.Tlog.Error("Error : only GET request can be prosess by the server!")
		http.Error(w, "Error : only GET request can be prosess by the server!", http.StatusMethodNotAllowed)
	}

}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		l.Tlog.Error("failed , method post must called to the server!")
		http.Error(w, "Post request must be configured!", http.StatusMethodNotAllowed)
		return
	}

	// var file file

	file, header, err := r.FormFile("file")
	if err != nil {
		l.Tlog.Error("Reading file body has Error , failed!")
		http.Error(w, "Error for decoding the file ", http.StatusBadRequest)
		return
	}
	// closing the file
	defer file.Close()

	errDir := os.MkdirAll("./uploads", 0755)
	if errDir != nil {
		l.Tlog.Error("Makind directory upload have an error !")
		http.Error(w, "upload directory not found !", http.StatusInternalServerError)
		return
	}

	filepath.Base(header.Filename)

	fileExt := filepath.Ext(header.Filename)
	fileID := uuid.New().String() + fileExt

	filePath := filepath.Join("./uploads", filepath.Base(header.Filename) + fileID )

	destFile, err := os.Create(filePath)
	if err != nil {
		l.Tlog.Error("file creation on Server has Error !")
		http.Error(w, "Could not create file on server", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	_, errFile := io.Copy(destFile, file)
	if errFile != nil {
		l.Tlog.Error("copying uploaded file failed")
		http.Error(w, "Could not save uploaded file", http.StatusInternalServerError)
		return
	}

	
	

	l.Tlog.Info("file recived Sudccessfully!")
}

func Server(file file) {

	fs := http.FileServer(http.Dir("./frontend/"))
	http.Handle("/", fs)
	http.HandleFunc("/upload", UploadHandler)

	port := "8085"
	fmt.Println("Server is now running on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("Server can't start properly!")
	}
}
