package toolbox

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
)

type RequestMethod byte

const (
	HEAD RequestMethod = iota
	GET
	POST
	PUT
	DELETE
)

func (method RequestMethod) String() string {
	switch method {
	case HEAD:
		return "HEAD"
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	default:
		return fmt.Sprintf("Unknown(%d)", method)
	}
}

func JoinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

func ParseURL(someURL string) string {
	u, err := url.Parse(someURL)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("\nScheme: %s\n  Host: %s\n  Path: %s\n Query: %s\n",
		u.Scheme, u.Hostname(), u.Path, u.RawQuery)
}

func BaseURL(url string) string {
	idx := strings.LastIndex(url, "/")
	if idx != -1 {
		url = url[:idx]
	}
	return url
}

func SendRequest(method RequestMethod, url, body string, headers map[string]string) ([]byte, error) {
	var b []byte

	// create a request
	req, err := http.NewRequest(method.String(), url, strings.NewReader(body))
	if err != nil {
		return b, err
	}

	// NOTE this !! -You need to set Req.Close to true (the defer on resp.Body.Close() syntax used in the examples is not enough)
	req.Close = true

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return b, err
	}

	respDump, _ := httputil.DumpResponse(resp, true)

	defer func(Body io.ReadCloser) {
		// Ignore error explicitly
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode > 300 {
		reqStr := strings.TrimSpace(string(reqDump))
		respStr := strings.TrimSpace(string(respDump))
		message := fmt.Sprintf("Non 200 status code: %d REQUEST: %s RESPONSE: %s", resp.StatusCode, reqStr, respStr)
		err = errors.New(message)
		return b, err
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	return b, nil
}
