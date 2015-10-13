package main

import (
	"log"
	"net/http"
	"strconv"
)

var cfg = LoadConfiguration()

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe( cfg.Server.Bindip + ":" + strconv.Itoa(cfg.Server.Port) , router))
}
