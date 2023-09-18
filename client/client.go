package client

import "net/http"

type Client struct {
	http_client *http.Client
}

func New() *Client {
	return &Client{
		http_client: &http.Client{},
	}
}
