package internal

import (
	"log"
	"math"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type user struct {
	limiter     *rate.Limiter
	lastRequest time.Time
	// a behaviorScore > 0 indicates bad behavior
	// a behaviorScore = 0 is a neutral or good behavior
	behaviorScore int
	ulock         sync.RWMutex
}

var (
	users        = make(map[string]*user)
	rwlock       sync.RWMutex
	cleaningDone = make(chan struct{})
)

func onSiteVisit(ip string, auth bool) *user {
	rwlock.RLock()
	usr, ok := users[ip]
	rwlock.RUnlock()
	if !ok {
		return onFreshVisit(ip, auth)
	}

	if auth {
		// increment distrust
		usr.ulock.Lock()

		usr.behaviorScore++
		if usr.behaviorScore%10 == 0 { // we should only sometimes cause distrust to decrease level
			timeframe := time.Hour
			trust := time.Duration(float32(10) * float32(1/usr.behaviorScore))
			if trust == 0 {
				timeframe /= 12
				trust = 1
			}
			limiter := rate.NewLimiter(rate.Every(timeframe/trust), 5)
			log.Printf("Changed trust threshhold for %s new rate %s with trust of %d\n", ip, time.Second/trust, trust)
			usr.limiter = limiter
		}
		usr.ulock.Unlock()
	}

	log.Printf("%s visited with behavior score of %d\n", ip, usr.behaviorScore)
	return usr
}

func onFreshVisit(ip string, auth bool) *user {
	var limiter *rate.Limiter
	behavior := 0
	if auth {
		// don't "trust" initial auth requests as much
		limiter = rate.NewLimiter(rate.Every(time.Second/5), 5)
		behavior += 2
	} else {
		limiter = rate.NewLimiter(rate.Every(time.Second/10), 5)
	}
	user := &user{limiter: limiter, lastRequest: time.Now(), behaviorScore: behavior}
	rwlock.Lock()
	users[ip] = user
	rwlock.Unlock()
	return user
}

func cleanLimiters() {
	for {
		select {
		case <-time.After(1 * time.Minute):
			rwlock.Lock()
			for ip, usr := range users {
				usr.ulock.Lock()
				usr.behaviorScore = int(math.Max(float64(usr.behaviorScore-5), 0))
				usr.ulock.Unlock()
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
