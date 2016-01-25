//go:generate hel --package ./...
package bdd_talk

import "net/http"

type ResponseStorer interface {
	Store(resp *http.Response, index int) error
}

type Client struct {
	serverURL string
	index     int

	respStorer ResponseStorer
}

func NewClient(URL string, respStorer ResponseStorer) *Client {
	return &Client{
		serverURL:  URL,
		respStorer: respStorer,
	}
}

func (c *Client) HitServer() error {
	resp, err := http.Get(c.serverURL)
	if err != nil {
		return err
	}

	c.index++
	return c.respStorer.Store(resp, c.index)
}
