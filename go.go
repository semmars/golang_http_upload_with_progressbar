package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.HandleFunc("/upload", uploadFile)
	//http.ListenAndServe(":2020", nil)

	log.Println("Listening on :6688...")
	err := http.ListenAndServe(":6688", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func uploadFile(w http.ResponseWriter, r *http.Request) {

	//fmt.Fprintf(w, "Uploading File")
	r.ParseMultipartForm(1000 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	//abs,err :=filepath.Abs("/public/uploads/")
	uploadPath := "./public/uploads/"

	newPath := filepath.Join(uploadPath, handler.Filename)
	fmt.Printf("#{newPath}\n")

	dst, err := os.Create(newPath)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jmsg := "{\"message\": \"" + handler.Filename + "\"}"

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jmsg))

	//jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")

	//fmt.Fprintf(w, "Successfully Uploaded File\n")

}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
