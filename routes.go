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
	Route{"Index"    , "GET"  , "/list"            , Index     },
	Route{"Create"   , "POST" , "/add"             , Create    },
	Route{"LeastLoad", "GET"  , "/nextnode/{group}", LeastLoad },
	Route{"Deploy"   , "POST" , "/deploy"          , Deploy    },
}
