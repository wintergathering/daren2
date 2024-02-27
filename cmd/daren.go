package main

import "net/http"

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	w.Write([]byte("Hello, World"))
}

func main() {
	s := &http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/", handleIndex)

	s.ListenAndServe()

}
