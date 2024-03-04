package main

import (
	"github.com/gorilla/mux"
	"github.com/wintergathering/daren2/frstr"
	"github.com/wintergathering/daren2/server"
)

const listenAddr = ":8080"

func main() {

	client := frstr.NewFirestoreClient()

	defer client.Close()

	ds := frstr.NewDareService(client)
	r := mux.NewRouter()

	s := server.NewServer(r, ds, listenAddr)

	s.Run()
}
