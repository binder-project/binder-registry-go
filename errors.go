package main

type jsonErr struct {
	// Assumed to be e.g. http.StatusNotFound
	Code int    `json:"code"`
	Text string `json:"text"`
}
