package internal

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type user struct {
	limiter     *rate.Limiter
	lastRequest time.Time
}

var (
	users        = make(map[string]*user)
	rwlock       sync.RWMutex
	cleaningDone = make(chan struct{})
)

func SetupHttpServer(store *FileStore) {
	go cleanLimiters()

	http.Handle("/", http.FileServer(http.Dir("./resources/assets/public")))
	http.HandleFunc("/api/v1/", func(response http.ResponseWriter, request *http.Request) {
		handlePathing(store, strings.ReplaceAll(request.URL.String(), "/api/v1/", ""), response, request)
	})

	registerEndpoints()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

	close(cleaningDone)
}

func handlePathing(store *FileStore, urlPart string, rsp http.ResponseWriter, rqst *http.Request) {
	ip, _, err := net.SplitHostPort(rqst.RemoteAddr)
	if err != nil {
		rsp.WriteHeader(http.StatusInternalServerError)
		return
	}

	limiter := onVisit(ip)
	if !limiter.Allow() {
		http.Error(rsp, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}

	result, ok := endpoint_handlers[urlPart]
	if !ok {
		http.Error(rsp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	rsp.WriteHeader(result(store, urlPart, rsp, rqst))
}

func onVisit(ip string) *rate.Limiter {
	rwlock.RLock()

	usr, ok := users[ip]
	rwlock.RUnlock()
	if !ok {
		rwlock.Lock()
		limiter := rate.NewLimiter(rate.Every(time.Second/10), 5)
		users[ip] = &user{limiter: limiter, lastRequest: time.Now()}
		rwlock.Unlock()
		return limiter
	}

	rwlock.Lock()
	usr.lastRequest = time.Now()
	rwlock.Unlock()

	return usr.limiter
}

func cleanLimiters() {
	for {
		select {
		case <-time.After(1 * time.Minute):
			rwlock.Lock()
			for ip, usr := range users {
				if time.Since(usr.lastRequest) >= 5*time.Minute {
					delete(users, ip)
				}
			}
			rwlock.Unlock()
		case <-cleaningDone:
			return
		}
	}
}
