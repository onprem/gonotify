package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Client represents the GoNotify API client
type Client struct {
	base  *url.URL
	token string
	hc    *http.Client
	ep    struct {
		login    *url.URL
		register *url.URL
		send     *url.URL
	}
}

// NewClient return an instance of Client
func NewClient(baseURL string, token string) (*Client, error) {
	base, err := url.Parse(baseURL + "/api/v1/")
	if err != nil {
		return nil, err
	}

	c := &Client{
		base:  base,
		token: token,
		hc:    &http.Client{},
	}
	err = c.register()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) register() error {
	u, err := url.Parse("login")
	if err != nil {
		return err
	}
	c.ep.login = c.base.ResolveReference(u)

	u, err = url.Parse("register")
	if err != nil {
		return err
	}
	c.ep.register = c.base.ResolveReference(u)

	u, err = url.Parse("send")
	if err != nil {
		return err
	}
	c.ep.send = c.base.ResolveReference(u)

	return nil
}

// Login returns a token given valid credentials
func (c *Client) Login(phone, password string) (string, error) {
	var token string

	if phone == "" || password == "" {
		return token, errors.New("Number and password cannot be empty")
	}

	values := map[string]string{
		"phone":    phone,
		"password": password,
	}

	v, err := json.Marshal(values)
	if err != nil {
		return token, err
	}

	req, err := http.NewRequest("POST", c.ep.login.String(), bytes.NewBuffer(v))
	if err != nil {
		return token, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()

	res := map[string]string{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return token, err
	}

	msg, ok := res["error"]
	if ok {
		return token, errors.New(msg)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return token, errors.New("Request failed. Status: " + resp.Status)
	}

	token = res["token"]
	return token, nil
}

// Send sends a message to given group
func (c *Client) Send(body, group string) error {
	if c.token == "" {
		return errors.New("You are not logged in")
	}

	if body == "" {
		return errors.New("Cannot send empty message")
	}

	values := map[string]string{
		"body":  body,
		"group": group,
	}

	v, err := json.Marshal(values)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.ep.send.String(), bytes.NewBuffer(v))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	res := map[string]string{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}

	errMsg, ok := res["error"]
	if ok {
		return errors.New(errMsg)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("Request failed. Status: " + resp.Status)
	}

	return nil
}
