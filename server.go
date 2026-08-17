package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	l "github.com/Guploader/log"
	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"uuid"`
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type FileMapper map[string]string

var fileMapper = make(FileMapper)

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

	l.Tlog.Info("checking method of the Reqeust!")
	if r.Method != http.MethodPost {
		l.Tlog.Error("failed , method post must called to the server!")
		http.Error(w, "Post request must be configured!", http.StatusMethodNotAllowed)
		return
	}

	l.Tlog.Info("method checked Successfully!")
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

	fileExt := filepath.Ext(header.Filename)
	fileName := filepath.Base(strings.TrimSuffix(header.Filename, fileExt))
	fileID := uuid.New().String()
	filePath := filepath.Join("./uploads", fileName+"_"+fileID+fileExt)

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

	// adding key value of the id and filepath to our map
	fileMapper[fileID] = filePath

	// generating randomize link for user
	url := fmt.Sprintf("http://localhost:8085/d/%s", fileID)

	// setting up header for response
	w.Header().Set("Content-Type", "application/json")

	res := FileMapper{
		"url": url,
	}
	encodeErr := json.NewEncoder(w).Encode(res)
	if encodeErr != nil {
		l.Tlog.Error("Can't Encode url to json object for sending to frontend!")
		return
	}

	// URLchecker := make(map[string]string)

}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// validating the URL‌ request method

	if r.Method != http.MethodGet {
		l.Tlog.Error("Method of the URL must be GET in this route!")
		http.Error(w, "Error for bad Request ", http.StatusBadRequest)
		return
	}
	// getting url path first

	urlPath := r.URL.Path
	reqFileID := strings.TrimPrefix(urlPath, "/d/")

	idValue, exist := fileMapper[reqFileID]
	if !exist {
		l.Tlog.Error("file does not exists on the server !")
		http.Error(w, "File does not exists on our server !", http.StatusNotFound)
		return
	}

	fmt.Println("idValue:", idValue)
	fmt.Println("downloadName:", filepath.Base(idValue))


	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(idValue)))

	// if file exists sending to the server
	http.ServeFile(w, r, idValue)

}
func Server() {

	fs := http.FileServer(http.Dir("./frontend/"))
	http.Handle("/", fs)
	// /upload
	http.HandleFunc("/upload", UploadHandler)
	// /d
	http.HandleFunc("/d/", DownloadHandler)
	port := "8085"
	fmt.Println("Server is now running on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("Server can't start properly!")
	}
}
