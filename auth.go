package bla

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

func fetchIP(ra string) string {

	i := strings.LastIndex(ra, ":")
	if i == -1 {
		return "unknown"
	}
	return ra[:i]
}

type authRateByIPHandler struct {
	origin http.Handler
	ticker *time.Ticker
	record map[string]int
	mu     sync.RWMutex

	username, password string
	limit              int
	realm              string
}

func NewAuthRateByIPHandler(realm string, origin http.Handler, username, password string, limit int) *authRateByIPHandler {

	ticker := time.NewTicker(time.Minute)

	a := &authRateByIPHandler{origin,
		ticker,
		map[string]int{},
		sync.RWMutex{},

		username,
		password,
		limit,
		realm,
	}

	go func() {
		for {
			<-a.ticker.C
			a.mu.Lock()
			for k, _ := range a.record {
				delete(a.record, k)
			}
			a.mu.Unlock()
		}
	}()
	return a
}

func (a *authRateByIPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ip := fetchIP(r.RemoteAddr)
	a.mu.RLock()
	rec := a.record[ip]
	a.mu.RUnlock()
	if rec > a.limit {
		fmt.Fprintf(w, "<h1>Too many request</h1>")
		w.WriteHeader(429)
		return
	}

	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="webfs@%s"`, a.realm))
	if !a.checkAndLimit(w, r) {
		w.WriteHeader(401)
		return
	}
	a.origin.ServeHTTP(w, r)
}

func (a *authRateByIPHandler) checkAndLimit(w http.ResponseWriter, r *http.Request) (result bool) {

	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(s) != 2 {
		return
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return
	}

	result = (pair[0] == a.username && pair[1] == a.password)

	if !result {
		ip := fetchIP(r.RemoteAddr)
		a.mu.Lock()
		rec, ok := a.record[ip]
		if !ok {
			rec = 0
		}
		a.record[ip] = rec + 1
		a.mu.Unlock()
	}
	return result
}
