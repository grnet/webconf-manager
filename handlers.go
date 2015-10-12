package main

import (
	"encoding/json"
//	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	nodes := GetNodesFromDb()
	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		panic(err)
	}

}

func Create(w http.ResponseWriter, r *http.Request) {
	var nodes Nodes
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &nodes); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	for _, node := range nodes{
		if len(node.Name) == 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(jsonError{Code: http.StatusBadRequest, Message: "Name not provided in json"}); err != nil {
				panic(err)
			}
			return
		}
		if len(node.Type) == 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(jsonError{Code: http.StatusBadRequest, Message: "Type not provided in json"}); err != nil {
				panic(err)
			}
			return
		}
		if node.Type != "bigbluebutton" {
			if node.Type != "transcoding" {
				if node.Type != "testing" {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusBadRequest)
					if err := json.NewEncoder(w).Encode(jsonError{Code: http.StatusBadRequest, Message: "Type value not supported. Select one of {bigbluebutton, transcoding}"}); err != nil {
						panic(err)
					}
					return
				}
			}
		}
		if node.StoragePath == 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(jsonError{Code: http.StatusBadRequest, Message: "Pair id not provided in json or provided but equal is to zero"}); err != nil {
				panic(err)
			}
			return
		}
		if len(node.InternalIP) == 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(jsonError{Code: http.StatusBadRequest, Message: "internal_ip not provided in json"}); err != nil {
				panic(err)
			}
			return
		}
	}

	var created_nodes Nodes
	for _, node := range nodes{
		new_node := AddNodeToDb(node)
		created_nodes = append(created_nodes, new_node)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created_nodes); err != nil {
		panic(err)
	}

}

