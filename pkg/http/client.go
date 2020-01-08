package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)


type Client struct {
	client *http.Client
	request *http.Request
	response *http.Response
	Body []byte
}

func New() *Client {
	return &Client{
		client:  &http.Client{},
	}
}

func (c *Client)SetRequest (method, uri string) *Client {
	req, err := http.NewRequest(method, uri, &bytes.Buffer{})
	if err != nil {
		return nil
	}

	c.Body = []byte{}
	c.request = req

	return c
}

func (c *Client) Do () error {
	resp, err := c.client.Do(c.request)
	if err != nil {
		return err
	}

	c.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UnmarshalJSON(data interface{}) error {
	err := json.Unmarshal(c.Body, data)
	if err != nil {
		return err
	}

	return nil
}