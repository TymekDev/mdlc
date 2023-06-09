package main

type Link struct {
	Filename    string `json:"filename"`
	Destination string `json:"destination"`
	Count       int    `json:"count"`
	StatusCode  int    `json:"status_code"`
	ErrMsg      string `json:"error,omitempty"`
}
