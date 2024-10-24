package client

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client interface {
	Add(key, value string, ttl time.Duration) error
	Set(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

// client implementation over HTTP
type HTTPClient struct {
	http http.Client
}

func NewHTTP() *HTTPClient {
	return &HTTPClient{}
}

func (c *HTTPClient) Add(key, value string, ttl time.Duration) error {
	req, err := c.createAddSetRequest("POST", "http://localhost:8080/api/v1/add", key, value, ttl)

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(res)

	return err
}

func (c *HTTPClient) Set(key, value string, ttl time.Duration) error {
	req, err := c.createAddSetRequest("PUT", "http://localhost:8080/api/v1/set", key, value, ttl)

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(res)

	return err
}

func (c *HTTPClient) Get(key string) (string, error) {
	req, err := c.createGetDelRequest("GET", "http://localhost:8080/api/v1/get", key)
	if err != nil {
		return "", err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}

	return c.handleResponse(res)
}

func (c *HTTPClient) Del(key string) error {
	req, err := c.createGetDelRequest("DELETE", "http://localhost:8080/api/v1/del", key)
	if err != nil {
		return err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(res)

	return err
}

func (c *HTTPClient) handleResponse(res *http.Response) (msg string, err error) {
	// read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return msg, err
	}

	msg = string(body)

	// check for server-side handled errors
	if res.Status != "200 OK" {
		return msg, errors.New(msg)
	}

	return msg, err
}

func (c *HTTPClient) createAddSetRequest(method, url, key, value string, ttl time.Duration) (*http.Request, error) {
	var expirationTime string
	if ttl == 0 {
		expirationTime = ""
	} else {
		expirationTime = time.Now().Add(ttl).Format(time.RFC3339)
	}

	postForm := c.createPostForm(key, value, expirationTime)
	req, err := http.NewRequest(method, url, strings.NewReader(postForm.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (c *HTTPClient) createPostForm(key, value, expiresAt string) url.Values {
	form := url.Values{}
	form.Add("key", key)
	form.Add("value", value)
	form.Add("expires_at", expiresAt)

	return form
}

func (c *HTTPClient) createGetDelRequest(method, url, key string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	queryParams.Add("key", key)
	req.URL.RawQuery = queryParams.Encode()

	return req, nil
}

// aggregate over client and parser
type ClientApp struct {
	cl Client
	p  Parser
}

func New(client Client) *ClientApp {
	return &ClientApp{
		cl: client,
		p:  Parser{},
	}
}

// run the entire client app
func (app ClientApp) Run() {
	// receiving user input
	op, key, val, ttl, err := app.p.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var res string

	// mapping operation to its corresponding handler
	switch *op {
	case "add":
		err = app.cl.Add(*key, *val, *ttl)
	case "set":
		err = app.cl.Set(*key, *val, *ttl)
	case "get":
		res, err = app.cl.Get(*key)
	case "del":
		err = app.cl.Del(*key)
	default:
		err = ErrOpUnsupported
	}

	if err != nil {
		fmt.Print("Error: ", err)
	} else {
		if res == "" {
			fmt.Println("Success!")
		} else {
			fmt.Println(res)
		}
	}
}

// incoming flag params parser
type Parser struct{}

func (p Parser) Parse() (op, key, val *string, ttl *time.Duration, err error) {
	op = flag.String("op", "", "operation to be executed")
	key = flag.String("key", "", "key to be inserted")
	val = flag.String("val", "", "value to be paired with the key")
	ttl = flag.Duration("ttl", 0, "key's time to live in the object storage")

	flag.Parse()

	if *key == "" {
		err = ErrEmptyKey
	}

	if *op == "" {
		err = ErrOpNotProvided
	}

	return op, key, val, ttl, err
}
