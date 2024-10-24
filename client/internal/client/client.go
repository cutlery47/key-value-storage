package client

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	Add(key, value string, ttl time.Duration) error
	Set(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

type HTTPClient struct {
	http http.Client
}

func NewHTTP() *HTTPClient {
	return &HTTPClient{}
}

func (c *HTTPClient) Add(key, value string, ttl time.Duration) error {
	if len(value) == 0 {
		return ErrEmptyVal
	}

	// calculating expiration time
	expiresAt := time.Now().Add(ttl)

	req, err := http.NewRequest("POST", "http:/localhost:8080/api/v1/add", nil)
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Add("key", key)
	form.Add("value", value)
	form.Add("expires_at", expiresAt.String())
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = c.http.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *HTTPClient) Set(key, value string, ttl time.Duration) error {
	if len(value) == 0 {
		return ErrEmptyVal
	}

	// calculating expiration time
	expiresAt := time.Now().Add(ttl)

	req, err := http.NewRequest("POST", "http:/localhost:8080/api/v1/add", nil)
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Add("key", key)
	form.Add("value", value)
	form.Add("expires_at", expiresAt.String())
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = c.http.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *HTTPClient) Get(key string) (string, error) {
	req, err := http.NewRequest("GET", "http:/localhost:8080/api/v1/add", nil)
	if err != nil {
		return "", err
	}

	req.URL.Query().Add("key", key)

	_, err = c.http.Do(req)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *HTTPClient) Del(key string) error {

	req, err := http.NewRequest("DELETE", "http:/localhost:8080/api/v1/add", nil)
	if err != nil {
		return err
	}

	req.URL.Query().Add("key", key)

	_, err = c.http.Do(req)
	if err != nil {
		return err
	}

	return nil
}

type ClientApp struct {
	cl Client
	p  parser
}

func New(client Client) *ClientApp {
	return &ClientApp{
		cl: client,
		p:  parser{},
	}
}

func (app ClientApp) Run() {
	op, key, val, ttl, err := app.p.parse()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var res string

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
		fmt.Println("Error:", err)
	} else {
		if res == "" {
			fmt.Println("Success!")
		} else {
			fmt.Println(res)
		}
	}

}

type parser struct{}

func (p parser) parse() (op, key, val *string, ttl *time.Duration, err error) {
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
