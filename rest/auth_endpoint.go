package rest

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func AuthEndpoints(store *FileStore, urlPart string, rspn http.ResponseWriter, rqst *http.Request) {
	switch rqst.Method {
	case http.MethodGet:
		verifyAuth(store, urlPart, rspn, rqst)
	case http.MethodDelete:

	default:
		werr(rspn, http.StatusBadRequest)
	}
}

func verifyAuth(_ *FileStore, _ string, rspn http.ResponseWriter, rqst *http.Request) {
	valid := false

	if rqst.Header.Get("X-API-Key") == os.Getenv("ADMIN_SECRET") {
		valid = true
		// valid can decrement hatred counter :3
		ip, _, err := net.SplitHostPort(rqst.RemoteAddr)
		if err != nil {
			rspn.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("Valid login from %s\n", ip)
	}

	rslt := fmt.Sprintf(`{"valid": %v}`, valid)
	rspn.Header().Set("Content-Type", "application/json")
	rspn.WriteHeader(http.StatusOK)
	rspn.Write([]byte(rslt))
}
