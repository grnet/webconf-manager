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
