package httpserv

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Y2Kwastaken/gdn/rest"
)

var (
	endpoint_handlers = make(map[string]func(*FileStore, string, http.ResponseWriter, *http.Request))
)

func registerEndpoints() {
	endpoint_handlers["photos"] = rest.PhotoEndpoints
	endpoint_handlers["auth"] = rest.AuthEndpoints
}

func SetupHttpServer(store *FileStore) {
	go cleanLimiters()

	http.Handle("/", http.FileServer(http.Dir("./resources/assets/public")))
	http.HandleFunc("/api/v1/", func(response http.ResponseWriter, request *http.Request) {
		handlePathing(store, strings.ReplaceAll(request.URL.String(), "/api/v1/", ""), response, request)
	})

	registerEndpoints()

	log.Println("GDN open on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}

	close(cleaningDone)
}

func handlePathing(store *FileStore, urlPart string, rsp http.ResponseWriter, rqst *http.Request) {
	ip, _, err := net.SplitHostPort(rqst.RemoteAddr)
	if err != nil {
		rsp.WriteHeader(http.StatusInternalServerError)
		return
	}

	auth := strings.Contains(urlPart, "auth")
	usr := onSiteVisit(ip, auth)
	usr.ulock.RLock()
	if usr.behaviorScore >= 50 {
		usr.ulock.RUnlock()
		http.Error(rsp, "Temporarily Banned", http.StatusForbidden)
		return
	}
	if !usr.limiter.Allow() {
		usr.ulock.RUnlock()
		http.Error(rsp, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}
	usr.ulock.RUnlock()

	pathParts := strings.Split(strings.Trim(rqst.URL.Path, "/"), "/")
	if len(pathParts) == 0 {
		http.Error(rsp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resource := pathParts[2] // 0 api 1 v1 2 [resource]

	result, ok := endpoint_handlers[resource]
	if !ok {
		http.Error(rsp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var remainderPath string
	if len(pathParts) > 1 {
		remainderPath = strings.Join(pathParts[3:], "/")
	} else {
		remainderPath = ""
	}

	result(store, remainderPath, rsp, rqst)
}
