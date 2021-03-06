package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func GetAllNodesFromDb() Nodes {

	db, err := sql.Open("sqlite3", cfg.Sql.Path )
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		fmt.Println("Failed to keep connection alive")
	}

	rows, err := db.Query("select * from hosts")
	if err != nil {
		panic(err)
	}

	var nodes Nodes
	for rows.Next() {
		var node Node
		err = rows.Scan(&node.Id, &node.Type, &node.Name, &node.StoragePath, &node.InternalIP)
		nodes = append(nodes, node)
	}
	db.Close()

	return nodes

}

func QueryTypeFromDb(typeFilter string) Nodes {

	db, err := sql.Open("sqlite3", cfg.Sql.Path )
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		fmt.Println("Failed to keep connection alive")
	}

	rows, err := db.Query("select * from hosts where type=?", typeFilter)
	if err != nil {
		panic(err)
	}

	var nodes Nodes
	for rows.Next() {
		var node Node
		err = rows.Scan(&node.Id, &node.Type, &node.Name, &node.StoragePath, &node.InternalIP)
		nodes = append(nodes, node)
	}
	db.Close()

	return nodes

}

func AddNodeToDb(node Node) Node {

	db, err := sql.Open("sqlite3", cfg.Sql.Path )
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		fmt.Println("Failed to keep connection alive")
	}

	stmt, err := db.Prepare("insert into hosts (type, hostname, pair_id, internal_ip) values (?,?,?,?)")
	if err != nil {
		panic(err)
	}
	
	res , err := stmt.Exec(node.Type, node.Name, node.StoragePath, node.InternalIP)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	node.Id = id

	db.Close()

	return node

}
