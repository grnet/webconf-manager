package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bsphere/celery"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

func Index(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	nodes := GetAllNodesFromDb()
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
	for _, node := range nodes {
		if len(node.Name) == 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(jsonResponse{Code: http.StatusBadRequest, Message: "Name not provided in json"}); err != nil {
				panic(err)
			}
			return
		}
		if len(node.Type) == 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(jsonResponse{Code: http.StatusBadRequest, Message: "Type not provided in json"}); err != nil {
				panic(err)
			}
			return
		}
		if node.Type != "bigbluebutton" {
			if node.Type != "transcoding" {
				if node.Type != "testing" {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusBadRequest)
					if err := json.NewEncoder(w).Encode(jsonResponse{Code: http.StatusBadRequest, Message: "Type value not supported. Select one of {bigbluebutton, transcoding}"}); err != nil {
						panic(err)
					}
					return
				}
			}
		}
		if node.StoragePath == 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(jsonResponse{Code: http.StatusBadRequest, Message: "Pair id not provided in json or provided but is equal to zero"}); err != nil {
				panic(err)
			}
			return
		}
		if len(node.InternalIP) == 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(jsonResponse{Code: http.StatusBadRequest, Message: "internal_ip not provided in json"}); err != nil {
				panic(err)
			}
			return
		}
	}

	var created_nodes Nodes
	for _, node := range nodes {
		new_node := AddNodeToDb(node)
		created_nodes = append(created_nodes, new_node)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created_nodes); err != nil {
		panic(err)
	}

}

func LeastLoad(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	group := vars["group"]
	nodes := QueryTypeFromDb(group)
	if len(nodes) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonResponse{Code: http.StatusNotFound, Message: "No nodes with given type found"}); err != nil {
			panic(err)
		}
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := cfg.MonHost.Method + "://" + cfg.MonHost.IP + "/icingaweb2/monitoring/list/services?modifyFilter=1&service=load&format=json"
	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(cfg.MonHost.Username, cfg.MonHost.Password)

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	var AllOuts icingaOuts
	if err := json.Unmarshal(body, &AllOuts); err != nil {
		panic(err)
	}

	var validOuts icingaOuts

	for _, out := range AllOuts {
		for _, node := range nodes {
			if node.Name == out.HostName {
				tokenized := strings.Split(out.ServicePerfdata, " ")
				loads1 := strings.Split(tokenized[0], "=")
				loads5 := strings.Split(tokenized[1], "=")
				loads15 := strings.Split(tokenized[2], "=")
				load1 := strings.Split(loads1[1], ";")
				load5 := strings.Split(loads5[1], ";")
				load15 := strings.Split(loads15[1], ";")
				out.Load1, _ = strconv.ParseFloat(load1[0], 64)
				out.Load5, _ = strconv.ParseFloat(load5[0], 64)
				out.Load15, _ = strconv.ParseFloat(load15[0], 64)
				validOuts = append(validOuts, out)
			}
		}
	}

	var nextnode icingaOut
	var out icingaOut

	if len(validOuts) > 0 {
		nextnode = validOuts[0]
		load := nextnode.Load1
		for _, out = range validOuts {
			if out.Load1 < load {
				nextnode = out
				load = out.Load1
			}
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, nextnode.HostName)
	}

	return

}

func Deploy(w http.ResponseWriter, r *http.Request) {

	conn, err := amqp.Dial("amqp://localhost://")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	var args []string
	args = append(args, "deploy.yml")

	task, err := celery.NewTask("tasks.run_playbook", args, nil)
	if err != nil {
		panic(err)
	}

	err = task.Publish(ch, "", "celery")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(jsonResponse{Code: http.StatusAccepted, Message: "Request to deploy via Ansible has been accepted"}); err != nil {
		panic(err)
	}

}
