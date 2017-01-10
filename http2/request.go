package http2

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func NewRequest(method, url string, params url.Values) (*http.Request, error) {
	var body io.Reader
	body = strings.NewReader(params.Encode())
	return http.NewRequest(method, url, body)
}

func DoRequest(c *http.Client, req *http.Request) (*http.Response, []byte, error) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	rep, err := c.Do(req)
	if err != nil {
		return rep, nil, err
	}
	defer rep.Body.Close()

	repBody, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return rep, nil, err
	}
	return rep, repBody, err
}

func Request(method, url string, params url.Values) ([]byte, error) {
	req, err := NewRequest(method, url, params)
	if err != nil {
		return nil, err
	}
	_, repBody, err := DoRequest(http.DefaultClient, req)
	return repBody, err
}

func DoJSONRequest(c *http.Client, req *http.Request, result interface{}) (*http.Response, error) {
	rep, repBody, err := DoRequest(c, req)
	if err != nil {
		return rep, err
	}
	err = json.Unmarshal(repBody, &result)
	return rep, err
}

func JSONRequest(method, url string, params url.Values) (result map[string]interface{}, err error) {
	req, err := NewRequest(method, url, params)
	if err != nil {
		return nil, err
	}
	_, err = DoJSONRequest(http.DefaultClient, req, &result)
	return result, err
}