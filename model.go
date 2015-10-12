package main

type Node struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	StoragePath int    `json:"storage_path"`
	InternalIP  string `json:"internal_ip"`
}

type Nodes []Node

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type icingaOut struct {
	HostName           string  `json:"host_name"`
	HostState          string  `json:"host_state"`
	ServiceDescription string  `json:"service_description"`
	ServiceState       string  `json:"service_state"`
	ServiceOutput      string  `json:"service_output"`
	ServicePerfdata    string  `json:"service_perfdata"`
	Load1              float64 `json:"load1"`
	Load5              float64 `json:"load5"`
	Load15             float64 `json:"load15"`
}

type icingaOuts []icingaOut
