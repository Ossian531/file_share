package api

import (
	"fmt"
	"net/http"
    "encoding/json"

	"file_share/storage"
)


type PresignResponse struct {
    URL string `json:"url"`
}


func health(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "API upp and running :)")
}


func upload(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("filename")
    if filename == "" {
		fmt.Println("Missing filename")
        http.Error(w, "Missing filename", http.StatusBadRequest)
        return
    }

    url, err := storage.GeneratePresignedUploadURL(r.Context(), filename)
    if err != nil {
		fmt.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(PresignResponse{URL: url})
}


func download(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("filename")
    if filename == "" {
        http.Error(w, "Missing filename", http.StatusBadRequest)
        return
    }

    url, err := storage.GeneratePresignedDownloadURL(r.Context(), filename)
    if err != nil {
		fmt.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(PresignResponse{URL: url})
}


func listObjects(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	files, err := storage.ListObjects(r.Context(), prefix)

	if err != nil {
		fmt.Println(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(files)
}


func Mux() *http.ServeMux {

	api := http.NewServeMux()

	api.HandleFunc("/health", health)

	api.HandleFunc("/upload", upload)
	api.HandleFunc("/download", download)
	api.HandleFunc("/list", listObjects)

	return api
}
