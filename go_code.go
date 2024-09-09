package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Handler for the / route
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "Welcome to the Go Microservice!")
}

// Handler for the /healthz route
func healthz(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/healthz" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Service is healthy")
}

// Handler for the /mockoon route
func mockoon() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mockoon" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		resp, err := http.Get(os.Getenv("MOCKURL"))
		if err != nil {
			http.Error(w, "Failed to get response from Mockoon", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Querying Mockoon API from ENV VAR:", os.Getenv("MOCKURL"))
		fmt.Fprintln(w, string(body))
	})
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/healthz", healthz)
	http.Handle("/mockoon", mockoon())

	log.Fatal(http.ListenAndServe(":8090", nil))
}