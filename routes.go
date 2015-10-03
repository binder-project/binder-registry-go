package main

import "net/http"

// Route is a simple HTTP Method, Pattern, and Handler pairing
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TemplateIndex",
		"GET",
		"/templates",
		TemplateIndex,
	},
	Route{
		"TemplateShow",
		"GET",
		"/templates/{templateName}",
		TemplateShow,
	},
	Route{
		"TemplateCreate",
		"POST",
		"/templates",
		TemplateCreate,
	},
}
