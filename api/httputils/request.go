package httputils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	neturl "net/url"
	"reflect"
	"time"
)

type FD map[string]string

type HttpRequest struct {
	c http.Client
}

// NewRequest
func NewRequest(timeout time.Duration, keepAlive time.Duration, handshakeTimeout time.Duration) *HttpRequest {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: keepAlive,
		}).Dial,
		TLSHandshakeTimeout: handshakeTimeout,
	}
	c := http.Client{Transport: transport}

	return &HttpRequest{
		c: c,
	}
}

// SendForm
func (r *HttpRequest) SendForm(url string, method string, data FD) (int, []byte, error) {
	form := neturl.Values{}
	for key, value := range data {
		form.Add(key, value)
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, nil, err
	}
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := r.c.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return 0, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, body, nil
}

// SendJSON
func (r *HttpRequest) SendJSON(url string, method string, data interface{}, reply interface{}) (int, interface{}, error) {
	rawData, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}

	req, err := http.NewRequest(
		method, url, bytes.NewBuffer(rawData),
	)

	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return 0, nil, err
	}

	resp, err := r.c.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return 0, nil, err
	}

	var result interface{}
	result = reflect.New(reflect.TypeOf(reply)).Interface()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, result, nil
}
