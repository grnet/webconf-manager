package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"Index" ,"GET" ,"/",Index },
	Route{"Create","POST","/",Create},
	Route{"LeastLoad", "GET", "/nextnode/{group}", LeastLoad },
//	Route{
//		"RemoveNode",
//		"DELETE",
//		"/nodes",
//		RemoveNode,
//	},
}
