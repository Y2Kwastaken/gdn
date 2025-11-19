package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func PhotoEndpoints(store *FileStore, urlPart string, rspn http.ResponseWriter, rqst *http.Request) {
	switch rqst.Method {
	case http.MethodPut:
		putPhoto(store, urlPart, rspn, rqst)
	case http.MethodGet:
		getPhoto(store, urlPart, rspn, rqst)
	case http.MethodDelete:
		delPhoto(store, urlPart, rspn, rqst)
	default:
		werr(rspn, http.StatusBadRequest)
	}

}

func delPhoto(store *FileStore, urlPart string, rspn http.ResponseWriter, rqst *http.Request) {
	if rqst.Header.Get("X-API-Key") != os.Getenv("ADMIN_SECRET") {
		werr(rspn, http.StatusUnauthorized)
		return
	}

	uuid, err := uuid.Parse(urlPart)
	if err != nil {
		http.Error(rspn, "Unable to parse photo uuid from"+urlPart+" wrong endpoint?", http.StatusBadRequest)
		return
	}
	uuidstr := uuid.String()

	err = store.Client.RemoveObject(rqst.Context(), "images", uuidstr, minio.RemoveObjectOptions{})
	if err != nil {
		werr(rspn, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = store.Database.DeleteImage(uuid)
	if err != nil {
		werr(rspn, http.StatusInternalServerError)
		log.Println("[SEVERE] Couldn't remove image of uuid ", uuid, "Manual removal might be required ", err)
		// this is bad we should try atleast 10 more times otherwise note this in logs

		for i := range 10 {
			err = store.Database.DeleteImage(uuid)
			if err != nil {
				log.Println("[SEVERE] ", i, "/10", " Couldn't remove image of uuid ", uuid, "Manual removal might be required ", err)
			}
		}

		log.Println("[SEVERE] unable to delete the uuid ", uuid, " from the database Manual removal IS REQUIRED")
		return
	}

	wstd(rspn, http.StatusOK)
}

func getPhoto(store *FileStore, urlPart string, rspn http.ResponseWriter, rqst *http.Request) {
	if urlPart == "" {
		getPhotoIds(store, urlPart, rspn, rqst)
		return
	}
	uuid, err := uuid.Parse(urlPart)
	if err != nil {
		http.Error(rspn, "Unlabe to parse photo uuid from "+urlPart+" wrong endpoint?", http.StatusBadRequest)
		return
	}
	uuidstr := uuid.String()

	meta, err := store.Database.QueryImage(uuid)
	if err != nil {
		werr(rspn, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if meta == nil {
		http.Error(rspn, "No image with uuid "+uuidstr, http.StatusBadRequest)
		log.Println("No image with uuid", uuidstr)
		return
	}

	path := meta.ImageName
	err = store.Client.FGetObject(rqst.Context(), "images", uuidstr, path, minio.GetObjectOptions{})
	if err != nil {
		werr(rspn, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer os.Remove(path)
	defer log.Println("disposed " + path)
	log.Println("serving ", path)

	rspn.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, path))
	http.ServeFile(rspn, rqst, path)
}

func getPhotoIds(store *FileStore, _ string, rspn http.ResponseWriter, rqst *http.Request) {
	query := rqst.URL.Query()
	limit := -1
	offset := 0
	entries := -1

	if query.Has("limit") {
		rslt, err := strconv.Atoi(query.Get("limit"))
		if err != nil {
			werr(rspn, http.StatusBadRequest)
			log.Println(err)
			return
		}

		if rslt < 1 || rslt > 20 {
			werr(rspn, http.StatusBadRequest)
			log.Printf("limit out of bounds: %d\n", rslt)
			return
		}
		limit = rslt
	} else {
		limit = 20
	}

	if query.Has("offset") {
		rslt, err := strconv.Atoi(query.Get("offset"))
		if err != nil {
			werr(rspn, http.StatusBadRequest)
			log.Println(err)
			return
		}

		if rslt < 0 {
			werr(rspn, http.StatusBadRequest)
			log.Printf("offset out of bounds: %d\n", offset)
			return
		}
		offset = rslt
	}

	if query.Has("entries") {
		rslt, err := strconv.Atoi(query.Get("entries"))
		if err != nil {
			werr(rspn, http.StatusBadRequest)
			log.Println(err)
			return
		}

		if rslt != 1 {
			werr(rspn, http.StatusBadRequest)
			log.Println("entries in request is not 1")
			return
		}

		entries = rslt
	}

	var uuids []uuid.UUID
	if limit != -1 {
		rslt, err := store.Database.QueryIds(limit, offset)
		if err != nil {
			werr(rspn, http.StatusInternalServerError)
			log.Println(err)
			return
		}
		uuids = rslt
	}

	if entries != -1 {
		rslt, err := store.Database.CountEntries()
		if err != nil {
			werr(rspn, http.StatusInternalServerError)
			log.Println(err)
			return
		}
		entries = rslt

		if offset >= entries {
			http.Error(rspn, "offset greater than or equal to total entry length", http.StatusBadRequest)
			return
		}
	}

	response := IdResponse{Ids: uuids, Entries: entries}

	rspn.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rspn).Encode(response); err != nil {
		werr(rspn, http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func putPhoto(store *FileStore, _ string, rspn http.ResponseWriter, rqst *http.Request) {
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

	if metadata.Title == "" {
		http.Error(rspn, "Title Required", http.StatusBadRequest)
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

	log.Println("Uploaded Image", metadata.Title)
	wstd(rspn, http.StatusOK)
}
