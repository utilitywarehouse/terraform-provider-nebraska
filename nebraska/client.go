package nebraska

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// FlatcarApplicationID is the id that the default Flatcar application is
	// created with
	FlatcarApplicationID = "e96281a6-d1af-4bde-9a0a-97b76e56dc57"
)

var (
	// ErrNotFound is returned when the client receives a 404 from Nebraska
	ErrNotFound = fmt.Errorf("nebraska: not found")
)

// Client communicates with a Nebraska server
type Client struct {
	BaseURL string

	c         *http.Client
	userAgent string
}

// New returns a new client for the given Nebraska server URL
func New(baseURL, userAgent string) *Client {
	return &Client{
		BaseURL:   baseURL,
		c:         &http.Client{},
		userAgent: userAgent,
	}

}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.BaseURL, path), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (c *Client) do(req *http.Request, data interface{}) error {
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("Couldn't retrieve response from request")
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var body []byte
		if resp.Body != nil {
			body, _ = ioutil.ReadAll(resp.Body)
		}
		return fmt.Errorf("Bad response: req_uri=%s, response_code=%d, response=%s", req.URL.String(), resp.StatusCode, string(body))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &data); err != nil {
			return err
		}
	}

	return nil
}
