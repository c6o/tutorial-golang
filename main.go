package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("x-c6o-userid")

		log.Printf("Request from %s, userID=%s", r.RemoteAddr, userID)

		req, err := http.NewRequest("GET", "http://service-c:8080", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		req.Header.Add("x-c6o-userid", userID)

		client := &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if len(body) > 20 {
			body = body[:20]
		}
		w.Write(body)

		log.Printf("Sent response to %s, userID=%s: %s", r.RemoteAddr, userID, body)
	})

	log.Println("Starting service-b on port 8080")

	http.ListenAndServe("127.0.0.1:8080", nil)
}
