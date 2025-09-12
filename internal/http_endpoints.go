package internal

import (
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type Metadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	ImageType   string
}

var endpoint_handlers = make(map[string]func(*FileStore, string, http.ResponseWriter, *http.Request))

func registerEndpoints() {
	endpoint_handlers["photos"] = endpointPhotos
}

func endpointPhotos(store *FileStore, urlPart string, rspn http.ResponseWriter, rqst *http.Request) {
	if rqst.Method == http.MethodPost {
		putPhoto(store, urlPart, rspn, rqst)
	}

	if rqst.Method == http.MethodGet {
		getPhoto(store, urlPart, rspn, rqst)
	}
}

func getPhoto(store *FileStore, urlPart string, rspn http.ResponseWriter, rqst *http.Request) {
	if urlPart == "" {
		getPhotoIds(store, urlPart, rspn, rqst)
		return
	}
	uuid, err := uuid.Parse(urlPart)
	uuidstr := uuid.String()
	if err != nil {
		http.Error(rspn, "Unlabe to parse photo uuid from "+urlPart+" wrong endpoint?", http.StatusBadRequest)
		return
	}

	err = store.Client.FGetObject(rqst.Context(), "images", uuidstr, uuidstr, minio.GetObjectOptions{})
	if err != nil {
		werr(rspn, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer os.Remove(uuidstr)
	http.ServeFile(rspn, rqst, uuidstr)
}

func getPhotoIds(store *FileStore, urlPart string, rspn http.ResponseWriter, rqst *http.Request) {
	println("This endpoint gets a photo id!")
}

func putPhoto(store *FileStore, urlPart string, rspn http.ResponseWriter, rqst *http.Request) {
	if rqst.Header.Get("X-API-Key") != os.Getenv("ADMIN_SECRET") {
		werr(rspn, http.StatusUnauthorized)
		return
	}

	ctype := rqst.Header.Get("Content-Type")
	mtype, _, err := mime.ParseMediaType(ctype)
	if err != nil {
		werr(rspn, http.StatusUnsupportedMediaType)
		return
	}

	if mtype != "multipart/form-data" {
		werr(rspn, http.StatusUnsupportedMediaType)
		return
	}

	reader, err := rqst.MultipartReader()
	if err != nil {
		http.Error(rspn, err.Error(), http.StatusBadRequest)
		return
	}

	part, err := reader.NextPart()
	if err != nil {
		http.Error(rspn, "Unable to Process form parts", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if part.Header.Get("Content-Type") != "application/json" {
		http.Error(rspn, "No Json Content found header "+part.Header.Get("Content-Type"), http.StatusBadRequest)
		return
	}

	limreader := io.LimitReader(part, 25<<10) // json size must not exceed 25kbs
	data, err := io.ReadAll(limreader)
	if err != nil {
		http.Error(rspn, "Data read failed check to ensure your json file doesn't exceed 25kbs", http.StatusBadRequest)
		log.Println(rqst.RemoteAddr, " attempted to send json file of size ", len(data)/(1024), "kbs")
	}

	var metadata Metadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		http.Error(rspn, "Invalid JSON", http.StatusBadRequest)
		return
	}
	part.Close()

	part, err = reader.NextPart()
	if err != nil {
		http.Error(rspn, "Unable to Process form parts", http.StatusBadRequest)
		log.Println(err)
		return
	}

	imageType := part.Header.Get("Content-Type")
	if !strings.Contains(imageType, "image") {
		http.Error(rspn, "Content Type Not Image", http.StatusBadRequest)
		return
	}

	metadata.ImageType = imageType

	limreader = io.LimitReader(part, 50<<20)
	err = store.UploadFS(rqst.Context(), "images", &metadata, limreader)
	if err != nil {
		werr(rspn, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	part.Close()

	// var metadata Metadata
	// if err := json.Unmarshal([]byte(rqst.PostFormValue("metadata")), &metadata); err != nil {
	// 	http.Error(rspn, "Invalid JSON", http.StatusBadRequest)
	// 	return
	// }

	// println(metadata.Name)
	// println(metadata.Description)
	// println(strings.Join(metadata.Tags, ","))

	wstd(rspn, http.StatusOK)
}

func wstd(rspn http.ResponseWriter, code int) {
	rspn.WriteHeader(code)
}

func werr(rspn http.ResponseWriter, code int) {
	http.Error(rspn, http.StatusText(code), code)
}
